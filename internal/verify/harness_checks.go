package verify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Group-Active-IA/active-stack/internal/harness/config/permissions"
	"github.com/Group-Active-IA/active-stack/internal/model"
	toml "github.com/pelletier/go-toml/v2"
)

// sddOrchestratorSectionID is the marker section written by the config harness
// installer (internal/harness/config). It is the string used by
// filemerge.InjectMarkdownSection when injecting the orchestrator block.
const sddOrchestratorSectionID = "sdd-orchestrator"

// ChecksForHarness derives the verify.Check slice for a single harness,
// one check per (harness × adapter). It dispatches by HarnessType:
//
//   - skill    → SKILL.md exists and is non-empty in adapter.SkillsDir
//   - config   → idempotent marker present exactly once in InstructionsPath
//     (special case "permissions" → checks SettingsPath for key)
//   - external → MCP config parseable (hard) + binary in PATH (Soft)
//
// Paths are ALWAYS resolved via adapters — never hardcoded.
func ChecksForHarness(h model.Harness, adapters []Adapter, homeDir string) []Check {
	var checks []Check
	for _, adapter := range adapters {
		checks = append(checks, checksForHarnessAdapter(h, adapter, homeDir)...)
	}
	return checks
}

// checksForHarnessAdapter derives checks for one (harness, adapter) pair.
func checksForHarnessAdapter(h model.Harness, adapter Adapter, homeDir string) []Check {
	switch h.Type {
	case model.HarnessSkill:
		return checkSkill(h, adapter, homeDir)
	case model.HarnessConfig:
		if h.ID == "permissions" {
			return checkPermissions(h, adapter, homeDir)
		}
		return checkConfig(h, adapter, homeDir)
	case model.HarnessExternal:
		return checkExternal(h, adapter, homeDir)
	default:
		return nil
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// skill checks
// ─────────────────────────────────────────────────────────────────────────────

func checkSkill(h model.Harness, adapter Adapter, homeDir string) []Check {
	skillsDir := adapter.SkillsDir(homeDir)
	agentID := string(adapter.Agent())

	return []Check{
		{
			ID:          fmt.Sprintf("skill:%s:%s", h.ID, agentID),
			Description: fmt.Sprintf("SKILL.md present and non-empty for %s/%s", agentID, h.ID),
			// A best-effort harness that failed to install is a warning, not a
			// hard failure: the operator was already warned at install time, and
			// the verify hook must not roll back a successful (soft) install.
			Soft: h.BestEffort,
			Run: func(_ context.Context) error {
				skillMD := filepath.Join(skillsDir, h.ID, "SKILL.md")
				info, err := os.Stat(skillMD)
				if err != nil {
					if os.IsNotExist(err) {
						return fmt.Errorf("SKILL.md not found at %q", skillMD)
					}
					return fmt.Errorf("stat SKILL.md: %w", err)
				}
				if info.Size() == 0 {
					return fmt.Errorf("SKILL.md is empty at %q", skillMD)
				}
				return nil
			},
		},
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// config checks
// ─────────────────────────────────────────────────────────────────────────────

// checkConfig verifies that the injected block's marker appears exactly once
// in the agent's instructions file. Uses the same sectionID the config installer
// writes (internal/harness/config/inject.go: "sdd-orchestrator").
func checkConfig(h model.Harness, adapter Adapter, homeDir string) []Check {
	// Primary-agent delivery (e.g. OpenCode) writes to the settings JSON, not
	// the instructions file — verify it where it actually lives.
	if adapter.ConfigDelivery() == model.ConfigDeliveryPrimaryAgent {
		return checkPrimaryAgentConfig(h, adapter, homeDir)
	}

	instrPath := adapter.InstructionsPath(homeDir)
	agentID := string(adapter.Agent())
	openMarker := "<!-- active-stack:" + sddOrchestratorSectionID + " -->"

	return []Check{
		{
			ID:          fmt.Sprintf("config:%s:%s", h.ID, agentID),
			Description: fmt.Sprintf("idempotent marker present exactly once in %s instructions (%s)", agentID, h.ID),
			Run: func(_ context.Context) error {
				data, err := os.ReadFile(instrPath)
				if err != nil {
					if os.IsNotExist(err) {
						return fmt.Errorf("instructions file not found at %q", instrPath)
					}
					return fmt.Errorf("read instructions file: %w", err)
				}
				count := strings.Count(string(data), openMarker)
				switch {
				case count == 0:
					return fmt.Errorf("marker %q not found in %q (harness not installed?)", openMarker, instrPath)
				case count > 1:
					return fmt.Errorf("marker %q appears %d times in %q (idempotency violation)", openMarker, count, instrPath)
				}
				return nil
			},
		},
	}
}

// permissionsKeyFor returns the top-level JSON key that the permissions
// installer writes for the given agent.
//
// Schema source: internal/harness/config/permissions/overlays.go
//
//   - Claude  → "permissions" (plural)  — claudeCodeOverlayJSON L13
//   - OpenCode → "permission" (singular) — openCodeOverlayJSON L32
//   - default  → "permissions" (plural)  — matches every other overlay that
//     either uses the same schema (Gemini uses "general", VS Code uses
//     "chat.tools.autoApprove") or injects no overlay at all. The default
//     covers future agents until their overlay lands — a failing verify check
//     will catch any mismatch.
func permissionsKeyFor(agent model.Agent) string {
	switch agent {
	case model.AgentOpenCode:
		return "permission"
	default:
		return "permissions"
	}
}

// checkPermissions verifies that the permissions harness is actually applied in
// the agent's settings file (written by internal/harness/config/permissions).
// The key name is resolved per-agent via permissionsKeyFor to match the exact
// JSON key each agent's overlay writes (e.g. "permission" for OpenCode,
// "permissions" for Claude). The error message names the searched key so that
// future schema drift is self-diagnosing (D4 in design.md).
//
// Beyond the key, the check reads the CONTENT and asserts the security floor:
// the permissions harness is the only one the installer declares non-disableable,
// so it is the one place where "the key is there" is not a good enough answer.
// Until 2026-09-08 this check parsed the settings into map[string]json.RawMessage
// and only asserted that the key existed. A Claude settings.json with
// "permissions": {"deny": []}, or an opencode.json whose bash block had been
// reduced to {"*": "allow"}, passed verify green with the deny floor emptied —
// the check that guards the product's only non-negotiable guarantee confirmed
// nothing but the presence of a key with the right name.
//
// The floor is per-agent and comes from the installer, never invented here:
// Claude gets permissions.SecurityFloorDeny(), OpenCode gets
// permissions.OpencodeBashFloorDeny() plus the last-match-wins ordering rule.
// Agents whose overlay writes no deny floor at all (gemini, vscode) keep the
// key-presence check: verifying a floor they never had would be a fabricated
// guarantee, which is the same disease in the other direction.
func checkPermissions(h model.Harness, adapter Adapter, homeDir string) []Check {
	settingsPath := adapter.SettingsPath(homeDir)
	agentID := string(adapter.Agent())
	key := permissionsKeyFor(adapter.Agent())

	description := fmt.Sprintf("permissions key present in %s settings", agentID)
	switch adapter.Agent() {
	case model.AgentClaude, model.AgentOpenCode:
		description = fmt.Sprintf("security floor applied in %s settings", agentID)
	}

	checks := []Check{
		{
			ID:          fmt.Sprintf("permissions:%s:%s", h.ID, agentID),
			Description: description,
			Run: func(_ context.Context) error {
				data, err := os.ReadFile(settingsPath)
				if err != nil {
					if os.IsNotExist(err) {
						return fmt.Errorf("settings file not found at %q", settingsPath)
					}
					return fmt.Errorf("read settings file: %w", err)
				}
				if adapter.Agent() == model.AgentCodex {
					var root map[string]any
					if err := toml.Unmarshal(data, &root); err != nil {
						return fmt.Errorf("parse settings file %q: %w", settingsPath, err)
					}
					for _, expected := range []string{"approval_policy", "sandbox_mode"} {
						if _, ok := root[expected]; !ok {
							return fmt.Errorf("expected key %q not found in %q", expected, settingsPath)
						}
					}
					return nil
				}
				var raw map[string]json.RawMessage
				if err := json.Unmarshal(data, &raw); err != nil {
					return fmt.Errorf("parse settings file %q: %w", settingsPath, err)
				}
				block, ok := raw[key]
				if !ok {
					return fmt.Errorf("expected key %q not found in %q", key, settingsPath)
				}
				switch adapter.Agent() {
				case model.AgentClaude:
					return verifyClaudeSecurityFloor(block, settingsPath)
				case model.AgentOpenCode:
					return verifyOpenCodeSecurityFloor(block, settingsPath)
				default:
					// gemini/vscode overlays carry no deny floor to verify.
					return nil
				}
			},
		},
	}
	if adapter.Agent() == model.AgentCodex {
		if home, err := os.UserHomeDir(); err == nil {
			machineConfig := filepath.Join(home, ".codex", "config.toml")
			if filepath.Clean(settingsPath) != filepath.Clean(machineConfig) {
				checks = append(checks, checkCodexProjectTrust(settingsPath, machineConfig))
			}
		}
	}
	return checks
}

// verifyClaudeSecurityFloor asserts that every rule the installer writes into
// permissions.deny is actually present. Claude evaluates deny → ask → allow, so
// the deny list is the whole guarantee: a rule that is not in it does not apply,
// no matter what defaultMode says. The floor is a MINIMUM, not an equality —
// extra user rules are fine, missing floor rules are not.
func verifyClaudeSecurityFloor(block json.RawMessage, settingsPath string) error {
	var shape struct {
		Deny []string `json:"deny"`
	}
	if err := json.Unmarshal(block, &shape); err != nil {
		return fmt.Errorf("parse %q permissions block in %q: %w", "permissions", settingsPath, err)
	}
	present := make(map[string]struct{}, len(shape.Deny))
	for _, rule := range shape.Deny {
		present[rule] = struct{}{}
	}
	var missing []string
	for _, want := range permissions.SecurityFloorDeny() {
		if _, ok := present[want]; !ok {
			missing = append(missing, want)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf(
			"security floor not applied in %q: permissions.deny is missing %d of %d required rule(s): %s",
			settingsPath, len(missing), len(permissions.SecurityFloorDeny()), strings.Join(missing, ", "))
	}
	return nil
}

// verifyOpenCodeSecurityFloor asserts the opencode floor, which is a different
// shape and a different failure mode than Claude's.
//
// opencode resolves permission.bash with LAST-MATCH-WINS — the opposite of
// Claude. A deny rule placed BEFORE a broader pattern is silently overwritten by
// it, so "the rule is in the file" and "the rule applies" are two different
// questions. This checks both: the floor rules are denied, AND no later entry
// whose pattern also matches them downgrades them to ask/allow.
func verifyOpenCodeSecurityFloor(block json.RawMessage, settingsPath string) error {
	var shape struct {
		Bash json.RawMessage `json:"bash"`
	}
	if err := json.Unmarshal(block, &shape); err != nil {
		return fmt.Errorf("parse %q permission block in %q: %w", "permission", settingsPath, err)
	}
	if len(shape.Bash) == 0 {
		return fmt.Errorf(
			"security floor not applied in %q: permission.bash is absent, so none of the deny rules (%s) apply",
			settingsPath, strings.Join(permissions.OpencodeBashFloorDeny(), ", "))
	}
	pairs, err := orderedStringObject(shape.Bash)
	if err != nil {
		return fmt.Errorf("parse permission.bash in %q: %w", settingsPath, err)
	}

	for _, want := range permissions.OpencodeBashFloorDeny() {
		idx := -1
		for i, p := range pairs {
			if p.key == want {
				idx = i
			}
		}
		if idx < 0 || pairs[idx].value != "deny" {
			return fmt.Errorf(
				"security floor not applied in %q: permission.bash does not deny %q",
				settingsPath, want)
		}
		for _, later := range pairs[idx+1:] {
			if later.value == "deny" {
				continue
			}
			if globsIntersect(later.key, want) {
				return fmt.Errorf(
					"security floor not applied in %q: permission.bash denies %q but the later entry %q=%q overrides it (opencode is last-match-wins)",
					settingsPath, want, later.key, later.value)
			}
		}
	}
	return nil
}

// orderedPair is one entry of a JSON object, in file order.
type orderedPair struct {
	key   string
	value string
}

// orderedStringObject decodes a flat JSON object of string values PRESERVING
// key order. encoding/json's map decoding loses order, and for opencode the
// order IS the semantics (last match wins), so a map would silently discard the
// very thing being verified.
func orderedStringObject(raw json.RawMessage) ([]orderedPair, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if delim, ok := tok.(json.Delim); !ok || delim != '{' {
		return nil, fmt.Errorf("expected a JSON object, got %v", tok)
	}
	var pairs []orderedPair
	for dec.More() {
		keyTok, err := dec.Token()
		if err != nil {
			return nil, err
		}
		key, ok := keyTok.(string)
		if !ok {
			return nil, fmt.Errorf("expected a string key, got %v", keyTok)
		}
		var value string
		if err := dec.Decode(&value); err != nil {
			return nil, fmt.Errorf("value of key %q is not a string: %w", key, err)
		}
		pairs = append(pairs, orderedPair{key: key, value: value})
	}
	return pairs, nil
}

// globsIntersect reports whether two opencode bash patterns have at least one
// command in common — i.e. whether the languages they denote intersect.
//
// This is the question the last-match-wins check actually needs, and it is NOT
// the same as matching one pattern against the other's literal text. Until
// 2026-09-08 this was `pattern matches the floor rule's text`, which answered
// "no overlap" for `rm -rf /home/*` against the floor's `rm -rf /*` — so an
// entry appended by hand at the end of opencode.json neutralised the floor for
// everything under /home and verify still reported green. That is the most
// likely way for the floor to break in the field, and it was the one case the
// check could not see.
//
// Only "*" is special on either side, and it matches any run of characters
// ("/" included). Both patterns are walked in lockstep over the hypothetical
// common string with memoised recursion: every transition advances at least one
// index, so the state graph is acyclic and O(len(a)*len(b)) states suffice.
func globsIntersect(a, b string) bool {
	const (
		unknown int8 = iota
		no
		yes
	)
	memo := make([][]int8, len(a)+1)
	for i := range memo {
		memo[i] = make([]int8, len(b)+1)
	}

	// solve reports whether some string is generated by both a[i:] and b[j:].
	var solve func(i, j int) bool
	solve = func(i, j int) bool {
		if m := memo[i][j]; m != unknown {
			return m == yes
		}
		var res bool
		switch {
		case i == len(a) && j == len(b):
			// Both exhausted: the empty string is common to both.
			res = true
		case i < len(a) && a[i] == '*':
			// The star ends here, or it absorbs one more character — which the
			// other side must also account for, hence the advance on j.
			res = solve(i+1, j) || (j < len(b) && solve(i, j+1))
		case j < len(b) && b[j] == '*':
			res = solve(i, j+1) || (i < len(a) && solve(i+1, j))
		case i == len(a) || j == len(b):
			// One side still demands literal characters the other cannot supply.
			res = false
		default:
			res = a[i] == b[j] && solve(i+1, j+1)
		}
		if res {
			memo[i][j] = yes
		} else {
			memo[i][j] = no
		}
		return res
	}
	return solve(0, 0)
}

func checkCodexProjectTrust(projectConfig, machineConfig string) Check {
	projectRoot := filepath.Dir(filepath.Dir(projectConfig))
	return Check{
		ID:          "codex:project-trust:" + projectRoot,
		Description: "Codex project configuration is trusted",
		Soft:        true,
		Run: func(context.Context) error {
			raw, err := os.ReadFile(machineConfig)
			if err != nil {
				return fmt.Errorf("cannot confirm Codex trust for %q", projectRoot)
			}
			var root struct {
				Projects map[string]struct {
					TrustLevel string `toml:"trust_level"`
				} `toml:"projects"`
			}
			if err := toml.Unmarshal(raw, &root); err != nil {
				return fmt.Errorf("cannot parse Codex machine config to confirm trust: %w", err)
			}
			for path, project := range root.Projects {
				if filepath.Clean(path) == filepath.Clean(projectRoot) && project.TrustLevel == "trusted" {
					return nil
				}
			}
			return fmt.Errorf("Codex will ignore %q until the project is trusted", projectConfig)
		},
	}
}

// checkPrimaryAgentConfig verifies that a config harness delivered as a primary
// agent is registered under agent.<id> with mode:primary in the settings JSON
// (e.g. opencode.json) — the entry that makes it a tab-able agent. Mirrors the
// install-time registration in internal/harness/config/primary_agent.go.
func checkPrimaryAgentConfig(h model.Harness, adapter Adapter, homeDir string) []Check {
	settingsPath := adapter.SettingsPath(homeDir)
	agentID := string(adapter.Agent())

	return []Check{
		{
			ID:          fmt.Sprintf("config:%s:%s", h.ID, agentID),
			Description: fmt.Sprintf("primary agent %q registered (mode:primary) in %s settings", h.ID, agentID),
			Run: func(_ context.Context) error {
				data, err := os.ReadFile(settingsPath)
				if err != nil {
					if os.IsNotExist(err) {
						return fmt.Errorf("settings file not found at %q", settingsPath)
					}
					return fmt.Errorf("read settings file: %w", err)
				}
				var root struct {
					Agent map[string]struct {
						Mode string `json:"mode"`
					} `json:"agent"`
				}
				if err := json.Unmarshal(data, &root); err != nil {
					return fmt.Errorf("parse settings file %q: %w", settingsPath, err)
				}
				entry, ok := root.Agent[h.ID]
				if !ok {
					return fmt.Errorf("agent %q not found in %q (harness not installed?)", h.ID, settingsPath)
				}
				if entry.Mode != "primary" {
					return fmt.Errorf("agent %q has mode %q, want \"primary\" (not tab-able) in %q", h.ID, entry.Mode, settingsPath)
				}
				return nil
			},
		},
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// external checks
// ─────────────────────────────────────────────────────────────────────────────

func checkExternal(h model.Harness, adapter Adapter, homeDir string) []Check {
	if h.External == nil {
		return nil
	}

	var checks []Check

	switch h.External.Method {
	case "mcp":
		// Hard check: MCP config file is present and parseable JSON.
		checks = append(checks, checkMCPConfig(h, adapter, homeDir))
		// Soft check: remote endpoint is reachable (warning only).
		if h.External.URL != "" {
			checks = append(checks, checkMCPEndpoint(h))
		}

	case "homebrew", "npm", "go-install", "download":
		if adapter.Agent() == model.AgentCodex && h.External.MCP != nil {
			checks = append(checks, checkMCPConfig(h, adapter, homeDir))
		}
		// Soft check: binary in PATH after install.
		// Network/exec failures → warning, not hard failure.
		checks = append(checks, checkBinaryInPATH(h))
	}

	return checks
}

func checkMCPConfig(h model.Harness, adapter Adapter, homeDir string) Check {
	mcpPath := adapter.MCPConfigPath(homeDir, h.ID)
	agentID := string(adapter.Agent())
	serverName := h.ID
	if h.External != nil && h.External.MCP != nil && h.External.MCP.Name != "" {
		serverName = h.External.MCP.Name
	}

	return Check{
		ID:          fmt.Sprintf("external:mcp-config:%s:%s", h.ID, agentID),
		Description: fmt.Sprintf("MCP config exists and is valid JSON at %s for %s", mcpPath, agentID),
		Soft:        false, // hard check: if we wrote it, it must be there
		Run: func(_ context.Context) error {
			data, err := os.ReadFile(mcpPath)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("MCP config not found at %q", mcpPath)
				}
				return fmt.Errorf("read MCP config: %w", err)
			}
			if adapter.Agent() == model.AgentCodex {
				var root struct {
					MCPServers map[string]any `toml:"mcp_servers"`
				}
				if err := toml.Unmarshal(data, &root); err != nil {
					return fmt.Errorf("MCP config at %q is not valid TOML: %w", mcpPath, err)
				}
				if _, ok := root.MCPServers[serverName]; !ok {
					return fmt.Errorf("MCP server %q not found in %q", serverName, mcpPath)
				}
				return nil
			}
			var raw json.RawMessage
			if err := json.Unmarshal(data, &raw); err != nil {
				return fmt.Errorf("MCP config at %q is not valid JSON: %w", mcpPath, err)
			}
			return nil
		},
	}
}

func checkMCPEndpoint(h model.Harness) Check {
	return Check{
		ID:          fmt.Sprintf("external:mcp-endpoint:%s", h.ID),
		Description: fmt.Sprintf("MCP endpoint %s is reachable", h.External.URL),
		Soft:        true, // warning only — network may be down
		Run:         nil,  // skipped: network checks are out of scope for unit verify
	}
}

func checkBinaryInPATH(h model.Harness) Check {
	binaryName := binaryNameFromHarness(h)

	return Check{
		ID:          fmt.Sprintf("external:binary:%s", h.ID),
		Description: fmt.Sprintf("binary %q in PATH after install", binaryName),
		Soft:        true, // warning only — PATH may differ across shells
		Run: func(_ context.Context) error {
			_, err := exec.LookPath(binaryName)
			if err != nil {
				return fmt.Errorf("binary %q not found in PATH: %w", binaryName, err)
			}
			return nil
		},
	}
}

// binaryNameFromHarness derives the expected binary name from the harness.
func binaryNameFromHarness(h model.Harness) string {
	if h.External == nil || h.External.Pkg == "" {
		return h.ID
	}
	// Strip npm scope prefix: "@scope/name" → "name".
	pkg := h.External.Pkg
	if strings.HasPrefix(pkg, "@") {
		parts := strings.SplitN(pkg, "/", 2)
		if len(parts) == 2 {
			return filepath.Base(parts[1])
		}
	}
	return filepath.Base(pkg)
}
