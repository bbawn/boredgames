package ram

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/bbawn/boredgames/internal/dao/errors"
	"github.com/bbawn/boredgames/internal/games/set"
)

type Sets struct {
	m sync.RWMutex
	// sets stores json-serialized set Games
	// This avoids shared-object confusion if we used unserialized Games
	sets map[uuid.UUID][]byte
}

func NewSets() *Sets {
	return &Sets{sets: make(map[uuid.UUID][]byte)}
}

func (s *Sets) List() ([]*set.Game, error) {
	var gs []*set.Game
	s.m.Lock()
	defer s.m.Unlock()
	for _, jGame := range s.sets {
		var g *set.Game
		err := json.Unmarshal(jGame, &g)
		if err != nil {
			return []*set.Game{}, errors.InternalError{fmt.Sprintf("Could not Unmarshal json game: %s", jGame)}
		}
		gs = append(gs, g)
	}
	return gs, nil
}

func (s *Sets) Insert(g *set.Game) error {
	s.m.Lock()
	defer s.m.Unlock()
	jGame, ok := s.sets[g.ID]
	if ok {
		return errors.AlreadyExistsError{g.ID.String()}
	}
	jGame, err := json.Marshal(g)
	if err != nil {
		return errors.InternalError{fmt.Sprintf("Could not Marshal json game: %s err %s", g.ID, err)}
	}
	s.sets[g.ID] = jGame
	return nil
}

func (s *Sets) Get(uuid uuid.UUID) (*set.Game, error) {
	s.m.Lock()
	defer s.m.Unlock()
	jGame, ok := s.sets[uuid]
	if !ok {
		return nil, errors.NotFoundError{uuid.String()}
	}
	var g *set.Game
	err := json.Unmarshal(jGame, &g)
	if err != nil {
		return nil, errors.InternalError{fmt.Sprintf("Could not Unmarshal json game: %s err: %s", jGame, err)}
	}
	return g, nil
}

func (s *Sets) Update(g *set.Game) error {
	s.m.Lock()
	defer s.m.Unlock()
	jGame, ok := s.sets[g.ID]
	if !ok {
		return errors.NotFoundError{g.ID.String()}
	}
	jGame, err := json.Marshal(g)
	if err != nil {
		return errors.InternalError{fmt.Sprintf("Could not Marshal json game: %s", g.ID)}
	}
	s.sets[g.ID] = jGame
	return nil
}

func (s *Sets) Delete(uuid uuid.UUID) error {
	s.m.Lock()
	defer s.m.Unlock()
	if _, ok := s.sets[uuid]; !ok {
		return errors.NotFoundError{uuid.String()}
	}
	delete(s.sets, uuid)
	return nil
}
