package dao

import (
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/bbawn/boredgames/internal/dao/errors"
	"github.com/bbawn/boredgames/internal/dao/ram"
	"github.com/bbawn/boredgames/internal/rooms"
)

// TestRamRooms tests the ram implementation of Rooms
func TestRamRooms(t *testing.T) {
	ram := ram.NewRooms()
	testRooms(t, ram)
}

// testRooms tests the given implementor of Rooms
func testRooms(t *testing.T, rms Rooms) {
	// Empty list
	expRs := []*rooms.Room{}
	rs, err := rms.List()
	if err != nil {
		t.Errorf("List returned unexpected err %#v", err)
	}
	if !roomsEqual(rs, expRs) {
		t.Errorf("List returned %#v, expected %#v", rs, expRs)
	}

	// Insert room with players
	r0 := rooms.NewRoom("r0", map[string]bool{"p0": true, "p1": true})
	err = rms.Insert(r0)
	if err != nil {
		t.Errorf("Unexpected err %s on Insert", err)
	}

	// Duplicate Insert fails
	err = rms.Insert(r0)
	_, ok := err.(errors.AlreadyExistsError)
	if !ok {
		t.Errorf("Expected err %s to be of type AlreadyExistsError", err)
	}

	// Insert room without players
	r1 := rooms.NewRoom("r1", map[string]bool{})
	err = rms.Insert(r1)
	if err != nil {
		t.Errorf("Unexpected err %s on Insert", err)
	}

	// Retrieve existing room
	r, err := rms.Get(r0.Name)
	if err != nil {
		t.Errorf("Unexpected err %s on Get", err)
	}
	if !reflect.DeepEqual(r, r0) {
		t.Errorf("Get returned %#v, expected %#v", r, r0)
	}

	// Retrieve existing room
	r, err = rms.Get(r1.Name)
	if err != nil {
		t.Errorf("Unexpected err %s on Get", err)
	}
	if !reflect.DeepEqual(r, r1) {
		t.Errorf("Get returned %#v, expected %#v", r, r1)
	}

	// List all rooms
	expRs = []*rooms.Room{r0, r1}
	rs, err = rms.List()
	if err != nil {
		t.Errorf("List returned unexpected err %#v", err)
	}
	if !roomsEqual(rs, expRs) {
		t.Errorf("List returned %#v, expected %#v", rs, expRs)
	}

	// Add player to room
	r0.Usernames["p2"] = true
	r, err = rms.AddPlayer(r0.Name, "p2")
	if err != nil {
		t.Errorf("Unexpected err %s on AddPlayer", err)
	}
	if !reflect.DeepEqual(r, r0) {
		t.Errorf("Get returned %#v, expected %#v", r, r0)
	}

	// Retrieve existing room
	r, err = rms.Get(r0.Name)
	if err != nil {
		t.Errorf("Unexpected err %s on Get", err)
	}
	if !reflect.DeepEqual(r, r0) {
		t.Errorf("Get returned %#v, expected %#v", r, r0)
	}

	// Remove player from room
	delete(r0.Usernames, "p2")
	r, err = rms.DeletePlayer(r0.Name, "p2")
	if err != nil {
		t.Errorf("Unexpected err %s on DeletePlayer", err)
	}
	if !reflect.DeepEqual(r, r0) {
		t.Errorf("Get returned %#v, expected %#v", r, r0)
	}

	// Retrieve existing room
	r, err = rms.Get(r0.Name)
	if err != nil {
		t.Errorf("Unexpected err %s on Get", err)
	}
	if !reflect.DeepEqual(r, r0) {
		t.Errorf("Get returned %#v, expected %#v", r, r0)
	}

	// Set room's game
	r0.GameType = rooms.Set
	r0.GameID = uuid.New()
	r, err = rms.SetGame(r0.Name, r0.GameType, r0.GameID)
	if err != nil {
		t.Errorf("Unexpected err %s on SetTame", err)
	}
	if !reflect.DeepEqual(r, r0) {
		t.Errorf("Get returned %#v, expected %#v", r, r0)
	}

	// Retrieve existing room
	r, err = rms.Get(r0.Name)
	if err != nil {
		t.Errorf("Unexpected err %s on Get", err)
	}
	if !reflect.DeepEqual(r, r0) {
		t.Errorf("Get returned %#v, expected %#v", r, r0)
	}

	// Delete existing room
	err = rms.Delete(r0.Name)
	if err != nil {
		t.Errorf("Unexpected err %s on Delete", err)
	}

	// Retrieve non-existing room
	_, err = rms.Get(r0.Name)
	_, ok = err.(errors.NotFoundError)
	if !ok {
		t.Errorf("Expected Get err %s to be of type NotFoundError", err)
	}

	// Delete non-existing room
	err = rms.Delete(r0.Name)
	_, ok = err.(errors.NotFoundError)
	if !ok {
		t.Errorf("Expected Delete err %s to be of type NotFoundError", err)
	}
}

func roomsEqual(rs1, rs2 []*rooms.Room) bool {
	if len(rs1) != len(rs2) {
		return false
	}
	m1 := make(map[string]*rooms.Room)
	m2 := make(map[string]*rooms.Room)
	for i := 0; i < len(rs1); i++ {
		m1[rs1[i].Name] = rs1[i]
		m2[rs2[i].Name] = rs2[i]
	}
	return reflect.DeepEqual(m1, m2)
}
