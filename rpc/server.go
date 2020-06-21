package rpc

import (
	"fmt"
	"github.com/smallnest/rpcx/server"
)

type Server struct {
	address   string
	port      string
	rpcServer *server.Server
}

func NewServer(address string, port string) (*Server, error) {
	s := &Server{
		rpcServer: server.NewServer(),
		address:   address,
		port:      port,
	}
	if err := s.rpcServer.Register(new(PSO), ""); err != nil {
		return nil, err
	}
	return s, nil
}
func (s *Server) Serve() error {
	if err := s.rpcServer.Serve("tcp", fmt.Sprintf("%s:%s", s.address, s.port)); err != nil {
		return err
	}
	return nil
}
