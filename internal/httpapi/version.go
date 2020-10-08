package httpapi

import (
	"net/http"
	"time"
)

// these global variables are set by the build pipeline
// see build/package/service/Dockerfile.
var (
	revision   string //nolint:gochecknoglobals
	pipelineID string //nolint:gochecknoglobals
	buildDate  string //nolint:gochecknoglobals
	launchDate string //nolint:gochecknoglobals
)

// versionResponse is the struct representing version responses
type versionResponse struct {
	Revision   string `json:"revision"`
	PipelineID string `json:"pipelineId"`
	BuildDate  string `json:"buildDate"`
	LaunchDate string `json:"launchDate"`
}

// handleVersion returns an HTTP handler for the version endpoint
func (s *Service) handleVersion() http.HandlerFunc {
	// Custom initialisation for handlers happens here
	launchDate = time.Now().UTC().Format(time.RFC3339)

	return func(w http.ResponseWriter, r *http.Request) {
		response := versionResponse{
			Revision:   revision,
			PipelineID: pipelineID,
			BuildDate:  buildDate,
			LaunchDate: launchDate,
		}
		s.respond(w, r, response, http.StatusOK)
	}
}
