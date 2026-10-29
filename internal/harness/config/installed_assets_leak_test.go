package config

import (
	"io/fs"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/assets"
	"github.com/Group-Active-IA/active-stack/internal/harness/roledefs"
)

// Everything below is copied into OTHER people's agents. On 2026-10-29 two
// reviews of this repo found facts of this project inside it: "Juani's written
// decision" in five legacy fragments and in /opsx:verify (an agent elsewhere
// would escalate to someone who is not its operator), and "For Active Stack
// require ci.yml" in the lean ship contract. This guard scans every installable
// asset so the next one fails here instead of in someone's project.
var projectOnlyFacts = []string{"Juani", "For Active Stack require", `E:\ESCRITORIO`}

func TestInstalledAssets_CarryNoProjectOnlyFacts(t *testing.T) {
	check := func(origin, name, content string) {
		for _, fact := range projectOnlyFacts {
			if strings.Contains(content, fact) {
				t.Errorf("%s %s ships project-only fact %q", origin, name, fact)
			}
		}
	}
	walk := func(origin string, fsys fs.FS) {
		err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			raw, readErr := fs.ReadFile(fsys, p)
			if readErr != nil {
				return readErr
			}
			check(origin, p, string(raw))
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", origin, err)
		}
	}
	walk("orchestrator fragment", embeddedAssets)
	walk("command", assets.CommandsFS)
	walk("embedded skill", assets.SkillsFS)
	for _, constant := range []string{leanGovernance, leanEvidence, leanKnowledge, leanMemory, leanShip, leanClaudeDispatch} {
		check("lean contract", "lean.go", constant)
	}
	roles, err := roledefs.DefaultResources()
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range roles {
		check("role", r.Name, string(r.Content))
	}
}
