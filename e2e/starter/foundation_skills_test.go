package starter

// Foundation-skill install post-conditions.
//
// These tests exercise the REAL embedded catalog (catalog.Load), not a fixture
// catalog: a harness that is missing from harnesses.yaml fails here. What they
// assert is the thing that matters to a user — the skill's SKILL.md is on disk
// in the agent's skills directory after an install — rather than the presence
// of a line in a YAML file, which no user ever reads.
//
// Hermetic: each skill's source.repo is redirected to a local git fixture
// served over file://, so no network is touched.

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/cmd/active-stack/headless"
	"github.com/Group-Active-IA/active-stack/internal/catalog"
	"github.com/Group-Active-IA/active-stack/internal/install"
	"github.com/Group-Active-IA/active-stack/internal/model"
)

// fixtureBody is the SKILL.md body written for a given harness id. It embeds
// the id so a test can tell WHICH fixture landed at a path, not merely that
// something did.
func fixtureBody(id string) string {
	return "# " + id + "\n\nfixture-body-for:" + id + "\n"
}

// redirectSkillRepo points a real-catalog skill harness at a local fixture git
// repo. It fails loudly when the harness is absent from the embedded catalog,
// or when it is not a clone-able skill.
//
// Source is a *model.Source, and catalog.parse stores the same pointer in both
// c.Harnesses and the ByID index, so mutating it here is visible on both paths.
func redirectSkillRepo(t *testing.T, c *catalog.Catalog, id, fileURL string) {
	t.Helper()
	for i := range c.Harnesses {
		h := &c.Harnesses[i]
		if h.ID != id {
			continue
		}
		if h.Type != model.HarnessSkill {
			t.Fatalf("harness %q: type = %q, want %q", id, h.Type, model.HarnessSkill)
		}
		if h.Source == nil {
			t.Fatalf("harness %q: source is nil, want a clone source", id)
		}
		h.Source.Repo = fileURL
		h.Source.Ref = ""
		h.Source.Path = ""
		return
	}
	t.Fatalf("harness %q is not in the embedded catalog (internal/catalog/harnesses.yaml)", id)
}

// installFoundationSkills runs the real headless install path for the given
// catalog skill ids with their repos redirected to fixtures, and returns the
// resolved Claude skills directory.
func installFoundationSkills(t *testing.T, ids ...string) string {
	t.Helper()
	requireGit(t)

	c, err := catalog.Load()
	if err != nil {
		t.Fatalf("catalog.Load(): %v", err)
	}
	for _, id := range ids {
		repo := buildFixtureRepo(t, fixtureBody(id), "")
		redirectSkillRepo(t, c, id, repo.FileURL)
	}

	homeDir := t.TempDir()
	params := headless.ParsedFlags{
		Yes:           true,
		HomeDir:       homeDir,
		NoSelfInstall: true,
		Intent: install.Intent{
			Agents: []model.Agent{model.AgentClaude},
			Mode:   model.ModeCustom,
			Custom: ids,
		},
	}

	var out bytes.Buffer
	if code := headless.RunHeadless(params, c, newClaudeRegistry(), &out); code != 0 {
		t.Fatalf("RunHeadless exited %d; output:\n%s", code, out.String())
	}
	return claudeSkillsDir(homeDir)
}

// assertFixtureLanded asserts that the SKILL.md on disk for id is the fixture
// written for THAT id — not merely that some file exists there.
func assertFixtureLanded(t *testing.T, skillsDir, id string) {
	t.Helper()
	path := filepath.Join(skillsDir, id, "SKILL.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("skill %q: SKILL.md not on disk at %s: %v", id, path, err)
		return
	}
	if want := "fixture-body-for:" + id; !strings.Contains(string(data), want) {
		t.Errorf("skill %q: SKILL.md at %s does not contain %q; got:\n%s", id, path, want, string(data))
	}
}

// TestFoundation_ProjectScaffold_LandsOnDisk is the RED test for the
// project-scaffold catalog entry: installing it must put its SKILL.md in the
// agent's skills directory.
func TestFoundation_ProjectScaffold_LandsOnDisk(t *testing.T) {
	skillsDir := installFoundationSkills(t, "project-scaffold")
	assertFixtureLanded(t, skillsDir, "project-scaffold")
}

// TestFoundation_DomainProcessPlan_AllThreeLandOnDisk triangulates with a
// different input: the three foundation skills installed in one run, each
// pointed at its OWN fixture repo. Asserting the per-id body catches an
// installer that lands one skill under three names, which the single-skill
// test above cannot see.
func TestFoundation_DomainProcessPlan_AllThreeLandOnDisk(t *testing.T) {
	// Domain → process → plan, the order they run in.
	ids := []string{"kb-creator", "project-scaffold", "roadmap-generator"}

	skillsDir := installFoundationSkills(t, ids...)
	for _, id := range ids {
		assertFixtureLanded(t, skillsDir, id)
	}
}
