package snapshot

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"github.com/getmitran/mitran/server/apierr"
)

type Component string

const (
	ComponentMemory   Component = "memory"
	ComponentSessions Component = "sessions"
	ComponentConfig   Component = "config"
	ComponentArtifacts Component = "artifacts"
	ComponentCrons    Component = "crons"
)

var AllComponents = []Component{ComponentMemory, ComponentSessions, ComponentConfig, ComponentArtifacts, ComponentCrons}

type SnapshotMeta struct {
	ID         string      `json:"id"`
	CreatedAt  time.Time   `json:"created_at"`
	Components []Component `json:"components"`
	Size       int64       `json:"size"`
	Path       string      `json:"path"`
}

type RestoreRequest struct {
	Path       string      `json:"path"`
	Components []Component `json:"components"`
}

var dataDir string

func SetDataDir(dir string) { dataDir = dir }

func componentDir(c Component) string {
	return filepath.Join(dataDir, string(c))
}

func CreateSnapshot(outputDir string) (*SnapshotMeta, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, err
	}

	id := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("snapshot-%s.tar.gz", id)
	path := filepath.Join(outputDir, filename)

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gw := gzip.NewWriter(f)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	var components []Component
	for _, c := range AllComponents {
		dir := componentDir(c)
		if _, err := os.Stat(dir); err == nil {
			components = append(components, c)
			if err := addDirToTar(tw, dir, string(c)); err != nil {
				return nil, fmt.Errorf("archiving %s: %w", c, err)
			}
		}
	}

	tw.Close()
	gw.Close()
	f.Close()

	info, _ := os.Stat(path)
	meta := &SnapshotMeta{
		ID:         id,
		CreatedAt:  time.Now(),
		Components: components,
		Size:       info.Size(),
		Path:       path,
	}

	metaPath := path + ".meta.json"
	metaBytes, _ := json.Marshal(meta)
	os.WriteFile(metaPath, metaBytes, 0644)

	return meta, nil
}

func RestoreSnapshot(path string, components []Component) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	allowed := make(map[string]bool)
	for _, c := range components {
		allowed[string(c)] = true
	}

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		parts := strings.SplitN(hdr.Name, "/", 2)
		if len(parts) == 0 || !allowed[parts[0]] {
			continue
		}

		target := filepath.Join(dataDir, hdr.Name)
		switch hdr.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(target, 0755)
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(target), 0755)
			out, err := os.Create(target)
			if err != nil {
				return err
			}
			io.Copy(out, tr)
			out.Close()
		}
	}
	return nil
}

func ListSnapshots(dir string) ([]SnapshotMeta, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var snapshots []SnapshotMeta
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".meta.json") {
			data, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				continue
			}
			var meta SnapshotMeta
			if json.Unmarshal(data, &meta) == nil {
				snapshots = append(snapshots, meta)
			}
		}
	}
	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].CreatedAt.After(snapshots[j].CreatedAt)
	})
	return snapshots, nil
}

func addDirToTar(tw *tar.Writer, srcDir, prefix string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(srcDir, path)
		name := filepath.Join(prefix, rel)

		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		hdr.Name = name

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if !info.IsDir() {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(tw, f)
			return err
		}
		return nil
	})
}

// REST Handlers

func HandleCreateSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		return
	}
	snapshotDir := filepath.Join(dataDir, ".snapshots")
	meta, err := CreateSnapshot(snapshotDir)
	if err != nil {
		apierr.WriteError(w, apierr.Internal(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meta)
}

func HandleRestore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		return
	}
	var req RestoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierr.WriteError(w, apierr.BadRequest("invalid request body"))
		return
	}
	if req.Path == "" {
		apierr.WriteError(w, apierr.BadRequest("path required"))
		return
	}
	if len(req.Components) == 0 {
		req.Components = AllComponents
	}
	if err := RestoreSnapshot(req.Path, req.Components); err != nil {
		apierr.WriteError(w, apierr.Internal(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "restored"})
}

func HandleListSnapshots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		return
	}
	snapshotDir := filepath.Join(dataDir, ".snapshots")
	list, err := ListSnapshots(snapshotDir)
	if err != nil {
		apierr.WriteError(w, apierr.Internal(err.Error()))
		return
	}
	if list == nil {
		list = []SnapshotMeta{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/snapshot", HandleCreateSnapshot)
	mux.HandleFunc("/api/v1/restore", HandleRestore)
	mux.HandleFunc("/api/v1/snapshots", HandleListSnapshots)
}
