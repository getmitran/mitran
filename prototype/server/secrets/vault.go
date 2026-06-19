package secrets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// SecretStore abstracts secret retrieval and storage.
type SecretStore interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

// EnvSecretStore reads secrets from environment variables prefixed with MITRAN_SECRET_.
// Set() writes to an internal map; Get() checks the map first, then falls back to os.Getenv.
type EnvSecretStore struct {
	mu    sync.RWMutex
	store map[string]string
}

func (e *EnvSecretStore) Get(key string) (string, error) {
	envKey := "MITRAN_SECRET_" + strings.ToUpper(key)
	e.mu.RLock()
	if v, ok := e.store[envKey]; ok {
		e.mu.RUnlock()
		return v, nil
	}
	e.mu.RUnlock()
	v := os.Getenv(envKey)
	if v == "" {
		return "", fmt.Errorf("secret %q not found", key)
	}
	return v, nil
}

func (e *EnvSecretStore) Set(key, value string) error {
	envKey := "MITRAN_SECRET_" + strings.ToUpper(key)
	e.mu.Lock()
	if e.store == nil {
		e.store = make(map[string]string)
	}
	e.store[envKey] = value
	e.mu.Unlock()
	return nil
}

// FileSecretStore reads secrets from a directory (one file per secret).
type FileSecretStore struct {
	Dir string
}

func (f *FileSecretStore) Get(key string) (string, error) {
	data, err := os.ReadFile(filepath.Join(f.Dir, key))
	if err != nil {
		return "", fmt.Errorf("secret %q: %w", key, err)
	}
	return strings.TrimSpace(string(data)), nil
}

func (f *FileSecretStore) Set(key, value string) error {
	if err := os.MkdirAll(f.Dir, 0700); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(f.Dir, key), []byte(value), 0600)
}

// NewSecretStore returns a SecretStore based on MITRAN_SECRETS_BACKEND (env|file, default env).
func NewSecretStore() SecretStore {
	switch os.Getenv("MITRAN_SECRETS_BACKEND") {
	case "file":
		dir := os.Getenv("MITRAN_SECRETS_DIR")
		if dir == "" {
			dir = "/etc/mitran/secrets"
		}
		return &FileSecretStore{Dir: dir}
	default:
		return &EnvSecretStore{}
	}
}
