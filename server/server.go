package server

import (
	"context"
	"errors"
	"fmt"
	"net"
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
	mux.HandleFunc("GET /healthz", api.Healthz)

	mux.HandleFunc("GET /api/tables", api.Tables(drv))
	mux.HandleFunc("GET /api/tables/{table}/data", api.TableData(drv))

	return &Server{
		srv: &http.Server{
			Addr:              addr,
			Handler:           logRequest(mux),
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
	}
}

func (s *Server) RunAndShut(ctx context.Context, shutdownTimeout time.Duration) error {
	l, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	go func() { errCh <- s.srv.Serve(l) }()

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve %s: %w", s.srv.Addr, err)
		}
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		_ = s.srv.Close()
		return fmt.Errorf("graceful shutdown: %w", err)
	}

	return nil
}
