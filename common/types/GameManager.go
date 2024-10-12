package types

import (
	"fmt"
	"sync"
)

// GameManager struct manages all active game rooms
type GameManager struct {
	Rooms map[int]*Room // Manages generic rooms
	Lock  sync.Mutex    // Mutex for concurrent access
}

// NewGameManager creates a new GameManager
func NewGameManager() *GameManager {
	return &GameManager{
		Rooms: make(map[int]*Room), // Initialize the rooms map
	}
}

// RemoveRoom removes the room from the GameManager once the game ends
func (gm *GameManager) RemoveRoom(roomID int) {
	gm.Lock.Lock()
	defer gm.Lock.Unlock()

	if _, exists := gm.Rooms[roomID]; exists {
		delete(gm.Rooms, roomID)
		fmt.Printf("Room with ID %d has been removed from GameManager\n", roomID)
	} else {
		fmt.Printf("Room with ID %d does not exist\n", roomID)
	}
}

// JoinRoom allows a player to join a room. If the room doesn't exist, it creates a new one.
func (gm *GameManager) JoinRoom(player *User, roomID int, maxPlayers int) error {
	gm.Lock.Lock()
	defer gm.Lock.Unlock()

	// Check if the room exists, if not create it
	room, exists := gm.Rooms[roomID]
	if !exists {
		// Create a new Room (e.g., for TicTacMemoRoom, max players = 2)
		room = CreateRoom(maxPlayers)
		gm.Rooms[roomID] = room
		fmt.Printf("Created new Room with Room ID %d\n", roomID)
	}

	// Let the player join the room
	err := room.JoinRoom(player)
	if err != nil {
		return fmt.Errorf("error joining room: %v", err)
	}

	return nil
}
