package types

import (
	"time"
)

// Room struct represents a room in the game containing two players
type Room struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	Player1ID  int       `json:"player1Id"`
	Player1    User      `json:"player1" gorm:"foreignKey:Player1ID"`
	Player2ID  int       `json:"player2Id"`
	Player2    User      `json:"player2" gorm:"foreignKey:Player2ID"`
	MatchEnded bool      `json:"matchEnded"`
	WinnerID   *int      `json:"winnerId"`
	Winner     User      `json:"winner" gorm:"foreignKey:WinnerID"`
	CreatedAt  time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

func CreateRoom(player1 Player, player2 Player) Room {
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
