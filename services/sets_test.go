package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/bbawn/boredgames/internal/dao/ram"
	"github.com/bbawn/boredgames/internal/games/set"
	"github.com/bbawn/boredgames/internal/router"
)

func TestSets(t *testing.T) {
	ram := ram.NewSets()
	tr := new(router.TableRouter)
	SetsAddRoutes(ram, tr)

	t.Log("List with no games")
	r := httptest.NewRequest("GET", "http://example.com/sets", nil)
	w := httptest.NewRecorder()
	tr.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusOK, resp.StatusCode)
	}
	expBody := "[]"
	if strings.TrimSpace(string(body)) != expBody {
		t.Errorf("Expected body %s, got %s", expBody, string(body))
	}

	t.Log("Create a game with invalid payload - TODO")
	t.Log("Create a couple of games")
	d := `{ "usernames": [ "p1", "p2", "p3" ] }`
	r = httptest.NewRequest("POST", "http://example.com/sets", bytes.NewReader([]byte(d)))
	w = httptest.NewRecorder()
	tr.ServeHTTP(w, r)
	resp = w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusOK, resp.StatusCode)
	}
	var g1 *set.Game
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&g1)
	if err != nil {
		t.Errorf("Failed to decode g1: %s", err)
		return
	}
	if err := checkNewGame(g1, "p1", "p2", "p3"); err != nil {
		t.Errorf("checkNewGame failed for g1: %s", err)
	}

	d = `{ "usernames": [ "p2", "p0" ] }`
	r = httptest.NewRequest("POST", "http://example.com/sets", bytes.NewReader([]byte(d)))
	w = httptest.NewRecorder()
	tr.ServeHTTP(w, r)
	resp = w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusOK, resp.StatusCode)
	}
	var g2 *set.Game
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&g2)
	if err != nil {
		t.Errorf("Failed to decode g2: %s", err)
		return
	}
	if err := checkNewGame(g2, "p2", "p0"); err != nil {
		t.Errorf("checkNewGame failed for g2: %s", err)
	}

	t.Log("Fail to Get non-existent game")
	uid := uuid.New()
	r = httptest.NewRequest("GET", "http://example.com/sets/"+uid.String(), nil)
	w = httptest.NewRecorder()
	tr.ServeHTTP(w, r)
	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
	expBody = fmt.Sprintf("Failed to get game from datastore: Key %s not found in datastore\n", uid.String())
	if string(body) != expBody {
		t.Errorf("Expected body: %s got %s", expBody, string(body))
	}

	t.Log("Get each game")
	r = httptest.NewRequest("GET", "http://example.com/sets/"+g1.ID.String(), nil)
	w = httptest.NewRecorder()
	tr.ServeHTTP(w, r)
	resp = w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusOK, resp.StatusCode)
	}
	var g *set.Game
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&g)
	if err != nil {
		t.Errorf("Failed to decode g: %s", err)
		return
	}
	if !reflect.DeepEqual(g, g1) {
		t.Errorf("Expected game from get: %v to equal inserted game %v", g, g1)
	}

	r = httptest.NewRequest("GET", "http://example.com/sets/"+g2.ID.String(), nil)
	w = httptest.NewRecorder()
	tr.ServeHTTP(w, r)
	resp = w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusOK, resp.StatusCode)
	}
	g = nil
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&g)
	if err != nil {
		t.Errorf("Failed to decode g: %s", err)
		return
	}
	if !reflect.DeepEqual(g, g2) {
		t.Errorf("Expected game from get: %v to equal inserted game %v", g, g2)
	}

	t.Log("List the games")
	r = httptest.NewRequest("GET", "http://example.com/sets", nil)
	w = httptest.NewRecorder()
	tr.ServeHTTP(w, r)
	resp = w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusOK, resp.StatusCode)
	}
	dec = json.NewDecoder(resp.Body)
	var gs []*set.Game
	err = dec.Decode(&gs)
	if err != nil {
		t.Errorf("failed to decode games: %s", err)
		return
	}
	expGs := gameMap(g1, g2)
	actualGs := gameMap(gs...)
	if !reflect.DeepEqual(expGs, actualGs) {
		t.Errorf("Expected games from get: %v to equal inserted games %v", actualGs, expGs)
	}

	t.Log("Fail to Delete non-existent game")
	uid = uuid.New()
	r = httptest.NewRequest("DEL", "http://example.com/sets/"+uid.String(), nil)
	w = httptest.NewRecorder()
	tr.ServeHTTP(w, r)
	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
	expBody = fmt.Sprintf("Failed to delete game from datastore: Key %s not found in datastore\n", uid.String())
	if string(body) != expBody {
		t.Errorf("Expected body: %s got %s", expBody, string(body))
	}

	t.Log("Delete a game")
	r = httptest.NewRequest("DEL", "http://example.com/sets/"+g2.ID.String(), nil)
	w = httptest.NewRecorder()
	tr.ServeHTTP(w, r)
	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if string(body) != "" {
		t.Errorf("Expected empty body: got %s", string(body))
	}

	t.Log("Next move in invalid state")
	r = httptest.NewRequest("POST", "http://example.com/sets/"+g1.ID.String()+"/next", nil)
	w = httptest.NewRecorder()
	tr.ServeHTTP(w, r)
	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusConflict, resp.StatusCode)
	}
	expBody = fmt.Sprintf("Failed to advance game to next round: Invalid method: NextRound detail: round not yet claimed\n")
	if string(body) != expBody {
		t.Errorf("Expected body: %s got %s", expBody, string(body))
	}

	t.Log("Claim a set with invalid data")
	cd := claimData{}
	payload, err := json.Marshal(&cd)
	if err != nil {
		t.Errorf("Unexpected Marshal err: %s", err)
		return
	}
	r = httptest.NewRequest("POST", "http://example.com/sets/"+g1.ID.String()+"/claim", bytes.NewReader(payload))
	w = httptest.NewRecorder()
	tr.ServeHTTP(w, r)
	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusConflict, resp.StatusCode)
	}
	expBody = fmt.Sprintf("Failed to claim set in game: Invalid value:  for arg: username\n")
	if string(body) != expBody {
		t.Errorf("Expected body: %s got %s", expBody, string(body))
	}

	t.Log("Claim a set")
	s1 := g1.FindExpandSet()
	cd = claimData{
		Username: "p1",
		Card1:    s1[0],
		Card2:    s1[1],
		Card3:    s1[2],
	}
	payload, err = json.Marshal(&cd)
	if err != nil {
		t.Errorf("Unexpected Marshal err: %s", err)
		return
	}
	r = httptest.NewRequest("POST", "http://example.com/sets/"+g1.ID.String()+"/claim", bytes.NewReader(payload))
	w = httptest.NewRecorder()
	tr.ServeHTTP(w, r)
	resp = w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusOK, resp.StatusCode)
	}
	var g1Claimed *set.Game
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&g1Claimed)
	if err != nil {
		t.Errorf("Failed to decode g1Claimed: %s", err)
		return
	}
	err = checkNextGame(g1, g1Claimed)
	if err != nil {
		t.Errorf("Invalid games after Claim:  %#v, %#v: err %s", g1, g1Claimed, err)
	}

	t.Log("Claim a set in invalid game state")
	r = httptest.NewRequest("POST", "http://example.com/sets/"+g1.ID.String()+"/claim", bytes.NewReader(payload))
	w = httptest.NewRecorder()
	tr.ServeHTTP(w, r)
	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusConflict, resp.StatusCode)
	}
	expBody = fmt.Sprintf("Failed to claim set in game: Invalid method: ClaimSet detail: round already claimed by p1\n")
	if string(body) != expBody {
		t.Errorf("Expected body: %s got %s", expBody, string(body))
	}

	t.Log("Valid Next round request")
	r = httptest.NewRequest("POST", "http://example.com/sets/"+g1.ID.String()+"/next", nil)
	w = httptest.NewRecorder()
	tr.ServeHTTP(w, r)
	resp = w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusOK, resp.StatusCode)
	}
	var g1Next *set.Game
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&g1Next)
	if err != nil {
		t.Errorf("Failed to decode g1Next: %s", err)
		return
	}
	err = checkNextGame(g1Claimed, g1Next)
	if err != nil {
		t.Errorf("Invalid games after Next:  %#v, %#v: err %s", g1Claimed, g1Next, err)
	}
}

// checkNewGame validates that the given game is in a valid initial state
func checkNewGame(g *set.Game, usernames ...string) error {
	if g.ID.URN() == "" {
		return fmt.Errorf("Invalid ID: %s", g.ID)
	}
	if len(g.Players) != len(usernames) {
		return fmt.Errorf("Expected %d players, got %d", len(g.Players), len(usernames))
	}
	for _, u := range usernames {
		var (
			p  *set.Player
			ok bool
		)
		if p, ok = g.Players[u]; !ok {
			return fmt.Errorf("Player with username %s not found", u)
		}
		if p.Username != u {
			return fmt.Errorf("Expected player Username %s, got %s", u, p.Username)
		}
		if len(p.Sets) != 0 {
			return fmt.Errorf("Expected empty Sets for player %s", u)
		}
	}
	return nil
}

// checkNextGame validates that g1 is a valid next state of g0
func checkNextGame(g0, g1 *set.Game) error {
	if g0.ID != g1.ID {
		return fmt.Errorf("Expected g0 ID %s to equal g1 ID %s", g0.ID, g1.ID)
	}
	if g0.GetState() == set.Playing {
		if g1.GetState() != set.SetClaimed {
			return fmt.Errorf("Expected g1 State to be SetClaimed")
		}
	} else {
		if g1.GetState() != set.Playing {
			return fmt.Errorf("Expected g1 State to be Playing")
		}
	}
	// TODO? we could check invariants on players, scores, deck size, etc...
	return nil
}

func gameMap(gs ...*set.Game) map[uuid.UUID]*set.Game {
	m := make(map[uuid.UUID]*set.Game)
	for _, g := range gs {
		m[g.ID] = g
	}
	return m
}
