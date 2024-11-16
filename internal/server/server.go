package server

import (
	"context"
	"net/http"
	"time"

	"github.com/paudarco/referral-api/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(cfg config.Server, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           cfg.Host + ":" + cfg.Port,
		MaxHeaderBytes: 1 << 20, // 1 MB
		Handler:        handler,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
