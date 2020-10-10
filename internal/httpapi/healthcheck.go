package httpapi

import (
	"net/http"
	"time"

	"github.com/etherlabsio/healthcheck"
)

type healthcheckResponse struct {
	Status string            `json:"status,omitempty"`
	Errors map[string]string `json:"errors,omitempty"`
}

// handleLiveness returns an http.HandlerFunc which performs all checks that determine
// whether the service is alive.
//
// If this check fails, the container for this service will be destroyed and recreated. Only fail
// these checks if recreating is what you want, otherwise use the readiness endpoint.
func (s *Service) handleLiveness() http.HandlerFunc {
	return healthcheck.HandlerFunc(
		// WithTimeout allows you to set a max overall timeout.
		healthcheck.WithTimeout(5*time.Second),

		// Checkers will fail the status in case of an error.
		// Since we're talking about a SQLite database, it makes sense to kill the container
		// in this case and have it create a new, empty database.
		healthcheck.WithChecker(
			"database", &databaseChecker{service: s},
		),

		// Observers (as opposed to checkers) do not fail the status in case of an error.
		// healthcheck.WithObserver(
		// 	"example", customChecker{},
		// ),
	)
}

// handleReadiness returns an http.HandlerFunc which performs all checks that determine
// whether the service is ready to take traffic.
//
// Checks here should fail e.g. when dependencies are down, but where destroying and recreating this
// container won't help.
func (s *Service) handleReadiness() http.HandlerFunc {
	return healthcheck.HandlerFunc(
		// WithTimeout allows you to set a max overall timeout.
		healthcheck.WithTimeout(5*time.Second),

		// Checkers will fail the status if they return an error
		healthcheck.WithChecker(
			"example", customChecker{},
		),

		// Observers (as opposed to checkers) do not fail the status in case of an error.
		// healthcheck.WithObserver(
		// 	"example", customChecker{},
		// ),
	)
}
