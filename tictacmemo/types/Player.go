package types

import (
	"game-server/common/types"
)

type Player struct {
	types.User             // Embedding User struct
	WaitlistId      string // Temporary id, client subscribes to this while we find a match
	InitialGameData *InitialGameData
}
