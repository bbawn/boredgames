package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/bbawn/boredgames/services"
)

var addr = flag.String("addr", ":8080", "http service address")

func newServeMux() *http.ServeMux {
	// XXX See https://golang.org/doc/articles/wiki/ for next steps...
	mux := http.NewServeMux()
	mux.Handle("/set", http.HandlerFunc(services.SetContainerHandler))
	mux.Handle("/set/", http.HandlerFunc(services.SetObjectHandler))

	return mux
}

func main() {
	flag.Parse()
	mux := newServeMux()
	srv := &http.Server{Addr: *addr, Handler: mux}

	log.Printf("INFO: ListenAndServe(): addr: %s", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("WARN: api: ListenAndServe() failed: %s", err)
	}
}
