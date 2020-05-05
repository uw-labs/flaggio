package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"github.com/victorkt/flaggio/internal/service"
)

// NewServer returns a new server object
func NewServer(
	router chi.Router,
	flagsService service.Flag,
	logger *logrus.Entry,
) *Server {
	srv := &Server{
		router:       router,
		flagsService: flagsService,
		logger:       logger,
	}
	srv.routes()
	return srv
}

// Server handles evaluation requests
type Server struct {
	router       chi.Router
	flagsService service.Flag
	logger       *logrus.Entry
}

// ServeHTTP responds to an HTTP request
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Setup all routes
func (s *Server) routes() {
	// API version 1
	s.router.Route("/v1", func(r chi.Router) {
		r.Post("/evaluate", s.handleEvaluateAll)
		r.Post("/evaluate/{key}", s.handleEvaluate)
	})
}
