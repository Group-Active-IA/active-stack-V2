package verify_test

import (
	"os"
	"path/filepath"
	"testing"

	roledefsinstaller "github.com/Group-Active-IA/active-stack/internal/harness/roledefs"
	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/verify"
)

// agentsDirFakeAdapter adds AgentsDir to the existing fakeAdapter — same
// low-blast-radius pattern as internal/install and internal/uninstall's own
// role_definitions_wiring_test.go.
type agentsDirFakeAdapter struct {
	fakeAdapter
	dir string
}

func (a agentsDirFakeAdapter) AgentsDir(string) string { return a.dir }

func roleDefinitionsHarnessForVerify() model.Harness {
	return model.Harness{ID: "role-definitions", Type: model.HarnessConfig}
}

// TestCheckForRoleDefinitions_Pass installs the real seven roles via
// roledefs.Install, then confirms verify.ChecksForHarness reports healthy.
func TestCheckForRoleDefinitions_Pass(t *testing.T) {
	home := t.TempDir()
	agentsDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")

	res, err := roledefsinstaller.DefaultResources()
	if err != nil {
		t.Fatalf("DefaultResources: %v", err)
	}
	adapter := roledefsAgentStub{agent: model.AgentClaude, dir: agentsDir}
	if _, err := roledefsinstaller.Install([]roledefsinstaller.AgentAdapter{adapter}, home, backupDir, res); err != nil {
		t.Fatalf("roledefs.Install (setup): %v", err)
	}

	fa := agentsDirFakeAdapter{fakeAdapter: fakeAdapter{agent: model.AgentClaude}, dir: agentsDir}
	checks := verify.ChecksForHarness(roleDefinitionsHarnessForVerify(), []verify.Adapter{fa}, home)
	if len(checks) != 1 {
		t.Fatalf("got %d checks, want 1: %+v", len(checks), checks)
	}
	results := runChecks(t, checks)
	for _, r := range results {
		if r.Status != verify.CheckStatusPassed {
			t.Errorf("check %q status = %q on a healthy install, want passed (error: %s)", r.ID, r.Status, r.Error)
		}
	}
}

// TestCheckForRoleDefinitions_FailMissingResource is RED first via mutant
// (documented in the apply report): removing the `h.ID == "role-definitions"`
// branch in checksForHarnessAdapter makes this pass for the wrong reason —
// zero checks is not the same as a check that actually caught the missing
// file — so this asserts a real failure with the missing role named.
func TestCheckForRoleDefinitions_FailMissingResource(t *testing.T) {
	home := t.TempDir()
	agentsDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")

	res, err := roledefsinstaller.DefaultResources()
	if err != nil {
		t.Fatalf("DefaultResources: %v", err)
	}
	adapter := roledefsAgentStub{agent: model.AgentClaude, dir: agentsDir}
	if _, err := roledefsinstaller.Install([]roledefsinstaller.AgentAdapter{adapter}, home, backupDir, res); err != nil {
		t.Fatalf("roledefs.Install (setup): %v", err)
	}
	if err := os.Remove(filepath.Join(agentsDir, "role-dba.md")); err != nil {
		t.Fatal(err)
	}

	fa := agentsDirFakeAdapter{fakeAdapter: fakeAdapter{agent: model.AgentClaude}, dir: agentsDir}
	checks := verify.ChecksForHarness(roleDefinitionsHarnessForVerify(), []verify.Adapter{fa}, home)
	results := runChecks(t, checks)
	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	if results[0].Status != verify.CheckStatusFailed {
		t.Fatalf("expected the check to fail when role-dba.md is missing, got status %q (error: %s)", results[0].Status, results[0].Error)
	}
}

// TestCheckForRoleDefinitions_UnsupportedAdapter_NoCheck confirms "not
// applicable" is not reported as a failure.
func TestCheckForRoleDefinitions_UnsupportedAdapter_NoCheck(t *testing.T) {
	fa := fakeAdapter{agent: model.AgentCodex} // no AgentsDir method at all
	checks := verify.ChecksForHarness(roleDefinitionsHarnessForVerify(), []verify.Adapter{fa}, t.TempDir())
	if len(checks) != 0 {
		t.Fatalf("got %d checks for an adapter with no AgentsDir capability, want 0", len(checks))
	}
}

// roledefsAgentStub is a minimal roledefs.AgentAdapter for setup.
type roledefsAgentStub struct {
	agent model.Agent
	dir   string
}

func (a roledefsAgentStub) Agent() model.Agent      { return a.agent }
func (a roledefsAgentStub) AgentsDir(string) string { return a.dir }
