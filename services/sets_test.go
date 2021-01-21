package services

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bbawn/boredgames/internal/dao/ram"
	"github.com/bbawn/boredgames/internal/games/set"
)

func TestSets(t *testing.T) {
	d := ram.NewSets()
	s := NewSets(d)
	w := httptest.NewRecorder()

	g, _ := set.NewGame()
	r := httptest.NewRequest("POST", "http://example.com/sets", nil)
	d.Insert(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if string(body) != "" {
		t.Errorf("Expected empty body, got %s", string(body))
	}
	// List with no games
	/*
		r := httptest.NewRequest("GET", "http://example.com/sets", nil)
		s.List(w, r)
		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

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
