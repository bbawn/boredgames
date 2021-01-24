package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"

	"github.com/bbawn/boredgames/internal/dao/ram"
	"github.com/bbawn/boredgames/internal/router"
	"github.com/bbawn/boredgames/services"
)

var addr = flag.String("addr", ":8080", "http service address")

func logHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		log.Println(fmt.Sprintf("%q", x))
		rec := httptest.NewRecorder()
		fn(rec, r)
		log.Println(fmt.Sprintf("%q", rec.Body))

		// this copies the recorded response to the response writer
		for k, v := range rec.HeaderMap {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		rec.Body.WriteTo(w)
	}
}

func newTableRouter() *router.TableRouter {
	daoSets := ram.NewSets()
	tr := new(router.TableRouter)
	services.NewSets(daoSets, tr)
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
