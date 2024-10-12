package types

import (
	"fmt"
	commonTypes "game-server/common/types"

	"github.com/google/uuid"
)

type TicTacMemoGameManager struct {
	*commonTypes.GameManager[TicTacMemoRoom]
}

// NewTicTacMemoGameManager creates a new TicTacMemoGameManager
func NewTicTacMemoGameManager() *TicTacMemoGameManager {
	return &TicTacMemoGameManager{
		GameManager: commonTypes.NewGameManager[TicTacMemoRoom](),
	}
}

// CreateRoom creates a new room and returns the room's UUID
func (gm *TicTacMemoGameManager) CreateRoom(maxPlayers int) (uuid.UUID, *TicTacMemoRoom, error) {
	gm.Lock.Lock()
	defer gm.Lock.Unlock()

	// Generate a new UUID for the room
	roomID := uuid.New()

	room := CreateRoom(maxPlayers, roomID)
	room.Room.OnStartGame = room

	// Store the room in the GameManager's map
	gm.Rooms[roomID] = room
	fmt.Printf("Created new Room with Room ID %s\n", roomID)

	// Return the roomID and the room itself
	return roomID, room, nil
}
