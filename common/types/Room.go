package types

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid" // Import the UUID package
)

// Define an interface that both Base and Derived will adhere to
type IStartGame interface {
	StartGame()
}

// Room struct that represents a room that can hold any number of players
type Room struct {
	ID         uuid.UUID `json:"room_id" gorm:"primaryKey"`
	Players    []*User   `json:"players"`    // Dynamic list of players
	MaxPlayers int       `json:"maxPlayers"` // Max number of players allowed in the room // This can be removed
	CreatedAt  time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	Mutex      sync.Mutex

	OnStartGame IStartGame `json:"-"`
}

// CreateRoom initializes a new room with a specified max number of players
func CreateRoom(maxPlayers int, roomID uuid.UUID) *Room {
	return &Room{
		ID:         roomID,
		MaxPlayers: maxPlayers,
		Players:    []*User{},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// JoinRoom adds a player to the room if there is space
func (room *Room) JoinRoom(player *User) error {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	if len(room.Players) >= room.MaxPlayers {
		return fmt.Errorf("room is full")
	}

	room.Players = append(room.Players, player)
	fmt.Printf("Player (%s) joined Room %d\n", player.Username, room.ID)

	room.UpdatedAt = time.Now()

	// If the room is full, you can start the game or notify the players
	if len(room.Players) == room.MaxPlayers {
		room.StartGame()
	}

	return nil
}

// StartGame starts the game when all players have joined
func (room *Room) StartGame() {
	fmt.Printf("Begning to start %d with %d players\n", room.ID, len(room.Players))
	// Game-specific logic will go in the derived game room

	room.OnStartGame.StartGame()
}
