package dao

import (
	"github.com/bbawn/boredgames/internal/rooms"
)

// Rooms provides persistences operations for game rooms
type Rooms interface {
	List() ([]*rooms.Room, error)
	Insert(r *rooms.Room) error
	Get(name string) (*rooms.Room, error)
	Delete(name string) error
	AddPlayer(name, username string) (*rooms.Room, error)
	DeletePlayer(name, username string) (*rooms.Room, error)
}
