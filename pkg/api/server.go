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
	router.HandleFunc("/linuxkit/{name}/build/{format}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		createBuild(vars["name"], vars["format"], w, r)
	}).Methods("POST")

	return &Server{
		router: router,
		port:   port,
	}
}

// Serve starts the API server
func (s *Server) Serve() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router)
}
