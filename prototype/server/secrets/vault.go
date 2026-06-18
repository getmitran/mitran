package secrets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SecretStore abstracts secret retrieval and storage.
type SecretStore interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

// EnvSecretStore reads secrets from environment variables prefixed with MITRAN_SECRET_.
type EnvSecretStore struct{}

func (e *EnvSecretStore) Get(key string) (string, error) {
	v := os.Getenv("MITRAN_SECRET_" + strings.ToUpper(key))
	if v == "" {
		return "", fmt.Errorf("secret %q not found", key)
	}
	return v, nil
}

func (e *EnvSecretStore) Set(key, value string) error {
	return os.Setenv("MITRAN_SECRET_"+strings.ToUpper(key), value)
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
