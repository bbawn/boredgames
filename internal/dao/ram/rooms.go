package ram

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/bbawn/boredgames/internal/dao/errors"
	"github.com/bbawn/boredgames/internal/rooms"
)

// Rooms is the collection of fake dao roomms
type Rooms struct {
	m sync.RWMutex
	// rooms stores json-serialized set Rooms
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
	for _, jGame := range rms.rooms {
		var r *rooms.Room
		err := json.Unmarshal(jGame, &r)
		if err != nil {
			return nil, errors.InternalError{fmt.Sprintf("Could not Unmarshal json game: %s", jGame)}
		}
		rs = append(rs, r)
	}
	return rs, nil
}

func (rms *Rooms) Insert(r *rooms.Room) error {
	rms.m.Lock()
	defer rms.m.Unlock()
	jGame, ok := rms.rooms[r.Name]
	if ok {
		return errors.AlreadyExistsError{r.Name}
	}
	jGame, err := json.Marshal(r)
	if err != nil {
		return errors.InternalError{fmt.Sprintf("Could not Marshal json game: %s err %s", r.Name, err)}
	}
	rms.rooms[r.Name] = jGame
	return nil
}

func (rms *Rooms) Get(name string) (*rooms.Room, error) {
	rms.m.Lock()
	defer rms.m.Unlock()
	jGame, ok := rms.rooms[name]
	if !ok {
		return nil, errors.NotFoundError{name}
	}
	var r *rooms.Room
	err := json.Unmarshal(jGame, &r)
	if err != nil {
		return nil, errors.InternalError{fmt.Sprintf("Could not Unmarshal json game: %s err: %s", jGame, err)}
	}
	return r, nil
}

func (rms *Rooms) Delete(name string) error {
	rms.m.Lock()
	defer rms.m.Unlock()
	if _, ok := rms.rooms[name]; !ok {
		return errors.NotFoundError{name}
	}
	delete(rms.rooms, name)
	return nil
}

func (rms *Rooms) AddPlayer(name, username string) (*rooms.Room, error) {
	rms.m.Lock()
	defer rms.m.Unlock()
	jGame, ok := rms.rooms[name]
	if !ok {
		return nil, errors.NotFoundError{name}
	}
	var r *rooms.Room
	err := json.Unmarshal(jGame, &r)
	if err != nil {
		return nil, errors.InternalError{fmt.Sprintf("Could not Unmarshal json game: %s err: %s", jGame, err)}
	}
	if _, ok := r.Usernames[username]; ok {
		return nil, errors.AlreadyExistsError{username}
	}
	r.Usernames[username] = true
	jGame, err = json.Marshal(r)
	if err != nil {
		return r, errors.InternalError{fmt.Sprintf("Could not Marshal json game: %s err %s", r.Name, err)}
	}
	rms.rooms[r.Name] = jGame
	return r, nil
}

func (rms *Rooms) DeletePlayer(name, username string) (*rooms.Room, error) {
	rms.m.Lock()
	defer rms.m.Unlock()
	jGame, ok := rms.rooms[name]
	if !ok {
		return nil, errors.NotFoundError{name}
	}
	var r *rooms.Room
	err := json.Unmarshal(jGame, &r)
	if err != nil {
		return nil, errors.InternalError{fmt.Sprintf("Could not Unmarshal json game: %s err: %s", jGame, err)}
	}
	if _, ok := r.Usernames[username]; !ok {
		return nil, errors.NotFoundError{username}
	}
	// Remove element from players
	delete(r.Usernames, username)
	jGame, err = json.Marshal(r)
	if err != nil {
		return r, errors.InternalError{fmt.Sprintf("Could not Marshal json game: %s err %s", r.Name, err)}
	}
	rms.rooms[r.Name] = jGame
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