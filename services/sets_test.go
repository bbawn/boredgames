package services

import (
	"bytes"
	// "fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bbawn/boredgames/internal/dao/ram"
	// "github.com/bbawn/boredgames/internal/games/set"
)

func TestSets(t *testing.T) {
	ram := ram.NewSets()
	s := NewSets(ram)
	w := httptest.NewRecorder()

	// List with no games
	r := httptest.NewRequest("GET", "http://example.com/sets", nil)
	s.List(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusOK, resp.StatusCode)
	}
	expBody := "[]"
	if strings.TrimSpace(string(body)) != expBody {
		t.Errorf("Expected body %s, got %s", expBody, string(body))
	}

	// g, _ := set.NewGame()
	d := `{ "usernames": [ "p1", "p2", "3" ] }`
	r = httptest.NewRequest("POST", "http://example.com/sets", bytes.NewReader([]byte(d)))
	s.Create(w, r)
	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if string(body) != "" {
		t.Errorf("Expected empty body, got %s", string(body))
	}
	/*
		fmt.Println("d:", d)
		fmt.Println("statusCode:", resp.StatusCode)
		fmt.Println("resp.Header:", resp.Header.Get("Content-Type"))
		fmt.Println("body", string(body))

		// List with games
		g, _ := set.NewGame()
		d.Insert(g)
		g, _ = set.NewGame("p0_g0", "p1_g1")
		d.Insert(g)

		s.List(w, r)
		resp = w.Result()
		body, _ = ioutil.ReadAll(resp.Body)

		fmt.Printf("d: %#v\n", d)
		fmt.Println("statusCode:", resp.StatusCode)
		fmt.Println("resp.Header:", resp.Header.Get("Content-Type"))
		fmt.Println("body", string(body))
	*/
}
