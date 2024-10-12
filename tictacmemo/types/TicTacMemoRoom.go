package types

import (
	"fmt"
	"game-server/common/types"
	"time"

	"github.com/google/uuid"
)

type TicTacMemoRoom struct {
	types.Room
}

// StartGame starts the game when all players have joined
func (room *TicTacMemoRoom) StartGame() {
	fmt.Printf("Starting the game in Room %d with %d players\n", room.ID, len(room.Players))
	// Game-specific logic will go in the derived game room
}

func CreateRoom(maxPlayers int, roomID uuid.UUID) *TicTacMemoRoom {
	return &TicTacMemoRoom{
		types.Room{
			ID:         roomID,
			MaxPlayers: maxPlayers,
			Players:    []*types.User{},
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}
}
