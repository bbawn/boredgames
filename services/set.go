package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bbawn/boredgames/internal/router"
	"github.com/google/uuid"
)

func ListSets(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>ListSets</h1><div>foo</div>")
}

func CreateSet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>CreateSet</h1><div>foo</div>")
}

func GetSet(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(router.GetField(r, 0))
	log.Printf("parseObjectURL: uuid %v, err %v", uuid, err)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid set uuid %s: %v", router.GetField(r, 0), err), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "<h1>GetSet</h1><div>%s</div>", uuid)
}

func DeleteSet(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(router.GetField(r, 0))
	log.Printf("parseObjectURL: uuid %v, err %v", uuid, err)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid set uuid %s: %v", router.GetField(r, 0), err), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "<h1>DeleteSet</h1><div>%s</div>", uuid)
}

func ClaimSet(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(router.GetField(r, 0))
	log.Printf("parseObjectURL: uuid %v, err %v", uuid, err)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid set uuid %s: %v", router.GetField(r, 0), err), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "<h1>ClaimSet</h1><div>%s</div>", uuid)
}

func NextSet(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(router.GetField(r, 0))
	log.Printf("parseObjectURL: uuid %v, err %v", uuid, err)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid set uuid %s: %v", router.GetField(r, 0), err), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "<h1>NextSet</h1><div>%s</div>", uuid)
}
