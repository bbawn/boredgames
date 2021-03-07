package dao

import (
	"github.com/bbawn/boredgames/internal/rooms"
)

// Rooms provides persistences operations for game rooms
type Rooms interface {
	List() ([]*rooms.Room, error)
	Insert(g *rooms.Room) error
	Get(name string) (*rooms.Room, error)
	Update(g *rooms.Room) error
	Delete(name string) error
}
