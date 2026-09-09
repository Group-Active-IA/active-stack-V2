package install_test

import (
	"path/filepath"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/backup"
	"github.com/Group-Active-IA/active-stack/internal/install"
	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/pipeline"
	"github.com/Group-Active-IA/active-stack/internal/system"
)

// TestBuildPlan_SelfInstallGoosMatchesSnapshotGoos pins the two halves of the
// self-install path to the same source of truth.
//
// They were not: plan.go passed opts.Profile.OS when collecting the path to
// snapshot, and left the step's goos empty so it resolved runtime.GOOS at
// Run() time. On any run where those two disagree the snapshot backs up one
// file and the step overwrites another, so Rollback() restores nothing and the
// user's previous binary is gone with no trace.
//
// The datum that discriminates: this test declares a Profile OS that is NOT
// this machine's. Every call to the bin-dir resolver must be handed that OS.
// If some call arrives with runtime.GOOS instead, the two halves disagreed.
//
// It also stops TestSelfInstall_WritesBinaryUnderGivenHome from depending on
// the machine it runs on — that test set Profile{OS: "linux"} believing it
// pinned the platform, and CI proved on 2026-09-09 that it did not: green on
// ubuntu-latest, red on windows-latest.
func TestBuildPlan_SelfInstallGoosMatchesSnapshotGoos(t *testing.T) {
	const declaredOS = "windows"

	var seen []string
	restoreDir := install.SetSelfInstallBinaryInstallDirFn(func(goos, homeDir string) string {
		seen = append(seen, goos)
		return filepath.Join(homeDir, "bin")
	})
	defer restoreDir()

	restorePath := install.SetAddToUserPathFn(func(string) error { return nil })
	defer restorePath()

	restoreSnap := install.SetSnapshotCreate(func(string, []string) (backup.Manifest, error) {
		return backup.Manifest{}, nil
	})
	defer restoreSnap()

	restoreExt := install.SetExternalInstallFn(fakeExternalSuccess)
	defer restoreExt()

	reg := &fakeRegistry{adapters: map[model.Agent]install.AgentAdapter{
		model.AgentClaude: fakeAdapter{agent: model.AgentClaude},
	}}
	opts := buildOptions(t.TempDir(), reg, nil)
	opts.Profile = system.PlatformProfile{OS: declaredOS}

	intent := install.Intent{Agents: []model.Agent{model.AgentClaude}, Mode: model.ModeFull}
	cat := &fakeCatalog{harnesses: []model.Harness{externalHarness("inert", "npm", "inert-pkg")}}
	plan, err := install.BuildPlan(cat, intent, opts)
	if err != nil {
		t.Fatalf("BuildPlan() error = %v", err)
	}

	orch := pipeline.NewOrchestrator(pipeline.DefaultRollbackPolicy())
	if result := orch.Execute(plan.StagePlan); result.Err != nil {
		t.Fatalf("Execute() error = %v", result.Err)
	}

	if len(seen) < 2 {
		t.Fatalf("bin dir resolved %d time(s) (%v); the snapshot half and the step half must both resolve it, or this test proves nothing", len(seen), seen)
	}
	for i, goos := range seen {
		if goos != declaredOS {
			t.Errorf("bin dir resolution #%d used goos %q, want %q (all resolutions: %v)", i+1, goos, declaredOS, seen)
		}
	}
}
