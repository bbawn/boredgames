package dao

import (
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/bbawn/boredgames/internal/dao/ram"
	"github.com/bbawn/boredgames/internal/games/set"
)

func TestSets(t *testing.T) {
	ram := ram.NewSets()
	daoTest(t, ram)
	// sqllite := sqllite.NewSets()
	// daoTest(T, sqllite)
	// postgres := postgres.NewSets()
	// daoTest(T, postgres)
}

func daoTest(t *testing.T, s Sets) {
	// Empty list
	expGs := []*set.Game{}
	gs, err := s.List()
	if !gamesEqual(gs, expGs) {
		t.Errorf("Expected %#v to equal %#v", gs, expGs)
	}

	// Insert game with players
	g0, _ := set.NewGame("p0", "p1")
	err = s.Insert(g0)
	if err != nil {
		t.Errorf("Unexpected err %s on Insert", err)
	}

	// Duplicate Insert fails
	err = s.Insert(g0)
	_, ok := err.(AlreadyExistsError)
	if !ok {
		t.Errorf("Expected err %s to be of type AlreadyExistsError", err)
	}

	// Insert game without players
	g1, _ := set.NewGame()
	err = s.Insert(g1)
	if err != nil {
		t.Errorf("Unexpected err %s on Insert", err)
	}

	// Retrieve existing game
	g, err := s.Get(g0.ID)
	if err != nil {
		t.Errorf("Unexpected err %s on Get", err)
	}
	if !reflect.DeepEqual(g, g0) {
		t.Errorf("Get returned %#v, expected equal to %#v", g, g0)
	}

	// Retrieve existing game
	g, err = s.Get(g1.ID)
	if err != nil {
		t.Errorf("Unexpected err %s on Get", err)
	}
	if !reflect.DeepEqual(g, g1) {
		t.Errorf("Get returned %#v, expected equal to %#v", g, g1)
	}

	// List all games
	expGs = []*set.Game{g0, g1}
	gs, err = s.List()
	if !gamesEqual(gs, expGs) {
		t.Errorf("Expected %#v to equal %#v", gs, expGs)
	}

	// Update game
	set := g0.FindExpandSet()
	g0.ClaimSet("p0", set[0], set[1], set[2])
	err = s.Update(g0)
	if err != nil {
		t.Errorf("Unexpected err %s on Insert", err)
	}

	// Retrieve existing game
	g, err = s.Get(g0.ID)
	if err != nil {
		t.Errorf("Unexpected err %s on Get", err)
	}
	if !reflect.DeepEqual(g, g0) {
		t.Errorf("Get returned %#v, expected equal to %#v", g, g0)
	}

	// Delete existing game
	err = s.Delete(g0.ID)
	if err != nil {
		t.Errorf("Unexpected err %s on Delete", err)
	}

	// Retrieve non-existing game
	g, err = s.Get(g0.ID)
	_, ok = err.(NotFoundError)
	if !ok {
		t.Errorf("Expected Get err %s to be of type NotFoundError", err)
	}

	// Update non-existing game
	err = s.Update(g0)
	_, ok = err.(NotFoundError)
	if !ok {
		t.Errorf("Expected Update err %s to be of type NotFoundError", err)
	}

	// Delete non-existing game
	err = s.Delete(g0.ID)
	_, ok = err.(NotFoundError)
	if !ok {
		t.Errorf("Expected Delete err %s to be of type NotFoundError", err)
	}
}

func gamesEqual(gs1, gs2 []*set.Game) bool {
	if len(gs1) != len(gs2) {
		return false
	}
	m1 := make(map[uuid.UUID]*set.Game)
	m2 := make(map[uuid.UUID]*set.Game)
	for i := 0; i < len(gs1); i++ {
		m1[gs1[i].ID] = gs1[i]
		m2[gs2[i].ID] = gs2[i]
	}
	return reflect.DeepEqual(m1, m2)
}
