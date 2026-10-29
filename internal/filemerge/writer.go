package filemerge

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type WriteResult struct {
	Changed bool
	Created bool
}

// WriteFileAtomic replaces path with content via a temp file and rename.
//
// perm == 0 keeps the mode of an existing file and uses 0644 for a new one.
// On 2026-10-29 installing rewrote a user's 0600 opencode.json and
// ~/.codex/config.toml as 0644 because every config writer passed 0o644;
// those files may hold API keys. An explicit perm still wins: restore passes
// the mode recorded in the backup manifest.
func WriteFileAtomic(path string, content []byte, perm fs.FileMode) (WriteResult, error) {
	keepMode := perm == 0
	if keepMode {
		perm = 0o644
	}

	created := false
	existing, err := os.ReadFile(path)
	if err == nil {
		if bytes.Equal(existing, content) {
			return WriteResult{}, nil
		}
		if keepMode {
			if info, statErr := os.Stat(path); statErr == nil {
				perm = info.Mode().Perm()
			}
		}
	} else if !os.IsNotExist(err) {
		return WriteResult{}, fmt.Errorf("read existing file %q: %w", path, err)
	} else {
		created = true
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return WriteResult{}, fmt.Errorf("create parent directories for %q: %w", path, err)
	}
	// Ensure the directory is writable — it may have been created with
	// restricted permissions by a previous installer version.
	// MkdirAll succeeds on existing dirs but does not fix their permissions.
	if err := os.Chmod(dir, 0o755); err != nil {
		return WriteResult{}, fmt.Errorf("set write permission on directory for %q: %w", path, err)
	}

	tmp, err := os.CreateTemp(dir, ".active-stack-*.tmp")
	if err != nil {
		return WriteResult{}, fmt.Errorf("create temp file for %q: %w", path, err)
	}

	tmpPath := tmp.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err := tmp.Write(content); err != nil {
		_ = tmp.Close()
		return WriteResult{}, fmt.Errorf("write temp file for %q: %w", path, err)
	}

	if err := tmp.Chmod(perm); err != nil {
		_ = tmp.Close()
		return WriteResult{}, fmt.Errorf("set permissions on temp file for %q: %w", path, err)
	}

	if err := tmp.Close(); err != nil {
		return WriteResult{}, fmt.Errorf("close temp file for %q: %w", path, err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return WriteResult{}, fmt.Errorf("replace %q atomically: %w", path, err)
	}

	cleanup = false
	return WriteResult{Changed: true, Created: created}, nil
}
