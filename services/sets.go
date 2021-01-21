package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bbawn/boredgames/internal/dao"
	"github.com/bbawn/boredgames/internal/dao/errors"
	"github.com/bbawn/boredgames/internal/games/set"
	"github.com/bbawn/boredgames/internal/router"
	"github.com/google/uuid"
)

// Sets provides the REST API for the Set board game
type Sets struct {
	dao dao.Sets
}

func NewSets(dao dao.Sets) *Sets {
	return &Sets{dao}
}

func (s *Sets) List(w http.ResponseWriter, r *http.Request) {
	games, err := s.dao.List()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load games from datastore: %s", err), errors.httpStatus(err))
		return
	}
	enc := json.NewEncoder(w)
	err := enc.Encode(games)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode games from datastore: %s", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, payload)
}

func (s *Sets) Insert(w http.ResponseWriter, r *http.Request) {
	var game *set.Game
	dec := json.NewDecoder(r.Body)
	dec.Decode(&game)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal game from data %s: %s", r.Body, err), errors.httpStatus(err))
		return
	}
	err := s.dao.Insert(game)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert game into datastore: %s", err), errors.httpStatus(err))
		return
	}
}

func (s *Sets) Get(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(router.GetField(r, 0))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid set uuid %s: %s", router.GetField(r, 0), err), http.StatusNotFound)
		return
	}
	game, err := s.dao.Get(uuid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get game from datastore: %s", err), errors.httpStatus(err))
		return
	}
	enc := json.NewEncoder(w)
	err := enc.Encode(game)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode game from datastore: %s", err), http.StatusInternalServerError)
		return
	}
}

func (s *Sets) Delete(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(router.GetField(r, 0))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid set uuid %s: %s", router.GetField(r, 0), err), errors.httpStatus(err))
		return
	}
	err := s.dao.Delete(uuid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete game from datastore: %s", err), errors.httpStatus(err))
		return
	}
	// TODO: some APIs return the resource on delete - should we?
}

type claimData struct {
	username string
	set      []*set.Card
}

func (s *Sets) Claim(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(router.GetField(r, 0))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid set uuid %s: %s", router.GetField(r, 0), err), errors.httpStatus(err))
		return
	}
	fmt.Fprintf(w, "<h1>ClaimSet</h1><div>%s</div>", uuid)
}

func (s *Sets) Next(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(router.GetField(r, 0))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid set uuid %s: %s", router.GetField(r, 0), err), http.StatusNotFound)
		return
	}
	game, err := s.dao.Get(uuid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get game from datastore: %s", err), errors.httpStatus(err))
		return
	}
	game.NextRound()
	err := s.dao.Update(game)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update game in datastore: %s", err), errors.httpStatus(err))
		return
	}
	enc := json.NewEncoder(w)
	err := enc.Encode(games)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode game next game: %s", err), http.StatusInternalServerError)
		return
	}
}
