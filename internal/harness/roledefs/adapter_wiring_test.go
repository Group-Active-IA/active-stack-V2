package roledefs

import (
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/agents/claude"
	"github.com/Group-Active-IA/active-stack/internal/agents/codex"
	"github.com/Group-Active-IA/active-stack/internal/agents/opencode"
)

// TestRealAdapters_SatisfyAgentAdapter confirms the three real adapters
// structurally satisfy roledefs.AgentAdapter — this is what actually wires
// Install/Uninstall/Health to internal/install and internal/uninstall
// without widening the shared agents.Adapter interface (low blast radius,
// same pattern internal/harness/skill.AgentAdapter already uses).
func TestRealAdapters_SatisfyAgentAdapter(t *testing.T) {
	var _ AgentAdapter = claude.NewAdapter()
	var _ AgentAdapter = codex.NewAdapter()
	var _ AgentAdapter = opencode.NewAdapter()
}

// TestOnlyClaudeSupportsAgentsDirToday triangula el capability real: Claude
// resuelve a una ruta no vacía bajo el home dado; Codex y OpenCode devuelven
// "" (J3/J4, TBD) y por lo tanto se saltean en Install/Uninstall/Health.
func TestOnlyClaudeSupportsAgentsDirToday(t *testing.T) {
	home := "/fake/home"
	if got := claude.NewAdapter().AgentsDir(home); got == "" {
		t.Error("claude.Adapter.AgentsDir returned empty — Claude Code sub-agents should be supported")
	}
	if got := codex.NewAdapter().AgentsDir(home); got != "" {
		t.Errorf("codex.Adapter.AgentsDir = %q, want \"\" (capability not confirmed yet, J3/J4)", got)
	}
	if got := opencode.NewAdapter().AgentsDir(home); got != "" {
		t.Errorf("opencode.Adapter.AgentsDir = %q, want \"\" (capability not confirmed yet, J3/J4)", got)
	}
}
