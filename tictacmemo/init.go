package tictacmemo

import "game-server/tictacmemo/core"

func InitMatchMaking() *core.MatchmakingSystem {
	// Initialize the global matchmaking system
	return core.NewMatchmakingSystem()
}
