package dao

import (
	"github.com/google/uuid"

	"github.com/bbawn/boredgames/internal/rooms"
)

// Rooms provides persistence operations for game rooms
type Rooms interface {
	List() ([]*rooms.Room, error)
	Insert(r *rooms.Room) error
	Get(name string) (*rooms.Room, error)
	Delete(name string) error
	AddPlayer(name, username string) (*rooms.Room, error)
	DeletePlayer(name, username string) (*rooms.Room, error)
	SetGame(name string, typ rooms.GameType, id uuid.UUID) (*rooms.Room, error)
}
