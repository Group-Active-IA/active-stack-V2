package headless_test

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/cmd/active-stack/headless"
	"github.com/Group-Active-IA/active-stack/internal/backup"
	"github.com/Group-Active-IA/active-stack/internal/harness/roledefs"
	"github.com/Group-Active-IA/active-stack/internal/install"
	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/uninstall"
)

func runRoleDefinitions(t *testing.T, conflicts []roledefs.ResourceResult) (int, string) {
	t.Helper()
	restoreSnap := install.SetSnapshotCreate(func(string, []string) (backup.Manifest, error) { return backup.Manifest{}, nil })
	t.Cleanup(restoreSnap)
	restoreRoles := install.SetRoleDefinitionsInstallFn(func([]install.AgentAdapter, string, string) (roledefs.Result, error) {
		return roledefs.Result{Resources: conflicts}, nil
	})
	t.Cleanup(restoreRoles)
	h := model.Harness{ID: "role-definitions", Type: model.HarnessConfig, InstallModes: []model.InstallMode{model.ModeFull}, Agents: []model.Agent{model.AgentClaude}}
	cat := &fakeExecCatalog{harnesses: []model.Harness{h}}
	reg := &fakeExecRegistry{adapters: map[model.Agent]install.AgentAdapter{model.AgentClaude: fakeExecAdapter{agent: model.AgentClaude}}}
	var out bytes.Buffer
	code := headless.RunHeadless(headless.ParsedFlags{
		Yes: true, HomeDir: t.TempDir(), NoSelfInstall: true,
		Intent: install.Intent{Agents: []model.Agent{model.AgentClaude}, Mode: model.ModeFull},
	}, cat, reg, &out)
	return code, out.String()
}

// A BLOQUEANTE must not end in "ready to start building", and the empty
// "0 passed, 0 failed" report must not pose as verification (isolated-home
// install, 2026-10-29). Exit code stays 0 by decision: scripts rely on it.
func TestHeadlessExecutor_FinalStatusWithBlockingConflict(t *testing.T) {
	code, out := runRoleDefinitions(t, []roledefs.ResourceResult{{
		Agent: model.AgentClaude, Resource: "role-auditor.md", Outcome: roledefs.OutcomeConflict,
		Path: "/h/.claude/agents/role-auditor.md", ConflictPath: "/h/.claude/agents/role-auditor.md.active-stack.md",
	}})
	if code != 0 {
		t.Fatalf("exit %d, want 0 (decision: blocking stays a warning):\n%s", code, out)
	}
	if strings.Contains(out, "ready to start building") {
		t.Errorf("blocking install still claims ready:\n%s", out)
	}
	if strings.Contains(out, "0 passed, 0 failed") {
		t.Errorf("empty report still printed as verification:\n%s", out)
	}
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	tail := lines[len(lines)-1]
	if !strings.Contains(tail, "BLOCKING") {
		t.Errorf("last line must state the blocking result, got %q:\n%s", tail, out)
	}
}

func TestHeadlessExecutor_FinalStatusCleanRun(t *testing.T) {
	code, out := runRoleDefinitions(t, nil)
	if code != 0 {
		t.Fatalf("exit %d:\n%s", code, out)
	}
	if strings.Contains(out, "0 passed, 0 failed") {
		t.Errorf("empty report still printed as verification:\n%s", out)
	}
	if !strings.Contains(out, "Post-install checks passed") || !strings.Contains(out, "ready to start building") {
		t.Errorf("clean run must report passed checks and readiness:\n%s", out)
	}
}

// Restore returns snapshotted files only; binaries downloaded by the install
// (engram under ~/.local/bin) stay, and the user must be told.
func TestRunHeadlessUninstall_RestoreWarnsDownloadedBinariesStay(t *testing.T) {
	homeDir := t.TempDir()
	cat, reg := minimalCatalogAndRegistry(homeDir)
	m, err := backup.NewSnapshotter().Create(filepath.Join(homeDir, "snap"), []string{filepath.Join(homeDir, "absent.md")})
	if err != nil {
		t.Fatal(err)
	}
	params := minimalParams(homeDir)
	params.Intent.Strategy = uninstall.StrategyRestore
	params.RestoreManifestPath = filepath.Join(m.RootDir, "manifest.json")
	var out bytes.Buffer
	if code := headless.RunHeadlessUninstall(params, cat, reg, &out); code != 0 {
		t.Fatalf("exit %d:\n%s", code, out.String())
	}
	if !strings.Contains(out.String(), "downloaded binaries") {
		t.Errorf("restore does not warn that downloaded binaries stay:\n%s", out.String())
	}
}
