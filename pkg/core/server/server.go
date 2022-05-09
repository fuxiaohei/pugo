package server

import (
	"fmt"
	"net/http"
	"pugo/pkg/utils/zlog"
)

// Server is the server.
type Server struct {
	opt ServerOption
}

type ServerOption struct {
	Port int
	Dir  string
}

// New returns a new server.
func New(opt ServerOption) *Server {
	return &Server{
		opt: opt,
	}
}

// Run runs the server.
func (s *Server) Run() error {
	http.Handle("/", http.FileServer(http.Dir(s.opt.Dir)))
	zlog.Infof("listening on port %d, serving %s", s.opt.Port, s.opt.Dir)
	return http.ListenAndServe(":"+fmt.Sprintf("%d", s.opt.Port), nil)
}
