package skill_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/harness/skill"
	"github.com/Group-Active-IA/active-stack/internal/model"
)

// ──────────────────────────────────────────────────────────────────────────────
// Review 2026-10-29 (8.1): every install cloned every skill over the network
// and only then compared SKILL.md. A runner that can return output asks
// `git ls-remote` for the remote commit first; when it matches the commit
// recorded at the last install and the skill is on disk, the clone is skipped.
// ──────────────────────────────────────────────────────────────────────────────

// outputRunner is a stubRunner that also answers ls-remote.
type outputRunner struct {
	stubRunner
	remoteSHA string
	lsErr     error
	lsCalls   int
}

func (r *outputRunner) Output(_ context.Context, args []string) ([]byte, error) {
	r.lsCalls++
	if r.lsErr != nil {
		return nil, r.lsErr
	}
	return []byte(r.remoteSHA + "\tHEAD\n"), nil
}

func cloneWritesSkill(args []string) {
	dest := args[len(args)-1]
	_ = os.MkdirAll(dest, 0o755)
	_ = os.WriteFile(filepath.Join(dest, "SKILL.md"), []byte("# cached-skill"), 0o644)
}

func clones(r *outputRunner) int {
	n := 0
	for _, c := range r.called {
		if len(c) > 1 && c[1] == "clone" {
			n++
		}
	}
	return n
}

func installCached(t *testing.T, r *outputRunner, skillsDir string) skill.Result {
	t.Helper()
	h := makeHarness("cached-skill", "clone", "owner/cached-skill", false)
	res, err := skill.NewInstaller(r, nil).Install(context.Background(), h,
		[]skill.AgentAdapter{fakeAdapter{agent: model.AgentClaude, skillsDir: skillsDir}}, t.TempDir(), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return res[0]
}

func TestClone_SkipsWhenRemoteCommitUnchanged(t *testing.T) {
	skillsDir := filepath.Join(t.TempDir(), "skills")
	r := &outputRunner{stubRunner: stubRunner{sideEffect: cloneWritesSkill}, remoteSHA: "aaa111"}

	installCached(t, r, skillsDir)
	if clones(r) != 1 {
		t.Fatalf("first install must clone once, got %d", clones(r))
	}
	second := installCached(t, r, skillsDir)
	if clones(r) != 1 {
		t.Errorf("same remote commit: second install cloned again (%d clones)", clones(r))
	}
	if !second.AlreadyInstalled {
		t.Error("skipped clone must report AlreadyInstalled")
	}
}

func TestClone_ClonesWhenRemoteCommitChanged(t *testing.T) {
	skillsDir := filepath.Join(t.TempDir(), "skills")
	r := &outputRunner{stubRunner: stubRunner{sideEffect: cloneWritesSkill}, remoteSHA: "aaa111"}
	installCached(t, r, skillsDir)
	r.remoteSHA = "bbb222"
	installCached(t, r, skillsDir)
	if clones(r) != 2 {
		t.Errorf("new remote commit must clone, got %d clones", clones(r))
	}
}

func TestClone_ClonesWhenSkillMissingFromDisk(t *testing.T) {
	skillsDir := filepath.Join(t.TempDir(), "skills")
	r := &outputRunner{stubRunner: stubRunner{sideEffect: cloneWritesSkill}, remoteSHA: "aaa111"}
	installCached(t, r, skillsDir)
	_ = os.Remove(filepath.Join(skillsDir, "cached-skill", "SKILL.md"))
	installCached(t, r, skillsDir)
	if clones(r) != 2 {
		t.Errorf("a deleted SKILL.md must be reinstalled, got %d clones", clones(r))
	}
}

func TestClone_LsRemoteFailureFallsBackToClone(t *testing.T) {
	skillsDir := filepath.Join(t.TempDir(), "skills")
	r := &outputRunner{stubRunner: stubRunner{sideEffect: cloneWritesSkill}, lsErr: errors.New("offline")}
	installCached(t, r, skillsDir)
	installCached(t, r, skillsDir)
	if clones(r) != 2 {
		t.Errorf("without ls-remote the clone path must stay as before, got %d clones", clones(r))
	}
}

// A restore or a hand edit leaves the stamp but changes the content; the stamp
// must not vouch for a SKILL.md that is no longer the one it installed.
func TestClone_ClonesWhenInstalledSkillChanged(t *testing.T) {
	skillsDir := filepath.Join(t.TempDir(), "skills")
	r := &outputRunner{stubRunner: stubRunner{sideEffect: cloneWritesSkill}, remoteSHA: "aaa111"}
	installCached(t, r, skillsDir)
	if err := os.WriteFile(filepath.Join(skillsDir, "cached-skill", "SKILL.md"), []byte("# restored personal version"), 0o644); err != nil {
		t.Fatal(err)
	}
	installCached(t, r, skillsDir)
	if clones(r) != 2 {
		t.Errorf("a changed SKILL.md must not be skipped by an old stamp, got %d clones", clones(r))
	}
}
