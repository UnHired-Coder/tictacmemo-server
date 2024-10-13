package types

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// GameManager struct manages all active game rooms using generics
type GameManager[T any] struct {
	Rooms map[uuid.UUID]*T // Manages generic rooms using UUID as keys
	Lock  sync.Mutex       // Mutex for concurrent access
}

func NewGameManager[T any]() *GameManager[T] {
	return &GameManager[T]{
		Rooms: make(map[uuid.UUID]*T),
	}
}

// CreateRoom creates a new room and returns the room's UUID
func (gm *GameManager[T]) CreateRoom(roomID uuid.UUID, room *T) (uuid.UUID, *T, error) {
	gm.Lock.Lock()
	defer gm.Lock.Unlock()

	// Store the Room in GameManager
	gm.Rooms[roomID] = room
	fmt.Printf("Created new Room with Room ID %s\n", roomID)

	return roomID, room, nil
}

// RemoveRoom removes the room from the GameManager once the game ends
func (gm *GameManager[T]) RemoveRoom(roomID uuid.UUID) {
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
func (gm *GameManager[T]) JoinRoom(player *User, roomID uuid.UUID, joinFunc func(room *T, player *User) error) (*T, error) {
	gm.Lock.Lock()
	defer gm.Lock.Unlock()

	// Check if the room exists, if not return an error
	room, exists := gm.Rooms[roomID]
	if !exists {
		return nil, errors.New(fmt.Sprintf("room with ID %s does not exist", roomID))
	}

	// Let the player join the room using the provided join function
	err := joinFunc(room, player)
	if err != nil {
		return nil, fmt.Errorf("error joining room: %v", err)
	}

	return room, nil
}
