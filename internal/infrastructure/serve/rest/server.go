package rest

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	defaultServerReadTimeout  = 5 * time.Second
	defaultServerWriteTimeout = 55 * time.Second
)

type server struct {
	server *http.Server
	log    *log.Logger
}

func NewServer(handler http.Handler, host, port string, log *log.Logger) *server {
	return &server{
		server: &http.Server{
			Addr:         host + ":" + port,
			ReadTimeout:  defaultServerReadTimeout,
			WriteTimeout: defaultServerWriteTimeout,
			Handler:      handler,
		},
		log: log,
	}
}

func (s *server) ListenAndServe() {
	go func() {
		s.log.Printf("HTTP server running on %s!", s.server.Addr)

		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Printf("error on listen and serve: %v", err)
		}
	}()
}

func (s *server) Shutdown(d time.Duration) error {
	s.log.Printf("shutting down HTTP server running on %s!", s.server.Addr)

	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("could not shutdown HTTP server: %w", err)
	}

	return nil
}
