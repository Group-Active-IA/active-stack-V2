package external

// lean-agent-workflow 5.3: si Claude ya tiene el plugin engram@engram
// habilitado, el MCP engram NO se registra otra vez en ~/.claude.json (el
// agente terminaría con dos servidores engram: mcp__engram__* y
// mcp__plugin_engram_engram__*). Ante cualquier duda — archivo ausente,
// corrupto, plugin deshabilitado o no instalado — se registra como antes.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/model"
)

// pluginFakeAdapter es claudeFixAdapter + las rutas de plugins, resueltas
// bajo el home del test como lo hace el adapter real.
type pluginFakeAdapter struct {
	claudeFixAdapter
}

func (p *pluginFakeAdapter) SettingsPath(home string) string {
	return filepath.Join(home, ".claude", "settings.json")
}
func (p *pluginFakeAdapter) InstalledPluginsPath(home string) string {
	return filepath.Join(home, ".claude", "plugins", "installed_plugins.json")
}
func (p *pluginFakeAdapter) MCPPluginID(server string) string {
	if server == "engram" {
		return "engram@engram"
	}
	return ""
}

const (
	settingsEnabled   = `{"enabledPlugins":{"engram@engram":true,"codex@openai-codex":true}}`
	settingsDisabled  = `{"enabledPlugins":{"engram@engram":false}}`
	settingsNoPlugins = `{"theme":"dark"}`
	installedEngram   = `{"version":2,"plugins":{"engram@engram":[{"scope":"user","installPath":"/x","version":"0.1.0"}]}}`
	installedOther    = `{"version":2,"plugins":{"codex@openai-codex":[{"scope":"user"}]}}`
	installedEmpty    = `{"version":2,"plugins":{"engram@engram":[]}}`
)

func writeHomeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func engramServers(t *testing.T, home string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(home, ".claude.json"))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	var obj map[string]any
	if err := json.Unmarshal(data, &obj); err != nil {
		t.Fatalf("~/.claude.json not JSON: %v", err)
	}
	servers, _ := obj["mcpServers"].(map[string]any)
	return servers
}

func registerEngram(t *testing.T, home string) Result {
	t.Helper()
	origSnap := snapshotterCreate
	snapshotterCreate = func(string, []string) error { return nil }
	t.Cleanup(func() { snapshotterCreate = origSnap })
	mcp := model.MCP{Name: "engram", Command: "engram", Args: []string{"mcp"}}
	res, err := registerStdioMCP(mcp, []AgentAdapter{&pluginFakeAdapter{}}, home, "engram")
	if err != nil {
		t.Fatalf("registerStdioMCP: %v", err)
	}
	return res
}

// TestRegisterStdioMCP_EngramPluginEnabled_SkipsAndReports es el caso real
// de la máquina de Juani: plugin habilitado en settings.json e instalado en
// installed_plugins.json → no se toca ~/.claude.json y el salto se informa.
func TestRegisterStdioMCP_EngramPluginEnabled_SkipsAndReports(t *testing.T) {
	home := t.TempDir()
	a := &pluginFakeAdapter{}
	writeHomeFile(t, a.SettingsPath(home), settingsEnabled)
	writeHomeFile(t, a.InstalledPluginsPath(home), installedEngram)

	res := registerEngram(t, home)

	if _, err := os.Stat(filepath.Join(home, ".claude.json")); !os.IsNotExist(err) {
		t.Fatalf("~/.claude.json must not be written when the plugin provides engram (stat err=%v)", err)
	}
	if len(res.ConfigFiles) != 0 {
		t.Fatalf("ConfigFiles = %v, want none", res.ConfigFiles)
	}
	if len(res.Notes) != 1 || !strings.Contains(res.Notes[0], "engram@engram") || !strings.Contains(res.Notes[0], "engram") {
		t.Fatalf("Notes = %q, want one note naming the plugin engram@engram", res.Notes)
	}
}

// TestRegisterStdioMCP_EngramPluginEnabled_ExistingEntryIsLeftAlone: si
// ~/.claude.json ya tenía un mcpServers.engram (duplicado de una instalación
// anterior), no se borra — el instalador nunca borra lo que no puede probar
// propio — ni se reescribe; la nota lo dice para que se limpie a mano.
func TestRegisterStdioMCP_EngramPluginEnabled_ExistingEntryIsLeftAlone(t *testing.T) {
	home := t.TempDir()
	a := &pluginFakeAdapter{}
	writeHomeFile(t, a.SettingsPath(home), settingsEnabled)
	writeHomeFile(t, a.InstalledPluginsPath(home), installedEngram)
	before := `{"mcpServers":{"engram":{"type":"stdio","command":"engram","args":["mcp"]}},"other":1}`
	writeHomeFile(t, filepath.Join(home, ".claude.json"), before)

	res := registerEngram(t, home)

	got, _ := os.ReadFile(filepath.Join(home, ".claude.json"))
	if string(got) != before {
		t.Fatalf("~/.claude.json was modified:\n%s", got)
	}
	if len(res.Notes) != 1 || !strings.Contains(res.Notes[0], "mcpServers.engram") {
		t.Fatalf("note must flag the pre-existing duplicate entry for manual cleanup, got %q", res.Notes)
	}
}

// TestRegisterStdioMCP_EngramPluginNotProviding_RegistersAsBefore triangula
// cada caso de duda: se registra igual que antes del cambio.
func TestRegisterStdioMCP_EngramPluginNotProviding_RegistersAsBefore(t *testing.T) {
	cases := []struct {
		name      string
		settings  string // "" = archivo ausente
		installed string // "" = archivo ausente
		wantNote  bool   // archivo corrupto: se registra, pero se declara
	}{
		{name: "sin-archivos"},
		{name: "settings-sin-enabledPlugins", settings: settingsNoPlugins, installed: installedEngram},
		{name: "plugin-deshabilitado", settings: settingsDisabled, installed: installedEngram},
		{name: "habilitado-pero-no-instalado", settings: settingsEnabled, installed: installedOther},
		{name: "habilitado-instalado-vacio", settings: settingsEnabled, installed: installedEmpty},
		{name: "habilitado-sin-installed_plugins", settings: settingsEnabled},
		{name: "settings-corrupto", settings: `{"enabledPlugins": {`, installed: installedEngram, wantNote: true},
		{name: "installed-corrupto", settings: settingsEnabled, installed: `not json`, wantNote: true},
		{name: "enabled-no-booleano", settings: `{"enabledPlugins":{"engram@engram":"yes"}}`, installed: installedEngram},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			a := &pluginFakeAdapter{}
			if tc.settings != "" {
				writeHomeFile(t, a.SettingsPath(home), tc.settings)
			}
			if tc.installed != "" {
				writeHomeFile(t, a.InstalledPluginsPath(home), tc.installed)
			}
			res := registerEngram(t, home)
			if _, ok := engramServers(t, home)["engram"]; !ok {
				t.Fatalf("engram MCP must still be registered in this case")
			}
			if tc.wantNote != (len(res.Notes) == 1) {
				t.Fatalf("Notes = %q, wantNote=%v", res.Notes, tc.wantNote)
			}
		})
	}
}

// TestRegisterStdioMCP_AdapterWithoutPluginCapability_Unchanged: un adapter
// sin PluginMCPAdapter (Codex, OpenCode, el wrapper de target proyecto)
// registra exactamente como antes, aunque el home tenga el plugin.
func TestRegisterStdioMCP_AdapterWithoutPluginCapability_Unchanged(t *testing.T) {
	home := t.TempDir()
	a := &pluginFakeAdapter{}
	writeHomeFile(t, a.SettingsPath(home), settingsEnabled)
	writeHomeFile(t, a.InstalledPluginsPath(home), installedEngram)

	origSnap := snapshotterCreate
	snapshotterCreate = func(string, []string) error { return nil }
	defer func() { snapshotterCreate = origSnap }()
	mcp := model.MCP{Name: "engram", Command: "engram", Args: []string{"mcp"}}
	res, err := registerStdioMCP(mcp, []AgentAdapter{&claudeFixAdapter{}}, home, "engram")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := engramServers(t, home)["engram"]; !ok || len(res.Notes) != 0 {
		t.Fatalf("adapter without plugin capability must register as before (notes %q)", res.Notes)
	}
}

// TestMaybeRegisterStdioMCP_PropagatesNotes: la nota sobrevive al merge con
// el resultado del binario, que es lo que ve el step del instalador.
func TestMaybeRegisterStdioMCP_PropagatesNotes(t *testing.T) {
	home := t.TempDir()
	a := &pluginFakeAdapter{}
	writeHomeFile(t, a.SettingsPath(home), settingsEnabled)
	writeHomeFile(t, a.InstalledPluginsPath(home), installedEngram)
	h := model.Harness{ID: "engram", External: &model.External{
		Method: "homebrew",
		MCP:    &model.MCP{Name: "engram", Command: "engram", Args: []string{"mcp"}},
	}}
	res, err := maybeRegisterStdioMCP(Result{BinaryPath: "/bin/engram", AlreadyInstalled: true}, h, []AgentAdapter{a}, home)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Notes) != 1 || res.BinaryPath != "/bin/engram" {
		t.Fatalf("got %+v, want the plugin note and the binary path kept", res)
	}
}
