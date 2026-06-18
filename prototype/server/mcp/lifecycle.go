//go:build !windows

package mcp

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type ServerConfig struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Command   string            `json:"command"`
	Args      []string          `json:"args"`
	Env       map[string]string `json:"env,omitempty"`
	Enabled   bool              `json:"enabled"`
	AutoStart bool              `json:"auto_start"`
}

type ServerStatus struct {
	ID        string    `json:"id"`
	Running   bool      `json:"running"`
	PID       int       `json:"pid,omitempty"`
	StartedAt time.Time `json:"started_at,omitempty"`
	Restarts  int       `json:"restarts"`
	LastError string    `json:"last_error,omitempty"`
}

type managedServer struct {
	config    ServerConfig
	cmd       *exec.Cmd
	cancel    context.CancelFunc
	status    ServerStatus
	mu        sync.Mutex
}

type LifecycleManager struct {
	servers map[string]*managedServer
	mu      sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	healthInterval time.Duration
}

func NewLifecycleManager(healthInterval time.Duration) *LifecycleManager {
	if healthInterval == 0 {
		healthInterval = 10 * time.Second
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &LifecycleManager{
		servers:        make(map[string]*managedServer),
		ctx:            ctx,
		cancel:         cancel,
		healthInterval: healthInterval,
	}
}

func (lm *LifecycleManager) Start(config ServerConfig) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	if _, exists := lm.servers[config.ID]; exists {
		return fmt.Errorf("server %s already running", config.ID)
	}

	srv := &managedServer{
		config: config,
		status: ServerStatus{ID: config.ID},
	}
	lm.servers[config.ID] = srv

	if err := lm.startProcess(srv); err != nil {
		delete(lm.servers, config.ID)
		return err
	}

	lm.wg.Add(1)
	go lm.monitor(srv)
	return nil
}

func (lm *LifecycleManager) Stop(id string) error {
	lm.mu.RLock()
	srv, exists := lm.servers[id]
	lm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("server %s not found", id)
	}

	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.cancel != nil {
		srv.cancel()
	}
	if srv.cmd != nil && srv.cmd.Process != nil {
		_ = srv.cmd.Process.Signal(os.Interrupt)
		done := make(chan error, 1)
		go func() { done <- srv.cmd.Wait() }()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			_ = srv.cmd.Process.Kill()
		}
	}
	srv.status.Running = false

	lm.mu.Lock()
	delete(lm.servers, id)
	lm.mu.Unlock()
	return nil
}

func (lm *LifecycleManager) Restart(id string) error {
	lm.mu.RLock()
	srv, exists := lm.servers[id]
	lm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("server %s not found", id)
	}

	srv.mu.Lock()
	if srv.cancel != nil {
		srv.cancel()
	}
	if srv.cmd != nil && srv.cmd.Process != nil {
		_ = srv.cmd.Process.Kill()
		_ = srv.cmd.Wait()
	}
	srv.status.Restarts++
	err := lm.startProcess(srv)
	srv.mu.Unlock()
	return err
}

func (lm *LifecycleManager) HealthCheck(id string) (*ServerStatus, error) {
	lm.mu.RLock()
	srv, exists := lm.servers[id]
	lm.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("server %s not found", id)
	}

	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.cmd != nil && srv.cmd.Process != nil {
		if err := srv.cmd.Process.Signal(syscall.Signal(0)); err != nil {
			srv.status.Running = false
		} else {
			srv.status.Running = true
		}
	}
	status := srv.status
	return &status, nil
}

func (lm *LifecycleManager) ListRunning() []ServerStatus {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	result := make([]ServerStatus, 0, len(lm.servers))
	for _, srv := range lm.servers {
		srv.mu.Lock()
		result = append(result, srv.status)
		srv.mu.Unlock()
	}
	return result
}

func (lm *LifecycleManager) AutoStart(configs []ServerConfig) []error {
	var errs []error
	for _, cfg := range configs {
		if cfg.AutoStart {
			if err := lm.Start(cfg); err != nil {
				errs = append(errs, fmt.Errorf("auto-start %s: %w", cfg.ID, err))
			}
		}
	}
	return errs
}

func (lm *LifecycleManager) Shutdown() {
	lm.cancel()
	lm.mu.RLock()
	ids := make([]string, 0, len(lm.servers))
	for id := range lm.servers {
		ids = append(ids, id)
	}
	lm.mu.RUnlock()

	for _, id := range ids {
		_ = lm.Stop(id)
	}
	lm.wg.Wait()
}

func (lm *LifecycleManager) startProcess(srv *managedServer) error {
	ctx, cancel := context.WithCancel(lm.ctx)
	cmd := exec.CommandContext(ctx, srv.config.Command, srv.config.Args...)
	env := os.Environ()
	for k, v := range srv.config.Env {
		env = append(env, k+"="+v)
	}
	cmd.Env = env

	if err := cmd.Start(); err != nil {
		cancel()
		srv.status.LastError = err.Error()
		return fmt.Errorf("start %s: %w", srv.config.ID, err)
	}

	srv.cmd = cmd
	srv.cancel = cancel
	srv.status.Running = true
	srv.status.PID = cmd.Process.Pid
	srv.status.StartedAt = time.Now()
	return nil
}

func (lm *LifecycleManager) monitor(srv *managedServer) {
	defer lm.wg.Done()

	for {
		select {
		case <-lm.ctx.Done():
			return
		case <-time.After(lm.healthInterval):
			srv.mu.Lock()
			if srv.cmd == nil || srv.cmd.Process == nil {
				srv.mu.Unlock()
				continue
			}
			// Check if process exited
			if srv.cmd.ProcessState != nil && srv.cmd.ProcessState.Exited() {
				srv.status.Running = false
				srv.status.LastError = "process exited unexpectedly"
				srv.status.Restarts++
				_ = lm.startProcess(srv)
			}
			srv.mu.Unlock()
		}
	}
}
