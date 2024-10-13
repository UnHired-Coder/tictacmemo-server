package types

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Define Action type and enums
type Action string

const (
	ActionJoinRoom Action = "join-room"
	ActionMakeMove Action = "make-move"
)

type GameEvent struct {
	Action Action          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type JoinRoomData struct {
	PlayerID int       `json:"playerID"`
	RoomID   uuid.UUID `json:"roomID"`
}

type MakeMoveData struct {
	PlayerID int       `json:"playerID"`
	RoomID   uuid.UUID `json:"roomID"`
	PosX     int       `json:"posX"`
	PosY     int       `json:"posY"`
}
