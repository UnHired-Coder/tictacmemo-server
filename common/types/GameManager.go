package types

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid" // Import the UUID package
)

// GameManager struct manages all active game rooms
type GameManager struct {
	Rooms map[uuid.UUID]*Room // Manages generic rooms using UUID as keys
	Lock  sync.Mutex          // Mutex for concurrent access
}

// NewGameManager creates a new GameManager
func NewGameManager() *GameManager {
	return &GameManager{
		Rooms: make(map[uuid.UUID]*Room), // Initialize the rooms map with UUID keys
	}
}

// CreateRoom creates a new room and returns the room's UUID
func (gm *GameManager) CreateRoom(maxPlayers int) (uuid.UUID, error) {
	gm.Lock.Lock()
	defer gm.Lock.Unlock()

	// Generate a new UUID for the room
	roomID := uuid.New()

	// Create and store the new room
	gm.Rooms[roomID] = CreateRoom(maxPlayers)
	fmt.Printf("Created new Room with Room ID %s\n", roomID)

	return roomID, nil
}

// RemoveRoom removes the room from the GameManager once the game ends
func (gm *GameManager) RemoveRoom(roomID uuid.UUID) {
	gm.Lock.Lock()
	defer gm.Lock.Unlock()

	if _, exists := gm.Rooms[roomID]; exists {
		delete(gm.Rooms, roomID)
		fmt.Printf("Room with ID %s has been removed from GameManager\n", roomID)
	} else {
		fmt.Printf("Room with ID %s does not exist\n", roomID)
	}
}

// JoinRoom allows a player to join a room. Throws an error if the room doesn't exist.
func (gm *GameManager) JoinRoom(player *User, roomID uuid.UUID) error {
	gm.Lock.Lock()
	defer gm.Lock.Unlock()

	// Check if the room exists, if not return an error
	room, exists := gm.Rooms[roomID]
	if !exists {
		return errors.New(fmt.Sprintf("room with ID %s does not exist", roomID))
	}

	// Let the player join the room
	err := room.JoinRoom(player)
	if err != nil {
		return fmt.Errorf("error joining room: %v", err)
	}

	return nil
}
