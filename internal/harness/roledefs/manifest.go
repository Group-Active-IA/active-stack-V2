package roledefs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Group-Active-IA/active-stack/internal/model"
)

// ownershipManifest records, per resource filename, the sha256 (hex) of the
// content THIS harness last wrote there. It is the only way install/update/
// uninstall can tell "ours, unmodified since we wrote it" apart from
// "personal, or ours but edited by the user" — comparing current bytes
// against the CANDIDATE content is not enough, because the candidate changes
// across versions of active-stack and a file matching last year's candidate
// but not this year's is still ours.
//
// This is the "índice" doc.go talks about: JSON state, not a SKILL.md, and
// deliberately outside ~/.claude/agents/ — see doc.go for why.
type ownershipManifest struct {
	Resources map[string]string `json:"resources"`
}

// manifestPath returns where the ownership manifest for one agent lives.
// Mirrors the existing ".active-stack/backups/..." convention (see
// internal/i18n for the precedent of non-backup state under .active-stack/).
func manifestPath(homeDir string, agent model.Agent) string {
	return filepath.Join(homeDir, ".active-stack", "state", "role-definitions", string(agent)+".json")
}

// loadManifest reads the ownership manifest for one agent. A missing file is
// not an error — it means "never installed on this machine", which is a
// legitimate state on a fresh home.
func loadManifest(homeDir string, agent model.Agent) (ownershipManifest, error) {
	data, err := os.ReadFile(manifestPath(homeDir, agent))
	if err != nil {
		if os.IsNotExist(err) {
			return ownershipManifest{Resources: map[string]string{}}, nil
		}
		return ownershipManifest{}, fmt.Errorf("roledefs: read ownership manifest for %q: %w", agent, err)
	}
	var m ownershipManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return ownershipManifest{}, fmt.Errorf("roledefs: parse ownership manifest for %q: %w", agent, err)
	}
	if m.Resources == nil {
		m.Resources = map[string]string{}
	}
	return m, nil
}

// saveManifest writes the ownership manifest for one agent, atomically
// enough for this purpose (os.WriteFile truncate+write — the manifest is
// pure bookkeeping, never the source of truth for what's on disk; a torn
// write here is recovered by the next Install re-deriving ownership from
// actual file hashes, never by trusting a corrupt manifest blindly).
func saveManifest(homeDir string, agent model.Agent, m ownershipManifest) error {
	path := manifestPath(homeDir, agent)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("roledefs: create state dir for %q manifest: %w", agent, err)
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("roledefs: marshal ownership manifest for %q: %w", agent, err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("roledefs: write ownership manifest for %q: %w", agent, err)
	}
	return nil
}
