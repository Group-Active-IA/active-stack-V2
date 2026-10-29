package external

import (
	"encoding/json"
	"fmt"
	"os"
)

// PluginMCPAdapter is implemented by adapters whose agent can receive an MCP
// server through a plugin (Claude Code today). Paths come from the adapter —
// never hardcoded here — and MCPPluginID names the plugin that ships
// serverName ("" when none does).
//
// It is an optional capability, discovered by type assertion like
// AgentsDir in internal/install: an adapter without it (Codex, OpenCode, the
// project-target wrapper) keeps registering exactly as before.
type PluginMCPAdapter interface {
	SettingsPath(homeDir string) string
	InstalledPluginsPath(homeDir string) string
	MCPPluginID(serverName string) string
}

// pluginProvidesMCP reports whether an enabled plugin already provides
// serverName for this adapter (lean-agent-workflow 5.3). Registering the MCP
// again on top would leave the agent with two copies of the same server
// (e.g. mcp__engram__* and mcp__plugin_engram_engram__*).
//
// "Provided" requires BOTH, because either alone is not proof:
//   - settings.json → enabledPlugins[<id>] == true (the user enabled it), and
//   - installed_plugins.json → plugins[<id>] is a non-empty list (it exists).
//
// Every other case falls back to the previous behaviour (register): missing
// files, no enabledPlugins key, plugin disabled, enabled but not installed.
// A file that exists but cannot be read or parsed also registers, and is
// declared through note so the user knows why the duplicate may appear.
func pluginProvidesMCP(adapter AgentAdapter, homeDir, serverName string) (pluginID string, provided bool, note string) {
	pa, ok := adapter.(PluginMCPAdapter)
	if !ok {
		return "", false, ""
	}
	pluginID = pa.MCPPluginID(serverName)
	if pluginID == "" {
		return "", false, ""
	}

	var settings struct {
		EnabledPlugins map[string]any `json:"enabledPlugins"`
	}
	settingsPath := pa.SettingsPath(homeDir)
	if found, err := readJSONFile(settingsPath, &settings); err != nil {
		return pluginID, false, fmt.Sprintf("no se pudo leer %s (%v): se registra el MCP %q igual, aunque el plugin %s podría estar habilitado", settingsPath, err, serverName, pluginID)
	} else if !found {
		return pluginID, false, ""
	}
	if enabled, _ := settings.EnabledPlugins[pluginID].(bool); !enabled {
		return pluginID, false, ""
	}

	var installed struct {
		Plugins map[string][]json.RawMessage `json:"plugins"`
	}
	installedPath := pa.InstalledPluginsPath(homeDir)
	if found, err := readJSONFile(installedPath, &installed); err != nil {
		return pluginID, false, fmt.Sprintf("no se pudo leer %s (%v): se registra el MCP %q igual, aunque el plugin %s figura habilitado", installedPath, err, serverName, pluginID)
	} else if !found {
		return pluginID, false, ""
	}
	if len(installed.Plugins[pluginID]) == 0 {
		return pluginID, false, ""
	}
	return pluginID, true, ""
}

// readJSONFile decodes path into v. found=false with nil error means the
// file does not exist — a legitimate state, not a failure.
func readJSONFile(path string, v any) (found bool, err error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(data, v); err != nil {
		return false, err
	}
	return true, nil
}

// pluginSkipNote is the message reported when registration is skipped. It
// also flags a pre-existing entry for serverName in configPath: that one is
// left alone (never delete what we cannot prove is ours) and needs manual
// cleanup to stop the duplicate.
func pluginSkipNote(pluginID, serverName, configPath string) string {
	note := fmt.Sprintf("MCP %q no registrado en %s: el plugin %s ya está habilitado y lo provee", serverName, configPath, pluginID)
	var cfg struct {
		MCPServers map[string]json.RawMessage `json:"mcpServers"`
	}
	if found, err := readJSONFile(configPath, &cfg); err == nil && found {
		if _, dup := cfg.MCPServers[serverName]; dup {
			note += fmt.Sprintf("; ya existe mcpServers.%s en ese archivo (duplicado de antes): quedó intacto, sacarlo a mano si sobra", serverName)
		}
	}
	return note
}
