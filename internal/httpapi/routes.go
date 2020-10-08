package httpapi

// routes is the function where routes and their handlers are added. It is meant to be used as the
// one place for all the routes to make it easy to see what's happening.
func (s *Service) routes() {
	// These are the three default routes that we must keep
	s.router.HandleFunc("/version", s.handleVersion())
	s.router.HandleFunc("/liveness", s.handleLiveness())
	s.router.HandleFunc("/readiness", s.handleReadiness())

	// New routes go here
	s.router.HandleFunc("/", s.logRequest(s.handleHomePage()))
}
