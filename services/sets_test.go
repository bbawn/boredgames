package services

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/bbawn/boredgames/internal/dao/ram"
	"github.com/bbawn/boredgames/internal/games/set"
)

func TestSets(t *testing.T) {
	d := ram.NewSets()
	s := NewSets(d)
	w := httptest.NewRecorder()

	// List with no games
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
}
