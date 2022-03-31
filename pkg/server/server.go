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

func (s Server) GetMux() *http.ServeMux {
	return s.mux
}
