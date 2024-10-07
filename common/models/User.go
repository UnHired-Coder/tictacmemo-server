package models

import (
	"time"
)

// User struct represents the users table in the database
type User struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	Email     string `gorm:"uniqueIndex"`
	CreatedAt time.Time
}
