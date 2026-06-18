package plugins

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type SandboxPlugin struct {
	Name    string
	Path    string
	Timeout time.Duration
}

func (p *SandboxPlugin) Execute(input []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, p.Path)
	cmd.Stdin = bytes.NewReader(input)

	out, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		return nil, fmt.Errorf("plugin %s timed out after %v", p.Name, p.Timeout)
	}
	if err != nil {
		return nil, fmt.Errorf("plugin %s failed: %w", p.Name, err)
	}
	return out, nil
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
