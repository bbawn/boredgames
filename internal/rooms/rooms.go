package rooms

import (
	"github.com/google/uuid"
)

type GameType int

const (
	None GameType = iota
	Set
	Boggle
	RRobots
)

//go:generate stringer -type=GameType

// Room is an instance of a game room
type Room struct {
	// Name is the human-readable unique identifier of the Room
	Name string
	// Usernames are the usernames of players in the room
	Usernames []string
	// GameType is the type of game currently being played in the room
	GameType GameType
	// GameID is the identifier of the current game being played in the room
	GameID uuid.UUID
}

func NewRoom(name string, usernames ...string) *Room {
	r := new(Room)
	r.Name = name
	r.Usernames = usernames
	return r
}
