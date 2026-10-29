package install_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/backup"
	"github.com/Group-Active-IA/active-stack/internal/install"
	"github.com/Group-Active-IA/active-stack/internal/model"
)

// 2026-10-29, closing instalacion-sana with OpenCode: every install wrote its
// snapshot to backups/install, so a second install overwrote the first and
// restore could only go back to the state after the previous install, never
// to the user's original config. Each install gets its own snapshot dir.
func TestBuildPlan_EachInstallHasItsOwnSnapshot(t *testing.T) {
	homeDir := t.TempDir()
	var dirs []string
	restore := install.SetSnapshotCreate(func(dir string, _ []string) (backup.Manifest, error) {
		dirs = append(dirs, dir)
		return backup.Manifest{}, nil
	})
	defer restore()
	cat := &fakeCatalog{harnesses: []model.Harness{{
		ID: "cfg-h", Type: model.HarnessConfig, InstallModes: []model.InstallMode{model.ModeLite}, Agents: []model.Agent{model.AgentClaude},
	}}}
	reg := &fakeRegistry{adapters: map[model.Agent]install.AgentAdapter{model.AgentClaude: fakeAdapter{agent: model.AgentClaude}}}
	for i := 0; i < 2; i++ {
		plan, err := install.BuildPlan(cat, install.Intent{Agents: []model.Agent{model.AgentClaude}, Mode: model.ModeLite},
			install.Options{HomeDir: homeDir, Registry: reg, NoSelfInstall: true})
		if err != nil {
			t.Fatal(err)
		}
		for _, s := range plan.Prepare {
			_ = s.Run()
		}
	}
	if len(dirs) != 2 || dirs[0] == dirs[1] {
		t.Fatalf("two installs must get two snapshot dirs, got %v", dirs)
	}
	for _, d := range dirs {
		if filepath.Dir(d) != filepath.Join(homeDir, ".active-stack", "backups") || !strings.HasPrefix(filepath.Base(d), "install-") {
			t.Errorf("snapshot dir %q, want <home>/.active-stack/backups/install-<time>", d)
		}
	}
}
