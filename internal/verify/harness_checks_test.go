package verify_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/harness/config/permissions"
	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/verify"
)

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

// newFakeAdapter builds a fakeAdapter wired to a temp dir.
func newFakeAdapter(t *testing.T, homeDir string, agent model.Agent) fakeAdapter {
	t.Helper()
	return fakeAdapter{
		agent:            agent,
		skillsDir:        filepath.Join(homeDir, "skills"),
		instructionsPath: filepath.Join(homeDir, "CLAUDE.md"),
		settingsPath:     filepath.Join(homeDir, "settings.json"),
		mcpConfigPath:    filepath.Join(homeDir, "mcp", "server.json"),
	}
}

// runCheck is a small helper that runs the first check for a harness and
// returns the result's status/error.
func runChecks(t *testing.T, checks []verify.Check) []verify.CheckResult {
	t.Helper()
	return verify.RunChecks(context.Background(), checks)
}

// ─────────────────────────────────────────────────────────────────────────────
// 2.2 — HarnessType == skill
// ─────────────────────────────────────────────────────────────────────────────

func TestCheckForSkillHarness_Pass(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentClaude)

	h := model.Harness{
		ID:   "active-orchestrator",
		Type: model.HarnessSkill,
	}

	// Create the expected SKILL.md with content.
	skillDir := filepath.Join(fa.SkillsDir(homeDir), h.ID)
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# Skill content"), 0o644); err != nil {
		t.Fatal(err)
	}

	checks := verify.ChecksForHarness(h, []verify.Adapter{fa}, homeDir)
	if len(checks) == 0 {
		t.Fatal("expected at least one check for skill harness")
	}

	results := runChecks(t, checks)
	for _, r := range results {
		if r.Status != verify.CheckStatusPassed {
			t.Errorf("check %q: status = %q, error = %q, want passed", r.ID, r.Status, r.Error)
		}
	}
}

func TestCheckForSkillHarness_FailMissingSKILLMD(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentClaude)

	h := model.Harness{
		ID:   "active-orchestrator",
		Type: model.HarnessSkill,
	}

	// Do NOT create SKILL.md — should fail.
	checks := verify.ChecksForHarness(h, []verify.Adapter{fa}, homeDir)
	if len(checks) == 0 {
		t.Fatal("expected at least one check for skill harness")
	}

	results := runChecks(t, checks)
	anyFailed := false
	for _, r := range results {
		if r.Status == verify.CheckStatusFailed {
			anyFailed = true
		}
	}
	if !anyFailed {
		t.Error("expected at least one failed check when SKILL.md is missing")
	}
}

func TestCheckForSkillHarness_FailEmptySKILLMD(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentClaude)

	h := model.Harness{
		ID:   "active-orchestrator",
		Type: model.HarnessSkill,
	}

	// Create an empty SKILL.md — should fail (not useful as a skill).
	skillDir := filepath.Join(fa.SkillsDir(homeDir), h.ID)
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}

	checks := verify.ChecksForHarness(h, []verify.Adapter{fa}, homeDir)
	results := runChecks(t, checks)
	anyFailed := false
	for _, r := range results {
		if r.Status == verify.CheckStatusFailed {
			anyFailed = true
		}
	}
	if !anyFailed {
		t.Error("expected at least one failed check when SKILL.md is empty")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// 2.3 — HarnessType == config
// ─────────────────────────────────────────────────────────────────────────────

func TestCheckForConfigHarness_Pass(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentClaude)

	h := model.Harness{
		ID:   "sdd-orchestrator",
		Type: model.HarnessConfig,
	}

	// Write instructions file with the idempotent marker exactly once.
	marker := "<!-- active-stack:sdd-orchestrator -->"
	closeMarker := "<!-- /active-stack:sdd-orchestrator -->"
	content := "# Instructions\n\n" + marker + "\nsome content\n" + closeMarker + "\n"
	if err := os.WriteFile(fa.InstructionsPath(homeDir), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	checks := verify.ChecksForHarness(h, []verify.Adapter{fa}, homeDir)
	results := runChecks(t, checks)
	for _, r := range results {
		if r.Status != verify.CheckStatusPassed {
			t.Errorf("check %q: status = %q, error = %q, want passed", r.ID, r.Status, r.Error)
		}
	}
}

// TestCheckForConfigHarness_PrimaryAgent_Pass verifies that for an agent that
// delivers config as a primary agent (OpenCode), the health check looks in the
// settings JSON for agent.<id> with mode:primary — NOT the AGENTS.md marker,
// which is intentionally absent in this delivery mode.
func TestCheckForConfigHarness_PrimaryAgent_Pass(t *testing.T) {
	homeDir := t.TempDir()
	fa := fakeAdapter{
		agent:            model.AgentOpenCode,
		instructionsPath: filepath.Join(homeDir, "AGENTS.md"), // intentionally has no marker
		settingsPath:     filepath.Join(homeDir, "opencode.json"),
		delivery:         model.ConfigDeliveryPrimaryAgent,
	}

	h := model.Harness{ID: "sdd-orchestrator", Type: model.HarnessConfig}

	settings := `{"agent":{"sdd-orchestrator":{"mode":"primary","prompt":"x"},"build":{"mode":"primary"}}}`
	if err := os.WriteFile(fa.SettingsPath(homeDir), []byte(settings), 0o644); err != nil {
		t.Fatal(err)
	}

	checks := verify.ChecksForHarness(h, []verify.Adapter{fa}, homeDir)
	results := runChecks(t, checks)
	if len(results) == 0 {
		t.Fatal("expected at least one check")
	}
	for _, r := range results {
		if r.Status != verify.CheckStatusPassed {
			t.Errorf("check %q: status = %q, error = %q, want passed", r.ID, r.Status, r.Error)
		}
	}
}

// TestCheckForConfigHarness_PrimaryAgent_FailAbsent verifies the check fails
// when the primary agent entry is missing from the settings JSON.
func TestCheckForConfigHarness_PrimaryAgent_FailAbsent(t *testing.T) {
	homeDir := t.TempDir()
	fa := fakeAdapter{
		agent:            model.AgentOpenCode,
		instructionsPath: filepath.Join(homeDir, "AGENTS.md"),
		settingsPath:     filepath.Join(homeDir, "opencode.json"),
		delivery:         model.ConfigDeliveryPrimaryAgent,
	}

	h := model.Harness{ID: "sdd-orchestrator", Type: model.HarnessConfig}

	// Settings exist but only have an unrelated agent — orchestrator absent.
	if err := os.WriteFile(fa.SettingsPath(homeDir), []byte(`{"agent":{"build":{"mode":"primary"}}}`), 0o644); err != nil {
		t.Fatal(err)
	}

	checks := verify.ChecksForHarness(h, []verify.Adapter{fa}, homeDir)
	results := runChecks(t, checks)
	anyFailed := false
	for _, r := range results {
		if r.Status == verify.CheckStatusFailed {
			anyFailed = true
		}
	}
	if !anyFailed {
		t.Error("expected failure when primary-agent entry is absent")
	}
}

func TestCheckForConfigHarness_FailMarkerAbsent(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentClaude)

	h := model.Harness{
		ID:   "sdd-orchestrator",
		Type: model.HarnessConfig,
	}

	// Write instructions file WITHOUT the marker.
	if err := os.WriteFile(fa.InstructionsPath(homeDir), []byte("# Instructions\nno marker here\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	checks := verify.ChecksForHarness(h, []verify.Adapter{fa}, homeDir)
	results := runChecks(t, checks)
	anyFailed := false
	for _, r := range results {
		if r.Status == verify.CheckStatusFailed {
			anyFailed = true
		}
	}
	if !anyFailed {
		t.Error("expected failure when marker is absent")
	}
}

func TestCheckForConfigHarness_FailDuplicateMarker(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentClaude)

	h := model.Harness{
		ID:   "sdd-orchestrator",
		Type: model.HarnessConfig,
	}

	// Write instructions file with duplicate markers (idempotency violated).
	marker := "<!-- active-stack:sdd-orchestrator -->"
	closeMarker := "<!-- /active-stack:sdd-orchestrator -->"
	content := marker + "\ncontent\n" + closeMarker + "\n" +
		marker + "\ncontent again\n" + closeMarker + "\n"
	if err := os.WriteFile(fa.InstructionsPath(homeDir), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	checks := verify.ChecksForHarness(h, []verify.Adapter{fa}, homeDir)
	results := runChecks(t, checks)
	anyFailed := false
	for _, r := range results {
		if r.Status == verify.CheckStatusFailed {
			anyFailed = true
		}
	}
	if !anyFailed {
		t.Error("expected failure when marker appears more than once (idempotency violation)")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// 2.3 — special case: permissions harness → settings path
// ─────────────────────────────────────────────────────────────────────────────

func TestCheckForPermissionsHarness_Pass(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentClaude)

	h := model.Harness{
		ID:   "permissions",
		Type: model.HarnessConfig,
	}

	// Write a settings JSON as the installer writes it: the permissions key
	// AND the security floor. A fixture with only the key used to pass here,
	// which is exactly the hole C-32 closed — see the section at the bottom.
	settings := map[string]interface{}{
		"permissions": map[string]interface{}{
			"allow": []string{"Read"},
			"deny":  permissions.SecurityFloorDeny(),
		},
	}
	data, _ := json.Marshal(settings)
	if err := os.WriteFile(fa.SettingsPath(homeDir), data, 0o644); err != nil {
		t.Fatal(err)
	}

	checks := verify.ChecksForHarness(h, []verify.Adapter{fa}, homeDir)
	results := runChecks(t, checks)
	for _, r := range results {
		if r.Status != verify.CheckStatusPassed {
			t.Errorf("check %q: status = %q, error = %q, want passed", r.ID, r.Status, r.Error)
		}
	}
}

func TestCheckForPermissionsHarness_FailMissingSettingsFile(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentClaude)

	h := model.Harness{
		ID:   "permissions",
		Type: model.HarnessConfig,
	}

	// Settings file does not exist.
	checks := verify.ChecksForHarness(h, []verify.Adapter{fa}, homeDir)
	results := runChecks(t, checks)
	anyFailed := false
	for _, r := range results {
		if r.Status == verify.CheckStatusFailed {
			anyFailed = true
		}
	}
	if !anyFailed {
		t.Error("expected failure when settings file is missing")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// C-18 — per-agent permissions key resolution
// ─────────────────────────────────────────────────────────────────────────────

// TestCheckPermissionsOpenCodeSingularKey is the RED test for C-18.
// OpenCode writes the key "permission" (singular) per overlays.go L32-33.
// Before the fix, checkPermissions searches "permissions" (plural) and
// MUST fail with a "permissions key not found" style error.
// After the fix it MUST pass.
func TestCheckPermissionsOpenCodeSingularKey(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentOpenCode)

	h := model.Harness{
		ID:   "permissions",
		Type: model.HarnessConfig,
	}

	// Write an opencode.json-style settings file with "permission" (singular),
	// exactly as the permissions installer writes for OpenCode.
	// Raw JSON, not a map: opencode is last-match-wins, so the deny floor has
	// to come after the wildcard and Go maps do not preserve key order.
	settings := `{"permission":{"read":"allow","bash":{"*":"allow","rm -rf /*":"deny"}}}`
	if err := os.WriteFile(fa.SettingsPath(homeDir), []byte(settings), 0o644); err != nil {
		t.Fatal(err)
	}

	checks := verify.ChecksForHarness(h, []verify.Adapter{fa}, homeDir)
	if len(checks) == 0 {
		t.Fatal("expected at least one check for permissions harness")
	}

	results := runChecks(t, checks)
	for _, r := range results {
		if r.Status != verify.CheckStatusPassed {
			t.Errorf("check %q: status = %q, error = %q, want passed (opencode uses singular \"permission\" key)", r.ID, r.Status, r.Error)
		}
	}
}

// TestCheckPermissionsClaudePluralKey guards against regressions on Claude:
// Claude writes "permissions" (plural); the fix must not break it.
func TestCheckPermissionsClaudePluralKey(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentClaude)

	h := model.Harness{
		ID:   "permissions",
		Type: model.HarnessConfig,
	}

	settings := map[string]interface{}{
		"permissions": map[string]interface{}{
			"defaultMode": "acceptEdits",
			"deny":        permissions.SecurityFloorDeny(),
		},
	}
	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fa.SettingsPath(homeDir), data, 0o644); err != nil {
		t.Fatal(err)
	}

	checks := verify.ChecksForHarness(h, []verify.Adapter{fa}, homeDir)
	results := runChecks(t, checks)
	for _, r := range results {
		if r.Status != verify.CheckStatusPassed {
			t.Errorf("check %q: status = %q, error = %q, want passed (claude uses plural \"permissions\" key)", r.ID, r.Status, r.Error)
		}
	}
}

// TestCheckPermissionsFailureMentionsCorrectKey asserts that when the
// agent-specific key is absent the error message names the searched key,
// not a hardcoded string. This covers D4 (self-describing failure message).
func TestCheckPermissionsFailureMentionsCorrectKey(t *testing.T) {
	tests := []struct {
		name         string
		agent        model.Agent
		settingsJSON string
		wantKeyInMsg string
	}{
		{
			name:         "opencode missing singular key",
			agent:        model.AgentOpenCode,
			settingsJSON: `{"someOtherKey": true}`,
			wantKeyInMsg: "permission",
		},
		{
			name:         "claude missing plural key",
			agent:        model.AgentClaude,
			settingsJSON: `{"someOtherKey": true}`,
			wantKeyInMsg: "permissions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir := t.TempDir()
			fa := newFakeAdapter(t, homeDir, tt.agent)

			h := model.Harness{
				ID:   "permissions",
				Type: model.HarnessConfig,
			}

			if err := os.WriteFile(fa.SettingsPath(homeDir), []byte(tt.settingsJSON), 0o644); err != nil {
				t.Fatal(err)
			}

			checks := verify.ChecksForHarness(h, []verify.Adapter{fa}, homeDir)
			results := runChecks(t, checks)

			anyFailed := false
			for _, r := range results {
				if r.Status == verify.CheckStatusFailed {
					anyFailed = true
					if r.Error == "" {
						t.Errorf("check %q: failed but Error is empty", r.ID)
						continue
					}
					if !containsKey(r.Error, tt.wantKeyInMsg) {
						t.Errorf("check %q: error = %q, want it to contain key name %q", r.ID, r.Error, tt.wantKeyInMsg)
					}
				}
			}
			if !anyFailed {
				t.Errorf("expected at least one failed check when key %q is absent from settings", tt.wantKeyInMsg)
			}
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// C-19: best-effort harness — tasks 4.1, 4.2
// ─────────────────────────────────────────────────────────────────────────────

// Task 4.1a — checkSkill for a best-effort harness emits a Check with Soft == true.
func TestCheckForSkillHarness_BestEffort_CheckIsSoft(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentClaude)

	h := model.Harness{
		ID:         "find-skill",
		Type:       model.HarnessSkill,
		BestEffort: true,
	}

	checks := verify.ChecksForHarness(h, []verify.Adapter{fa}, homeDir)
	if len(checks) == 0 {
		t.Fatal("expected at least one check for best-effort skill harness")
	}

	for _, c := range checks {
		if !c.Soft {
			t.Errorf("check %q: Soft = false, want true for best-effort harness", c.ID)
		}
	}
}

// Task 4.1b — checkSkill for a NON-best-effort harness emits a Check with Soft == false (regression).
func TestCheckForSkillHarness_NonBestEffort_CheckIsHard(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentClaude)

	h := model.Harness{
		ID:         "active-orchestrator",
		Type:       model.HarnessSkill,
		BestEffort: false,
	}

	checks := verify.ChecksForHarness(h, []verify.Adapter{fa}, homeDir)
	if len(checks) == 0 {
		t.Fatal("expected at least one check for non-best-effort skill harness")
	}

	for _, c := range checks {
		if c.Soft {
			t.Errorf("check %q: Soft = true, want false for non-best-effort harness", c.ID)
		}
	}
}

// Task 4.2 — a best-effort skill check that fails (SKILL.md missing) produces a
// warning (CheckStatusWarning) — not a hard failure — and the report stays Ready.
func TestCheckForSkillHarness_BestEffort_FailIsWarningNotFailure(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentClaude)

	h := model.Harness{
		ID:         "find-skill",
		Type:       model.HarnessSkill,
		BestEffort: true,
	}

	// Intentionally do NOT create SKILL.md — the check should warn, not fail.
	checks := verify.ChecksForHarness(h, []verify.Adapter{fa}, homeDir)
	results := runChecks(t, checks)

	for _, r := range results {
		if r.Status == verify.CheckStatusFailed {
			t.Errorf("check %q: status = Failed, want Warning for best-effort harness", r.ID)
		}
	}

	// The report must still be Ready (warnings don't flip Ready).
	report := verify.BuildReport(results)
	if !report.Ready {
		t.Errorf("report.Ready = false, want true when only best-effort skill check fails (should be warning)")
	}
}

// containsKey reports whether s contains the quoted or unquoted form of key.
func containsKey(s, key string) bool {
	return containsStr(s, `"`+key+`"`) || containsStr(s, key)
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && strings.Contains(s, substr))
}

// ─────────────────────────────────────────────────────────────────────────────
// 2.4 — HarnessType == external
// ─────────────────────────────────────────────────────────────────────────────

func TestCheckForExternalHarness_MCP_Pass(t *testing.T) {
	homeDir := t.TempDir()
	mcpDir := filepath.Join(homeDir, "mcp")
	if err := os.MkdirAll(mcpDir, 0o755); err != nil {
		t.Fatal(err)
	}

	fa := fakeAdapter{
		agent:         model.AgentClaude,
		mcpConfigPath: filepath.Join(mcpDir, "server.json"),
	}

	h := model.Harness{
		ID:       "context7",
		Type:     model.HarnessExternal,
		External: &model.External{Method: "mcp", URL: "https://mcp.context7.com"},
	}

	// Write a valid JSON MCP config.
	mcpConfig := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"context7": map[string]interface{}{
				"url": "https://mcp.context7.com",
			},
		},
	}
	data, _ := json.Marshal(mcpConfig)
	if err := os.WriteFile(fa.MCPConfigPath(homeDir, h.ID), data, 0o644); err != nil {
		t.Fatal(err)
	}

	checks := verify.ChecksForHarness(h, []verify.Adapter{fa}, homeDir)
	// MCP config check is hard; binary check is Soft.
	hasSoftCheck := false
	for _, c := range checks {
		if c.Soft {
			hasSoftCheck = true
		}
	}
	if !hasSoftCheck {
		t.Error("expected at least one Soft check (binary/endpoint) for external harness")
	}

	results := runChecks(t, checks)
	report := verify.BuildReport(results)
	// MCP file exists and is valid JSON, so the hard check passes.
	// The soft check may warn (binary not in PATH in test env) — that's OK.
	if report.Failed > 0 {
		t.Errorf("expected no hard failures for valid MCP config, got %d failures: %+v", report.Failed, results)
	}
}

func TestCheckForExternalHarness_MCP_FailInvalidJSON(t *testing.T) {
	homeDir := t.TempDir()
	mcpDir := filepath.Join(homeDir, "mcp")
	if err := os.MkdirAll(mcpDir, 0o755); err != nil {
		t.Fatal(err)
	}

	fa := fakeAdapter{
		agent:         model.AgentClaude,
		mcpConfigPath: filepath.Join(mcpDir, "server.json"),
	}

	h := model.Harness{
		ID:       "context7",
		Type:     model.HarnessExternal,
		External: &model.External{Method: "mcp", URL: "https://mcp.context7.com"},
	}

	// Write invalid JSON.
	if err := os.WriteFile(fa.MCPConfigPath(homeDir, h.ID), []byte("{invalid json"), 0o644); err != nil {
		t.Fatal(err)
	}

	checks := verify.ChecksForHarness(h, []verify.Adapter{fa}, homeDir)
	results := runChecks(t, checks)

	anyFailed := false
	for _, r := range results {
		if r.Status == verify.CheckStatusFailed {
			anyFailed = true
		}
	}
	if !anyFailed {
		t.Error("expected hard failure when MCP config JSON is invalid")
	}
}

func TestCheckForExternalHarness_SoftCheckProducesWarningNotFailure(t *testing.T) {
	homeDir := t.TempDir()

	fa := fakeAdapter{
		agent:         model.AgentClaude,
		mcpConfigPath: filepath.Join(homeDir, "mcp", "engram.json"),
	}

	// Engram uses homebrew, not mcp — its binary check is Soft.
	h := model.Harness{
		ID:       "engram",
		Type:     model.HarnessExternal,
		External: &model.External{Method: "homebrew", Pkg: "engram"},
	}

	checks := verify.ChecksForHarness(h, []verify.Adapter{fa}, homeDir)

	// All checks for a homebrew harness should be Soft (binary in PATH).
	for _, c := range checks {
		if !c.Soft {
			t.Errorf("check %q should be Soft for homebrew external harness", c.ID)
		}
	}

	results := runChecks(t, checks)
	report := verify.BuildReport(results)

	// The binary is not installed in test env, so we get warnings — but NOT failures.
	if report.Failed > 0 {
		t.Errorf("Soft checks must not produce failures: got %d failures", report.Failed)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// 2.2 — Multiple adapters: one check per (harness × agent)
// ─────────────────────────────────────────────────────────────────────────────

func TestCheckForSkillHarness_MultipleAdapters(t *testing.T) {
	homeDir := t.TempDir()

	// Two adapters for two agents.
	claudeSkillsDir := filepath.Join(homeDir, "claude", "skills")
	opencodeSkillsDir := filepath.Join(homeDir, "opencode", "skills")

	claudeAdapter := fakeAdapter{
		agent:     model.AgentClaude,
		skillsDir: claudeSkillsDir,
	}
	opencodeAdapter := fakeAdapter{
		agent:     model.AgentOpenCode,
		skillsDir: opencodeSkillsDir,
	}

	h := model.Harness{
		ID:   "active-orchestrator",
		Type: model.HarnessSkill,
	}

	// Create SKILL.md for claude only.
	skillDir := filepath.Join(claudeSkillsDir, h.ID)
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# content"), 0o644); err != nil {
		t.Fatal(err)
	}
	// opencode's SKILL.md is missing.

	checks := verify.ChecksForHarness(h, []verify.Adapter{claudeAdapter, opencodeAdapter}, homeDir)

	// Should have checks for both agents.
	if len(checks) < 2 {
		t.Fatalf("expected at least 2 checks (one per adapter), got %d", len(checks))
	}

	results := runChecks(t, checks)
	report := verify.BuildReport(results)

	// Claude passes, opencode fails → report not ready.
	if report.Ready {
		t.Error("Ready = true, want false when opencode SKILL.md is missing")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// C-32: el piso de seguridad tiene que estar APLICADO, no sólo nombrado
//
// Hasta el 2026-09-08 checkPermissions parseaba el settings a
// map[string]json.RawMessage y sólo miraba que la clave existiera. Un
// settings.json de Claude con "permissions": {"deny": []} o un opencode.json
// con "permission": {"bash": {"*": "allow"}} pasaban el verify en verde con el
// piso de denegación vaciado. El harness permissions es el único que el
// instalador declara NO desactivable: si el verify no mira el contenido, la
// única garantía no negociable del stack no está verificada por nadie.
// ─────────────────────────────────────────────────────────────────────────────

// permissionsHarness is the harness under test in this section.
func permissionsHarnessFixture() model.Harness {
	return model.Harness{ID: "permissions", Type: model.HarnessConfig}
}

// writeSettings writes raw JSON to the adapter's settings path.
func writeSettings(t *testing.T, fa fakeAdapter, homeDir, body string) {
	t.Helper()
	if err := os.WriteFile(fa.SettingsPath(homeDir), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// assertPermissionsCheckFails runs the permissions checks and requires at least
// one hard failure whose message contains wantInMsg.
func assertPermissionsCheckFails(t *testing.T, fa fakeAdapter, homeDir, wantInMsg string) {
	t.Helper()
	checks := verify.ChecksForHarness(permissionsHarnessFixture(), []verify.Adapter{fa}, homeDir)
	if len(checks) == 0 {
		t.Fatal("expected at least one check for permissions harness")
	}
	results := runChecks(t, checks)
	for _, r := range results {
		if r.Status == verify.CheckStatusFailed {
			if !strings.Contains(r.Error, wantInMsg) {
				t.Errorf("check %q failed with error = %q, want it to mention %q", r.ID, r.Error, wantInMsg)
			}
			return
		}
	}
	t.Fatalf("expected the permissions check to FAIL: the security floor is not applied, but verify reported %+v", results)
}

// assertPermissionsCheckPasses runs the permissions checks and requires all to pass.
func assertPermissionsCheckPasses(t *testing.T, fa fakeAdapter, homeDir string) {
	t.Helper()
	checks := verify.ChecksForHarness(permissionsHarnessFixture(), []verify.Adapter{fa}, homeDir)
	if len(checks) == 0 {
		t.Fatal("expected at least one check for permissions harness")
	}
	for _, r := range runChecks(t, checks) {
		if r.Status != verify.CheckStatusPassed {
			t.Errorf("check %q: status = %q, error = %q, want passed", r.ID, r.Status, r.Error)
		}
	}
}

// claudeSettingsWithDeny renders a Claude settings.json with the given deny list.
func claudeSettingsWithDeny(t *testing.T, deny []string) string {
	t.Helper()
	body, err := json.Marshal(map[string]any{
		"permissions": map[string]any{
			"defaultMode": "default",
			"deny":        deny,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

// Claude — el piso vaciado a mano tiene que dar rojo.
func TestCheckPermissionsClaudeRejectsEmptiedDenyFloor(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentClaude)
	writeSettings(t, fa, homeDir, claudeSettingsWithDeny(t, []string{}))
	assertPermissionsCheckFails(t, fa, homeDir, "Read(.env)")
}

// Claude — el piso ausente por completo (clave "deny" no escrita) tiene que dar rojo.
func TestCheckPermissionsClaudeRejectsMissingDenyKey(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentClaude)
	writeSettings(t, fa, homeDir, `{"permissions":{"defaultMode":"bypassPermissions","allow":["Bash"]}}`)
	assertPermissionsCheckFails(t, fa, homeDir, "deny")
}

// Claude — que le saquen UNA regla del piso también tiene que dar rojo,
// y el error tiene que nombrar exactamente cuál falta.
func TestCheckPermissionsClaudeRejectsPartialDenyFloor(t *testing.T) {
	full := permissions.SecurityFloorDeny()
	for i, dropped := range full {
		t.Run(dropped, func(t *testing.T) {
			partial := make([]string, 0, len(full)-1)
			partial = append(partial, full[:i]...)
			partial = append(partial, full[i+1:]...)

			homeDir := t.TempDir()
			fa := newFakeAdapter(t, homeDir, model.AgentClaude)
			writeSettings(t, fa, homeDir, claudeSettingsWithDeny(t, partial))
			assertPermissionsCheckFails(t, fa, homeDir, dropped)
		})
	}
}

// Claude — el piso completo (lo que escribe el instalador) tiene que dar verde,
// incluso con reglas extra del usuario (el piso es un mínimo, no un igual).
func TestCheckPermissionsClaudeAcceptsFloorPlusExtras(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentClaude)
	deny := append(permissions.SecurityFloorDeny(), "Bash(curl:*)")
	writeSettings(t, fa, homeDir, claudeSettingsWithDeny(t, deny))
	assertPermissionsCheckPasses(t, fa, homeDir)
}

// OpenCode — bash sin ninguna regla del piso tiene que dar rojo.
func TestCheckPermissionsOpenCodeRejectsEmptiedBashFloor(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentOpenCode)
	writeSettings(t, fa, homeDir, `{"permission":{"bash":{"*":"allow"},"read":"allow"}}`)
	assertPermissionsCheckFails(t, fa, homeDir, "rm -rf /*")
}

// OpenCode — el piso presente pero SOMBREADO por un comodín posterior.
// opencode evalúa last-match-wins: un "*": "allow" después del deny lo anula.
// Sintácticamente la regla está; semánticamente no aplica. Tiene que dar rojo.
func TestCheckPermissionsOpenCodeRejectsFloorShadowedByLaterWildcard(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentOpenCode)
	writeSettings(t, fa, homeDir, `{"permission":{"bash":{"rm -rf /*":"deny","*":"allow"}}}`)
	assertPermissionsCheckFails(t, fa, homeDir, "rm -rf /*")
}

// OpenCode — el piso en el orden correcto (comodín primero, deny último)
// tiene que dar verde.
func TestCheckPermissionsOpenCodeAcceptsDenyLast(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentOpenCode)
	writeSettings(t, fa, homeDir, `{"permission":{"bash":{"*":"ask","go test *":"allow","rm -rf /*":"deny"}}}`)
	assertPermissionsCheckPasses(t, fa, homeDir)
}

// El contrato de verdad: lo que el instalador escribe, en cualquier tier y para
// cualquier agente, tiene que pasar el verify. Ata el check al instalador real
// en vez de a un fixture escrito a mano que puede derivar.
func TestCheckPermissionsAcceptsRealInstallerOutput(t *testing.T) {
	agents := []model.Agent{model.AgentClaude, model.AgentOpenCode}
	tiers := []model.PermissionTier{model.TierEstricto, model.TierBalanceado, model.TierBypass}

	for _, agent := range agents {
		for _, tier := range tiers {
			t.Run(string(agent)+"/"+string(tier), func(t *testing.T) {
				homeDir := t.TempDir()
				fa := newFakeAdapter(t, homeDir, agent)

				if _, err := permissions.Install(homeDir, []permissions.PermissionsAdapter{fa}, tier); err != nil {
					t.Fatalf("permissions.Install: %v", err)
				}
				assertPermissionsCheckPasses(t, fa, homeDir)
			})
		}
	}
}

// OpenCode — el sombreado PARCIAL. Una regla posterior mas angosta que se
// solapa con el piso tambien lo anula, y es la forma mas probable de romperlo:
// alguien agregando una linea al final del archivo a mano. Bajo last-match-wins
// `rm -rf /home/proyectos` matchea el deny del piso Y el allow posterior, y
// gana el posterior: el piso deja de aplicar para todo /home.
func TestCheckPermissionsOpenCodeRejectsFloorShadowedByNarrowerLaterRule(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentOpenCode)
	writeSettings(t, fa, homeDir, `{"permission":{"bash":{"*":"ask","rm -rf /*":"deny","rm -rf /home/*":"allow"}}}`)
	assertPermissionsCheckFails(t, fa, homeDir, "rm -rf /home/*")
}

// OpenCode — el sombreado en la direccion contraria: el piso es angosto y la
// regla posterior es mas ancha. Tambien lo anula.
func TestCheckPermissionsOpenCodeRejectsFloorShadowedByBroaderLaterRule(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentOpenCode)
	writeSettings(t, fa, homeDir, `{"permission":{"bash":{"*":"ask","rm -rf /home/*":"deny","rm -rf /*":"allow"}}}`)
	// El piso real de opencode se sigue exigiendo aparte, asi que este fixture
	// falla igual por "rm -rf /*" ausente. Lo que importa aca es que tambien
	// exista el caso simetrico en el test del predicado (glob_internal_test.go).
	assertPermissionsCheckFails(t, fa, homeDir, "rm -rf /*")
}

// OpenCode — una regla posterior que NO se solapa con el piso tiene que dar
// VERDE. El check no puede prohibirle al usuario agregar reglas propias: solo
// las que anulan el piso. Sin este test, un predicado que devuelva siempre true
// pasaria por bueno.
func TestCheckPermissionsOpenCodeAcceptsUnrelatedLaterRule(t *testing.T) {
	homeDir := t.TempDir()
	fa := newFakeAdapter(t, homeDir, model.AgentOpenCode)
	writeSettings(t, fa, homeDir, `{"permission":{"bash":{"*":"ask","rm -rf /*":"deny","git status":"allow"}}}`)
	assertPermissionsCheckPasses(t, fa, homeDir)
}
