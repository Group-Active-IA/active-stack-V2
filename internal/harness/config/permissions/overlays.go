package permissions

import (
	"encoding/json"

	"github.com/Group-Active-IA/active-stack/internal/model"
)

// securityFloorDeny is the single source of truth for the C-21 deny rules that
// CLAUDE CODE receives in all three of its tiers, bypass included. Claude
// evaluates deny → ask → allow, so deny ALWAYS wins regardless of defaultMode.
// This slice must NEVER be modified by tier composition — it is what makes the
// security guarantee real for Claude, and TestSecurityFloorInAllClaudeTiers
// reads the composed overlay back to keep that true.
//
// SCOPE, corrected 2026-09-08: this comment used to say the floor "applies to
// ALL tiers (including bypass)" with no subject, and it read as a promise about
// the whole installer. It is not one. opencode has a different floor of its own
// (opencodeBashFloorDeny). codex has NO floor, in any tier — see
// mergeCodexTOMLFile in install.go for what was looked for and why nothing was
// invented. gemini and vscode carry none either. An agent has the floor its own
// installer writes, and nothing more.
//
// The five rules are the normative contract tested in permissions_test.go.
// The implementation is a superset (extra rules are allowed per design.md §D2).
var securityFloorDeny = []string{
	// ─── Normative floor (must-have, tested) ───
	"Read(.env)",
	"Read(.env.*)",
	"Bash(rm -rf /)",
	"Bash(rm -rf ~)",
	"Bash(git push --force:*)",
	// ─── Superset extras (defense in depth) ───
	"Bash(sudo rm -rf /)",
	"Bash(sudo rm -rf ~)",
	"Edit(.env)",
	"Edit(.env.*)",
}

// SecurityFloorDeny returns a copy of the security floor deny rules.
// Exported for testing; the internal slice must not be mutated.
func SecurityFloorDeny() []string {
	out := make([]string, len(securityFloorDeny))
	copy(out, securityFloorDeny)
	return out
}

// opencodeBashFloorDeny is the opencode counterpart of securityFloorDeny: the
// bash patterns that composeOpencodeOverlay denies in EVERY tier (estricto,
// balanceado and bypass). opencode has no cross-tool deny-list — its floor is
// expressed as pattern-object entries inside "permission.bash" — so the two
// floors are different sets on purpose. Do NOT copy Claude's rules here: an
// agent only has the floor its installer actually writes.
//
// Kept honest by TestOpencodeBashFloorIsDeniedInEveryTier, which reads back the
// composed overlay for each tier instead of trusting this list.
var opencodeBashFloorDeny = []string{
	"rm -rf /*",
}

// OpencodeBashFloorDeny returns a copy of the opencode bash deny floor.
// Exported so internal/verify can check the installed opencode.json against the
// same list the installer writes, instead of hardcoding a second copy that can
// drift away from this one.
func OpencodeBashFloorDeny() []string {
	out := make([]string, len(opencodeBashFloorDeny))
	copy(out, opencodeBashFloorDeny)
	return out
}

// ── Claude Code overlay composition ─────────────────────────────────────────

// claudePermissionsShape is the struct that serializes to the Claude Code
// "permissions" block in settings.json.
//
// Allow deliberately has NO omitempty. MergeJSONObjects only adds or overwrites
// keys — it never deletes one — so a key the overlay omits survives from
// whatever was installed before. On 2026-09-08 that meant installing balanceado
// and then estricto left balanceado's allow-list alive: the user hardened their
// config and kept auto-approving Read, Edit and Bash(go test:*). A tier is only
// a tier if it declares every key it owns, including the ones it wants empty.
type claudePermissionsShape struct {
	DefaultMode string   `json:"defaultMode"`
	Allow       []string `json:"allow"`
	Deny        []string `json:"deny"`
}

// claudeSettingsShape is the top-level settings.json for Claude Code.
type claudeSettingsShape struct {
	Permissions claudePermissionsShape `json:"permissions"`
}

// claudeAllowList is the curated allow-list for tier balanceado.
// Kept as a package-level var so it can be referenced in tests if needed.
var claudeAllowList = []string{
	"Read",
	"Edit",
	"Bash(go test:*)",
	"Bash(go build:*)",
	"Bash(git status)",
	"Bash(git diff:*)",
}

// composeClaudeOverlay returns the JSON overlay bytes for Claude Code given a
// permission tier. Zero-value tier normalizes to TierBalanceado (never TierBypass).
//
// Claude semantics: deny → ask → allow; deny ALWAYS wins.
// The security floor is present in every tier.
func composeClaudeOverlay(tier model.PermissionTier) []byte {
	tier = tier.Normalize()

	var shape claudeSettingsShape
	shape.Permissions.Deny = SecurityFloorDeny()

	switch tier {
	case model.TierEstricto:
		shape.Permissions.DefaultMode = "default"
		// Empty, not absent: an explicit [] overwrites the previous tier's list.
		shape.Permissions.Allow = []string{}
	case model.TierBypass:
		shape.Permissions.DefaultMode = "bypassPermissions"
		// bypassPermissions auto-approves everything not denied, so the list is
		// redundant — but it still has to be emitted to clear the old one.
		shape.Permissions.Allow = []string{}
	default: // TierBalanceado (and any unknown after normalize)
		shape.Permissions.DefaultMode = "default"
		shape.Permissions.Allow = claudeAllowList
	}

	b, err := json.Marshal(shape)
	if err != nil {
		// This should never happen with a static struct — panic so it surfaces in tests.
		panic("permissions: composeClaudeOverlay: marshal failed: " + err.Error())
	}
	return b
}

// ── OpenCode overlay composition (LAST-MATCH-WINS) ─────────────────────────
//
// CRITICAL: opencode evaluates pattern-objects with LAST-MATCH-WINS semantics —
// the OPPOSITE of Claude Code. If a deny rule appears BEFORE a wildcard "*", the
// wildcard overwrites it and the deny is silently ignored. The deny floor MUST be
// serialized LAST in every bash pattern-object.
//
// Because Go maps do NOT preserve insertion order, we use orderedPairs (a slice
// of key-value structs) that serializes to a JSON object in insertion order via a
// custom MarshalJSON. This is the ONLY safe way to guarantee deny-last.

// orderedPair is a single key-value entry in an ordered JSON object.
type orderedPair struct {
	Key   string
	Value string
}

// orderedBashBlock is an ordered collection of bash pattern-object entries.
// It serializes to a JSON object preserving insertion order (deny-last guarantee).
type orderedBashBlock []orderedPair

// MarshalJSON serializes the bash block as a JSON object in insertion order.
func (b orderedBashBlock) MarshalJSON() ([]byte, error) {
	buf := []byte{'{'}
	for i, p := range b {
		key, err := json.Marshal(p.Key)
		if err != nil {
			return nil, err
		}
		val, err := json.Marshal(p.Value)
		if err != nil {
			return nil, err
		}
		buf = append(buf, key...)
		buf = append(buf, ':')
		buf = append(buf, val...)
		if i < len(b)-1 {
			buf = append(buf, ',')
		}
	}
	buf = append(buf, '}')
	return buf, nil
}

// opencodePermissionShape is the top-level opencode "permission" object.
// Fields use omitempty because the zero value is the empty string, which is not
// a legal opencode verb — writing "external_directory": "" would be worse than
// omitting it.
//
// The catch, found on 2026-09-08: omitempty plus a merge that never deletes keys
// means a key only ONE tier sets stays on disk after switching to a tier that
// does not. external_directory was set to "deny" by estricto alone, so going
// estricto → balanceado kept denying external directories forever. Every tier
// must therefore set EVERY field of this struct explicitly, and
// TestOpencodeTiersDeclareTheSameKeys is what keeps that true when a field is
// added.
type opencodePermissionShape struct {
	Read              string           `json:"read,omitempty"`
	Edit              string           `json:"edit,omitempty"`
	Glob              string           `json:"glob,omitempty"`
	Grep              string           `json:"grep,omitempty"`
	List              string           `json:"list,omitempty"`
	Bash              orderedBashBlock `json:"bash"`
	Task              string           `json:"task,omitempty"`
	ExternalDirectory string           `json:"external_directory,omitempty"`
	Webfetch          string           `json:"webfetch,omitempty"`
	Websearch         string           `json:"websearch,omitempty"`
	LSP               string           `json:"lsp,omitempty"`
}

// composeOpencodeOverlay returns the JSON overlay bytes for opencode given a
// permission tier. Zero-value tier normalizes to TierBalanceado (never TierBypass).
//
// LAST-MATCH-WINS: deny rules MUST be the last entries in the bash pattern-object.
// orderedBashBlock guarantees serialization order (deny-last).
//
// REPLACE SCOPE, 2026-10-23: only "permission.bash" is wrapped in the
// __replace__ sentinel (same mechanism already used for MCP server entries —
// see buildStdioOverlay in internal/harness/external/mcp.go). The OTHER keys
// of "permission" (read, edit, external_directory, ...) stay as plain
// overwrites, which was already correct: every tier declares every scalar
// field explicitly (TestOpencodeTiersDeclareTheSameKeys), so a plain merge
// already discards the previous tier's scalar value.
//
// Only "bash" needed __replace__, because it is the one key whose value is
// itself a pattern-object: MergeJSONObjects deep-merges nested objects key by
// key, so a per-command override left by a previous tier (e.g. balanceado's
// "go test *": "allow") survived a switch to a stricter tier inside that
// nested object, where a plain scalar overwrite does not apply. Until
// 2026-10-23 that survival was accepted as a documented, provably-inert
// residue (see git history of this file and
// TestOpencodeTierSwitchFullyReplacesBash, which replaced the test that used
// to assert the residue on purpose).
//
// REJECTED ALTERNATIVE: wrapping the WHOLE "permission" object in
// __replace__ (done on another branch, commit ddeb94f) fixes the same bug but
// also destroys any sibling key inside "permission" that no tier here manages
// — a field a future opencode version adds, or one a human added by hand.
// Scoping the sentinel to "bash" alone keeps those siblings alive; see
// TestOpencodeTierSwitchPreservesPersonalSiblingKeyInPermission.
//
// ACCEPTED LOSS: a personal pattern added by hand INSIDE permission.bash
// itself (not a sibling of bash, but an entry of the object bash replaces)
// does NOT survive a reinstall — see
// TestOpencodeTierSwitchReplacesPersonalBashRule. The propose accepts this
// ("perdida deliberada de bash personal exige diff/aprobación y snapshot")
// and leaves the diff/approval gate as TBD (J6); this function does not
// implement that gate — see the apply report's deuda nueva.
func composeOpencodeOverlay(tier model.PermissionTier) []byte {
	tier = tier.Normalize()

	var shape opencodePermissionShape

	switch tier {
	case model.TierEstricto:
		// Everything → ask; external_directory → deny for extra hardening.
		// bash: wildcard ask first, then deny floor LAST (last-wins guarantee).
		shape = opencodePermissionShape{
			Read:              "ask",
			Edit:              "ask",
			Glob:              "ask",
			Grep:              "ask",
			List:              "ask",
			ExternalDirectory: "deny",
			Webfetch:          "ask",
			Websearch:         "ask",
			LSP:               "ask",
			Bash: orderedBashBlock{
				{Key: "*", Value: "ask"},
				// Deny floor — MUST come after the wildcard (last-wins).
				{Key: "rm -rf *", Value: "deny"},
				{Key: "rm -rf /*", Value: "deny"},
			},
		}

	case model.TierBypass:
		// Everything → allow; only the most catastrophic commands denied LAST.
		shape = opencodePermissionShape{
			Read:              "allow",
			Edit:              "allow",
			Glob:              "allow",
			Grep:              "allow",
			List:              "allow",
			ExternalDirectory: "allow",
			Webfetch:          "allow",
			Websearch:         "allow",
			LSP:               "allow",
			Bash: orderedBashBlock{
				{Key: "*", Value: "allow"},
				// Catastrophic deny floor LAST (last-wins guarantee).
				{Key: "rm -rf /*", Value: "deny"},
				{Key: "rm -rf ~/", Value: "deny"},
			},
		}

	default: // TierBalanceado (and zero-value after normalize)
		// Safe operations auto-allowed; deny floor LAST in bash.
		shape = opencodePermissionShape{
			Read:              "allow",
			Edit:              "ask",
			Glob:              "allow",
			Grep:              "allow",
			List:              "allow",
			ExternalDirectory: "ask",
			Webfetch:          "ask",
			Websearch:         "ask",
			LSP:               "allow",
			Bash: orderedBashBlock{
				{Key: "*", Value: "ask"},
				// Safe operations allowed.
				{Key: "go test *", Value: "allow"},
				{Key: "go build *", Value: "allow"},
				{Key: "git status", Value: "allow"},
				{Key: "git diff *", Value: "allow"},
				// Deny floor MUST be last (last-wins guarantee).
				{Key: "rm -rf *", Value: "deny"},
				{Key: "rm -rf /*", Value: "deny"},
			},
		}
	}

	// Every scalar field is a plain value: MergeJSONObjects overwrites it
	// outright because every tier that sets it sets it on every tier (no
	// omitempty gap for Read/Edit/Glob/Grep/List/ExternalDirectory/Webfetch/
	// Websearch/LSP). "bash" is the one field wrapped in
	// {"__replace__": ...} so the merge discards the previous tier's
	// pattern-object instead of merging into it; filemerge strips the
	// sentinel key before anything reaches disk (see
	// TestMergeJSONObjectsReplaceSentinelNotInOutput), so no "__replace__"
	// literal ever lands in opencode.json.
	//
	// "task" is omitted when empty, same as the struct's original
	// `json:"task,omitempty"` tag: NO tier sets it today, and "" is not a
	// legal opencode verb (see the comment on opencodePermissionShape) — it
	// would be worse than not writing the key at all. Found the hard way:
	// TestOpencodeNeverWritesEmptyPermissionVerb caught "task": "" landing in
	// a real opencode.json during the "Primer uso real" check for this fix,
	// after the switch from the struct (omitempty) to this map literal
	// (no omitempty) dropped that guarantee for the one field that relied on
	// it.
	permission := map[string]any{
		"read":               shape.Read,
		"edit":               shape.Edit,
		"glob":               shape.Glob,
		"grep":               shape.Grep,
		"list":               shape.List,
		"external_directory": shape.ExternalDirectory,
		"webfetch":           shape.Webfetch,
		"websearch":          shape.Websearch,
		"lsp":                shape.LSP,
		"bash": map[string]any{
			"__replace__": shape.Bash,
		},
	}
	if shape.Task != "" {
		permission["task"] = shape.Task
	}

	overlay := map[string]any{"permission": permission}

	b, err := json.Marshal(overlay)
	if err != nil {
		panic("permissions: composeOpencodeOverlay: marshal failed: " + err.Error())
	}
	return b
}

// ── Static overlays (gemini, vscode — tier-agnostic, TBD per design.md §D5) ──

// geminiCLIOverlayJSON sets Gemini CLI to "auto_edit" mode.
// Gemini CLI does not expose a deny-list equivalent in its schema;
// this is a known limitation. Tier mapeo: (TBD) — see design.md §D5.
var geminiCLIOverlayJSON = []byte(`{
  "general": {
    "defaultApprovalMode": "auto_edit"
  }
}
`)

// vscodeCopilotOverlayJSON enables auto-approve for VS Code Copilot chat tools.
// VS Code Copilot's settings.json does not support deny-list granularity.
// Tier mapeo: (TBD) — see design.md §D5.
var vscodeCopilotOverlayJSON = []byte(`{
  "chat.tools.autoApprove": true
}
`)

// agentOverlay returns the permission overlay bytes for the given agent and tier.
// Returns nil for agents with no injectable settings file (explicit no-op).
//
// Only claude and opencode differentiate by tier. gemini/vscode use static overlays
// regardless of tier (TBD per design.md §D5). Other agents are explicit no-ops.
//
// (TBD) Tier mapping for gemini: only exposes defaultApprovalMode (default|auto_edit|yolo)
// with no deny-list equivalent — see design.md §D5 for why this is deferred.
//
// (TBD) Tier mapping for vscode: only exposes chat.tools.autoApprove (bool, binary)
// with no deny-list granularity — see design.md §D5 for why this is deferred.
func agentOverlay(agent model.Agent, tier model.PermissionTier) []byte {
	switch agent {
	case model.AgentClaude:
		return composeClaudeOverlay(tier)
	case model.AgentOpenCode:
		return composeOpencodeOverlay(tier)
	case model.AgentGemini:
		return geminiCLIOverlayJSON
	case model.AgentVSCode:
		return vscodeCopilotOverlayJSON
	case model.AgentCursor:
		// Cursor manages permissions via cli-config.json, not settings.json.
		return nil
	case model.AgentCodex:
		// Codex uses config.toml and is handled by mergeCodexTOMLFile.
		return nil
	case model.AgentAntigravity:
		// Antigravity manages permissions via IDE UI. No injectable settings.json.
		return nil
	case model.AgentWindsurf:
		// Windsurf manages permissions via IDE UI. No injectable settings.json schema.
		return nil
	default:
		return nil
	}
}
