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

func SetsAddRoutes(dao dao.Sets, router *router.TableRouter) {
	s := &Sets{dao}
	router.AddRoute("GET", "/sets", http.HandlerFunc(s.List))
	router.AddRoute("POST", "/sets", http.HandlerFunc(s.Create))
	router.AddRoute("GET", "/sets/([^/]+)", http.HandlerFunc(s.Get))
	router.AddRoute("DEL", "/sets/([^/]+)", http.HandlerFunc(s.Delete))
	router.AddRoute("POST", "/sets/([^/]+)/claim", http.HandlerFunc(s.Claim))
	router.AddRoute("POST", "/sets/([^/]+)/next", http.HandlerFunc(s.Next))
}

func (s *Sets) List(w http.ResponseWriter, r *http.Request) {
	games, err := s.dao.List()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load games from datastore: %s", err), httpStatus(err))
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
	var cd createData
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
		http.Error(w, fmt.Sprintf("Failed to insert game into datastore: %s", err), httpStatus(err))
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
		http.Error(w, fmt.Sprintf("Failed to get game from datastore: %s", err), httpStatus(err))
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
		http.Error(w, fmt.Sprintf("Invalid set uuid %s: %s", router.GetField(r, 0), err), httpStatus(err))
		return
	}
	err = s.dao.Delete(uuid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete game from datastore: %s", err), httpStatus(err))
		return
	}
}

type claimData struct {
	Username string
	Cards    set.CardTriple
}

func (s *Sets) Claim(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(router.GetField(r, 0))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid set uuid %s: %s", router.GetField(r, 0), err), httpStatus(err))
		return
	}
	var cd claimData
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&cd)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal claim data: %s", err), http.StatusBadRequest)
		return
	}
	game, err := s.dao.Get(uuid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get game from datastore: %s", err), httpStatus(err))
		return
	}
	err = game.ClaimSet(cd.Username, cd.Cards)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to claim set in game: %s", err), httpStatus(err))
		return
	}
	err = s.dao.Update(game)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update game in datastore: %s", err), httpStatus(err))
		return
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(game)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode game next game: %s", err), http.StatusInternalServerError)
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
		http.Error(w, fmt.Sprintf("Failed to get game from datastore: %s", err), httpStatus(err))
		return
	}
	err = game.NextRound()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to advance game to next round: %s", err), httpStatus(err))
		return
	}
	err = s.dao.Update(game)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update game in datastore: %s", err), httpStatus(err))
		return
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(game)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode game next game: %s", err), http.StatusInternalServerError)
		return
	}
}

func httpStatus(err error) int {
	switch err.(type) {
	case errors.AlreadyExistsError:
		return http.StatusConflict
	case errors.InternalError:
		return http.StatusInternalServerError
	case errors.NotFoundError:
		return http.StatusNotFound
	case set.InvalidArgError:
		return http.StatusBadRequest
	case set.InvalidStateError:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
