package roledefs

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Group-Active-IA/active-stack/internal/backup"
	"github.com/Group-Active-IA/active-stack/internal/model"
)

// installBackupCreate is the backing function for snapshotting a resource
// before an "own, intact" overwrite. A package-level var so tests can inject
// a failing fake without touching the real filesystem semantics.
var installBackupCreate = func(snapshotDir string, paths []string) (backup.Manifest, error) {
	return backup.NewSnapshotter().Create(snapshotDir, paths)
}

// restoreManifest is the backing function used to undo an "update" during a
// self-rollback (QA-2). Package-level var, same seam pattern as
// internal/install's restoreFn.
var restoreManifest = func(m backup.Manifest) error {
	return backup.RestoreService{}.Restore(m)
}

// testHookBeforeWrite, when non-nil, runs immediately before a resource's
// final write to disk, after the concurrency re-read. It exists only so a
// test can land a racing edit inside the read-decide-write window — see
// TestInstall_ConcurrentEditAbortsWrite.
var testHookBeforeWrite func(path string)

// undoAction reverses a single write made during this Install call.
type undoAction struct {
	kind           string // "create" | "update"
	path           string
	backupManifest *backup.Manifest // set only for "update"
}

// Install materializes resources (the seven canonical role-*.md files, in
// production always DefaultResources()) into each adapter's AgentsDir,
// following the lifecycle the propose requires: absent creates, identical is
// a no-op, ours-and-intact updates with a backup, and anything personal,
// unknown, or user-modified is preserved untouched with our version written
// to an alternate path instead.
//
// Adapters whose AgentsDir is empty are skipped (capability not supported).
//
// If any resource write fails partway through — including on a LATER
// adapter after an EARLIER one already succeeded (QA-2) — every write this
// call already made, across every adapter, is undone before returning the
// error. This does not rely on the pipeline's generic rollback: that one
// only replays Rollback() on steps whose StepResult.Status is
// StepStatusSucceeded (internal/pipeline/rollback.go), and a step whose
// Run() itself returns this error never reaches that status.
func Install(adapters []AgentAdapter, homeDir, backupDir string, resources []Resource) (Result, error) {
	if err := ValidateResources(resources); err != nil {
		return Result{}, err
	}

	var result Result
	var undo []undoAction
	newManifests := map[model.Agent]ownershipManifest{}

	for _, adapter := range adapters {
		dir := adapter.AgentsDir(homeDir)
		if dir == "" {
			result.Skipped = append(result.Skipped, adapter.Agent())
			continue
		}

		m, err := loadManifest(homeDir, adapter.Agent())
		if err != nil {
			rollbackActions(undo)
			return Result{}, err
		}
		newM := ownershipManifest{Resources: make(map[string]string, len(m.Resources))}
		for k, v := range m.Resources {
			newM.Resources[k] = v
		}

		for _, res := range resources {
			destPath := filepath.Join(dir, res.Name)
			rr, action, err := installOneResource(adapter.Agent(), dir, destPath, res, m, backupDir)
			if err != nil {
				rollbackActions(undo)
				return Result{}, fmt.Errorf("roledefs: install %q for agent %q: %w", res.Name, adapter.Agent(), err)
			}
			result.Resources = append(result.Resources, rr)
			if action != nil {
				undo = append(undo, *action)
			}
			if rr.Outcome == OutcomeCreated || rr.Outcome == OutcomeUpdated {
				newM.Resources[res.Name] = sha256Hex(res.Content)
			}
		}
		newManifests[adapter.Agent()] = newM
	}

	for agent, m := range newManifests {
		if err := saveManifest(homeDir, agent, m); err != nil {
			rollbackActions(undo)
			return Result{}, fmt.Errorf("roledefs: persist ownership manifest for %q: %w", agent, err)
		}
	}

	return result, nil
}

// installOneResource decides and applies the outcome for one (resource,
// agent) pair, returning the undo action needed to reverse it (nil when
// nothing was written, i.e. Noop or an idempotent conflict re-check).
func installOneResource(agent model.Agent, dir, destPath string, res Resource, m ownershipManifest, backupDir string) (ResourceResult, *undoAction, error) {
	base := ResourceResult{Agent: agent, Resource: res.Name, Path: destPath}

	lst, lerr := os.Lstat(destPath)
	switch {
	case lerr != nil && !os.IsNotExist(lerr):
		return base, nil, fmt.Errorf("stat %q: %w", destPath, lerr)
	case lerr == nil && lst.Mode()&os.ModeSymlink != 0:
		return installSymlinkTarget(base, dir, destPath, res)
	}

	existing, rerr := os.ReadFile(destPath)
	notExist := os.IsNotExist(rerr)
	if rerr != nil && !notExist {
		return base, nil, fmt.Errorf("read %q: %w", destPath, rerr)
	}

	if notExist {
		if err := writeFile(destPath, res.Content); err != nil {
			return base, nil, err
		}
		base.Outcome = OutcomeCreated
		return base, &undoAction{kind: "create", path: destPath}, nil
	}

	if bytes.Equal(existing, res.Content) {
		base.Outcome = OutcomeNoop
		return base, nil, nil
	}

	recorded, known := m.Resources[res.Name]
	currentHash := sha256Hex(existing)
	if known && recorded == currentHash {
		return installOwnIntactUpdate(base, destPath, existing, res, backupDir)
	}

	return writeConflictSideFile(base, destPath, res, "personal or user-modified — preserved untouched")
}

// installSymlinkTarget handles a destination that is a symlink. One pointing
// outside dir is rejected outright (hard error — "rechazar traversal/
// symlinks externos"); one pointing inside dir is never followed and
// overwritten, and is treated like any other content we do not own.
func installSymlinkTarget(base ResourceResult, dir, destPath string, res Resource) (ResourceResult, *undoAction, error) {
	target, rerr := os.Readlink(destPath)
	if rerr != nil {
		return base, nil, fmt.Errorf("readlink %q: %w", destPath, rerr)
	}
	resolved := target
	if !filepath.IsAbs(resolved) {
		resolved = filepath.Join(dir, resolved)
	}
	resolved = filepath.Clean(resolved)
	cleanDir := filepath.Clean(dir)
	if resolved != cleanDir && !strings.HasPrefix(resolved, cleanDir+string(filepath.Separator)) {
		return base, nil, fmt.Errorf(
			"refusing to touch %q: it is a symlink pointing outside %q (resolves to %q)",
			destPath, dir, resolved)
	}
	return writeConflictSideFile(base, destPath, res, "symlink inside the agents dir — never followed")
}

// installOwnIntactUpdate handles the "propio intacto" path: backup, a cheap
// concurrency re-read, then overwrite.
func installOwnIntactUpdate(base ResourceResult, destPath string, existing []byte, res Resource, backupDir string) (ResourceResult, *undoAction, error) {
	snapName := fmt.Sprintf("%s-%s", time.Now().UTC().Format("20060102-150405.000000000"), res.Name)
	snapDir := filepath.Join(backupDir, snapName)
	bm, berr := installBackupCreate(snapDir, []string{destPath})
	if berr != nil {
		return base, nil, fmt.Errorf("backup %q before update: %w", destPath, berr)
	}

	if testHookBeforeWrite != nil {
		testHookBeforeWrite(destPath)
	}

	// Concurrency check: re-read right before writing. Cheap and not a full
	// lock, but it catches the common case — an edit landing in the window
	// between our read-decide and our write — instead of clobbering it
	// silently. A race narrower than this window is the same open debt
	// internal/harness/config/permissions already carries (D-12, step 3).
	reread, rrerr := os.ReadFile(destPath)
	if rrerr != nil || !bytes.Equal(reread, existing) {
		return base, nil, fmt.Errorf(
			"concurrent change detected at %q between decision and write — aborting without overwriting", destPath)
	}

	if err := writeFile(destPath, res.Content); err != nil {
		return base, nil, err
	}
	base.Outcome = OutcomeUpdated
	bmCopy := bm
	return base, &undoAction{kind: "update", path: destPath, backupManifest: &bmCopy}, nil
}

// writeConflictSideFile preserves the existing file untouched and writes our
// candidate content to destPath + ".active-stack.md" instead — idempotent:
// if that side file already carries our content, nothing is written and no
// undo is needed.
func writeConflictSideFile(base ResourceResult, destPath string, res Resource, reason string) (ResourceResult, *undoAction, error) {
	altPath := destPath + ".active-stack.md"
	base.Outcome = OutcomeConflict
	base.ConflictPath = altPath

	existingAlt, aerr := os.ReadFile(altPath)
	switch {
	case aerr == nil && bytes.Equal(existingAlt, res.Content):
		return base, nil, nil
	case aerr == nil || os.IsNotExist(aerr):
		if err := writeFile(altPath, res.Content); err != nil {
			return base, nil, fmt.Errorf("write conflict side file %q (%s): %w", altPath, reason, err)
		}
		return base, &undoAction{kind: "create", path: altPath}, nil
	default:
		return base, nil, fmt.Errorf("read existing conflict side file %q: %w", altPath, aerr)
	}
}

func writeFile(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir %q: %w", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write %q: %w", path, err)
	}
	return nil
}

// rollbackActions reverses undo actions in reverse order (last write first).
// Best-effort: an error here is swallowed because the caller is already
// returning the original failure, and a self-rollback that half-fails has
// left more behind, not less — same category of residual debt as the
// pipeline's own ExecuteRollback (internal/pipeline/rollback.go).
func rollbackActions(actions []undoAction) {
	for i := len(actions) - 1; i >= 0; i-- {
		a := actions[i]
		switch a.kind {
		case "create":
			_ = os.Remove(a.path)
		case "update":
			if a.backupManifest != nil {
				_ = restoreManifest(*a.backupManifest)
			}
		}
	}
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
