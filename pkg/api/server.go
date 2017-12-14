package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router http.Handler
	port   int
}

// New creates new API server
func New(port int) *Server {
	router := mux.NewRouter()
	router.HandleFunc("/build", createBuild).Methods("POST")

	return &Server{
		router: router,
		port:   port,
	}
}

// Serve starts the API server
func (s *Server) Serve() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router)
}
