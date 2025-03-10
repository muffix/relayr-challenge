package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/muffix/relayr-challenge/internal/database"
	"github.com/muffix/relayr-challenge/internal/review"
)

// Service is the struct representing the service.
//
// Dependencies such as database connections should be added here. They can be made available to
// handler funcs by implementing the handlers with a service receiver (see the version, readiness,
// and liveness endpoints for an example).
//
// This pattern also avoids global state.
type Service struct {
	server *http.Server
	router *mux.Router

	offers   database.Offers
	reviewer review.Reviewer
}

// NewService returns a new service struct.
//
// This only sets up the routes, but shouldn't set up dependencies. That way, testing is easier
// since all of them can be mocked or added later.
func NewService(servicePort int) *Service {
	router := mux.NewRouter().StrictSlash(true)

	service := &Service{
		server: createServerWithRouter(router, servicePort),
		router: router,
	}

	service.routes()

	return service
}

// SetDatabase is a setter for the database
func (s *Service) SetDatabase(db database.Offers) {
	s.offers = db
}

// SetReviewer is a setter for a client of a reviews engine
func (s *Service) SetReviewer(r review.Reviewer) {
	s.reviewer = r
}

func createServerWithRouter(router http.Handler, port int) *http.Server {
	return &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

// Start starts the HTTP service
func (s *Service) Start() {
	go func() {
		log.Printf("Serving on port %s", s.server.Addr)
		err := s.server.ListenAndServe()
		if err != nil {
			log.Fatalf("Error from router %s", err.Error())
		}
	}()
	defer s.close()

	// Handle interrupts
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}

func (s *Service) close() {
	_ = s.server.Shutdown(nil)
}

// respond is a helper function to create a response for an encodable struct. It sets the content
// type and response code.
func (s *Service) respond(w http.ResponseWriter, _ *http.Request, data interface{}, status int) {
	if data == nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// decode is a helper function that decodes request data into a struct
func (s *Service) decode(_ http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
