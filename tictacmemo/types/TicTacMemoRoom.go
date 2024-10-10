package types

import (
	"game-server/common/types"
)

type TicTacMemoRoom struct {
	types.Room            // Embeds the Room struct
	Board       [3][3]int // A specific board for this game (could vary for different games)
	CurrentTurn int       // Tracks the current player’s turn
}
