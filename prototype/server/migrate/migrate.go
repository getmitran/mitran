package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Migration struct {
	Version int
	Name    string
	Up      func() error
}

const versionFile = "./data/migration_version"

func GetVersion() int {
	data, err := os.ReadFile(versionFile)
	if err != nil {
		return 0
	}
	v, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	return v
}

func SetVersion(v int) error {
	if err := os.MkdirAll(filepath.Dir(versionFile), 0755); err != nil {
		return err
	}
	return os.WriteFile(versionFile, []byte(strconv.Itoa(v)), 0644)
}

func Run(migrations []Migration, currentVersion int) error {
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})
	for _, m := range migrations {
		if m.Version > currentVersion {
			if err := m.Up(); err != nil {
				return fmt.Errorf("migration %d (%s) failed: %w", m.Version, m.Name, err)
			}
			if err := SetVersion(m.Version); err != nil {
				return fmt.Errorf("failed to set version %d: %w", m.Version, err)
			}
		}
	}
	return nil
}
