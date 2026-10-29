package command

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Group-Active-IA/active-stack/internal/backup"
)

// Installer materializes slash-command files from the embedded CommandsFS into
// each focused agent's command directory. It mirrors internal/harness/skill.Installer.
//
// Governance ALTO: snapshot before write (backup); idempotent content-hash skip.
type Installer struct {
	commandsFS fs.FS
}

// NewInstaller creates a new command Installer.
// commandsFS is the embedded FS (assets.CommandsFS) containing the per-agent
// command .md files under commands/<variantKey>/...
func NewInstaller(commandsFS fs.FS) *Installer {
	return &Installer{commandsFS: commandsFS}
}

// Install materializes EVERY command file embedded for each adapter that provides
// a non-empty command directory. An empty CommandsDir() skips the adapter silently.
//
// Parameters:
//   - adapters: per-agent adapters; adapters with empty CommandsDir are skipped.
//   - homeDir: user's home directory, passed to adapter.CommandsDir.
//   - backupDir: parent directory for pre-overwrite backup snapshots.
//
// The assets of an agent live under commands/<variantKey>/ in the embedded FS, and
// the destination path is derived from the asset path:
//
//	commands/<variantKey>/opsx/ship.md → <commandsDir>/opsx/ship.md
//
// One Result per file installed, so an adapter can contribute several.
func (ins *Installer) Install(
	adapters []AgentAdapter,
	homeDir, backupDir string,
) ([]Result, error) {
	if ins.commandsFS == nil {
		return nil, fmt.Errorf("command installer: commandsFS is nil")
	}

	var results []Result
	for _, adapter := range adapters {
		commandsDir := adapter.CommandsDir(homeDir)
		if commandsDir == "" {
			// This adapter does not support commands — skip silently.
			continue
		}

		adapterResults, err := ins.installForAdapter(adapter, commandsDir, backupDir)
		if err != nil {
			return nil, fmt.Errorf("agent %q: %w", adapter.Agent(), err)
		}
		results = append(results, adapterResults...)
	}
	return results, nil
}

// installForAdapter walks the adapter's variant directory in the embedded FS and
// writes every command it finds, idempotently and with backup-before-overwrite.
//
// 2026-09-09: antes leia UN archivo por variante, con la ruta hardcodeada en
// assetAndRelPathForVariant. Los assets commands/claude/opsx/ship.md,
// commands/claude/opsx/verify.md y commands/claude/active/init.md ya existian en
// el repo y no los instalaba nadie: en una maquina limpia, `active-stack install`
// inyectaba en el CLAUDE.md la doctrina de /opsx:ship y no instalaba el comando
// que la ejecuta. Caminar el FS hace que agregar un .md alcance, sin tocar Go.
func (ins *Installer) installForAdapter(
	adapter AgentAdapter,
	commandsDir, backupDir string,
) ([]Result, error) {
	assetPaths, err := commandAssetsForVariant(ins.commandsFS, adapter.VariantKey())
	if err != nil {
		return nil, err
	}

	root := variantRoot(adapter.VariantKey())
	var results []Result
	for _, assetPath := range assetPaths {
		// El path del destino se deriva del path del asset: commands/<variante>/X → <commandsDir>/X.
		relPath := filepath.FromSlash(strings.TrimPrefix(assetPath, root+"/"))
		result, installErr := ins.installOneFile(adapter, assetPath, relPath, commandsDir, backupDir)
		if installErr != nil {
			return nil, installErr
		}
		results = append(results, result)
	}
	return results, nil
}

// RelPathsForVariant returns the relative path (under commandsDir) of EVERY
// command embedded for the given variant, in lexical order, with the platform
// separator. A variant with no embedded directory returns (nil, nil).
//
// Es el mismo walk que usa el instalador, expuesto para que internal/uninstall
// borre exactamente el conjunto que se instalo. Si esta lista y la del install
// se separan, el uninstall deja huerfanos: por eso sale del mismo lugar.
func RelPathsForVariant(commandsFS fs.FS, variantKey string) ([]string, error) {
	assetPaths, err := commandAssetsForVariant(commandsFS, variantKey)
	if err != nil {
		return nil, err
	}
	root := variantRoot(variantKey)
	var relPaths []string
	for _, assetPath := range assetPaths {
		relPaths = append(relPaths, filepath.FromSlash(strings.TrimPrefix(assetPath, root+"/")))
	}
	return relPaths, nil
}

// commandAssetsForVariant lists every embedded file under commands/<variantKey>/,
// in lexical order. A variant with no directory in the FS returns (nil, nil) —
// the caller skips silently, as it did for an unknown variant key before.
func commandAssetsForVariant(commandsFS fs.FS, variantKey string) ([]string, error) {
	if variantKey == "" {
		return nil, nil
	}
	root := variantRoot(variantKey)

	var assetPaths []string
	err := fs.WalkDir(commandsFS, root, func(p string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		// Sin filtro por extension, a proposito. Filtrar a *.md se ve mas prolijo
		// y rompe la unica propiedad por la que existe este walk: que agregar un
		// archivo a assets/commands/<variante>/ alcance para que se instale, sin
		// tocar Go. El dia que un comando necesite un adjunto al lado (un script,
		// un ejemplo), un filtro lo dejaria afuera en silencio, que es el mismo
		// modo de fallar que este walk vino a cerrar en 2026-09-09.
		assetPaths = append(assetPaths, p)
		return nil
	})
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// Variante sin assets embebidos — nada que instalar.
			return nil, nil
		}
		return nil, fmt.Errorf("command: walk embedded assets at %q: %w", root, err)
	}
	return assetPaths, nil
}

// variantRoot returns the embedded FS directory for a variant. Always uses
// forward slashes: an embed.FS is not the OS filesystem, and on Windows a
// filepath.Join here yields a path that fs.WalkDir cannot open.
func variantRoot(variantKey string) string {
	return "commands/" + variantKey
}

// installOneFile writes a single command file: idempotent skip on identical
// content, snapshot before overwrite (governance ALTO), then write.
func (ins *Installer) installOneFile(
	adapter AgentAdapter,
	assetPath, relPath, commandsDir, backupDir string,
) (Result, error) {
	content, err := fs.ReadFile(ins.commandsFS, assetPath)
	if err != nil {
		return Result{}, fmt.Errorf("command: embedded asset not found at %q: %w", assetPath, err)
	}

	destFile := filepath.Join(commandsDir, relPath)

	// Idempotence check: skip if content is byte-identical.
	identical, err := checkIdempotent(destFile, content)
	if err != nil {
		return Result{}, err
	}
	if identical {
		return Result{CommandPath: destFile, AlreadyInstalled: true}, nil
	}

	// Backup existing file if it exists (governance ALTO).
	if _, statErr := os.Stat(destFile); statErr == nil {
		if backupErr := snapshotCommandFile(backupDir, destFile, adapter.VariantKey()); backupErr != nil {
			return Result{}, backupErr
		}
	}

	// Write the command file.
	if err := os.MkdirAll(filepath.Dir(destFile), 0o755); err != nil {
		return Result{}, fmt.Errorf("command: mkdir %q: %w", filepath.Dir(destFile), err)
	}
	if err := os.WriteFile(destFile, content, 0o644); err != nil {
		return Result{}, fmt.Errorf("command: write %q: %w", destFile, err)
	}

	return Result{CommandPath: destFile}, nil
}

// RelPathForVariant returns the relative path (under commandsDir) of the
// starter-add command for the given agent variant key:
//
//	claude   → "active/starter-add.md"  (namespaced; invoked as /active:starter-add)
//	opencode → "active-starter-add.md"  (flat, hyphenated; invoked as /active-starter-add)
//
// An unknown variantKey returns "" — the caller must skip silently.
//
// 2026-09-09: dejo de ser la fuente de verdad del instalador, que ahora camina
// commands/<variantKey>/ y deriva cada destino del path del asset. Sobrevive
// porque internal/uninstall todavia la usa (plan.go, steps.go) y con eso borra
// SOLO starter-add: los demas comandos instalados quedan huerfanos tras un
// uninstall. Simetrizar eso pide meterle el FS embebido al motor de uninstall,
// y es un cambio aparte.
func RelPathForVariant(variantKey string) string {
	switch variantKey {
	case "claude":
		return filepath.Join("active", "starter-add.md")
	case "opencode":
		return "active-starter-add.md"
	default:
		return ""
	}
}

// snapshotCommandFile backs up the existing command file at destFile into a
// timestamped sub-directory inside backupDir, before the file is overwritten.
// Mirrors snapshotSkillDir from internal/harness/skill/idempotent.go.
func snapshotCommandFile(backupDir, destFile, variantKey string) error {
	// Nanoseconds plus the file name keep one snapshot per file: a name by
	// second alone let two overwrites in the same second share a directory, and
	// the second snapshot replaced the first file's backup.
	snapName := fmt.Sprintf("%s-cmd-%s-%s", time.Now().UTC().Format("20060102-150405.000000000"), variantKey, filepath.Base(destFile))
	snapDir := filepath.Join(backupDir, snapName)

	s := backup.NewSnapshotter()
	if _, err := s.Create(snapDir, []string{destFile}); err != nil {
		return fmt.Errorf("command: snapshot %q: %w", destFile, err)
	}
	return nil
}
