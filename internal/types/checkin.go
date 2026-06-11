package types

import (
	"time"
)

// CheckIn represents a user's daily check-in
type CheckIn struct {
	Address string    `json:"address"`
	Date    time.Time `json:"date"`
}