package roledefs

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/roles"
)

// withShippedHashes reemplaza la lista embebida de versiones publicadas por
// una sintética durante el test. La lista real se valida aparte contra el
// historial git (TestShippedRoleHashes_MatchGitHistory).
func withShippedHashes(t *testing.T, m map[string][]string) {
	t.Helper()
	old := shippedRoleHashes
	shippedRoleHashes = m
	t.Cleanup(func() { shippedRoleHashes = old })
}

func resultFor(t *testing.T, res Result, name string) ResourceResult {
	t.Helper()
	for _, rr := range res.Resources {
		if rr.Resource == name {
			return rr
		}
	}
	t.Fatalf("no result for %q in %+v", name, res.Resources)
	return ResourceResult{}
}

// TestInstall_ShippedVersionWithoutManifest_IsBackedUpAndUpdated es 5.2: un
// archivo byte-idéntico a una versión YA PUBLICADA del asset es del stack
// aunque no haya manifiesto de ownership (máquina que instaló antes de que
// existiera, o manifiesto perdido). Se trata como "propio intacto": backup y
// actualización, nunca conflicto.
func TestInstall_ShippedVersionWithoutManifest_IsBackedUpAndUpdated(t *testing.T) {
	home := t.TempDir()
	agentsDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")
	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	old := []byte("---\nname: role-dba\n---\n\nversión vieja publicada por el stack\n")
	withShippedHashes(t, map[string][]string{roles.RoleDBA + ".md": {sha256Hex(old)}})
	dbaPath := filepath.Join(agentsDir, roles.RoleDBA+".md")
	if err := os.WriteFile(dbaPath, old, 0o644); err != nil {
		t.Fatal(err)
	}

	adapters := []AgentAdapter{fakeAdapter{agent: model.AgentClaude, dir: agentsDir}}
	res, err := Install(adapters, home, backupDir, mustResources(t))
	if err != nil {
		t.Fatalf("Install: %v", err)
	}

	rr := resultFor(t, res, roles.RoleDBA+".md")
	if rr.Outcome != OutcomeUpdated {
		t.Fatalf("role-dba.md outcome = %q, want %q (conflict path %q)", rr.Outcome, OutcomeUpdated, rr.ConflictPath)
	}
	got, err := os.ReadFile(dbaPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, resourceByRole(t, roles.RoleDBA).Content) {
		t.Fatalf("role-dba.md was not updated to the current stack content")
	}
	if _, err := os.Stat(dbaPath + ".active-stack.md"); !os.IsNotExist(err) {
		t.Fatalf("a shipped version must not produce a conflict side file (stat err=%v)", err)
	}
	entries, err := os.ReadDir(backupDir)
	if err != nil || len(entries) == 0 {
		t.Fatalf("expected a snapshot of the old version under %q, got entries=%v err=%v", backupDir, entries, err)
	}
	m, err := loadManifest(home, model.AgentClaude)
	if err != nil {
		t.Fatal(err)
	}
	if m.Resources[roles.RoleDBA+".md"] != sha256Hex(resourceByRole(t, roles.RoleDBA).Content) {
		t.Fatalf("manifest must record the new content after updating a shipped version, got %q", m.Resources[roles.RoleDBA+".md"])
	}
}

// TestInstall_NonShippedContent_StaysConflict triangula los bordes de 5.2:
// contenido desconocido, una versión publicada truncada (hash parcial) y una
// versión publicada de OTRO rol siguen siendo conflicto, nunca update.
func TestInstall_NonShippedContent_StaysConflict(t *testing.T) {
	shipped := []byte("---\nname: role-dba\n---\n\nversión vieja publicada por el stack\n")
	otherRole := []byte("---\nname: role-qa\n---\n\nversión vieja de QA\n")
	withShippedHashes(t, map[string][]string{
		roles.RoleDBA + ".md": {sha256Hex(shipped)},
		roles.RoleQA + ".md":  {sha256Hex(otherRole)},
	})

	cases := map[string][]byte{
		"desconocido":      []byte("mi rol personal de DBA\n"),
		"truncado":         shipped[:len(shipped)-1],
		"de-otro-rol":      otherRole,
		"con-byte-de-mas":  append(append([]byte{}, shipped...), '\n'),
		"vacio-no-shipped": {},
	}
	for name, content := range cases {
		t.Run(name, func(t *testing.T) {
			home := t.TempDir()
			agentsDir := filepath.Join(home, ".claude", "agents")
			backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")
			if err := os.MkdirAll(agentsDir, 0o755); err != nil {
				t.Fatal(err)
			}
			dbaPath := filepath.Join(agentsDir, roles.RoleDBA+".md")
			if err := os.WriteFile(dbaPath, content, 0o644); err != nil {
				t.Fatal(err)
			}
			res, err := Install([]AgentAdapter{fakeAdapter{agent: model.AgentClaude, dir: agentsDir}}, home, backupDir, mustResources(t))
			if err != nil {
				t.Fatalf("Install: %v", err)
			}
			if rr := resultFor(t, res, roles.RoleDBA+".md"); rr.Outcome != OutcomeConflict {
				t.Fatalf("outcome = %q, want %q", rr.Outcome, OutcomeConflict)
			}
			got, _ := os.ReadFile(dbaPath)
			if !bytes.Equal(got, content) {
				t.Fatalf("non-shipped content was modified")
			}
		})
	}
}

// TestInstall_SymlinkToShippedContent_IsNeverFollowed: un symlink dentro del
// directorio cuyo destino tiene una versión publicada sigue sin seguirse.
// La regla de 5.2 mira bytes de un archivo regular, nunca a través de un link.
func TestInstall_SymlinkToShippedContent_IsNeverFollowed(t *testing.T) {
	home := t.TempDir()
	agentsDir := filepath.Join(home, ".claude", "agents")
	backupDir := filepath.Join(home, ".active-stack", "backups", "role-definitions")
	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	shipped := []byte("---\nname: role-dba\n---\n\nversión vieja publicada por el stack\n")
	withShippedHashes(t, map[string][]string{roles.RoleDBA + ".md": {sha256Hex(shipped)}})
	target := filepath.Join(agentsDir, "mi-dba-real.md")
	if err := os.WriteFile(target, shipped, 0o644); err != nil {
		t.Fatal(err)
	}
	dbaPath := filepath.Join(agentsDir, roles.RoleDBA+".md")
	if err := os.Symlink("mi-dba-real.md", dbaPath); err != nil {
		t.Skipf("symlinks not supported in this environment: %v", err)
	}

	res, err := Install([]AgentAdapter{fakeAdapter{agent: model.AgentClaude, dir: agentsDir}}, home, backupDir, mustResources(t))
	if err != nil {
		t.Fatalf("Install: %v", err)
	}
	if rr := resultFor(t, res, roles.RoleDBA+".md"); rr.Outcome != OutcomeConflict {
		t.Fatalf("outcome = %q, want %q", rr.Outcome, OutcomeConflict)
	}
	got, _ := os.ReadFile(target)
	if !bytes.Equal(got, shipped) {
		t.Fatalf("symlink target was written through")
	}
}

// TestShippedRoleHashes_Anchor5ce2c0f fija, sin depender de git, la evidencia
// real: en la máquina de Juani los siete especialistas instalados coinciden
// byte a byte con 5ce2c0f. Si la lista generada pierde esos hashes, esa
// máquina vuelve a ver siete conflictos falsos.
func TestShippedRoleHashes_Anchor5ce2c0f(t *testing.T) {
	want := map[string]string{
		"role-analyst.md":   "7eea4cec3f6579b4bc0a572622efede34e02b58f9ad90c40503010b7146df159",
		"role-architect.md": "7fae1cba18e19be4096bb334ab27b84aa394b7727256e28461c327de4dfa0778",
		"role-backend.md":   "c3aea771ee24b5f4797f90332b9b2fac5e6c01474ee66d0d1330ab00b939e5de",
		"role-dba.md":       "def9b40b4f715ccf6882a569347bf6b938179ef5ab1268f9c3c55fbe6238cece",
		"role-devops.md":    "2bf177695803482371cbf14e075ca4467438b986fdd14206ec663f8e0148a576",
		"role-frontend.md":  "27e393eb16a2deeebcfa7f7582250526b96a41f945a2bd3b524752164d3a6d85",
		"role-qa.md":        "67cc3d878ab0fcbfcb38ec51e69726c8bcb7f320cc2d78336ffbde1f9d027af7",
	}
	for name, h := range want {
		if !isShippedVersion(name, h) {
			t.Errorf("%s: 5ce2c0f hash %s is not in the embedded shipped list", name, h)
		}
	}
	// Ningún hash publicado vale para otro nombre.
	if isShippedVersion("role-qa.md", want["role-dba.md"]) {
		t.Errorf("a shipped hash must only count for its own resource name")
	}
}

// TestShippedRoleHashes_MatchGitHistory es el "testeada" de la lista
// generada: recalcula desde el historial git cada versión de cada
// assets/role-*.md y exige igualdad exacta con shippedRoleHashes (ni falta
// ni sobra). La lista NO se lee de git en runtime; este test sólo vigila que
// no se desactualice. Si un commit nuevo toca los assets, este test falla
// hasta regenerar shipped.go — ver el comentario de shippedRoleHashes.
//
// Se saltea si no hay git o el paquete no está dentro de un repo (tarball).
func TestShippedRoleHashes_MatchGitHistory(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git not available: %v", err)
	}
	out, err := exec.Command("git", "log", "--all", "--format=%H", "--", "assets").Output()
	if err != nil {
		t.Skipf("not inside a git checkout: %v", err)
	}
	got := map[string]map[string]bool{}
	for _, commit := range strings.Fields(string(out)) {
		ls, err := exec.Command("git", "ls-tree", "--name-only", commit, "assets/").Output()
		if err != nil {
			t.Fatalf("git ls-tree %s: %v", commit, err)
		}
		for _, p := range strings.Fields(string(ls)) {
			name := filepath.Base(p)
			if !strings.HasPrefix(name, "role-") || !strings.HasSuffix(name, ".md") {
				continue
			}
			content, err := exec.Command("git", "show", commit+":./assets/"+name).Output()
			if err != nil {
				t.Fatalf("git show %s:%s: %v", commit, name, err)
			}
			sum := sha256.Sum256(content)
			if got[name] == nil {
				got[name] = map[string]bool{}
			}
			got[name][hex.EncodeToString(sum[:])] = true
		}
	}
	if len(got) == 0 {
		t.Skip("git history has no role assets (shallow clone?)")
	}

	flatten := func(m map[string]map[string]bool) []string {
		var out []string
		for name, hs := range m {
			for h := range hs {
				out = append(out, name+" "+h)
			}
		}
		sort.Strings(out)
		return out
	}
	embedded := map[string]map[string]bool{}
	for name, hs := range shippedRoleHashes {
		for _, h := range hs {
			if embedded[name] == nil {
				embedded[name] = map[string]bool{}
			}
			embedded[name][h] = true
		}
	}
	g, e := flatten(got), flatten(embedded)
	if strings.Join(g, "\n") != strings.Join(e, "\n") {
		t.Fatalf("shippedRoleHashes is out of date with git history.\nfrom git:\n%s\nembedded:\n%s",
			strings.Join(g, "\n"), strings.Join(e, "\n"))
	}
}
