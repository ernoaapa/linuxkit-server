package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	router http.Handler
	port   int
}

// New creates new API server
func New(port int) *Server {
	router := mux.NewRouter()
	router.Methods("POST").Path("/linuxkit/{name}/build/{format}").HandlerFunc(routeHandler)

	return &Server{
		router: router,
		port:   port,
	}
}

func routeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	output, err := getOutputFormat(r)
	if err != nil {
		http.Error(w, err.Error(), 400)
	} else {
		createBuild(vars["name"], vars["format"], output, w, r)
	}
}

func getOutputFormat(r *http.Request) (string, error) {
	switch output := r.URL.Query().Get("output"); output {
	case "img":
		return "img", nil
	case "tar", "":
		return "tar", nil
	default:
		return "", fmt.Errorf("Invalid 'output' value %s. Must be one of [tar, img]", output)
	}
}

// Serve starts the API server
func (s *Server) Serve() error {
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Start server on %s", addr)
	return http.ListenAndServe(addr, handlers.CompressHandler(s.router))
}
