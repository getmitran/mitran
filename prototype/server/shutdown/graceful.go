package shutdown

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func ListenAndServeGraceful(addr string, handler http.Handler, timeout time.Duration) error {
	srv := &http.Server{Addr: addr, Handler: handler}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	select {
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	case <-ctx.Done():
		stop()
		shutCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := srv.Shutdown(shutCtx); err != nil {
			return err
		}
		return nil
	}
}
