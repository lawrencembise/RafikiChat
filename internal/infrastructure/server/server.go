package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"rafikichat/internal/infrastructure/config"
)

// Server represents the HTTP server
type Server struct {
	httpServer *http.Server
}

// NewServer initializes a new Server with the given configuration
func NewServer(cfg config.Config) *Server {
	router := mux.NewRouter()

	// Define your routes here
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to RafikiChat!"))
	}).Methods("GET")

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
	}
}

// Start runs the HTTP server
func (s *Server) Start() error {
	log.Printf("Starting server on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}
