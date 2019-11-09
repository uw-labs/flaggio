package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/victorkt/flaggio/internal/server/api/service"
)

// NewServer returns a new server object
func NewServer(
	router chi.Router,
	flagsService service.FlagService,
) *Server {
	srv := &Server{
		router:       router,
		flagsService: flagsService,
	}
	srv.routes()
	return srv
}

// Server handles evaluation requests
type Server struct {
	router       chi.Router
	flagsService service.FlagService
}

// ServeHTTP responds to an HTTP request
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Setup all routes
func (s *Server) routes() {
	s.router.Post("/evaluate", s.handleEvaluateAll)
	s.router.Post("/evaluate/{key}", s.handleEvaluate)
}
