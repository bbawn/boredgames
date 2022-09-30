package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/bbawn/boredgames/internal/dao"
	"github.com/bbawn/boredgames/internal/rooms"
	"github.com/bbawn/boredgames/internal/router"
)

// Rooms provides the REST API for the game room resource
type Rooms struct {
	dao dao.Rooms
}

// RoomsAddRoutes adds the routes for this service to the given router
func RoomsAddRoutes(dao dao.Rooms, router *router.TableRouter) {
	rms := &Rooms{dao}
	router.AddRoute("GET", "/rooms", http.HandlerFunc(rms.List))
	router.AddRoute("POST", "/rooms", http.HandlerFunc(rms.Create))
	router.AddRoute("GET", "/rooms/([^/]+)", http.HandlerFunc(rms.Get))
	router.AddRoute("DEL", "/rooms/([^/]+)", http.HandlerFunc(rms.Delete))
	router.AddRoute("POST", "/rooms/([^/]+)/players", http.HandlerFunc(rms.AddPlayer))
	router.AddRoute("DEL", "/rooms/([^/]+)/players", http.HandlerFunc(rms.DeletePlayer))
	router.AddRoute("PUT", "/rooms/([^/]+)/game", http.HandlerFunc(rms.SetGame))
}

// List returns the list of all rooms
func (rms *Rooms) List(w http.ResponseWriter, r *http.Request) {
	rooms, err := rms.dao.List()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load rooms from datastore: %s", err), httpStatus(err))
		return
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(rooms)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode rooms from datastore: %s", err), http.StatusInternalServerError)
		return
	}
}

// postData is the data payload of the create room request
type postData struct {
	// Name is human-readable identifier of the Room
	Name string `json:"name"`
	// Usernames is the set of players in the room
	Usernames map[string]bool `json:"usernames"`
}

func (pd *postData) validate() error {
	if pd.Name == "" {
		return errors.New("non-empty Name is required")
	}
	for u := range pd.Usernames {
		if u == "" {
			return errors.New("usernames must be non-empty")
		}
	}
	return nil
}

// Create creates a new room
func (rms *Rooms) Create(w http.ResponseWriter, r *http.Request) {
	var pd postData
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&pd)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal create data: %s", err), http.StatusBadRequest)
		return
	}
	err = pd.validate()
	if err != nil {
		m := fmt.Sprintf("Invalid request payload err: %s", err)
		http.Error(w, m, http.StatusBadRequest)
		return
	}
	room := rooms.NewRoom(pd.Name, pd.Usernames)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create new room: %s", err), http.StatusBadRequest)
		return
	}
	err = rms.dao.Insert(room)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert room into datastore: %s", err), httpStatus(err))
		return
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(room)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode new room: %s", err), http.StatusInternalServerError)
		return
	}
}

// Get returns the room with requested name
func (rms *Rooms) Get(w http.ResponseWriter, r *http.Request) {
	name := router.GetField(r, 0)
	if name == "" {
		http.Error(w, fmt.Sprintf("Invalid room name %s:", router.GetField(r, 0)), http.StatusNotFound)
		return
	}
	game, err := rms.dao.Get(name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get room from datastore: %s", err), httpStatus(err))
		return
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(game)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode room from datastore: %s", err), http.StatusInternalServerError)
		return
	}
}

// Delete deletes the room with the requested name
func (rms *Rooms) Delete(w http.ResponseWriter, r *http.Request) {
	name := router.GetField(r, 0)
	if name == "" {
		http.Error(w, fmt.Sprintf("Invalid room name %s:", router.GetField(r, 0)), http.StatusNotFound)
		return
	}
	err := rms.dao.Delete(name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete room from datastore: %s", err), httpStatus(err))
		return
	}
}

// playerData is the payload of the post and delete room player requests
type playerData struct {
	username string
}

// AddPlayer adds a player to the room
func (rms *Rooms) AddPlayer(w http.ResponseWriter, r *http.Request) {
	name := router.GetField(r, 0)
	if name == "" {
		http.Error(w, fmt.Sprintf("Invalid room name %s:", router.GetField(r, 0)), http.StatusNotFound)
		return
	}
	var pd playerData
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&pd)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal player data: %s", err), http.StatusBadRequest)
		return
	}
	room, err := rms.dao.AddPlayer(name, pd.username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add player into datastore: %s", err), httpStatus(err))
		return
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(room)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode updated room: %s", err), http.StatusInternalServerError)
		return
	}
}

// DeletePlayer deletes a player from the room
func (rms *Rooms) DeletePlayer(w http.ResponseWriter, r *http.Request) {
	name := router.GetField(r, 0)
	if name == "" {
		http.Error(w, fmt.Sprintf("Invalid room name %s:", router.GetField(r, 0)), http.StatusNotFound)
		return
	}
	var pd playerData
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&pd)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal player data: %s", err), http.StatusBadRequest)
		return
	}
	room, err := rms.dao.DeletePlayer(name, pd.username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete player from datastore: %s", err), httpStatus(err))
		return
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(room)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode updated room: %s", err), http.StatusInternalServerError)
		return
	}
}

// gameData is the payload of the game update PUT request
type gameData struct {
	typ rooms.GameType
	id  uuid.UUID
}

// SetGame sets the game type and id for the room
func (rms *Rooms) SetGame(w http.ResponseWriter, r *http.Request) {
	name := router.GetField(r, 0)
	if name == "" {
		http.Error(w, fmt.Sprintf("Invalid room name %s:", router.GetField(r, 0)), http.StatusNotFound)
		return
	}
	var gd gameData
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&gd)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal game data: %s", err), http.StatusBadRequest)
		return
	}
	room, err := rms.dao.SetGame(name, gd.typ, gd.id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update game: %s", err), httpStatus(err))
		return
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(room)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode updated room: %s", err), http.StatusInternalServerError)
		return
	}
}
