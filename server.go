package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	srv *http.Server
	l   net.Listener
}

func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Handler: mux,
		},
		l: l,
	}
}

func (s *Server) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)
	// Start HTTP server in another goroutine
	eg.Go(func() error {
		// http.ErrServerClosed is returned when the server is closed normally by http.Sever.Shutdown()
		if err := s.srv.Serve(s.l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	// Wait for interrupt signal from the cannel
	<-ctx.Done()
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	// Wait for other goroutine to finish
	return eg.Wait()
}