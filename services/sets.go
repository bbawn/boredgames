package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/bbawn/boredgames/internal/dao"
	"github.com/bbawn/boredgames/internal/dao/errors"
	"github.com/bbawn/boredgames/internal/games/set"
	"github.com/bbawn/boredgames/internal/router"
)

// Sets provides the REST API for the Set board game
type Sets struct {
	dao dao.Sets
}

func NewSets(dao dao.Sets, router *router.TableRouter) *Sets {
	s := &Sets{dao}
	router.AddRoute("GET", "/sets", http.HandlerFunc(s.List))
	router.AddRoute("POST", "/sets", http.HandlerFunc(s.Create))
	router.AddRoute("GET", "/sets/([^/]+)", http.HandlerFunc(s.Get))
	router.AddRoute("DEL", "/sets/([^/]+)", http.HandlerFunc(s.Delete))
	router.AddRoute("POST", "/sets/([^/]+)/claim", http.HandlerFunc(s.Claim))
	router.AddRoute("POST", "/sets/([^/]+)/next", http.HandlerFunc(s.Next))
	return s
}

func (s *Sets) List(w http.ResponseWriter, r *http.Request) {
	games, err := s.dao.List()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load games from datastore: %s", err), errors.HttpStatus(err))
		return
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(games)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode games from datastore: %s", err), http.StatusInternalServerError)
		return
	}
}

type createData struct {
	Usernames []string
}

func (s *Sets) Create(w http.ResponseWriter, r *http.Request) {
	var cd *createData
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&cd)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal create data: %s", err), http.StatusBadRequest)
		return
	}
	game, err := set.NewGame(cd.Usernames...)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create new game: %s", err), http.StatusBadRequest)
		return
	}
	err = s.dao.Insert(game)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert game into datastore: %s", err), errors.HttpStatus(err))
		return
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(game)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode new game: %s", err), http.StatusInternalServerError)
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
		http.Error(w, fmt.Sprintf("Failed to get game from datastore: %s", err), errors.HttpStatus(err))
		return
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(game)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode game from datastore: %s", err), http.StatusInternalServerError)
		return
	}
}

func (s *Sets) Delete(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(router.GetField(r, 0))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid set uuid %s: %s", router.GetField(r, 0), err), errors.HttpStatus(err))
		return
	}
	err = s.dao.Delete(uuid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete game from datastore: %s", err), errors.HttpStatus(err))
		return
	}
}

type claimData struct {
	username            string
	card1, card2, card3 *set.Card
}

func (s *Sets) Claim(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(router.GetField(r, 0))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid set uuid %s: %s", router.GetField(r, 0), err), errors.HttpStatus(err))
		return
	}
	var cd *claimData
	dec := json.NewDecoder(r.Body)
	dec.Decode(&cd)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal claim data: %s", err), errors.HttpStatus(err))
		return
	}
	game, err := s.dao.Get(uuid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get game from datastore: %s", err), errors.HttpStatus(err))
		return
	}
	game.ClaimSet(cd.username, cd.card1, cd.card2, cd.card3)
	err = s.dao.Update(game)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update game in datastore: %s", err), errors.HttpStatus(err))
		return
	}
}

func (s *Sets) Next(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(router.GetField(r, 0))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid set uuid %s: %s", router.GetField(r, 0), err), http.StatusNotFound)
		return
	}
	game, err := s.dao.Get(uuid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get game from datastore: %s", err), errors.HttpStatus(err))
		return
	}
	game.NextRound()
	err = s.dao.Update(game)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update game in datastore: %s", err), errors.HttpStatus(err))
		return
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(game)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode game next game: %s", err), http.StatusInternalServerError)
		return
	}
}
