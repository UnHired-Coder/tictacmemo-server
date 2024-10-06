package types

import "time"

type MatchRequest struct {
	UserId *int       `json:"userId"`
	Time   *time.Time `json:"time"`
}
