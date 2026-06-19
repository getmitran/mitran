package shutdown

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"
)

func freePort() string {
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	defer l.Close()
	return l.Addr().String()
}

func TestListenAndServeGraceful_RespondsToRequests(t *testing.T) {
	addr := freePort()
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	})

	errCh := make(chan error, 1)
	go func() { errCh <- ListenAndServeGraceful(addr, mux, 5*time.Second) }()

	// Wait for server to be ready
	var resp *http.Response
	var err error
	for i := 0; i < 50; i++ {
		resp, err = http.Get("http://" + addr + "/ping")
		if err == nil {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if err != nil {
		t.Fatalf("server never became ready: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// Send SIGTERM to trigger shutdown
	p, _ := os.FindProcess(os.Getpid())
	p.Signal(syscall.SIGTERM)

	select {
	case e := <-errCh:
		if e != nil {
			t.Fatalf("unexpected error: %v", e)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("server did not shut down in time")
	}
}

func TestListenAndServeGraceful_ShutdownOnSignal(t *testing.T) {
	addr := freePort()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	errCh := make(chan error, 1)
	go func() { errCh <- ListenAndServeGraceful(addr, handler, 2*time.Second) }()

	// Wait for server ready
	for i := 0; i < 50; i++ {
		if conn, err := net.DialTimeout("tcp", addr, 50*time.Millisecond); err == nil {
			conn.Close()
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	// Signal shutdown
	p, _ := os.FindProcess(os.Getpid())
	p.Signal(syscall.SIGTERM)

	select {
	case e := <-errCh:
		if e != nil {
			t.Fatalf("shutdown returned error: %v", e)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("server did not shut down within timeout")
	}

	// Verify server no longer accepts connections
	_, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
	if err == nil {
		t.Fatal("server still accepting connections after shutdown")
	}

	// Suppress unused import
	_ = context.Background()
}
