package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	roledefsinstaller "github.com/Group-Active-IA/active-stack/internal/harness/roledefs"
	"github.com/Group-Active-IA/active-stack/internal/model"
)

// claudeOnlyAdapter is the minimal roledefs.AgentAdapter needed to set up
// fixtures for these tests — mirrors roledefsAgentStub used elsewhere.
type claudeOnlyAdapter struct{ dir string }

func (a claudeOnlyAdapter) Agent() model.Agent      { return model.AgentClaude }
func (a claudeOnlyAdapter) AgentsDir(string) string { return a.dir }

// TestRoles_ConvenedRoleWithInstalledDefinition_Passes is RED first:
// --home doesn't exist as a flag yet. Covers "clasificación vacía" staying
// distinct from "definición faltante/vacía" (propose, sección 4): a diff
// that convenes role-dba, with role-dba.md actually installed, must not be
// blocked.
func TestRoles_ConvenedRoleWithInstalledDefinition_Passes(t *testing.T) {
	home := t.TempDir()
	agentsDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")
	res, err := roledefsinstaller.DefaultResources()
	if err != nil {
		t.Fatalf("DefaultResources: %v", err)
	}
	if _, err := roledefsinstaller.Install([]roledefsinstaller.AgentAdapter{claudeOnlyAdapter{agentsDir}}, home, backupDir, res); err != nil {
		t.Fatalf("roledefs.Install (setup): %v", err)
	}

	var out bytes.Buffer
	code := runRolesDispatch([]string{"--home", home, "--files", "migrations/003_tenant.sql"}, t.TempDir(), &out)
	if code != 0 {
		t.Fatalf("exit = %d, salida:\n%s", code, out.String())
	}
}

// TestRoles_ConvenedRoleWithMissingDefinition_BlocksWithDistinctExitCode
// covers the actual gate: role-dba is convened by the diff, but role-dba.md
// was never installed (fresh home) — it must block with a distinct exit
// code and name the missing role, not silently classify as if nothing were
// wrong.
func TestRoles_ConvenedRoleWithMissingDefinition_BlocksWithDistinctExitCode(t *testing.T) {
	var out bytes.Buffer
	code := runRolesDispatch([]string{"--home", t.TempDir(), "--files", "migrations/003_tenant.sql"}, t.TempDir(), &out)
	if code == 0 {
		t.Fatalf("exit = 0 with role-dba convened and no definition installed; salida:\n%s", out.String())
	}
	if code == 1 {
		t.Fatalf("exit = 1 is the generic usage/git error code — a missing definition needs its own, distinguishable code; salida:\n%s", out.String())
	}
	if !contains(out.String(), "role-dba") {
		t.Errorf("the error does not name the missing role:\n%s", out.String())
	}
}

// TestRoles_EmptyClassification_NeverBlocksOnDefinitions is the other half:
// a TRULY empty classification (zero files — QA is convened "siempre que
// haya algo que revisar", so any non-empty --files list convenes at least
// QA; a genuinely empty Convened() only happens with zero files, which
// --files refuses as a malformed invocation) must never block on
// definitions, no matter how bare the home is — "clasificación vacía" and
// "definición faltante" are different claims. Reached via git on a clean
// repo with nothing uncommitted, not via --files.
func TestRoles_EmptyClassification_NeverBlocksOnDefinitions(t *testing.T) {
	repo := initCleanGitRepo(t)
	var out bytes.Buffer
	code := runRolesDispatch([]string{"--home", t.TempDir()}, repo, &out)
	if code != 0 {
		t.Fatalf("exit = %d for a clean repo with nothing to review; salida:\n%s", code, out.String())
	}
}

// TestRoles_HomeFlag_DefaultsWhenAbsent confirms --home is optional and the
// fallback actually resolves to a non-empty path (os.UserHomeDir()) instead
// of leaving the definitions check silently unable to run. Asserted at the
// parseRolesFlags level, not through the full dispatch: which role
// definitions happen to exist on whatever machine runs this suite is not
// something a deterministic test should depend on (the two tests above
// already cover the dispatch's actual blocking behavior with a controlled
// --home).
func TestRoles_HomeFlag_DefaultsWhenAbsent(t *testing.T) {
	_, _, homeDir, err := parseRolesFlags([]string{"--files", "README.md"})
	if err != nil {
		t.Fatalf("parseRolesFlags: %v", err)
	}
	if homeDir == "" {
		t.Error("homeDir is empty without --home — the os.UserHomeDir() fallback did not resolve")
	}
}

// initCleanGitRepo creates a temp repo with one commit and nothing
// uncommitted, so roles.Collect (no --files, no --base) returns zero files.
func initCleanGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	runGitT(t, dir, "init", "-q", "-b", "main")
	runGitT(t, dir, "config", "user.email", "test@example.com")
	runGitT(t, dir, "config", "user.name", "test")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGitT(t, dir, "add", ".")
	runGitT(t, dir, "commit", "-qm", "inicial")
	return dir
}

func runGitT(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func contains(haystack, needle string) bool {
	return bytes.Contains([]byte(haystack), []byte(needle))
}
