package external

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/system"
)

// cmdRunner abstracts os/exec so tests can inject a fake without running real commands.
type cmdRunner func(ctx context.Context, name string, args ...string) ([]byte, error)

// realCmdRunner runs name with args and returns combined stdout+stderr.
func realCmdRunner(ctx context.Context, name string, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, name, args...).CombinedOutput()
}

// runner is the package-level runner; replaced in tests.
var runner cmdRunner = realCmdRunner

// lookPath is a package-level var so tests can replace it.
var lookPath = exec.LookPath

func installNPM(ctx context.Context, h model.Harness, profile system.PlatformProfile) (Result, error) {
	pkg := h.External.Pkg
	if pkg == "" {
		return Result{}, fmt.Errorf("harness %q: npm method requires External.Pkg", h.ID)
	}

	// Ask before installing. On 2026-09-09 an install on a machine that already
	// had openspec 1.3.1 failed with "sudo: npm: command not found" and rolled
	// the whole installation back: nvm lives under the user's home, so `npm` is
	// on the interactive PATH and not on sudo's. A machine that ALREADY has the
	// tool working was the one place this could not install it.
	//
	// homebrew.go and the download method already reported AlreadyInstalled;
	// npm was the only one of the three that never asked.
	binaryName := pkgBinaryName(pkg)
	var upgradeNote string
	if existing, err := lookPath(binaryName); err == nil && existing != "" {
		// Presence is not version (review 2026-10-29, 8.2): an old openspec
		// passed as installed while lacking the subcommands the lean block
		// uses. Without min_version the presence check stands.
		min := h.External.MinVersion
		if min == "" {
			return Result{BinaryPath: existing, AlreadyInstalled: true}, nil
		}
		out, _ := runner(ctx, existing, "--version")
		current, ok := parseVersion(string(out))
		if !ok {
			return Result{BinaryPath: existing, AlreadyInstalled: true, Notes: []string{fmt.Sprintf(
				"could not read the %s version (minimum %s); left as is: check it with `%s --version`", binaryName, min, binaryName)}}, nil
		}
		required, _ := parseVersion(min)
		if !versionLess(current, required) {
			return Result{BinaryPath: existing, AlreadyInstalled: true}, nil
		}
		upgradeNote = fmt.Sprintf("%s %s is below the minimum %s", binaryName, strings.Join(current, "."), min)
	}

	var (
		name string
		args []string
	)

	// --ignore-scripts: npm packages can ship postinstall hooks that spawn
	// "node" through a lifecycle-script shell npm builds itself; on some
	// Windows setups that shell fails to resolve "node" on PATH even though
	// the interactive shell resolves it fine, turning an optional postinstall
	// script (e.g. openspec's shell-completion tip) into a hard install
	// failure. We never depend on install-time lifecycle scripts.
	if useSudo(profile) {
		name = "sudo"
		args = []string{"npm", "install", "-g", "--ignore-scripts", pkg}
	} else {
		name = "npm"
		args = []string{"install", "-g", "--ignore-scripts", pkg}
	}

	if out, err := runner(ctx, name, args...); err != nil {
		if upgradeNote != "" {
			// The old binary still works: a failed upgrade must not roll the
			// whole install back (the 2026-09-09 sudo failure above).
			binaryPath, _ := lookPath(binaryName)
			return Result{BinaryPath: binaryPath, AlreadyInstalled: true, Notes: []string{fmt.Sprintf(
				"%s; upgrade failed (%v): run `npm install -g %s` yourself", upgradeNote, err, pkg)}}, nil
		}
		return Result{}, fmt.Errorf("npm install -g %s: %w\n%s", pkg, err, out)
	}

	binaryPath, _ := lookPath(binaryName)
	res := Result{BinaryPath: binaryPath}
	if upgradeNote != "" {
		res.Notes = []string{upgradeNote + "; upgraded with npm install -g " + pkg}
	}
	return res, nil
}

// versionPattern finds the first dotted version in `--version` output, which
// may carry a prefix ("openspec version 1.14.3") or a trailing newline.
var versionPattern = regexp.MustCompile(`(\d+)\.(\d+)(?:\.(\d+))?`)

// parseVersion returns major, minor and patch as strings ("0" when absent).
func parseVersion(text string) ([]string, bool) {
	m := versionPattern.FindStringSubmatch(text)
	if m == nil {
		return nil, false
	}
	patch := m[3]
	if patch == "" {
		patch = "0"
	}
	return []string{m[1], m[2], patch}, true
}

// versionLess compares two parsed versions numerically, part by part.
func versionLess(a, b []string) bool {
	for i := 0; i < 3; i++ {
		x, _ := strconv.Atoi(a[i])
		y, _ := strconv.Atoi(b[i])
		if x != y {
			return x < y
		}
	}
	return false
}

// useSudo returns true when the install should be run under sudo.
// Windows never uses sudo; other platforms skip it when npm prefix is user-writable.
func useSudo(profile system.PlatformProfile) bool {
	if profile.OS == "windows" {
		return false
	}
	return !profile.NpmWritable
}

// pkgBinaryName extracts the binary name from an npm package identifier.
// "@scope/name" → "name", "plain" → "plain".
func pkgBinaryName(pkg string) string {
	if strings.HasPrefix(pkg, "@") {
		parts := strings.SplitN(pkg, "/", 2)
		if len(parts) == 2 {
			return filepath.Base(parts[1])
		}
	}
	return filepath.Base(pkg)
}
