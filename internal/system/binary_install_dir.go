package system

import (
	"os"
	"path/filepath"
)

// binaryInstallDirWritableFn is the backing function for the writable-dir check.
// It is a package-level variable so tests can inject a fake.
var binaryInstallDirWritableFn = isWritableDirSystem

// BinaryInstallDir returns the directory where active-stack binaries should be installed.
// Windows: %LOCALAPPDATA%\active-stack\bin
// Unix: /usr/local/bin when writable, otherwise ~/.local/bin
//
// The home directory is discovered with os.UserHomeDir(). Callers that were
// given an explicit home (the installer's --home flag) must use
// BinaryInstallDirForHome instead, or the flag is silently ignored.
func BinaryInstallDir(goos string) string {
	return BinaryInstallDirForHome(goos, "")
}

// BinaryInstallDirForHome is BinaryInstallDir with the home directory supplied
// by the caller instead of discovered from the environment.
//
// This exists because `--home` was structurally unreachable from the code that
// decides where a downloaded binary lands: BinaryInstallDir took only goos and
// resolved the home itself, so an install run with --home pointing at a
// throwaway directory still wrote into the user's real ~/.local/bin.
//
// An empty homeDir means "discover it", which is what BinaryInstallDir passes,
// so the default install location is unchanged when no home is supplied.
func BinaryInstallDirForHome(goos, homeDir string) string {
	return binaryInstallDirWithCandidate(goos, homeDir, "/usr/local/bin")
}

// binaryInstallDirWithCandidate is the testable core of BinaryInstallDir,
// accepting an injected candidate dir (for tests that can't use the real /usr/local/bin).
func binaryInstallDirWithCandidate(goos, homeDir, candidate string) string {
	home := homeDir
	if home == "" {
		// Error is deliberately swallowed: an unresolvable home is handled
		// per-branch below, the same way it was before homeDir was threaded in.
		home, _ = os.UserHomeDir()
	}

	if goos == "windows" {
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			localAppData = filepath.Join(home, "AppData", "Local")
		}
		return filepath.Join(localAppData, "active-stack", "bin")
	}

	if binaryInstallDirWritableFn(candidate) {
		return candidate
	}

	if home == "" {
		return candidate
	}
	return filepath.Join(home, ".local", "bin")
}

// isWritableDirSystem probes whether a directory is writable by creating a temp file.
func isWritableDirSystem(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return false
	}
	tmp, err := os.CreateTemp(dir, ".active-stack-write-test-*")
	if err != nil {
		return false
	}
	tmp.Close()
	os.Remove(tmp.Name())
	return true
}
