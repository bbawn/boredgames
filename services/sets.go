package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bbawn/boredgames/internal/dao"
	"github.com/bbawn/boredgames/internal/router"
	"github.com/google/uuid"
)

// Sets provides the REST API for the Set board game
type Sets struct {
	dao dao.Sets
}

func NewSets(dao dao.Sets) *Sets {
	return &Sets{dao}
}

func (s *Sets) List(w http.ResponseWriter, r *http.Request) {
	games, err := s.dao.List()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load games from datastore: %v", err), http.StatusInternalServerError)
		return
	}
	var payload []byte
	err = json.Unmarshal(payload, games)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal games from datastore: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, payload)
}

func (s *Sets) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>CreateSet</h1><div>foo</div>")
}

func (s *Sets) Get(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(router.GetField(r, 0))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid set uuid %s: %v", router.GetField(r, 0), err), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "<h1>GetSet</h1><div>%s</div>", uuid)
}

func (s *Sets) Delete(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(router.GetField(r, 0))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid set uuid %s: %v", router.GetField(r, 0), err), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "<h1>DeleteSet</h1><div>%s</div>", uuid)
}

func (s *Sets) Claim(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(router.GetField(r, 0))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid set uuid %s: %v", router.GetField(r, 0), err), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "<h1>ClaimSet</h1><div>%s</div>", uuid)
}

func (s *Sets) Next(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(router.GetField(r, 0))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid set uuid %s: %v", router.GetField(r, 0), err), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "<h1>NextSet</h1><div>%s</div>", uuid)
}
