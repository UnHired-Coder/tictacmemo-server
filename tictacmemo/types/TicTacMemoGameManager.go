package types

import (
	"fmt"
	commonTypes "game-server/common/types"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicTacMemoGameManager struct {
	*commonTypes.GameManager[TicTacMemoRoom]
}

func NewTicTacMemoGameManager() *TicTacMemoGameManager {
	return &TicTacMemoGameManager{
		GameManager: commonTypes.NewGameManager[TicTacMemoRoom](),
	}
}

func (gm *TicTacMemoGameManager) CreateRoom(maxPlayers int) (uuid.UUID, *TicTacMemoRoom, error) {
	roomID := uuid.New()
	room := CreateRoom(maxPlayers, roomID)
	gm.GameManager.CreateRoom(room.ID, room)
	room.Room.OnStartGame = room

	fmt.Printf("Created new Room with Room ID %s\n", roomID)
	return roomID, room, nil
}

func (gm *TicTacMemoGameManager) JoinRoom(db *gorm.DB, roomData JoinRoomData) (*TicTacMemoRoom, error) {
	var user commonTypes.User
	result := db.First(&user, roomData.PlayerID)
	if result.Error != nil {
		log.Fatal("Error fetching user:", result.Error)
		return nil, result.Error
	}

	joinFunc := func(room *TicTacMemoRoom, player *commonTypes.User) error {
		for i := range room.Players {
			if room.Players[i].ID == player.ID {
				return fmt.Errorf("Player already joined!")
			}
		}
		return room.JoinRoom(player)
	}

	room, err := gm.GameManager.JoinRoom(&user, roomData.RoomID, joinFunc)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return room, err
}
