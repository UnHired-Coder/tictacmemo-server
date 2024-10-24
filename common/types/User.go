package types

import (
	"time"
)

// User struct represents the users table in the database
type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"userId" gorm:"index"`
	Username  string    `json:"username" gorm:"index"`
	Email     string    `json:"email" gorm:"index"`
	AuthType  string    `json:"authType"`
	Rating    int       `json:"rating" gorm:"index;default:1000"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
