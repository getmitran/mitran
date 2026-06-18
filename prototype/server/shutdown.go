package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func ListenAndServeGraceful(addr string, handler http.Handler) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	srv := &http.Server{Addr: addr, Handler: handler}
	go func() {
		<-ctx.Done()
		log.Println("shutting down...")
		shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(shutCtx)
	}()
	err := srv.ListenAndServe()
	if err == http.ErrServerClosed {
		log.Println("shutdown complete")
		return nil
	}
	return err
}
