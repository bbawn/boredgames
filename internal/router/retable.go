package router

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

// TableRouter is regexp table-based http router
// Inspired by https://github.com/benhoyt/go-routing/tree/master/retable
type TableRouter struct {
	routes []route
}

// AddRoute adds given handler for route matching given pattern and method
func (tr *TableRouter) AddRoute(method, pattern string, handler http.HandlerFunc) {
	tr.routes = append(tr.routes, newRoute(method, pattern, handler))
}

// Serve routes the given request
func (tr *TableRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range tr.routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

type ctxKey struct{}

// GetField returns the matched field at given index from the URL path
func GetField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[index]
}
