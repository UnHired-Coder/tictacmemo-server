package types

import (
	"fmt"
	"game-server/common/types"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicTacMemoRoom struct {
	types.Room
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

func (room *TicTacMemoRoom) StartGame() {
	fmt.Printf("Starting the game in Room %d with %d players\n", room.ID, len(room.Players))

}

func (room *TicTacMemoRoom) MakeMove(db *gorm.DB, makeMoveData MakeMoveData) {
	log.Printf("Make move: posX: %d, posY: %d", makeMoveData.PosX, makeMoveData.PosY)
}
