package services

import (
	"fmt"
	"net/http"
)

func SetContainerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>SetContainerHandler</h1><div>foo</div>")
}

func SetObjectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>SetContainerHandler</h1><div>foo</div>")
}
