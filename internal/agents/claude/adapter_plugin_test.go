package claude_test

import (
	"path/filepath"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/agents/claude"
	"github.com/Group-Active-IA/active-stack/internal/harness/external"
)

// TestClaudeAdapter_PluginMCP (lean-agent-workflow 5.3): las rutas donde
// Claude Code declara plugins habilitados e instalados salen del adapter,
// nunca hardcodeadas en el installer, y el adapter sabe qué plugin trae el
// MCP engram.
func TestClaudeAdapter_PluginMCP(t *testing.T) {
	var a external.PluginMCPAdapter = claude.NewAdapter()
	if got, want := a.SettingsPath(testHome), filepath.Join(testHome, ".claude", "settings.json"); got != want {
		t.Errorf("SettingsPath = %q, want %q", got, want)
	}
	if got, want := a.InstalledPluginsPath(testHome), filepath.Join(testHome, ".claude", "plugins", "installed_plugins.json"); got != want {
		t.Errorf("InstalledPluginsPath = %q, want %q", got, want)
	}
	if got := a.MCPPluginID("engram"); got != "engram@engram" {
		t.Errorf("MCPPluginID(engram) = %q, want engram@engram", got)
	}
	if got := a.MCPPluginID("context7"); got != "" {
		t.Errorf("MCPPluginID(context7) = %q, want empty (no plugin provides it)", got)
	}
}
