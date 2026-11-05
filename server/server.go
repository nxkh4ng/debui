package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/nxkh4ng/debui/api"
	"github.com/nxkh4ng/debui/driver"
)

type Server struct {
	srv *http.Server
}

func New(addr string, drv driver.Driver) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", api.Heathz)

	mux.HandleFunc("GET /api/tables", api.Tables(drv))

	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

func (s *Server) RunAndShut(ctx context.Context, shutdownTimeout time.Duration) error {
	errCh := make(chan error, 1)

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		_ = s.srv.Close()
		return err
	}

	return nil
}
