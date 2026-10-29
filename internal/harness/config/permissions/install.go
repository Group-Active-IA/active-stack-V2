package permissions

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Group-Active-IA/active-stack/internal/backup"
	"github.com/Group-Active-IA/active-stack/internal/filemerge"
	"github.com/Group-Active-IA/active-stack/internal/model"
)

// snapshotterCreate is the function used to create a backup snapshot.
// It is a package-level variable so tests can inject a failing implementation
// to verify that backup failures abort the write (task 4.8).
var snapshotterCreate = func(snapshotDir string, paths []string) error {
	_, err := backup.NewSnapshotter().Create(snapshotDir, paths)
	return err
}

// SetSnapshotterCreate replaces the snapshotter function for testing and
// returns a restore function that resets it to the original.
// This is the only test-seam exposed by this package.
func SetSnapshotterCreate(fn func(snapshotDir string, paths []string) error) (restore func()) {
	old := snapshotterCreate
	snapshotterCreate = fn
	return func() { snapshotterCreate = old }
}

// osReadFile reads a file and returns its contents, or nil if the file does
// not exist. It is a package-level variable so tests can inject fakes.
var osReadFile = func(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read json file %q: %w", path, err)
	}
	return content, nil
}

// Install applies the permission overlay for each adapter in the list.
//
// For each adapter:
//  1. Resolve the settings file path via adapter.SettingsPath(homeDir).
//  2. Determine the overlay for the agent and tier. If nil → skip without error (no-op).
//  3. If the settings file exists → create a backup snapshot FIRST.
//     If the backup fails → abort; the settings file is NOT touched.
//  4. Merge the overlay into the existing settings (idempotent deep-merge).
//  5. Write atomically; skip if byte-identical (WriteFileAtomic).
//
// The tier controls which permission level is composed. A zero-value tier
// normalizes defensively to TierBalanceado — NEVER TierBypass.
//
// Governance: ALTO. The backup is mandatory and cannot be disabled.
func Install(homeDir string, adapters []PermissionsAdapter, tier model.PermissionTier) (Result, error) {
	// Defensive normalization: never let a zero-value produce bypass behavior.
	tier = tier.Normalize()

	var result Result

	for _, adapter := range adapters {
		changed, path, err := installOne(homeDir, adapter, tier)
		if err != nil {
			return Result{}, err
		}
		if changed {
			result.Changed = true
			result.Files = append(result.Files, path)
		}
	}

	return result, nil
}

func installOne(homeDir string, adapter PermissionsAdapter, tier model.PermissionTier) (changed bool, path string, err error) {
	settingsPath := adapter.SettingsPath(homeDir)
	if settingsPath == "" {
		// Explicit no-op: agent has no injectable settings path.
		return false, "", nil
	}

	if adapter.Agent() == model.AgentCodex {
		if err := backupIfExists(homeDir, adapter.Agent(), settingsPath); err != nil {
			return false, "", err
		}
		writeResult, err := mergeCodexTOMLFile(settingsPath, tier)
		if err != nil {
			return false, "", fmt.Errorf("permissions: merge settings for agent %q: %w", adapter.Agent(), err)
		}
		return writeResult.Changed, settingsPath, nil
	}

	overlay := agentOverlay(adapter.Agent(), tier)
	if overlay == nil {
		// Explicit no-op: agent does not support settings.json permission injection.
		return false, "", nil
	}

	// Backup BEFORE touching the file — only if the file already exists.
	if err := backupIfExists(homeDir, adapter.Agent(), settingsPath); err != nil {
		return false, "", err
	}

	var writeResult filemerge.WriteResult
	if adapter.Agent() == model.AgentOpenCode {
		// opencode-permission-replace (2026-10-23): the installer validates the
		// base BEFORE merging, instead of inheriting MergeJSONObjects's default
		// "malformed base becomes {}" tolerance. That default is fine for most
		// consumers — overlay prevalence beats aborting — but it means an
		// unparseable opencode.json would be silently discarded rather than
		// reported, and the user loses whatever was in it with no error. JSONC
		// (comments, trailing commas) still passes: ValidateJSONObject applies
		// the same normalization MergeJSONObjects already does for the base.
		writeResult, err = mergeOpencodeJSONFile(settingsPath, overlay)
	} else {
		writeResult, err = mergeJSONFile(settingsPath, overlay)
	}
	if err != nil {
		return false, "", fmt.Errorf("permissions: merge settings for agent %q: %w", adapter.Agent(), err)
	}

	return writeResult.Changed, settingsPath, nil
}

// mergeOpencodeJSONFile is mergeJSONFile plus a pre-merge validation of the
// base file. See the comment at its call site in installOne for why opencode
// specifically needs this instead of the shared mergeJSONFile.
func mergeOpencodeJSONFile(path string, overlay []byte) (filemerge.WriteResult, error) {
	baseJSON, err := osReadFile(path)
	if err != nil {
		return filemerge.WriteResult{}, err
	}
	if err := filemerge.ValidateJSONObject(baseJSON); err != nil {
		return filemerge.WriteResult{}, fmt.Errorf("opencode settings file %q is not valid JSON/JSONC: %w", path, err)
	}

	merged, err := filemerge.MergeJSONObjects(baseJSON, overlay)
	if err != nil {
		return filemerge.WriteResult{}, fmt.Errorf("merge json objects: %w", err)
	}

	return filemerge.WriteFileAtomic(path, merged, 0)
}

// mergeCodexTOMLFile writes codex's permission keys for the given tier.
//
// NO DENY FLOOR, AND THAT IS THE HONEST STATE (2026-09-08). Claude gets
// securityFloorDeny and opencode gets opencodeBashFloorDeny in every tier;
// codex gets no equivalent in any tier, and least of all in bypass, where
// approval_policy = "never" plus sandbox_mode = "danger-full-access" is total
// access with no counterweight at all.
//
// What was looked for, and what was found: the entire codex permission surface
// this repo knows is approval_policy, sandbox_mode and [sandbox_workspace_write]
// — three knobs that grade HOW MUCH is permitted. None of them denies a named
// command or path the way Claude's deny-list or opencode's pattern-objects do.
// The design doc records the same shape: its tier table gives bypass
// "Additional configuration: none"
// (docs/superpowers/specs/2026-07-01-codex-agent-support-design.md:134-138).
// internal/filemerge can already write any top-level key, table value or whole
// section, so the writer is not the limit — there is simply no key to write.
//
// Whether codex ADMITS some deny construct that nothing here happens to
// document is UNDETERMINED: codex is not installed on the machine that wrote
// this, and no checked-in artifact answers it. So no key was guessed. A key
// codex ignores would be WORSE than none, because internal/verify would read it
// back off disk and report a floor that does not exist — the same failure the
// verify check was narrowed for on 2026-09-08, when "security floor applied"
// was demoted to "permissions key present" for exactly this agent
// (internal/verify/harness_checks.go:183-187).
//
// To close this for real, someone has to verify a candidate key against a
// running codex and then write it here WITH the tier-by-tier test. Until then
// the gap is documented, not simulated.
func mergeCodexTOMLFile(path string, tier model.PermissionTier) (filemerge.WriteResult, error) {
	raw, err := osReadFile(path)
	if err != nil {
		return filemerge.WriteResult{}, err
	}
	// network is written in EVERY tier, not just balanceado. Until 2026-09-08 it
	// was written only there, and UpsertTOMLTableValues does not remove what it
	// is not asked to write: installing balanceado and then bypass left
	// network_access = false behind, a restriction from a tier that was gone.
	// Same class of bug as Claude's sticky allow-list — a tier that skips a key
	// inherits the previous tier's answer for it.
	//
	// Every tier writes false, bypass included, and that is deliberate.
	// UNVERIFIED ASSUMPTION (2026-09-08): under sandbox_mode =
	// "danger-full-access" this key is believed to be inert, because the sandbox
	// it configures is off. Nothing in this repo proves it, so the value does not
	// lean on it — the asymmetry decides. If the key is inert, false costs
	// nothing; if it is not, true would open network access nobody asked for.
	//
	// It also PERSISTS on disk. Someone who later runs codex with
	// --sandbox workspace-write would inherit whatever is written here from a
	// tier that no longer governs that dimension — exactly the stale-value
	// pattern this fix exists to kill.
	approval, sandbox, network := "on-request", "workspace-write", "false"
	switch tier.Normalize() {
	case model.TierEstricto:
		approval, sandbox, network = "untrusted", "read-only", "false"
	case model.TierBypass:
		approval, sandbox, network = "never", "danger-full-access", "false"
	}
	merged, err := filemerge.UpsertTOMLTopLevel(string(raw), map[string]string{
		"approval_policy": fmt.Sprintf("%q", approval),
		"sandbox_mode":    fmt.Sprintf("%q", sandbox),
	})
	if err != nil {
		return filemerge.WriteResult{}, err
	}
	merged, err = filemerge.UpsertTOMLTableValues(merged, "sandbox_workspace_write", map[string]string{
		"network_access": network,
	})
	if err != nil {
		return filemerge.WriteResult{}, err
	}
	return filemerge.WriteFileAtomic(path, []byte(merged), 0)
}

// backupIfExists creates a snapshot of settingsPath if the file exists.
// If the backup fails, it returns an error — the caller must NOT write.
func backupIfExists(homeDir string, agent model.Agent, settingsPath string) error {
	if _, err := os.Stat(settingsPath); err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet — nothing to back up.
			return nil
		}
		return fmt.Errorf("permissions: stat settings file %q: %w", settingsPath, err)
	}

	snapshotDir := filepath.Join(homeDir, ".active-stack", "backups", fmt.Sprintf("permissions-%s", string(agent)))
	if err := snapshotterCreate(snapshotDir, []string{settingsPath}); err != nil {
		return fmt.Errorf("permissions: backup settings for agent %q: %w", agent, err)
	}

	return nil
}

// mergeJSONFile reads the existing settings file (or treats it as empty if it
// doesn't exist), deep-merges the overlay, and writes atomically.
func mergeJSONFile(path string, overlay []byte) (filemerge.WriteResult, error) {
	baseJSON, err := osReadFile(path)
	if err != nil {
		return filemerge.WriteResult{}, err
	}

	merged, err := filemerge.MergeJSONObjects(baseJSON, overlay)
	if err != nil {
		return filemerge.WriteResult{}, fmt.Errorf("merge json objects: %w", err)
	}

	return filemerge.WriteFileAtomic(path, merged, 0)
}
