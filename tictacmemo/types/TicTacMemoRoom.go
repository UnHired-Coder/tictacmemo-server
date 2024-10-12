package types

import (
	"fmt"
	"game-server/common/types"
	"time"
)

// TicTacMemoRoom struct extends Room and implements Tic-Tac-Toe-specific logic
type TicTacMemoRoom struct {
	types.Room                     // Embedding Room to reuse the generic player management
	Board       [3][3]int          // Tic-Tac-Toe-specific board
	CurrentTurn int                // Tracks which player’s turn it is
	GameManager *types.GameManager // Reference to GameManager for room removal
}

// NewTicTacMemoRoom creates a new TicTacMemoRoom
func NewTicTacMemoRoom(int) *TicTacMemoRoom {
	return &TicTacMemoRoom{
		Room: *types.CreateRoom(2), // Only two players for Tic-Tac-Toe
	}
}

// StartGame method for TicTacMemoRoom, initializes the game board and sets the current turn
func (r *TicTacMemoRoom) StartGame() {
	// Call the generic StartGame to notify room is full
	r.Room.StartGame()

	// Initialize the board and set Player 1 to start
	fmt.Println("TicTacToe game starting...")
	r.Board = [3][3]int{} // Reset the board
	r.CurrentTurn = 1     // Player 1 starts
}

// MakeMove allows players to make a move on the Tic-Tac-Toe board
func (r *TicTacMemoRoom) MakeMove(player *types.User, x, y int) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	// Check if the game has ended
	if len(r.Players) != 2 {
		return fmt.Errorf("the game hasn't started yet")
	}

	if r.Board[x][y] != 0 {
		return fmt.Errorf("invalid move, position already taken")
	}

	// Check if it's the player's turn
	currentPlayerID := r.Players[r.CurrentTurn-1].ID
	if player.ID != currentPlayerID {
		return fmt.Errorf("it's not your turn")
	}

	// Make the move
	r.Board[x][y] = r.CurrentTurn
	fmt.Printf("Player %d (%s) made a move at (%d, %d)\n", r.CurrentTurn, player.Username, x, y)

	// Switch turn
	r.switchTurn()

	// You can add a winner check here if necessary

	r.UpdatedAt = time.Now()
	return nil
}

// Method to switch turns
func (r *TicTacMemoRoom) switchTurn() {
	if r.CurrentTurn == 1 {
		r.CurrentTurn = 2
	} else {
		r.CurrentTurn = 1
	}
}

func (r *TicTacMemoRoom) endGame() {
	// Call GameManager to remove the room
	r.GameManager.RemoveRoom(r.ID)
}

// CreateTicTacMemoRoom creates a TicTacMemoRoom for a two-player game
// func CreateTicTacMemoRoom(gm *types.GameManager, roomID int) *TicTacMemoRoom {
// 	gm.lock.Lock()
// 	defer gm.lock.Unlock()

// 	// Check if the room already exists
// 	if _, exists := gm.rooms[roomID]; exists {
// 		fmt.Printf("Room with ID %d already exists\n", roomID)
// 		return nil
// 	}

// 	// Create a new TicTacMemoRoom and register it in the manager
// 	ticTacToeRoom := NewTicTacMemoRoom(roomID)
// 	gm.rooms[roomID] = &ticTacToeRoom.Room

// 	fmt.Printf("Created new TicTacMemoRoom with Room ID %d\n", roomID)
// 	return ticTacToeRoom
// }
