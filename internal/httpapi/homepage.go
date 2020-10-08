package httpapi

import (
	"net/http"
)

// homePageResponse is the struct representing home page responses
type homePageResponse struct {
	Message string `json:"message"`
}

// handleHomePage returns an http.HandlerFunc for the home page endpoint
func (s *Service) handleHomePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := homePageResponse{
			Message: "Hello from Go",
		}
		s.respond(w, r, response, http.StatusOK)
	}
}
