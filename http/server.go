package http

import (
	"context"
	"net/http"
	"time"

	"github.com/mroobert/monorepo-tixer/http/mid"
)

// ServerConfig represents the web server configuration details.
type ServerConfig struct {
	IdleTimeout     time.Duration
	WriteTimeout    time.Duration
	ReadTimeout     time.Duration
	ShutdownTimeout time.Duration
	Addr            string
	DebugAddr       string
	MaxBodySize     int64
}

// Server represents an HTTP server. It is meant to wrap all HTTP functionality
// used by the application so that dependent packages (such as main) do not
// need to reference the "net/http" package at all.
type Server struct {
	server *http.Server
	router *http.ServeMux

	MaxBodySize int64
}

// NewServer creates a new server with the provided configuration.
func NewServer(cfg ServerConfig) *Server {
	s := &Server{
		server: &http.Server{
			Addr:         cfg.Addr,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		router: http.NewServeMux(),

		MaxBodySize: cfg.MaxBodySize,
	}

	s.registerTicketRoutes(s.router)

	s.server.Handler = mid.Panics(mid.Logger(s.router))
	return s
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Close immediately closes the server.
func (s *Server) Close() {
	s.server.Close()
}
