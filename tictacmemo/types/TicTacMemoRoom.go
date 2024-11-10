package types

import (
	"encoding/json"
	"fmt"
	"game-server/common/types"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type GameState struct {
	Board         [3][3]string `json:"board"`         // Actual board with X's and O's
	VisibleBoard  [3][3]string `json:"visible_board"` // Board with hidden moves for users
	CurrentPlayer string       `json:"current_player"`
	Winner        string       `json:"winner"`
	IsDraw        bool         `json:"is_draw"`
}

type TicTacMemoRoom struct {
	types.Room
	GameState         GameState
	CurrentTurn       string
	PlayerIDs         map[string]string // Mapping of X or O to playerID
	OpponentPlayerIDs map[string]string // Mapping of X's opponent to X or O's opponent to O to playerID

}

// CreateRoom initializes a TicTacMemoRoom with a given maxPlayers and roomID.
func CreateRoom(maxPlayers int, roomID uuid.UUID, playerXID string, playerOID string) *TicTacMemoRoom {
	// Empty 3x3 board for TicTacToe game
	board := [3][3]string{
		{"", "", ""},
		{"", "", ""},
		{"", "", ""},
	}
	visibleBoard := [3][3]string{
		{"", "", ""},
		{"", "", ""},
		{"", "", ""},
	}

	return &TicTacMemoRoom{
		Room: types.Room{
			ID:         roomID,
			MaxPlayers: maxPlayers,
			Players:    []*types.User{},
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		GameState: GameState{
			Board:         board,
			VisibleBoard:  visibleBoard,
			CurrentPlayer: "X", // X always starts
			Winner:        "",
			IsDraw:        false,
		},
		CurrentTurn: "X",
		PlayerIDs: map[string]string{
			"X": playerXID,
			"O": playerOID,
		},
		OpponentPlayerIDs: map[string]string{
			"X": playerOID,
			"O": playerXID,
		},
	}
}

// StartGame starts the TicTacToe game, resetting the board and setting the first turn.
func (room *TicTacMemoRoom) StartGame() {
	fmt.Printf("Starting the game in Room %d with %d players\n", room.ID, len(room.Players))
	room.GameState = GameState{
		Board:         [3][3]string{},
		VisibleBoard:  [3][3]string{},
		CurrentPlayer: "X",
		Winner:        "",
		IsDraw:        false,
	}
	room.CurrentTurn = "X"

	gameStart := gin.H{
		"event": "start-game",
		"data":  room.GameState,
	}

	room.BroadcastGameState(gameStart)
}

// MakeMove processes the move and updates the game state, checking for winners or draw.
func (room *TicTacMemoRoom) MakeMove(db *gorm.DB, makeMoveData MakeMoveData, playerID string) {
	posX, posY := makeMoveData.PosX, makeMoveData.PosY

	// Validate the current player
	if room.PlayerIDs[room.CurrentTurn] != playerID {
		log.Printf("Invalid move by player %s. Not your turn!", playerID)
		return
	}

	// Validate if the move is within the bounds of the board
	if posX < 0 || posX > 2 || posY < 0 || posY > 2 {
		log.Printf("Invalid move at position (%d, %d). Out of bounds.", posX, posY)
		return
	}

	// Validate if the position is already taken
	if room.GameState.Board[posX][posY] != "" {
		log.Printf("Invalid move at position (%d, %d). Spot already taken.", posX, posY)
		return
	}

	// Mark the move on the actual board
	room.GameState.Board[posX][posY] = room.CurrentTurn

	// Update the visible board (hide the actual move)
	room.GameState.VisibleBoard[posX][posY] = "?"

	log.Printf("Move made by player %s at position (%d, %d)", room.CurrentTurn, posX, posY)

	// Check for win or draw
	if room.checkWin() {
		room.GameState.Winner = room.CurrentTurn
	} else if room.checkDraw() {
		room.GameState.IsDraw = true
	} else {
		// Switch turn
		if room.CurrentTurn == "X" {
			room.CurrentTurn = "O"
		} else {
			room.CurrentTurn = "X"
		}
		room.GameState.CurrentPlayer = room.CurrentTurn
	}

	// Log the current state
	log.Printf("Current GameState: %+v", room.GameState)
}

const WEIGHT_MOVES = 2
const WEIGHT_TIME_ELAPSED = 2
const MAX_SCORING_TIME = 30

// MakeMove processes the move and updates the game state, checking for winners or draw.
func (room *TicTacMemoRoom) UpdateScore(db *gorm.DB, updateScoreData UpdateScoreData) {
	if room.GameState.Winner != "" {
		var userId = room.PlayerIDs[updateScoreData.AssignedLabel]
		var user types.User
		if err := db.Where("user_id = ?", userId).First(&user).Error; err != nil {
			log.Fatal("User does not existes", err)
		}

		var opponentUserId = room.OpponentPlayerIDs[updateScoreData.AssignedLabel]
		var opponentUser types.User
		if err := db.Where("user_id = ?", opponentUserId).First(&opponentUser).Error; err != nil {
			log.Fatal("User does not existes", err)
		}

		if room.GameState.IsDraw {
			// no change
			// Bind incoming request parameters
			var game GameHistory = GameHistory{
				UserID:           userId,
				OpponentUserID:   opponentUser.UserID,
				OpponentUsername: opponentUser.Username,
				RatingChange:     0,
			}
			db.Create(&game)
			return
		}

		var ratingChange int = room.calculateRating(updateScoreData)
		user.Rating = user.Rating + ratingChange

		log.Printf("Updated Score: %s, Time: %d, Moves: %d RatingChange: %d",
			room.GameState.Winner, updateScoreData.ElapsedTime, updateScoreData.MoveCount, ratingChange)

		// Bind incoming request parameters
		var game GameHistory = GameHistory{
			UserID:           userId,
			OpponentUserID:   opponentUser.UserID,
			OpponentUsername: opponentUser.Username,
			RatingChange:     ratingChange,
		}
		db.Create(&game)

		db.Save(&user)
	}
}

func (room *TicTacMemoRoom) calculateRating(updateScoreData UpdateScoreData) int {
	var ratingChangeTime = (MAX_SCORING_TIME - updateScoreData.ElapsedTime) * WEIGHT_TIME_ELAPSED
	var ratingChangeMoves = updateScoreData.MoveCount * WEIGHT_MOVES
	var ratingChange = (ratingChangeTime + ratingChangeMoves)

	if updateScoreData.AssignedLabel != room.GameState.Winner {
		ratingChange = -1 * ratingChange // looser
	}

	return ratingChange
}

// checkWin checks if the current player has won.
func (room *TicTacMemoRoom) checkWin() bool {
	board := room.GameState.Board
	// Check rows, columns, and diagonals for a win
	for i := 0; i < 3; i++ {
		if board[i][0] == room.CurrentTurn && board[i][1] == room.CurrentTurn && board[i][2] == room.CurrentTurn {
			return true
		}
		if board[0][i] == room.CurrentTurn && board[1][i] == room.CurrentTurn && board[2][i] == room.CurrentTurn {
			return true
		}
	}
	if board[0][0] == room.CurrentTurn && board[1][1] == room.CurrentTurn && board[2][2] == room.CurrentTurn {
		return true
	}
	if board[0][2] == room.CurrentTurn && board[1][1] == room.CurrentTurn && board[2][0] == room.CurrentTurn {
		return true
	}
	return false
}

// checkDraw checks if the board is full and there is no winner, resulting in a draw.
func (room *TicTacMemoRoom) checkDraw() bool {
	board := room.GameState.Board
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == "" {
				return false
			}
		}
	}
	return true
}

func (room *TicTacMemoRoom) BroadcastGameState(data any) {
	json, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling game state:", err)
		return
	}

	// Iterate through the list of clients and send the message
	for i := 0; i < len(room.Clients); i++ {
		client := room.Clients[i]
		if client == nil {
			continue
		}
		err := client.WriteMessage(websocket.TextMessage, json)
		if err != nil {
			log.Println("Error sending game state to client:", err)
			/*client.Close()
			// Remove the client from the list if sending fails
			room.Clients = append(room.Clients[:i], room.Clients[i+1:]...)
			i-- // Adjust index after removal*/
		}
	}
}
