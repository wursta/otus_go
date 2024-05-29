package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

var (
	ErrServerStarted    = errors.New("server already started")
	ErrServerNotStarted = errors.New("server not started")
)

type Server struct {
	host   string
	port   string
	server *http.Server
	mux    *http.ServeMux
}

type Logger interface {
	Debug(msg string, params ...any)
	Info(msg string)
	Error(msg string)
}

type Application interface {
	CreateEvent(
		ctx context.Context,
		id, title string,
		startDate time.Time,
		endDate time.Time,
		notifyBefore time.Duration,
	) error
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{
		mux: http.NewServeMux(),
	}
}

func (s *Server) Start(_ context.Context) error {
	if s.server != nil {
		return ErrServerStarted
	}
	address := net.JoinHostPort(s.host, s.port)
	s.server = &http.Server{
		Addr:    address,
		Handler: loggingMiddleware(s.mux),
	}

	err := s.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return ErrServerNotStarted
	}

	return s.server.Close()
}

func (s *Server) AddRoute(route string, handlerFunc http.HandlerFunc) {
	s.mux.HandleFunc(route, handlerFunc)
}
