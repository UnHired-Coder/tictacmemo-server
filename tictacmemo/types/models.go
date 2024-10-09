package types

import (
	"game-server/common/models"
	"time"
)

// Room struct represents a room in the game containing two players
type Room struct {
	ID         int         `json:"id" gorm:"primaryKey"`
	Player1ID  int         `json:"player1Id"`
	Player1    models.User `json:"player1" gorm:"foreignKey:Player1ID"`
	Player2ID  int         `json:"player2Id"`
	Player2    models.User `json:"player2" gorm:"foreignKey:Player2ID"`
	MatchEnded bool        `json:"matchEnded"`
	WinnerID   *int        `json:"winnerId"`
	Winner     models.User `json:"winner" gorm:"foreignKey:WinnerID"`
	CreatedAt  time.Time   `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt  time.Time   `json:"updatedAt" gorm:"autoUpdateTime"`
}

func CreateRoom(player1 models.Player, player2 models.Player) Room {
	room := Room{
		Player1ID:  player1.ID,
		Player1:    player1.User,
		Player2ID:  player2.ID,
		Player2:    player2.User,
		MatchEnded: false, // Default match not ended
		WinnerID:   nil,   // Default winner is nil
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return room
}
