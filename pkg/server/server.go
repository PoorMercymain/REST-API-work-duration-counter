package server

import (
	"context"
	"net/http"
	"time"
)

type Server interface {
	Run() error
	Shutdown(ctx context.Context) error
}

type server struct {
	httpServer *http.Server
}

func New(port string, handler http.Handler) *server {
	var s = new(server)
	s.httpServer = &http.Server{
		Addr:           "0.0.0.0:" + port,
		MaxHeaderBytes: 1 << 20,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s
}

func (s *server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
