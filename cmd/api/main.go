package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/bbawn/boredgames/internal/dao/ram"
	"github.com/bbawn/boredgames/internal/router"
	"github.com/bbawn/boredgames/services"
)

var addr = flag.String("addr", ":8080", "http service address")

func logHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request", r.Method, r.RequestURI)
		rec := httptest.NewRecorder()
		fn(rec, r)
		log.Println("Response StatusCode", rec.Result().StatusCode)

		// this copies the recorded response to the response writer
		for k, v := range rec.Result().Header {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		rec.Body.WriteTo(w)
	}
}

func newTableRouter() *router.TableRouter {
	daoRooms := ram.NewRooms()
	daoSets := ram.NewSets()
	tr := new(router.TableRouter)

	// API routes
	services.RoomsAddRoutes(daoRooms, tr)
	services.SetsAddRoutes(daoSets, tr)

	// static routes
	tr.AddRoute("GET", "/.*", http.StripPrefix("/", http.FileServer(http.Dir("ui"))))
	return tr
}

func main() {
	flag.Parse()
	tr := newTableRouter()
	srv := &http.Server{Addr: *addr, Handler: logHandler(tr.ServeHTTP)}

	log.Printf("INFO: ListenAndServe(): addr: %s", *addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("WARN: api: ListenAndServe() failed: %s", err)
	}
}
