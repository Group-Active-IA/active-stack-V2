package command_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/assets"
	"github.com/Group-Active-IA/active-stack/internal/harness/command"
	"github.com/Group-Active-IA/active-stack/internal/model"
)

// Verify is a shipped executable instruction asset: exercise its actual installer,
// rather than checking an unrelated documentation copy.
func TestInstalledVerifySelectsWorkflowBeforeLegacyInstructions(t *testing.T) {
	dir := t.TempDir()
	_, err := command.NewInstaller(assets.CommandsFS).Install([]command.AgentAdapter{fakeAdapter{agent: model.AgentClaude, commandsDir: dir, variantKey: "claude"}}, "", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, "opsx", "verify.md"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	lean := strings.Index(s, "## Lean workflow")
	legacy := strings.Index(s, "## Legacy workflow")
	if lean < 0 || legacy <= lean {
		t.Fatal("installed verify lacks explicit lean-first and legacy-only dispatch")
	}
	if !strings.Contains(s[lean:legacy], "role-breaker") || !strings.Contains(s[lean:legacy], "role-auditor") {
		t.Fatal("lean dispatch does not bind breaker and auditor contracts")
	}
	if !strings.Contains(s[lean:legacy], "Do not continue into the legacy steps") {
		t.Fatal("lean dispatch can fall through into the legacy loop")
	}
}

func TestInstalledVerifyFollowsDepthExecutorsAndReportIndex(t *testing.T) {
	dir := t.TempDir()
	_, err := command.NewInstaller(assets.CommandsFS).Install([]command.AgentAdapter{fakeAdapter{agent: model.AgentClaude, commandsDir: dir, variantKey: "claude"}}, "", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, "opsx", "verify.md"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	lean := s[strings.Index(s, "## Lean workflow"):strings.Index(s, "## Legacy workflow")]
	for _, required := range []string{
		"depth table",
		"LOW runs one breaker round",
		"Role executors",
		"informes/verify-<attempt>-breaker.md",
		"LEEME.md",
	} {
		if !strings.Contains(lean, required) {
			t.Errorf("lean verify missing %q", required)
		}
	}
}
