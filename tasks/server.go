package tasks

import (
	"errors"
	"fmt"
	"net/http"
)

type RouteMethod string

const (
	GET  RouteMethod = "GET"
	POST RouteMethod = "POST"
)

type Server struct {
	server *http.Server
	muxes  map[Route]*http.ServeMux
}

type RouteHandler = func(http.ResponseWriter, *http.Request)

type Route struct {
	Method RouteMethod
	Path   string
}

func NewServer() *Server {
	return &Server{
		nil,
		make(map[Route]*http.ServeMux),
	}
}

func (s *Server) Start(port int) error {
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.getHandler(),
	}
	return s.server.ListenAndServe()
}

func (s *Server) Stop() error {
	if s.server == nil {
		return errors.New("Server is not running")
	}
	return s.server.Shutdown(nil)
}

func (s *Server) ServeRoute(route Route, handler RouteHandler) {
	mux := http.NewServeMux()
	mux.Handle(string(route.Method)+" "+route.Path, http.HandlerFunc(handler))

	s.muxes[route] = mux
}

func (s *Server) ServeDir(dirPath string, rootUri string) {
	mux := http.NewServeMux()
	handler := http.FileServer(http.Dir(dirPath))
	mux.Handle(rootUri, handler)

	s.muxes[Route{GET, rootUri}] = mux
}

func (s *Server) getHandler() http.Handler {
	mux := http.NewServeMux()

	for route, handler := range s.muxes {
		mux.Handle(string(route.Method)+" "+route.Path, http.StripPrefix(route.Path, handler))
	}

	return mux
}
