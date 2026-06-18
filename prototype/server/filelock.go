package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// WriteJSONAtomic writes data as JSON to path atomically via tmp+rename.
func WriteJSONAtomic(path string, data any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if err := json.NewEncoder(f).Encode(data); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	f.Close()
	return os.Rename(tmp, path)
}
