package types

import (
	"fmt"
	"game-server/common/types"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GameState struct {
	Board         [3][3]string `json:"board"`
	CurrentPlayer string       `json:"current_player"`
	Winner        string       `json:"winner"`
	IsDraw        bool         `json:"is_draw"`
}

type TicTacMemoRoom struct {
	types.Room
	GameState   GameState
	CurrentTurn string
}

// CreateRoom initializes a TicTacMemoRoom with a given maxPlayers and roomID.
func CreateRoom(maxPlayers int, roomID uuid.UUID, players []*Player) *TicTacMemoRoom {
	// Empty 3x3 board for TicTacToe game
	board := [3][3]string{
		{"", "", ""},
		{"", "", ""},
		{"", "", ""},
	}

	return &TicTacMemoRoom{
		Room: types.Room{
			ID:         roomID,
			MaxPlayers: maxPlayers,
			Players:    mapPlayersToUsers(players),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		GameState: GameState{
			Board:         board,
			CurrentPlayer: "X", // X always starts
			Winner:        "",
			IsDraw:        false,
		},
		CurrentTurn: "X",
	}
}

func mapPlayersToUsers(playerList []*Player) []*types.User {
	var userList []*types.User
	for _, player := range playerList {
		user := player.User
		userList = append(userList, &user)
	}
	return userList
}

// StartGame starts the TicTacToe game, resetting the board and setting the first turn.
func (room *TicTacMemoRoom) StartGame() {
	fmt.Printf("Starting the game in Room %d with %d players\n", room.ID, len(room.Players))
	room.GameState = GameState{
		Board:         [3][3]string{},
		CurrentPlayer: "X",
		Winner:        "",
		IsDraw:        false,
	}
	room.CurrentTurn = "X"
}

// MakeMove processes the move and updates the game state, checking for winners or draw.
func (room *TicTacMemoRoom) MakeMove(db *gorm.DB, makeMoveData MakeMoveData) {
	// posX, posY := makeMoveData.PosX, makeMoveData.PosY

	// if room.GameState.Board[posX][posY] != "" {
	// 	log.Printf("Invalid move at position (%d, %d). Spot already taken.", posX, posY)
	// 	return
	// }

	// // Mark the move on the board
	// room.GameState.Board[posX][posY] = room.CurrentTurn
	// log.Printf("Move made by player %s at position (%d, %d)", room.CurrentTurn, posX, posY)

	// // Check for win or draw
	// if room.checkWin() {
	// 	room.GameState.Winner = room.CurrentTurn
	// } else if room.checkDraw() {
	// 	room.GameState.IsDraw = true
	// } else {
	// 	// Switch turn
	// 	if room.CurrentTurn == "X" {
	// 		room.CurrentTurn = "O"
	// 	} else {
	// 		room.CurrentTurn = "X"
	// 	}
	// 	room.GameState.CurrentPlayer = room.CurrentTurn
	// }

	// Log the current state
	log.Printf("Current GameState: %+v", room.GameState)
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
