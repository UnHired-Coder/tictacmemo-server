package types

import (
	"fmt"
	"sync"
	"time"
)

// Room struct that represents a room that can hold any number of players
type Room struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	Players    []*User   `json:"players"`    // Dynamic list of players
	MaxPlayers int       `json:"maxPlayers"` // Max number of players allowed in the room
	CreatedAt  time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	Mutex      sync.Mutex
}

// CreateRoom initializes a new room with a specified max number of players
func CreateRoom(maxPlayers int) *Room {
	return &Room{
		MaxPlayers: maxPlayers,
		Players:    []*User{},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// JoinRoom adds a player to the room if there is space
func (r *Room) JoinRoom(player *User) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if len(r.Players) >= r.MaxPlayers {
		return fmt.Errorf("room is full")
	}

	r.Players = append(r.Players, player)
	fmt.Printf("Player (%s) joined Room %d\n", player.Username, r.ID)

	r.UpdatedAt = time.Now()

	// If the room is full, you can start the game or notify the players
	if len(r.Players) == r.MaxPlayers {
		r.StartGame()
	}

	return nil
}

// StartGame starts the game when all players have joined
func (r *Room) StartGame() {
	fmt.Printf("Starting the game in Room %d with %d players\n", r.ID, len(r.Players))
	// Game-specific logic will go in the derived game room
}
