package server

import "net/http"

type Server struct {
	mux *http.ServeMux
}

// InitializeServer returns an initialized Server object, with a fresh
// http.ServeMux.
func (s Server) InitializeServer() *Server {
	server := new(Server)
	server.mux = http.NewServeMux()
	return server
}

// GetMux returns the Server's mux object
func (s *Server) GetMux() *http.ServeMux {
	return s.mux
}

// AssignHandler assigns a given route to the desired handler function
func (s *Server) AssignHandler(route string, handlerFn http.HandlerFunc) {
	s.mux.HandleFunc(route, handlerFn)
}
