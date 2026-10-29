package filemerge

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// 2026-10-29: installing rewrote a user's opencode.json (0600) and
// ~/.codex/config.toml (0600) as 0644, because every config writer passed an
// explicit 0o644. Files like these may carry API keys. perm == 0 now means
// "keep the existing mode; 0644 only when the file is new". An explicit perm
// still wins: restore uses it to put the recorded mode back.
func TestWriteFileAtomic_ZeroPermKeepsExistingMode(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX modes")
	}
	path := filepath.Join(t.TempDir(), "settings.json")
	if err := os.WriteFile(path, []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := WriteFileAtomic(path, []byte(`{"a":1}`), 0); err != nil {
		t.Fatal(err)
	}
	info, _ := os.Stat(path)
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("mode = %o, want the existing 600", info.Mode().Perm())
	}
}

func TestWriteFileAtomic_ZeroPermNewFileIs0644(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX modes")
	}
	path := filepath.Join(t.TempDir(), "new.json")
	if _, err := WriteFileAtomic(path, []byte("{}"), 0); err != nil {
		t.Fatal(err)
	}
	info, _ := os.Stat(path)
	if info.Mode().Perm() != 0o644 {
		t.Fatalf("mode = %o, want 644 for a new file", info.Mode().Perm())
	}
}

func TestWriteFileAtomic_ExplicitPermStillApplies(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX modes")
	}
	path := filepath.Join(t.TempDir(), "restored.json")
	if err := os.WriteFile(path, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := WriteFileAtomic(path, []byte(`{"b":2}`), 0o600); err != nil {
		t.Fatal(err)
	}
	info, _ := os.Stat(path)
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("mode = %o, want the explicit 600", info.Mode().Perm())
	}
}
