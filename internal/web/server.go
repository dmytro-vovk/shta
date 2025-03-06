package web

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

func New(listen string, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  30 * time.Second,
			Addr:         listen,
			Handler:      handler,
		},
	}
}

func (w *Server) Serve() (err error) {
	err = w.server.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	}

	return
}

func (w *Server) Stop(ctx context.Context) error {
	return w.server.Shutdown(ctx)
}
