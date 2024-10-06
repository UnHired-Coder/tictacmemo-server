package types

import "time"

type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"index"`
	Email     string    `json:"email" gorm:"index"`
	AuthType  string    `json:"authType"`
	Rating    int       `json:"rating" gorm:"index"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

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
