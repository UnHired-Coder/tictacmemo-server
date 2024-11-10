package types

import "time"

// User struct represents the users table in the database
type GameHistory struct {
	ID               int       `json:"id" gorm:"primaryKey"`
	UserID           string    `json:"userID" gorm:"index"`
	OpponentUserID   string    `json:"opponentUserID" gorm:"index"`
	OpponentUsername string    `json:"username" gorm:"index"`
	RatingChange     int       `json:"RatingChange" gorm:"index;default:1000"`
	CreatedAt        time.Time `json:"createdAt" gorm:"autoCreateTime"`
}
