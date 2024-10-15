package core

import (
	"errors"
	"game-server/tictacmemo/types"
	"sort"
	"sync"
	"time"
)

// MatchmakingSystem holds the sorted list of players and a mutex for thread safety.
type MatchmakingSystem struct {
	players []types.Player
	mutex   sync.Mutex
}

func NewMatchmakingSystem() *MatchmakingSystem {
	return &MatchmakingSystem{
		players: []types.Player{},
	}
}

func (ms *MatchmakingSystem) AddPlayer(player types.Player) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	ms.players = append(ms.players, player)
	sort.Slice(ms.players, func(i, j int) bool {
		return ms.players[i].Rating < ms.players[j].Rating
	})
}

// MatchPlayers tries to find a match for two players with similar ratings within a specified timeout.
func (ms *MatchmakingSystem) MatchPlayers(timeout time.Duration) (*types.Player, *types.Player, error) {
	expiry := time.Now().Add(timeout)

	for {
		ms.mutex.Lock()
		if len(ms.players) >= 2 {
			// Find the closest pair
			closestIndex := -1
			minDiff := int(^uint(0) >> 1) // Max int value
			for i := 0; i < len(ms.players)-1; i++ {
				diff := ms.players[i+1].Rating - ms.players[i].Rating
				if diff < minDiff {
					minDiff = diff
					closestIndex = i
				}
			}

			if closestIndex != -1 {
				player1 := ms.players[closestIndex]
				player2 := ms.players[closestIndex+1]
				// Remove matched players from the list
				ms.players = append(ms.players[:closestIndex], ms.players[closestIndex+2:]...)
				ms.mutex.Unlock()
				return &player1, &player2, nil
			}
		}
		ms.mutex.Unlock()

		if time.Now().After(expiry) {
			return nil, nil, errors.New("timeout: no match found within the specified duration")
		}

		time.Sleep(100 * time.Millisecond) // Sleep for a short duration before trying again
	}
}
