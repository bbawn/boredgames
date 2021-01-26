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

	// List with no games
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

	// Create a couple of games
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
		http.Error(w, fmt.Sprintf("Failed to unmarshal game: %s", err), http.StatusInternalServerError)
		return
	}
	// TODO: checkGame?
	// if err := checkGame(g1); err != nil {
	// 	t.Errorf("Expected empty body, got %s", string(body))
	// }

	// Fail to Get non-existent game
	uuid := uuid.New()
	r = httptest.NewRequest("GET", "http://example.com/sets/"+uuid.String(), nil)
	w = httptest.NewRecorder()
	tr.ServeHTTP(w, r)
	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
	expBody = fmt.Sprintf("Failed to get game from datastore: Key %s not found in datastore\n", uuid.String())
	fmt.Printf("expBody: %s\n", expBody)
	fmt.Printf("expBody: %v\n", []byte(expBody))
	fmt.Printf("got:     %s\n", string(body))
	fmt.Printf("got:     %v\n", body)
	// if bytes.Equal(body, []byte(expBody)) {
	if string(body) != expBody {
		t.Errorf("Expected body: %s got %s", expBody, string(body))
	}

	// Get each game
	r = httptest.NewRequest("GET", "http://example.com/sets/"+uuid.String(), nil)
	w = httptest.NewRecorder()
	tr.ServeHTTP(w, r)
	resp = w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusOK, resp.StatusCode)
	}
	var g *set.Game
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&g1)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal game: %s", err), http.StatusInternalServerError)
		return
	}
	if reflect.DeepEqual(g, g1) {
		t.Errorf("Expected game from get: %v to equal inserted game %v", g, g1)
	}
	// List the games
	// Fail to Delete non-existent game
	// Delete a game
	// Next move in invalid state
	// Claim a set
	// Claim a set in invalid game state
	// Next move
}
