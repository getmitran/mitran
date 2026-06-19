package plugins

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const defaultMaxOutputBytes = 10 * 1024 * 1024 // 10MB

type SandboxPlugin struct {
	Name           string
	Path           string
	Timeout        time.Duration
	MaxOutputBytes int64
}

func (p *SandboxPlugin) maxOutput() int64 {
	if p.MaxOutputBytes > 0 {
		return p.MaxOutputBytes
	}
	return defaultMaxOutputBytes
}

func (p *SandboxPlugin) Execute(input []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, p.Path)
	cmd.Stdin = bytes.NewReader(input)

	// Resource limits: SysProcAttr-based limits (rlimit, cgroups) are
	// OS-specific. For cross-platform safety we constrain via env and
	// output capping. On Linux, add unix.Rlimit via SysProcAttr if needed.
	cmd.Env = append(os.Environ(), "GOMAXPROCS=1")

	var stdout bytes.Buffer
	cmd.Stdout = &limitedWriter{w: &stdout, max: p.maxOutput()}

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("plugin %s timed out after %v", p.Name, p.Timeout)
		}
		return nil, fmt.Errorf("plugin %s failed: %w", p.Name, err)
	}
	return stdout.Bytes(), nil
}

// limitedWriter caps bytes written, implementing io.Writer.
type limitedWriter struct {
	w       io.Writer
	max     int64
	written int64
}

func (lw *limitedWriter) Write(p []byte) (int, error) {
	remaining := lw.max - lw.written
	if remaining <= 0 {
		return 0, fmt.Errorf("output exceeded %d bytes", lw.max)
	}
	if int64(len(p)) > remaining {
		p = p[:remaining]
	}
	n, err := lw.w.Write(p)
	lw.written += int64(n)
	return n, err
}

func LoadSandboxPlugins(dir string) ([]SandboxPlugin, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var plugins []SandboxPlugin
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.Mode()&0111 == 0 {
			continue
		}
		plugins = append(plugins, SandboxPlugin{
			Name:    e.Name(),
			Path:    filepath.Join(dir, e.Name()),
			Timeout: 30 * time.Second,
		})
	}
	return plugins, nil
}
