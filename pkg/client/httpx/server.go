package httpx

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
)

type Server struct {
	server *http.Server
}

func NewServer() *Server {
	return &Server{
		server: &http.Server{},
	}
}

func (s *Server) ListenAndServe(ctx context.Context, addr string, tlsConfig *tls.Config) error {
	s.server.Addr = addr
	s.server.TLSConfig = tlsConfig

	if err := s.server.ListenAndServe(); err != nil {
		return fmt.Errorf("ListenAndServe: %w", err)
	}

	return nil
}
