package models

import (
	"time"
)

// Player struct represents the table holding active players in the database
type Player struct {
	ID        uint      `gorm:"primaryKey"`  // Auto-generated primary key
	UserID    string    `gorm:"uniqueIndex"` // Holds the user ID (linked to User table)
	CreatedAt time.Time // Time when the player became active
}
