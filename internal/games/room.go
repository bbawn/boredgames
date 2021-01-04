package games

import (
	"github.com/google/uuid"
)

// Room is virtual room where players play games
type Room struct {
	Uuid  uuid.UUID
	Users []string
}
