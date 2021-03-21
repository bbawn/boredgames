package services

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/bbawn/boredgames/internal/router"
)

func doRequest(
	tr *router.TableRouter,
	method, target string,
	reqBody io.Reader,
) *http.Response {
	r := httptest.NewRequest(method, target, reqBody)
	w := httptest.NewRecorder()
	tr.ServeHTTP(w, r)
	return w.Result()
}
