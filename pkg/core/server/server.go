package server

import (
	"fmt"
	"net/http"
	"pugo/pkg/utils/zlog"
)

// Server is the server.
type Server struct {
	dir  string
	port int
}

// New returns a new server.
func New(dir string, port int) *Server {
	return &Server{
		dir:  dir,
		port: port,
	}
}

// Run runs the server.
func (s *Server) Run() error {
	http.Handle("/", http.FileServer(http.Dir(s.dir)))
	zlog.Infof("listening on port %d, serving %s", s.port, s.dir)
	return http.ListenAndServe(":"+fmt.Sprintf("%d", s.port), nil)
}
