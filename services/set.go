package services

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

var Uuid uuid.UUID

func SetContainerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("SetContainerHandler: r %#v", r)
	switch r.Method {
	case "GET":
		listSets()
	case "POST":
		createSet()
	default:
	}
	fmt.Fprintf(w, "<h1>SetContainerHandler</h1><div>foo</div>")
}

func SetObjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("SetObjectHandler: r %#v", r)
	id, verb, err := parseObjectURL(r.URL.Path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing set %s: %v", r.URL.Path, err), http.StatusBadRequest)
		return
	}
	log.Printf("SetObjectHandler: id %#v, verb %s", id, verb)
	switch r.Method {
	case "GET":
		err = getSet(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving set %s: %v", id, err), http.StatusBadRequest)
			return
		}
	case "DEL":
		err = deleteSet(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error deleting set %s: %v", id, err), http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("Invalid method %s: %v", r.Method, err), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "<h1>SetObjectHandler</h1><div>%s</div>", id)
}

func parseObjectURL(path string) (*uuid.UUID, string, error) {
	comps := strings.Split(path, "/")
	if len(comps) > 3 || len(comps) < 2 {
		return nil, "", fmt.Errorf("Invalid set path %s", path)
	}

	// XXX this parses anything: "foo" returns all 0 uuid. WTF? Consider bson.ObjectID...
	uuid, err := uuid.Parse(comps[1])
	if err == nil {
		return nil, "", err
	}
	var verb string
	if len(comps) > 2 {
		verb = comps[2]
	}
	return &uuid, verb, nil
}

func createSet() error {
	return nil
}

func deleteSet(uuid *uuid.UUID) error {
	return nil
}

func getSet(uuid *uuid.UUID) error {
	return nil
}

func listSets() {
}
