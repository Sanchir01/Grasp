package server

import (
	"context"
	"github.com/Sanchir01/Grasp/internal/config"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	config     *config.Config
}

func NewHttpServer(cfg *config.Config) *Server {
	srv := &http.Server{
		Addr:         cfg.HttpServer.Host + ":" + cfg.HttpServer.Port,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}
	return &Server{
		httpServer: srv,
		config:     cfg,
	}
}
func (s *Server) Run(handler http.Handler) error {
	s.httpServer.Handler = handler
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
