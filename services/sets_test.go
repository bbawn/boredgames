package services

import (
	"bytes"
	// "fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/bbawn/boredgames/internal/dao/ram"
	"github.com/bbawn/boredgames/internal/router"
	// "github.com/bbawn/boredgames/internal/games/set"
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
	body, _ = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if string(body) != "" {
		t.Errorf("Expected empty body, got %s", string(body))
	}

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
	// XXX This is the error message - that what we want? I think so...
	if string(body) != "" {
		t.Errorf("Expected empty body, got %s", string(body))
	}

	// Get each game
	// List the games
	// Fail to Delete non-existent game
	// Delete a game
	// Next move in invalid state
	// Claim a set
	// Claim a set in invalid game state
	// Next move
}
