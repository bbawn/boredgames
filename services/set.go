package services

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

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
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.Printf("SetObjectHandler: id %#v, verb %s", id, verb)
	switch r.Method {
	case "GET":
		getSet(id)
	case "DEL":
		deleteSet(id)
	default:
	}
	fmt.Fprintf(w, "<h1>SetObjectHandler</h1><div>bar</div>")
}

var validObjectPath = regexp.MustCompile(`^/set/(\w+)/(\w+)$`)

func parseObjectURL(path string) (*uuid.UUID, string, error) {
	m := validObjectPath.FindStringSubmatch(path)
	if m == nil {
		return nil, "", fmt.Errorf("invalid object URL path: %s", path)
	}
	// XXX this parses anything: "foo" returns all 0 uuid. WTF? Consider bson.ObjectID...
	uuid, err := uuid.Parse(m[1])
	if m == nil {
		return nil, "", err
	}
	verb := m[2]
	return &uuid, verb, nil
}

func createSet() {
}

func deleteSet(uuid *uuid.UUID) {
}

func getSet(uuid *uuid.UUID) {
}

func listSets() {
}
