package types

import "time"

// User struct represents the users table in the database
type GameHistory struct {
	ID                 int       `json:"id" gorm:"primaryKey"`
	UserID             string    `json:"user_id" gorm:"index"`
	OpponentUserID     string    `json:"opponent_user_id" gorm:"index"`
	OpponentUsername   string    `json:"username" gorm:"index"`
	OpponentAvatar     string    `json:"opponent_avatar" gorm:"index"`
	RatingChange       int       `json:"rating_change" gorm:"index;default:0"`
	RatingBeforeChange int       `json:"rating_before_change" gorm:"index;default:1000"`
	CreatedAt          time.Time `json:"createdAt" gorm:"autoCreateTime"`
}
