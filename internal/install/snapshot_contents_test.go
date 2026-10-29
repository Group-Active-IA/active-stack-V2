package install_test

// A preexisting directory is a restore NO-OP by design (backup/restore.go: its
// files are restored "by their own ManifestEntry records"). So the Prepare
// snapshot must list every FILE a step will overwrite inside a directory that
// already exists, not only the directory. Found in an isolated-home install
// (2026-10-29): restore did not bring back opsx/ship.md and opsx/verify.md.

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/Group-Active-IA/active-stack/internal/backup"
	"github.com/Group-Active-IA/active-stack/internal/harness/skill"
	"github.com/Group-Active-IA/active-stack/internal/install"
	"github.com/Group-Active-IA/active-stack/internal/model"
)

func capturePrepare(t *testing.T, h model.Harness, homeDir string) ([]string, map[string]bool) {
	t.Helper()
	var paths []string
	var hints map[string]bool
	restore := install.SetSnapshotCreateWithHints(func(_ string, p []string, d map[string]bool) (backup.Manifest, error) {
		paths, hints = p, d
		return backup.Manifest{}, nil
	})
	t.Cleanup(restore)
	cat := &fakeCatalog{harnesses: []model.Harness{h}}
	reg := &fakeRegistry{adapters: map[model.Agent]install.AgentAdapter{
		model.AgentClaude: fakeAdapter{agent: model.AgentClaude},
	}}
	plan, err := install.BuildPlan(cat, install.Intent{Agents: []model.Agent{model.AgentClaude}, Mode: model.ModeLite},
		install.Options{HomeDir: homeDir, Registry: reg, NoSelfInstall: true})
	if err != nil {
		t.Fatalf("BuildPlan() error: %v", err)
	}
	for _, s := range plan.Prepare {
		_ = s.Run()
	}
	return paths, hints
}

func hasPath(paths []string, want string) bool {
	for _, p := range paths {
		if filepath.Clean(p) == filepath.Clean(want) {
			return true
		}
	}
	return false
}

func TestSnapshot_ListsEveryCommandFileTheStepWrites(t *testing.T) {
	install.ResetEmbeddedCommandsFS()
	install.WithEmbeddedCommandsFS(fstest.MapFS{
		"commands/claude/opsx/ship.md":   {Data: []byte("ship")},
		"commands/claude/opsx/verify.md": {Data: []byte("verify")},
	})
	t.Cleanup(install.ResetEmbeddedCommandsFS)
	restoreCommand := install.SetCommandInstallFn(func(_ []install.AgentAdapter, _, _ string) error { return nil })
	t.Cleanup(restoreCommand)

	homeDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(homeDir, "commands", "opsx"), 0o755); err != nil {
		t.Fatal(err)
	}
	paths, _ := capturePrepare(t, commandHarness(), homeDir)
	for _, name := range []string{"ship.md", "verify.md"} {
		want := filepath.Join(homeDir, "commands", "opsx", name)
		if !hasPath(paths, want) {
			t.Errorf("snapshot does not list %s; restore could not bring it back. paths: %v", want, paths)
		}
	}
}

func TestSnapshot_ListsExistingSkillFilesAndHintsMissingSkillDir(t *testing.T) {
	restoreSkill := install.SetSkillInstallFn(func(_ interface{}, _ fs.FS, _ context.Context, _ model.Harness, _ []skill.AgentAdapter, _, _ string) ([]skill.Result, error) {
		return nil, nil
	})
	t.Cleanup(restoreSkill)
	h := model.Harness{
		ID: "skill-h", Name: "Skill H", Type: model.HarnessSkill,
		Source:       &model.Source{Repo: "owner/repo", Method: "clone"},
		InstallModes: []model.InstallMode{model.ModeLite, model.ModeFull},
		Agents:       []model.Agent{model.AgentClaude},
	}

	// Existing personal skill with the same id: each file must be snapshotted.
	homeDir := t.TempDir()
	skillDir := filepath.Join(homeDir, "skills", "skill-h")
	nested := filepath.Join(skillDir, "references", "notes.md")
	if err := os.MkdirAll(filepath.Dir(nested), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{filepath.Join(skillDir, "SKILL.md"), nested} {
		if err := os.WriteFile(f, []byte("mine"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	paths, _ := capturePrepare(t, h, homeDir)
	for _, want := range []string{filepath.Join(skillDir, "SKILL.md"), nested} {
		if !hasPath(paths, want) {
			t.Errorf("snapshot does not list existing skill file %s. paths: %v", want, paths)
		}
	}

	// Absent skill inside an existing skills dir: the skill dir is a dir hint,
	// so restore removes what the install created there.
	homeDir2 := t.TempDir()
	if err := os.MkdirAll(filepath.Join(homeDir2, "skills"), 0o755); err != nil {
		t.Fatal(err)
	}
	paths2, hints2 := capturePrepare(t, h, homeDir2)
	missing := filepath.Join(homeDir2, "skills", "skill-h")
	if !hasPath(paths2, missing) || !hints2[filepath.Clean(missing)] {
		t.Errorf("absent skill dir %s not recorded as a dir hint. paths: %v", missing, paths2)
	}
}
