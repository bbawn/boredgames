package dao

import (
	"github.com/google/uuid"

	"github.com/bbawn/boredgames/internal/games/set"
)

// Sets provides persistences operations for set games
type Sets interface {
	List() ([]*set.Game, error)
	Insert(g *set.Game) error
	Get(uuid uuid.UUID) (*set.Game, error)
	Update(g *set.Game) error
	Delete(uuid uuid.UUID) error
}
