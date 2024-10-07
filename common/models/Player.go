package models

type Player struct {
	User              // Embedding User struct
	WaitlistId string // Temporary id, client subscribes to this while we find a match
}
