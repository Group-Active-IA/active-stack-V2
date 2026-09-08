// Package permissions implements the permissions harness for JR Stack.
// It injects safe-by-default permission overlays into each supported agent's
// settings.json file, with mandatory backup before any write.
//
// Governance: ALTO — writes user-level agent configuration.
// Backup is MANDATORY and cannot be disabled. If backup fails, the settings
// file is NOT touched.
//
// Supported agents and their overlay behavior:
//   - claude:      acceptEdits mode + deny-list (.env, rm -rf /, etc.)
//   - opencode:    bash/read granularity; git destructive = ask
//   - gemini:      auto_edit mode
//   - vscode:      chat.tools.autoApprove = true
//   - cursor:      no-op (permissions via cli-config.json)
//   - codex:       approval_policy + sandbox_mode in config.toml; NO deny floor
//   - antigravity: no-op (permissions via IDE UI / Artifact Review Policy)
//
// Deny floors are per-agent and are NOT all the same. claude and opencode each
// get one in every tier, bypass included; codex gets none in any tier, because
// nothing here establishes a deny construct its config.toml would honour, and a
// key it ignored would make verify report a floor that does not exist. The
// reasoning is in mergeCodexTOMLFile. The codex line above said "no-op (no known
// settings.json)" until 2026-09-08, which stopped being true when
// mergeCodexTOMLFile started writing that file.
package permissions
