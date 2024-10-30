package types

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Define Action type and enums
type Action string

const (
	ActionJoinRoom    Action = "join-room"
	ActionMakeMove    Action = "make-move"
	ActionUpdateScore Action = "update-score"
)

type GameEvent struct {
	Action Action          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type JoinRoomData struct {
	PlayerID string    `json:"playerID"`
	RoomID   uuid.UUID `json:"roomID"`
}

type MakeMoveData struct {
	PlayerID string `json:"playerID"`
	PosX     int    `json:"posX"`
	PosY     int    `json:"posY"`
}

type UpdateScoreData struct {
	PlayerID      string `json:"playerID"`
	ElapsedTime   int    `json:"elapsedTime"`
	MoveCount     int    `json:"moveCount"`
	AssignedLabel string `json:"assignedLabel"`
}
