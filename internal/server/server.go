package server

import (
	"context"
	"errors"
	"net/http"

	"ozon/pkg/logger"
	"time"
)

const (
	httpPort = ":8080"
)

type Server struct {
	httpServer *http.Server
}

func New(handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    httpPort,
			Handler: handler,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	logs := logger.GetLogger()
	logs.Info("Starting server")

	err := s.httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return s.httpServer.Shutdown(ctx)
}
