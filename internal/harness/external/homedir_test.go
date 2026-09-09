package external

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/system"
)

// TestInstall_HomebrewWritesBinaryUnderGivenHome is the regression test for the
// real-world leak found by running
//
//	active-stack install --agent claude --mode full --home /tmp/probehome --yes
//
// on a Linux machine: --home was passed precisely so the run would not touch
// the user's real home, and the installer still dropped a 19 MB `engram`
// binary into the REAL ~/.local/bin while /tmp/probehome/.local/bin was never
// even created.
//
// The cause is that Install threads homeDir down to the MCP-registration path
// only. The binary path resolved through system.BinaryInstallDir(goos), which
// takes no home at all and calls os.UserHomeDir() internally — so the flag was
// structurally unreachable from the code that decides where the file lands.
//
// This test deliberately asserts on WHERE THE FILE ENDS UP, not on what a path
// helper returns: a helper-level test would keep passing if a caller forgot to
// thread the home through, which is exactly the bug.
//
// HOME is pointed at a second temp dir so that running this test while it is
// still red reproduces the bug without re-polluting the developer's real home.
func TestInstall_HomebrewWritesBinaryUnderGivenHome(t *testing.T) {
	const binaryContent = "fake-engram-binary"

	// The home the installer is TOLD to use (the --home flag).
	givenHome := t.TempDir()
	// The home the process would discover on its own (os.UserHomeDir()).
	// Nothing may be written here.
	ambientHome := t.TempDir()
	t.Setenv("HOME", ambientHome)
	t.Setenv("USERPROFILE", ambientHome)

	tarGzData := buildTarGz(t, "engram", []byte(binaryContent))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/releases") {
			_ = json.NewEncoder(w).Encode([]map[string]string{{"tag_name": "v1.0.0"}})
			return
		}
		_, _ = w.Write(tarGzData)
	}))
	defer srv.Close()

	origBase := githubBaseURL
	githubBaseURL = srv.URL
	defer func() { githubBaseURL = origBase }()

	origClient := httpClient
	httpClient = srv.Client()
	defer func() { httpClient = origClient }()

	// Never mutate the real user PATH from a test.
	origAdd := addToUserPath
	addToUserPath = func(string) error { return nil }
	defer func() { addToUserPath = origAdd }()

	// Pin the writable probe: otherwise this test resolves to /usr/local/bin on
	// a machine where that happens to be writable, and the result would depend
	// on the host instead of on the code.
	restoreWritable := system.SetBinaryInstallDirWritableFn(func(string) bool { return false })
	defer restoreWritable()

	h := harnessWithMethod("homebrew", "engram", "")
	h.External.Repo = "Gentleman-Programming/engram"
	profile := system.PlatformProfile{OS: "linux", PackageManager: "apt"}

	res, err := Install(context.Background(), h, profile, nil, givenHome, nil)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	wantPath := filepath.Join(givenHome, ".local", "bin", "engram")
	leakPath := filepath.Join(ambientHome, ".local", "bin", "engram")

	if _, err := os.Stat(leakPath); err == nil {
		t.Errorf("binary leaked into the ambient home: %q exists; --home was %q", leakPath, givenHome)
	}

	if _, err := os.Stat(wantPath); err != nil {
		t.Errorf("binary was not installed under the given home: stat %q: %v", wantPath, err)
	} else if data, readErr := os.ReadFile(wantPath); readErr != nil || string(data) != binaryContent {
		t.Errorf("binary content at %q = %q (err %v), want %q", wantPath, data, readErr, binaryContent)
	}

	if res.BinaryPath != wantPath {
		t.Errorf("Result.BinaryPath = %q, want %q", res.BinaryPath, wantPath)
	}

	// Nothing at all may have been created outside the given home.
	assertTreeEmpty(t, ambientHome)
}

// assertTreeEmpty fails the test if root contains any file, reporting each one.
func assertTreeEmpty(t *testing.T, root string) {
	t.Helper()
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			t.Errorf("unexpected file written outside the given home: %q", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %q: %v", root, err)
	}
}
