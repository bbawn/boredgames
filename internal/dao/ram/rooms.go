package ram

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"

	"github.com/bbawn/boredgames/internal/dao/errors"
	"github.com/bbawn/boredgames/internal/rooms"
)

// Rooms is the collection of fake dao roomms
type Rooms struct {
	m sync.RWMutex
	// rooms stores json-serialized set Rooms keyed on name
	// This avoids future shared-object confusion if we used unserialized Rooms
	rooms map[string][]byte
}

func NewRooms() *Rooms {
	return &Rooms{rooms: make(map[string][]byte)}
}

func (rms *Rooms) List() ([]*rooms.Room, error) {
	// Empty slice, not nil so we can always unmarshal to json array
	rs := []*rooms.Room{}
	rms.m.Lock()
	defer rms.m.Unlock()
	for _, jRoom := range rms.rooms {
		var r *rooms.Room
		err := json.Unmarshal(jRoom, &r)
		if err != nil {
			return nil, errors.InternalError{Details: fmt.Sprintf("Could not Unmarshal json game: %s", jRoom)}
		}
		rs = append(rs, r)
	}
	return rs, nil
}

func (rms *Rooms) Insert(r *rooms.Room) error {
	rms.m.Lock()
	defer rms.m.Unlock()
	jRoom, ok := rms.rooms[r.Name]
	if ok {
		return errors.AlreadyExistsError{Key: r.Name}
	}
	jRoom, err := json.Marshal(r)
	if err != nil {
		return errors.InternalError{Details: fmt.Sprintf("Could not Marshal json game: %s err %s", r.Name, err)}
	}
	rms.rooms[r.Name] = jRoom
	return nil
}

func (rms *Rooms) Get(name string) (*rooms.Room, error) {
	rms.m.Lock()
	defer rms.m.Unlock()
	jRoom, ok := rms.rooms[name]
	if !ok {
		return nil, errors.NotFoundError{Key: name}
	}
	var r *rooms.Room
	err := json.Unmarshal(jRoom, &r)
	if err != nil {
		return nil, errors.InternalError{Details: fmt.Sprintf("Could not Unmarshal json game: %s err: %s", jRoom, err)}
	}
	return r, nil
}

func (rms *Rooms) Delete(name string) error {
	rms.m.Lock()
	defer rms.m.Unlock()
	if _, ok := rms.rooms[name]; !ok {
		return errors.NotFoundError{Key: name}
	}
	delete(rms.rooms, name)
	return nil
}

func (rms *Rooms) AddPlayer(name, username string) (*rooms.Room, error) {
	rms.m.Lock()
	defer rms.m.Unlock()
	jRoom, ok := rms.rooms[name]
	if !ok {
		return nil, errors.NotFoundError{Key: name}
	}
	var r *rooms.Room
	err := json.Unmarshal(jRoom, &r)
	if err != nil {
		return nil, errors.InternalError{Details: fmt.Sprintf("Could not Unmarshal json game: %s err: %s", jRoom, err)}
	}
	if _, ok := r.Usernames[username]; ok {
		return nil, errors.AlreadyExistsError{Key: username}
	}
	r.Usernames[username] = true
	jRoom, err = json.Marshal(r)
	if err != nil {
		return r, errors.InternalError{Details: fmt.Sprintf("Could not Marshal json game: %s err %s", r.Name, err)}
	}
	rms.rooms[r.Name] = jRoom
	return r, nil
}

func (rms *Rooms) DeletePlayer(name, username string) (*rooms.Room, error) {
	rms.m.Lock()
	defer rms.m.Unlock()
	jRoom, ok := rms.rooms[name]
	if !ok {
		return nil, errors.NotFoundError{Key: name}
	}
	var r *rooms.Room
	err := json.Unmarshal(jRoom, &r)
	if err != nil {
		return nil, errors.InternalError{Details: fmt.Sprintf("Could not Unmarshal json game: %s err: %s", jRoom, err)}
	}
	if _, ok := r.Usernames[username]; !ok {
		return nil, errors.NotFoundError{Key: username}
	}
	// Remove element from players
	delete(r.Usernames, username)
	jRoom, err = json.Marshal(r)
	if err != nil {
		return r, errors.InternalError{Details: fmt.Sprintf("Could not Marshal json game: %s err %s", r.Name, err)}
	}
	rms.rooms[r.Name] = jRoom
	return r, nil
}

func (rms *Rooms) SetGame(name string, typ rooms.GameType, id uuid.UUID) (*rooms.Room, error) {
	rms.m.Lock()
	defer rms.m.Unlock()
	jRoom, ok := rms.rooms[name]
	if !ok {
		return nil, errors.NotFoundError{Key: name}
	}
	var r *rooms.Room
	err := json.Unmarshal(jRoom, &r)
	if err != nil {
		return nil, errors.InternalError{Details: fmt.Sprintf("Could not Unmarshal json game: %s err: %s", jRoom, err)}
	}
	r.GameType = typ
	r.GameID = id
	jRoom, err = json.Marshal(r)
	if err != nil {
		return r, errors.InternalError{Details: fmt.Sprintf("Could not Marshal json game: %s err %s", r.Name, err)}
	}
	rms.rooms[r.Name] = jRoom
	return r, nil
}

func (rms *Rooms) Dump() string {
	var b strings.Builder
	rms.m.Lock()
	defer rms.m.Unlock()
	for name, room := range rms.rooms {
		b.WriteString(fmt.Sprintf("name %s: room %s\n", name, room))
	}
	return b.String()
}

// find returns the index of given string in given slice of -1 for not found
func find(ss []string, s string) int {
	for i, elt := range ss {
		if elt == s {
			return i
		}
	}
	return -1
}
