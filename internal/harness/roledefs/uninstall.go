package roledefs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Group-Active-IA/active-stack/internal/model"
)

// Uninstall removes the resources this harness owns and leaves everything
// else alone. For each (resource, agent) pair:
//
//   - absent → nothing to do.
//   - present and the manifest records OUR hash for the current bytes
//     ("propio intacto") → removed.
//   - present but not ours (personal, unknown, or user-modified since we
//     wrote it) → preserved, reported as such. Never deleted "by glob nor by
//     name/manifest alone" — the current bytes are hashed and compared every
//     time, the manifest entry alone is never trusted.
//   - a symlink at the destination is never followed for deletion — treated
//     as not-ours and preserved, same conservative rule as Install.
//
// The agents directory itself is never removed, even when every resource
// this harness owns has been deleted from it: it is shared with the user's
// own personal agent definitions, and "container removal" (propose, sección
// 4) only applies to a directory this harness exclusively owns, which this
// one never is.
func Uninstall(adapters []AgentAdapter, homeDir string, resources []Resource) (Result, error) {
	var result Result

	for _, adapter := range adapters {
		dir := adapter.AgentsDir(homeDir)
		if dir == "" {
			result.Skipped = append(result.Skipped, adapter.Agent())
			continue
		}

		m, err := loadManifest(homeDir, adapter.Agent())
		if err != nil {
			return result, err
		}

		changed := false
		for _, res := range resources {
			destPath := filepath.Join(dir, res.Name)
			rr, removed, err := uninstallOneResource(adapter.Agent(), destPath, res, m)
			if err != nil {
				return result, fmt.Errorf("roledefs: uninstall %q for agent %q: %w", res.Name, adapter.Agent(), err)
			}
			if rr != nil {
				result.Resources = append(result.Resources, *rr)
			}
			if removed {
				delete(m.Resources, res.Name)
				changed = true
			}

			uninstallConflictSideFile(destPath, res)
		}

		if changed {
			if err := saveManifest(homeDir, adapter.Agent(), m); err != nil {
				return result, err
			}
		}
	}

	return result, nil
}

// uninstallOneResource decides and applies the outcome for the main
// destination path. Returns (nil, false, nil) when there was nothing there
// to begin with — no entry worth reporting. The bool return reports whether
// the manifest entry for this resource should be dropped.
func uninstallOneResource(agent model.Agent, destPath string, res Resource, m ownershipManifest) (*ResourceResult, bool, error) {
	lst, lerr := os.Lstat(destPath)
	switch {
	case lerr != nil && os.IsNotExist(lerr):
		return nil, false, nil
	case lerr != nil:
		return nil, false, fmt.Errorf("stat %q: %w", destPath, lerr)
	case lst.Mode()&os.ModeSymlink != 0:
		// Never follow a symlink for deletion — same conservative rule as
		// Install: treat it as not-ours.
		return &ResourceResult{Agent: agent, Resource: res.Name, Path: destPath, Outcome: OutcomePreserved}, false, nil
	}

	data, err := os.ReadFile(destPath)
	if err != nil {
		return nil, false, fmt.Errorf("read %q: %w", destPath, err)
	}

	recorded, known := m.Resources[res.Name]
	if known && recorded == sha256Hex(data) {
		if err := os.Remove(destPath); err != nil {
			return nil, false, fmt.Errorf("remove %q: %w", destPath, err)
		}
		return &ResourceResult{Agent: agent, Resource: res.Name, Path: destPath, Outcome: OutcomeDeleted}, true, nil
	}
	return &ResourceResult{Agent: agent, Resource: res.Name, Path: destPath, Outcome: OutcomePreserved}, false, nil
}

func uninstallConflictSideFile(destPath string, res Resource) {
	altPath := destPath + ".active-stack.md"
	data, err := os.ReadFile(altPath)
	if err != nil {
		return // absent or unreadable — nothing this function can safely do
	}
	if bytes.Equal(data, res.Content) {
		// Unambiguously ours: we only ever write our own candidate content
		// to this exact alternate path. Best-effort removal.
		_ = os.Remove(altPath)
	}
}
