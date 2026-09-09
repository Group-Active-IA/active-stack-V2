package install_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/backup"
	extinstaller "github.com/Group-Active-IA/active-stack/internal/harness/external"
	"github.com/Group-Active-IA/active-stack/internal/install"
	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/pipeline"
	"github.com/Group-Active-IA/active-stack/internal/system"
)

// externalHarness builds an external harness installed in every mode.
func externalHarness(id, method, pkg string) model.Harness {
	return model.Harness{
		ID:           id,
		Type:         model.HarnessExternal,
		External:     &model.External{Method: method, Pkg: pkg},
		InstallModes: []model.InstallMode{model.ModeLite, model.ModeFull},
	}
}

// TestExternalRollback_RemovesDownloadedBinary is the regression test for the
// second half of the real-world probe run: openspec failed AFTER engram had
// already been installed, the installer printed "Rollback: succeeded", and the
// 19 MB engram binary was still on disk.
//
// A rollback that reports success and leaves files behind is worse than one
// that fails: it leaves you believing the machine is back to how it was.
//
// The cause is independent of the --home leak. Rollback restores a snapshot of
// a path list predicted before Apply, and collectWritePaths only predicts paths
// for External harnesses with method "mcp" — a downloaded binary is never in
// the manifest. On top of that externalStep.Run discarded the installer Result
// (`_, err := externalInstallFn(...)`), so nothing downstream ever learned the
// binary's path. Putting the binary inside the right home would not have
// changed any of that.
func TestExternalRollback_RemovesDownloadedBinary(t *testing.T) {
	installDir := t.TempDir()
	binaryPath := filepath.Join(installDir, "engram")

	restore := install.SetExternalInstallFn(func(
		_ context.Context,
		h model.Harness,
		_ system.PlatformProfile,
		_ []extinstaller.AgentAdapter,
		_ string,
	) (extinstaller.Result, error) {
		switch h.ID {
		case "engram":
			// Stand in for downloadBinary actually writing the file.
			if err := os.WriteFile(binaryPath, []byte("19mb-of-engram"), 0o755); err != nil {
				t.Fatalf("seed downloaded binary: %v", err)
			}
			return extinstaller.Result{BinaryPath: binaryPath, Downloaded: true}, nil
		case "openspec":
			return extinstaller.Result{}, errors.New("sudo: npm: command not found")
		}
		return extinstaller.Result{}, nil
	})
	defer restore()

	result := runTwoHarnessPipeline(t, "homebrew")

	if result.Err == nil {
		t.Fatal("Execute() must fail: the openspec step returns an error")
	}
	if result.Rollback.Stage != pipeline.StageRollback {
		t.Fatal("rollback stage did not run")
	}

	if _, err := os.Stat(binaryPath); err == nil {
		t.Errorf("rollback left the downloaded binary on disk: %q still exists", binaryPath)
	}
	if !result.Rollback.Success {
		t.Errorf("Rollback.Success = false (%v), want true: removing our own download must succeed", result.Rollback.Err)
	}
	if len(result.Rollback.Unreverted) != 0 {
		t.Errorf("Unreverted = %v, want empty: the binary was actually removed", result.Rollback.Unreverted)
	}
}

// TestExternalRollback_ReportsUnrevertedNpmPackage covers the case rollback
// genuinely cannot undo. `npm install -g` scatters files a path snapshot never
// modelled, and deleting the shim would leave a half-removed package, so the
// honest outcome is not "succeeded" — it is "partial, and here is what stayed".
func TestExternalRollback_ReportsUnrevertedNpmPackage(t *testing.T) {
	restore := install.SetExternalInstallFn(func(
		_ context.Context,
		h model.Harness,
		_ system.PlatformProfile,
		_ []extinstaller.AgentAdapter,
		_ string,
	) (extinstaller.Result, error) {
		switch h.ID {
		case "engram":
			return extinstaller.Result{BinaryPath: "/usr/local/bin/openspec"}, nil
		case "openspec":
			return extinstaller.Result{}, errors.New("boom")
		}
		return extinstaller.Result{}, nil
	})
	defer restore()

	result := runTwoHarnessPipeline(t, "npm")

	if result.Err == nil {
		t.Fatal("Execute() must fail")
	}
	if len(result.Rollback.Unreverted) == 0 {
		t.Fatal("Unreverted is empty: a global npm install that rollback cannot undo must be reported")
	}
	joined := strings.Join(result.Rollback.Unreverted, "\n")
	if !strings.Contains(joined, "some-pkg") {
		t.Errorf("Unreverted = %v, want it to name the package left installed", result.Rollback.Unreverted)
	}
}

// runTwoHarnessPipeline builds and runs a plan with an "engram" step followed
// by an "openspec" step, with the snapshot and restore stubbed out.
func runTwoHarnessPipeline(t *testing.T, firstMethod string) pipeline.ExecutionResult {
	t.Helper()

	cat := &fakeCatalog{harnesses: []model.Harness{
		externalHarness("engram", firstMethod, "some-pkg"),
		externalHarness("openspec", "npm", "openspec"),
	}}
	reg := &fakeRegistry{adapters: map[model.Agent]install.AgentAdapter{
		model.AgentClaude: fakeAdapter{agent: model.AgentClaude},
	}}
	homeDir := t.TempDir()
	intent := install.Intent{Agents: []model.Agent{model.AgentClaude}, Mode: model.ModeLite}

	plan, err := install.BuildPlan(cat, intent, buildOptions(homeDir, reg, nil))
	if err != nil {
		t.Fatalf("BuildPlan() error = %v", err)
	}

	restoreSnap := install.SetSnapshotCreate(func(string, []string) (backup.Manifest, error) {
		return backup.Manifest{}, nil
	})
	defer restoreSnap()

	restoreRestore := install.SetRestoreFn(func(backup.Manifest) error { return nil })
	defer restoreRestore()

	orch := pipeline.NewOrchestrator(pipeline.DefaultRollbackPolicy())
	return orch.Execute(plan.StagePlan)
}
