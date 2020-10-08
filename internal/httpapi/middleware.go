package httpapi

import (
	"log"
	"net/http"
)

// logRequest logs the method and URL of a request
//
// It's mostly a demo middleware trying to be useful.
func (s *Service) logRequest(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL)
		h(w, r)
	}
}
