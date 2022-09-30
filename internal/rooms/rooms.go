package rooms

import (
	"github.com/google/uuid"
)

// GameType is the type of a game
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
	Name string `json:"name"`
	// Usernames are the usernames of players in the room
	Usernames map[string]bool `json:"usernames"`
	// GameType is the type of game currently being played in the room
	GameType GameType `json:"gameType"`
	// GameID is the identifier of the current game being played in the room
	GameID uuid.UUID `json:"gameID"`
}

// NewRoom creates a room with given name and players
func NewRoom(name string, usernames map[string]bool) *Room {
	r := new(Room)
	r.Name = name
	r.Usernames = usernames
	return r
}
