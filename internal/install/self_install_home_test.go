package install_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/backup"
	"github.com/Group-Active-IA/active-stack/internal/install"
	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/pipeline"
	"github.com/Group-Active-IA/active-stack/internal/system"
)

// TestSelfInstall_WritesBinaryUnderGivenHome is the regression test for the
// SECOND binary path that ignored --home. Fixing system.BinaryInstallDir for
// the downloaded-harness path (engram) left a second caller still on the
// home-less variant: selfInstallBinaryInstallDirFn, which copies active-stack
// itself into the bin dir.
//
// Observed by running the installer twice:
//
//	--mode lite --home /tmp/probe3  →  ~/.local/bin unchanged
//	--mode full --home /tmp/probe4  →  ~/.local/bin/active-stack rewritten
//
// Like the engram test, this drives BuildPlan+Execute and asserts on WHERE THE
// FILE ENDS UP. A test that only checked the path-resolution helper would have
// stayed green through this bug: the helper was already fixed, and the defect
// was a second caller passing no home to it.
//
// HOME points at a separate temp dir so a red run cannot touch the real home.
func TestSelfInstall_WritesBinaryUnderGivenHome(t *testing.T) {
	const binaryContent = "the-active-stack-binary"

	givenHome := t.TempDir()   // what --home says
	ambientHome := t.TempDir() // what os.UserHomeDir() would find
	t.Setenv("HOME", ambientHome)
	t.Setenv("USERPROFILE", ambientHome)

	// Stand in for the running binary that self-install copies.
	srcDir := t.TempDir()
	src := filepath.Join(srcDir, "active-stack")
	if err := os.WriteFile(src, []byte(binaryContent), 0o755); err != nil {
		t.Fatalf("seed source binary: %v", err)
	}

	// Pin the writable probe so the result does not depend on whether this
	// machine happens to have a writable /usr/local/bin.
	restoreWritable := system.SetBinaryInstallDirWritableFn(func(string) bool { return false })
	defer restoreWritable()

	restorePath := install.SetAddToUserPathFn(func(string) error { return nil })
	defer restorePath()

	restoreSnap := install.SetSnapshotCreate(func(string, []string) (backup.Manifest, error) {
		return backup.Manifest{}, nil
	})
	defer restoreSnap()

	// BuildPlan returns an empty plan when no harness is selected, and then the
	// self-install step is never appended. One inert harness keeps the plan real.
	restoreExt := install.SetExternalInstallFn(fakeExternalSuccess)
	defer restoreExt()

	reg := &fakeRegistry{adapters: map[model.Agent]install.AgentAdapter{
		model.AgentClaude: fakeAdapter{agent: model.AgentClaude},
	}}
	opts := buildOptions(givenHome, reg, nil)
	opts.SelfInstallBinaryPath = src
	opts.Profile = system.PlatformProfile{OS: "linux"}

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

	wantPath := filepath.Join(givenHome, ".local", "bin", "active-stack")
	leakPath := filepath.Join(ambientHome, ".local", "bin", "active-stack")

	if _, err := os.Stat(leakPath); err == nil {
		t.Errorf("self-install wrote into the ambient home: %q exists; --home was %q", leakPath, givenHome)
	}
	if data, err := os.ReadFile(wantPath); err != nil {
		t.Errorf("self-install did not write under the given home: read %q: %v", wantPath, err)
	} else if string(data) != binaryContent {
		t.Errorf("content at %q = %q, want %q", wantPath, data, binaryContent)
	}

	assertNoFilesUnder(t, ambientHome)
}

// assertNoFilesUnder fails the test if root contains any file.
func assertNoFilesUnder(t *testing.T, root string) {
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
