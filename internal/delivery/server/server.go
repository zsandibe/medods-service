package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/zsandibe/medods-service/config"
	logger "github.com/zsandibe/medods-service/pkg"
)

type Server struct {
	httpServer http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: http.Server{
			Addr:           fmt.Sprintf("0.0.0.0:%s", cfg.Server.Port),
			Handler:        handler,
			MaxHeaderBytes: 1024 * 1024,
			ReadTimeout:    15 * time.Second,
			WriteTimeout:   15 * time.Second,
		},
	}
}

func (s *Server) Run() error {
	logger.Info("Starting server on: ", s.httpServer.Addr)

	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down server...")
	if err := s.httpServer.Close(); err != nil {
		return err
	}
	return nil
}
