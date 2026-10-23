package reports

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReportsCapture_Happy_ReturnsComplete(t *testing.T) {
	repo := newRepo(t)
	in := baseInput(repo)

	res, err := Capture(in)
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if res.Status != StatusComplete {
		t.Fatalf("Status = %s, quería COMPLETE (cause=%s)", res.Status, res.Cause)
	}
	if res.Key != "opsx/fase-2-pasos-1-4/verify/round-1/backend" {
		t.Fatalf("Key = %q, no matchea la forma literal del propose", res.Key)
	}
	raw := filepath.Join(repo, res.Dir, "raw.txt")
	if b, err := os.ReadFile(raw); err != nil || string(b) != in.Source.Text {
		t.Fatalf("raw.txt no tiene el contenido esperado: err=%v body=%q", err, string(b))
	}
	if _, err := os.Stat(filepath.Join(repo, res.Dir, "informe.md")); err != nil {
		t.Fatalf("informe.md no se escribió: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repo, res.Dir, "manifest.json")); err != nil {
		t.Fatalf("manifest.json no se escribió: %v", err)
	}
}

func TestReportsCapture_SameEventTwice_IsNoOp(t *testing.T) {
	repo := newRepo(t)
	in := baseInput(repo)

	first, err := Capture(in)
	if err != nil {
		t.Fatalf("primera captura: %v", err)
	}
	modTime1, _ := os.Stat(filepath.Join(repo, first.Dir, "manifest.json"))

	second, err := Capture(in)
	if err != nil {
		t.Fatalf("segunda captura (reintento): %v", err)
	}
	if second.Status != StatusComplete || second.Dir != first.Dir || second.ContentHash != first.ContentHash {
		t.Fatalf("el reintento con el mismo evento tiene que ser un no-op idéntico: first=%+v second=%+v", first, second)
	}
	modTime2, _ := os.Stat(filepath.Join(repo, second.Dir, "manifest.json"))
	if modTime1.ModTime() != modTime2.ModTime() {
		t.Fatalf("el reintento reescribió manifest.json; tenía que ser no-op (mtime cambió)")
	}
}

func TestReportsCapture_ConflictingBytes_KeepsOriginalSeparate(t *testing.T) {
	repo := newRepo(t)
	in := baseInput(repo)

	first, err := Capture(in)
	if err != nil {
		t.Fatalf("primera captura: %v", err)
	}

	in2 := in
	in2.Source.Text = "un informe completamente distinto, mismos identificadores"
	second, err := Capture(in2)
	if err != nil {
		t.Fatalf("segunda captura (conflicto): %v", err)
	}
	if second.Status != StatusConflict {
		t.Fatalf("Status = %s, quería CONFLICT", second.Status)
	}
	if second.Dir == first.Dir {
		t.Fatalf("el conflicto quedó en el MISMO directorio que el original: %s", second.Dir)
	}
	// El original no se tocó.
	rawOriginal, err := os.ReadFile(filepath.Join(repo, first.Dir, "raw.txt"))
	if err != nil || string(rawOriginal) != in.Source.Text {
		t.Fatalf("el original fue pisado por el conflicto: err=%v body=%q", err, string(rawOriginal))
	}
}

func TestReportsCapture_TwoRoles_NeverShareKey(t *testing.T) {
	repo := newRepo(t)
	backend := baseInput(repo)
	qa := baseInput(repo)
	qa.Role = "qa"
	qa.DispatchID = "d-2"

	resBackend, err := Capture(backend)
	if err != nil {
		t.Fatalf("backend: %v", err)
	}
	resQA, err := Capture(qa)
	if err != nil {
		t.Fatalf("qa: %v", err)
	}
	if resBackend.Key == resQA.Key || resBackend.Dir == resQA.Dir {
		t.Fatalf("dos roles distintos comparten clave/directorio: backend=%+v qa=%+v", resBackend, resQA)
	}
	if resBackend.Status != StatusComplete || resQA.Status != StatusComplete {
		t.Fatalf("ambos tenían que ser COMPLETE: backend=%s qa=%s", resBackend.Status, resQA.Status)
	}
}

func TestReportsCapture_BlankReport_ReturnsPending_NeverComplete(t *testing.T) {
	repo := newRepo(t)
	in := baseInput(repo)
	in.Source.Text = "   \n\t  "

	res, err := Capture(in)
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if res.Status != StatusPending {
		t.Fatalf("Status = %s, quería PENDING para un informe en blanco", res.Status)
	}
}

func TestReportsCapture_UnknownContract_ReturnsUnsupported_NoWrites(t *testing.T) {
	repo := newRepo(t)
	in := baseInput(repo)
	in.Source.Contract = "claude-hook-v1" // no registrado: 2b no está implementado

	res, err := Capture(in)
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if res.Status != StatusUnsupported {
		t.Fatalf("Status = %s, quería UNSUPPORTED", res.Status)
	}
	if _, err := os.Stat(filepath.Join(repo, "openspec", "changes", in.Change)); err == nil {
		t.Fatalf("UNSUPPORTED no debería escribir nada bajo openspec/changes/%s", in.Change)
	}
}

func TestReportsCapture_AmbiguousIdentity_MissingDispatchID_ReturnsPending(t *testing.T) {
	repo := newRepo(t)
	in := baseInput(repo)
	in.DispatchID = ""

	res, err := Capture(in)
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if res.Status != StatusPending {
		t.Fatalf("Status = %s, quería PENDING por identidad ambigua (sin dispatch id)", res.Status)
	}
}

func TestReportsCapture_AmbiguousIdentity_MissingSHA_ReturnsPending(t *testing.T) {
	repo := newRepo(t)
	in := baseInput(repo)
	in.SHA = ""

	res, err := Capture(in)
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if res.Status != StatusPending {
		t.Fatalf("Status = %s, quería PENDING por identidad ambigua (sin SHA)", res.Status)
	}
}

func TestReportsCapture_NonPositiveRound_ReturnsPending(t *testing.T) {
	repo := newRepo(t)
	in := baseInput(repo)
	in.Round = 0

	res, err := Capture(in)
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if res.Status != StatusPending {
		t.Fatalf("Status = %s, quería PENDING: reintentos no incrementan vuelta, y vuelta 0 es ambigua", res.Status)
	}
}

func TestReportsCapture_MissingRequiredFields_ReturnsError(t *testing.T) {
	repo := newRepo(t)
	_, err := Capture(CaptureInput{RepoDir: repo})
	if err == nil {
		t.Fatalf("Capture sin Change/Phase/Role tenía que devolver error, no Result")
	}
}

func TestReportsCapture_PathTraversalInChange_Rejected(t *testing.T) {
	repo := newRepo(t)
	in := baseInput(repo)
	in.Change = "../../etc"

	if _, err := Capture(in); err == nil {
		t.Fatalf("Change con \"..\" tenía que rechazarse con error")
	}
}

func TestReportsCapture_PathTraversalInRole_Rejected(t *testing.T) {
	repo := newRepo(t)
	in := baseInput(repo)
	in.Role = "backend/../../escape"

	if _, err := Capture(in); err == nil {
		t.Fatalf("Role con separador de ruta tenía que rechazarse con error")
	}
}

func TestReportsCapture_ForeignWorktree_SubdirectoryRejected(t *testing.T) {
	repo := newRepo(t)
	sub := filepath.Join(repo, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	in := baseInput(sub) // un subdirectorio del worktree, no su raíz

	if _, err := Capture(in); err == nil {
		t.Fatalf("RepoDir que no es la raíz del worktree tenía que rechazarse")
	} else if !strings.Contains(err.Error(), "raíz del worktree") {
		t.Fatalf("error inesperado: %v", err)
	}
}

func TestReportsCapture_NotAGitRepo_Rejected(t *testing.T) {
	dir := t.TempDir()
	in := baseInput(dir)

	if _, err := Capture(in); err == nil {
		t.Fatalf("un directorio que no es un repo git tenía que rechazarse")
	}
}

func TestReportsCapture_SymlinkEscapesInformesRoot_Rejected(t *testing.T) {
	repo := newRepo(t)
	outside := t.TempDir()

	changesDir := filepath.Join(repo, "openspec", "changes")
	if err := os.MkdirAll(changesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// El directorio del change es un symlink que escapa del worktree.
	if err := os.Symlink(outside, filepath.Join(changesDir, "fase-2-pasos-1-4")); err != nil {
		t.Skipf("symlinks no disponibles en este entorno: %v", err)
	}

	in := baseInput(repo)
	if _, err := Capture(in); err == nil {
		t.Fatalf("un ancestro symlinkeado fuera de la raíz tenía que rechazarse")
	}
	if entries, _ := os.ReadDir(outside); len(entries) != 0 {
		t.Fatalf("se escribió algo afuera de la raíz a través del symlink: %v", entries)
	}
}

// TestReportsCapture_FakeSecretToken_ReturnsRetained_NeverLeaksTheToken es el
// caso fijado explícitamente por el propose §2a ("Rojo antes"): un informe
// con un token falso de prueba tiene que terminar RETENIDO, sin que el token
// aparezca en ningún archivo versionado ni en la salida.
func TestReportsCapture_FakeSecretToken_ReturnsRetained_NeverLeaksTheToken(t *testing.T) {
	repo := newRepo(t)
	const fakeToken = "ghp_FAKEFAKEFAKEFAKEFAKEFAKEFAKEFAKEFAKE"
	in := baseInput(repo)
	in.Source.Text = "Informe de revisión.\n\nCredencial usada en el test: " + fakeToken + "\n"

	res, err := Capture(in)
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if res.Status != StatusRetained {
		t.Fatalf("Status = %s, quería RETAINED", res.Status)
	}
	if res.Cause == "" || strings.Contains(res.Cause, fakeToken) {
		t.Fatalf("Cause no puede contener el token: %q", res.Cause)
	}
	for _, k := range res.SecretKinds {
		if strings.Contains(k, fakeToken) {
			t.Fatalf("SecretKinds no puede contener el token: %v", res.SecretKinds)
		}
	}
	if res.RetainedPath == "" {
		t.Fatalf("RetainedPath vacío para un informe RETAINED")
	}

	// Ningún byte del token en ningún archivo bajo el repo, salvo el
	// original local retenido (que está ignorado por git, verificado abajo).
	walkErr := filepath.WalkDir(repo, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.Contains(string(b), fakeToken) {
			rel, _ := filepath.Rel(repo, path)
			if rel != res.RetainedPath {
				t.Fatalf("el token apareció en un archivo que no es el original retenido: %s", rel)
			}
		}
		return nil
	})
	if walkErr != nil {
		t.Fatalf("recorrer el repo: %v", walkErr)
	}

	// El original retenido existe y SÍ tiene el token (ahí es donde tiene
	// que vivir, local, nunca versionado).
	retained, err := os.ReadFile(filepath.Join(repo, res.RetainedPath))
	if err != nil || !strings.Contains(string(retained), fakeToken) {
		t.Fatalf("el original retenido no se persistió correctamente: err=%v", err)
	}

	// git check-ignore: el original retenido SÍ está ignorado.
	if ignored, verbose := gitCheckIgnore(t, repo, res.RetainedPath); !ignored {
		t.Fatalf("git check-ignore no marcó como ignorado el original retenido: %q", verbose)
	}

	// manifest.json y LEEME.md (versionables) no contienen el token.
	manifestPath := filepath.Join(repo, "openspec", "changes", in.Change, "informes", in.Phase, "round-1", in.Role, "manifest.json")
	mb, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("leer manifest.json: %v", err)
	}
	if strings.Contains(string(mb), fakeToken) {
		t.Fatalf("manifest.json contiene el token")
	}
	leeme, err := os.ReadFile(filepath.Join(repo, "openspec", "changes", in.Change, "informes", "LEEME.md"))
	if err != nil {
		t.Fatalf("leer LEEME.md: %v", err)
	}
	if strings.Contains(string(leeme), fakeToken) {
		t.Fatalf("LEEME.md contiene el token")
	}
}

func TestReportsCapture_OriginalIsVersionable_GitCheckIgnoreReturnsNotIgnored(t *testing.T) {
	repo := newRepo(t)
	in := baseInput(repo)

	res, err := Capture(in)
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	rawRel := filepath.Join(res.Dir, "raw.txt")

	// git check-ignore devuelve 1 (nada impreso) para un original activo: no
	// está ignorado, tiene que poder versionarse.
	if ignored, verbose := gitCheckIgnore(t, repo, rawRel); ignored {
		t.Fatalf("git check-ignore marcó como ignorado un original activo: %q", verbose)
	}

	testGit(t, repo, "add", "-A")
	tracked := testGit(t, repo, "ls-files", rawRel)
	if strings.TrimSpace(tracked) != filepath.ToSlash(rawRel) {
		t.Fatalf("tras git add, el original no quedó en el árbol versionado: %q", tracked)
	}
}

func TestReportsCapture_PermissionDenied_ReturnsError_NeverComplete(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("corriendo como root: chmod 0 no bloquea escrituras, el caso no es simulable")
	}
	repo := newRepo(t)
	changesDir := filepath.Join(repo, "openspec", "changes")
	if err := os.MkdirAll(changesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(changesDir, 0o500); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(changesDir, 0o755)

	in := baseInput(repo)
	res, err := Capture(in)
	if err == nil {
		t.Fatalf("permiso denegado tenía que devolver error, no Result (Status=%s)", res.Status)
	}
}

func TestReportsCapture_RecoversFromPartialWrite_WithoutFalseConflict(t *testing.T) {
	repo := newRepo(t)
	in := baseInput(repo)

	// Simula un crash entre escribir raw.txt y escribir manifest.json: el
	// directorio existe, raw.txt existe, pero no hay manifest.json. La
	// próxima captura tiene que recuperarse (COMPLETE), no declarar
	// conflicto contra una escritura que nunca se confirmó.
	segs, err := dirSegments(in.Change, in.Phase, in.Round, in.Role)
	if err != nil {
		t.Fatal(err)
	}
	partialDir := filepath.Join(repo, filepath.Join(segs...))
	if err := os.MkdirAll(partialDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(partialDir, "raw.txt"), []byte("basura de un intento interrumpido"), 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := Capture(in)
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if res.Status != StatusComplete {
		t.Fatalf("Status = %s, quería COMPLETE: la escritura parcial sin manifest.json tenía que ser recuperable", res.Status)
	}
	body, err := os.ReadFile(filepath.Join(repo, res.Dir, "raw.txt"))
	if err != nil || string(body) != in.Source.Text {
		t.Fatalf("raw.txt no se recuperó con el contenido correcto: err=%v body=%q", err, string(body))
	}
}

func TestReportsCapture_LEEME_ListsEveryCapturedIdentity(t *testing.T) {
	repo := newRepo(t)
	backend := baseInput(repo)
	qa := baseInput(repo)
	qa.Role = "qa"
	qa.DispatchID = "d-2"

	if _, err := Capture(backend); err != nil {
		t.Fatal(err)
	}
	if _, err := Capture(qa); err != nil {
		t.Fatal(err)
	}

	leeme, err := os.ReadFile(filepath.Join(repo, "openspec", "changes", backend.Change, "informes", "LEEME.md"))
	if err != nil {
		t.Fatalf("leer LEEME.md: %v", err)
	}
	if !strings.Contains(string(leeme), "backend") || !strings.Contains(string(leeme), "qa") {
		t.Fatalf("LEEME.md no lista ambos roles capturados: %s", leeme)
	}
}

func TestReportsCapture_OversizedContent_RejectedWithoutTruncation(t *testing.T) {
	repo := newRepo(t)
	in := baseInput(repo)
	in.Source.Text = strings.Repeat("a", MaxContentBytes+1)

	if _, err := Capture(in); err == nil {
		t.Fatalf("contenido que excede MaxContentBytes tenía que rechazarse con error, nunca truncarse en silencio")
	}
}
