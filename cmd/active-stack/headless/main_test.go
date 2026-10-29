package headless_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestMain points the machine home at a throwaway directory before any test in
// this package runs. Several flows here resolve the home with os.UserHomeDir()
// when HomeDir is empty (starter installs, by design), and the self-install
// step then copied this test binary over the developer's real
// ~/.local/bin/active-stack: on 2026-10-29 `active-stack version` printed the
// headless test usage and "PASS". TestRunWindowsStartersInstall_Success was
// the one writing it. No test in this package may touch the real home.
//
// The Go locations are pinned to the ORIGINAL home first: the installer's
// dependency check runs `go version`, and with a fresh HOME the go command
// would download the whole toolchain into the throwaway dir.
func TestMain(m *testing.M) {
	realHome, _ := os.UserHomeDir()
	realCache, _ := os.UserCacheDir()
	pin := map[string]string{
		"GOPATH":  filepath.Join(realHome, "go"),
		"GOCACHE": filepath.Join(realCache, "go-build"),
	}
	for key, value := range pin {
		if os.Getenv(key) == "" && realHome != "" {
			_ = os.Setenv(key, value)
		}
	}
	if os.Getenv("GOMODCACHE") == "" && os.Getenv("GOPATH") != "" {
		_ = os.Setenv("GOMODCACHE", filepath.Join(os.Getenv("GOPATH"), "pkg", "mod"))
	}

	home, err := os.MkdirTemp("", "headless-test-home-")
	if err != nil {
		fmt.Fprintln(os.Stderr, "headless tests: create throwaway home:", err)
		os.Exit(1)
	}
	for _, key := range []string{"HOME", "USERPROFILE", "LOCALAPPDATA", "APPDATA"} {
		_ = os.Setenv(key, home)
	}

	code := m.Run()
	_ = os.RemoveAll(home)
	os.Exit(code)
}
