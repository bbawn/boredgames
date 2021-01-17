package ram

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"

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
		err := json.Unmarshal(jGame, g)
		if err != nil {
			return []*set.Game{}, fmt.Errorf("Could not Unmarshal json game: %s", jGame)
		}
		fmt.Printf("ram: s.List(): g: %v\n", g)
		gs = append(gs, g)
	}
	return gs, nil
}

func (s *Sets) Insert(g *set.Game) error {
	return s.upsert(g)
}

func (s *Sets) Get(uuid uuid.UUID) (*set.Game, error) {
	s.m.Lock()
	defer s.m.Unlock()
	jGame, ok := s.sets[uuid]
	if !ok {
		return nil, nil
	}
	var g *set.Game
	err := json.Unmarshal(jGame, g)
	if err != nil {
		return nil, fmt.Errorf("Could not Unmarshal json game: %s", jGame)
	}
	return g, nil
}

func (s *Sets) Update(g *set.Game) error {
	return s.upsert(g)
}

func (s *Sets) Delete(uuid uuid.UUID) error {
	s.m.Lock()
	defer s.m.Unlock()
	if _, ok := s.sets[uuid]; !ok {
		return fmt.Errorf("Set id %s not found", uuid)
	}
	delete(s.sets, uuid)
	return nil
}

func (s *Sets) upsert(g *set.Game) error {
	s.m.Lock()
	defer s.m.Unlock()
	jGame, err := json.Marshal(g)
	if err != nil {
		return fmt.Errorf("Could not Marshal game: %v", g)
	}
	s.sets[g.ID] = jGame
	return nil
}
