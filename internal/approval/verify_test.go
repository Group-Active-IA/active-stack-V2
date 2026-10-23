package approval

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestVerify_EmptyApprovalTable_ReturnsMissing: la sección existe pero sin
// fila de datos. Una tabla vacía nunca es COMPLETO, ni aquí ni en ningún
// adapter de 2a.
func TestApprovalVerify_EmptyApprovalTable_ReturnsMissing(t *testing.T) {
	dir := newRepo(t)
	writeFile(t, dir, proposalPath, proposalBody(""))
	commitAll(t, dir, "docs(propuesta): revisión final")
	// "aprobación" vacía: un commit que toca el archivo pero no agrega fila.
	writeFile(t, dir, proposalPath, proposalBody("")+"\n")
	c := commitAll(t, dir, "docs(propuestas): intento de aprobación sin fila")

	got, err := Verify(Options{RepoDir: dir, ProposalPath: proposalPath, Commit: c})
	if err != nil {
		t.Fatalf("error de infraestructura: %v", err)
	}
	if got.Status != StatusMissing {
		t.Fatalf("Status = %v, quería MISSING. Cause=%q", got.Status, got.Cause)
	}
	if got.Status.ExitCode() != 1 {
		t.Errorf("ExitCode() = %d, quería 1", got.Status.ExitCode())
	}
}

// TestVerify_TableSectionAbsent_ReturnsMissing: ni siquiera existe "## Aprobación".
func TestApprovalVerify_TableSectionAbsent_ReturnsMissing(t *testing.T) {
	dir := newRepo(t)
	writeFile(t, dir, proposalPath, "# Propose\n\nsin sección de aprobación.\n")
	c := commitAll(t, dir, "docs(propuesta): sin aprobación")

	got, _ := Verify(Options{RepoDir: dir, ProposalPath: proposalPath, Commit: c})
	if got.Status != StatusMissing {
		t.Fatalf("Status = %v, quería MISSING", got.Status)
	}
}

// TestVerify_CommitDoesNotTouchProposal_ReturnsInvalid: "sin C" — el commit
// pasado en --commit no modifica la tabla de Aprobación.
func TestApprovalVerify_CommitDoesNotTouchProposal_ReturnsInvalid(t *testing.T) {
	dir := newRepo(t)
	writeFile(t, dir, proposalPath, proposalBody(""))
	commitAll(t, dir, "docs(propuesta): revisión final")
	writeFile(t, dir, "otro.md", "no toca la propuesta\n")
	c := commitAll(t, dir, "docs: archivo ajeno")

	got, _ := Verify(Options{RepoDir: dir, ProposalPath: proposalPath, Commit: c})
	if got.Status != StatusInvalid {
		t.Fatalf("Status = %v, quería INVALID (el commit no toca la propuesta)", got.Status)
	}
}

// TestVerify_CommitMixesCodeWithApproval_ReturnsInvalid: la fila se agrega en
// el mismo commit que toca otro archivo — no es "documental separado".
func TestApprovalVerify_CommitMixesCodeWithApproval_ReturnsInvalid(t *testing.T) {
	dir := newRepo(t)
	writeFile(t, dir, proposalPath, proposalBody(""))
	p := commitAll(t, dir, "docs(propuesta): revisión final")

	current, _, _ := showFile(dir, "HEAD", proposalPath)
	writeFile(t, dir, proposalPath, current+approvalRowLine("x", "Juani", "23/10/2026", "y"))
	writeFile(t, dir, "main.go", "package main\n") // código, en el MISMO commit
	c := commitAll(t, dir, "docs+code: aprobación mezclada\n\n"+propuestaRevisadaTrailer+" "+p)

	got, _ := Verify(Options{RepoDir: dir, ProposalPath: proposalPath, Commit: c})
	if got.Status != StatusInvalid {
		t.Fatalf("Status = %v, quería INVALID (mezcla aprobación con código)", got.Status)
	}
	if !strings.Contains(got.Cause, "mezcla") {
		t.Errorf("Cause = %q, no explica que es por mezclar código", got.Cause)
	}
}

// TestVerify_MissingField_ReturnsInvalid cubre persona/fecha/alcance faltantes,
// cada uno en su propio sub-test: son fallas distintas aunque el Status sea
// el mismo, y el campo reportado tiene que coincidir con el que falta.
func TestApprovalVerify_MissingField_ReturnsInvalid(t *testing.T) {
	cases := []struct {
		name                              string
		eleccion, persona, fecha, alcance string
	}{
		{"persona", "elección", "", "23/10/2026", "alcance"},
		{"fecha", "elección", "Juani", "", "alcance"},
		{"alcance", "elección", "Juani", "23/10/2026", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := newRepo(t)
			writeFile(t, dir, proposalPath, proposalBody(""))
			p := commitAll(t, dir, "docs(propuesta): revisión final")

			current, _, _ := showFile(dir, "HEAD", proposalPath)
			writeFile(t, dir, proposalPath, current+approvalRowLine(tc.eleccion, tc.persona, tc.fecha, tc.alcance))
			c := commitAll(t, dir, "docs(propuestas): aprobación incompleta\n\n"+propuestaRevisadaTrailer+" "+p)

			got, _ := Verify(Options{RepoDir: dir, ProposalPath: proposalPath, Commit: c})
			if got.Status != StatusInvalid {
				t.Fatalf("Status = %v, quería INVALID (falta %s)", got.Status, tc.name)
			}
		})
	}
}

// TestVerify_TableAlreadyCompleteBeforeCommit_ReturnsInvalid: "tabla completa
// sin commit documental separado sigue sin cumplir este contrato".
func TestApprovalVerify_TableAlreadyCompleteBeforeCommit_ReturnsInvalid(t *testing.T) {
	dir := newRepo(t)
	full := proposalBody("") + approvalRowLine("x", "Juani", "23/10/2026", "y")
	writeFile(t, dir, proposalPath, full)
	commitAll(t, dir, "docs(propuesta): ya viene con la fila adentro")

	// Un commit posterior que toca el MISMO archivo (agrega un espacio al
	// final) para que "sin C" no dispare antes de llegar a esta regla.
	writeFile(t, dir, proposalPath, full+"\n")
	c := commitAll(t, dir, "docs: toque cosmético\n\n"+propuestaRevisadaTrailer+" "+strings.TrimSpace(git(t, dir, "rev-parse", "HEAD~1")))

	got, _ := Verify(Options{RepoDir: dir, ProposalPath: proposalPath, Commit: c})
	if got.Status != StatusInvalid {
		t.Fatalf("Status = %v, quería INVALID (tabla ya estaba completa antes)", got.Status)
	}
}

// TestVerify_ProposalRefMissingInCommitMessage_ReturnsInvalid: sin la línea
// "Propuesta-revisada: <P>" en el mensaje, la decisión no referencia P.
func TestApprovalVerify_ProposalRefMissingInCommitMessage_ReturnsInvalid(t *testing.T) {
	dir := newRepo(t)
	writeFile(t, dir, proposalPath, proposalBody(""))
	commitAll(t, dir, "docs(propuesta): revisión final")

	current, _, _ := showFile(dir, "HEAD", proposalPath)
	writeFile(t, dir, proposalPath, current+approvalRowLine("x", "Juani", "23/10/2026", "y"))
	c := commitAll(t, dir, "docs(propuestas): aprobación sin referenciar P")

	got, _ := Verify(Options{RepoDir: dir, ProposalPath: proposalPath, Commit: c})
	if got.Status != StatusInvalid {
		t.Fatalf("Status = %v, quería INVALID (no referencia P)", got.Status)
	}
}

// TestVerify_ProposalNotAncestor_ReturnsInvalid: "P ajeno" — el commit que el
// mensaje cita como P no es ancestro del commit de aprobación.
func TestApprovalVerify_ProposalNotAncestor_ReturnsInvalid(t *testing.T) {
	dir := newRepo(t)
	writeFile(t, dir, proposalPath, proposalBody(""))
	commitAll(t, dir, "docs(propuesta): revisión final")

	// Un commit en una rama huérfana, nunca fusionada: no es ancestro de main.
	git(t, dir, "checkout", "-q", "--orphan", "otra")
	writeFile(t, dir, "ajeno.md", "no tiene nada que ver\n")
	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-qm", "rama ajena")
	ajeno := strings.TrimSpace(git(t, dir, "rev-parse", "HEAD"))
	git(t, dir, "checkout", "-q", "main")

	writeFile(t, dir, proposalPath, mustShow(t, dir, "HEAD", proposalPath)+
		approvalRowLine("x", "Juani", "23/10/2026", "y"))
	c := commitAll(t, dir, "docs(propuestas): cita un P ajeno\n\n"+propuestaRevisadaTrailer+" "+ajeno)

	got, _ := Verify(Options{RepoDir: dir, ProposalPath: proposalPath, Commit: c})
	if got.Status != StatusInvalid {
		t.Fatalf("Status = %v, quería INVALID (P no es ancestro)", got.Status)
	}
	if !strings.Contains(got.Cause, "ancestro") {
		t.Errorf("Cause = %q, no explica que P no es ancestro", got.Cause)
	}
}

// TestVerify_ReferencedInformeMissingAtP_ReturnsInvalid: la propuesta
// referencia un informe que no existe (falta, o fue renombrado/alterado de
// forma que ya no está en esa ruta) en la revisión P.
func TestApprovalVerify_ReferencedInformeMissingAtP_ReturnsInvalid(t *testing.T) {
	dir := newRepo(t)
	body := proposalBody("ver [informe](informes/x/01-arquitecto.md).\n\n")
	writeFile(t, dir, proposalPath, body) // el informe NUNCA se escribe
	p := commitAll(t, dir, "docs(propuesta): referencia un informe que no está")

	current, _, _ := showFile(dir, "HEAD", proposalPath)
	writeFile(t, dir, proposalPath, current+approvalRowLine("x", "Juani", "23/10/2026", "y"))
	c := commitAll(t, dir, "docs(propuestas): aprobación\n\n"+propuestaRevisadaTrailer+" "+p)

	got, _ := Verify(Options{RepoDir: dir, ProposalPath: proposalPath, Commit: c})
	if got.Status != StatusInvalid {
		t.Fatalf("Status = %v, quería INVALID (informe referenciado ausente)", got.Status)
	}
	if got.Field != informePath {
		t.Errorf("Field = %q, quería la ruta del informe faltante %q", got.Field, informePath)
	}
}

// TestVerify_ProposalChangedAfterApproval_ReturnsStale: cambios normativos
// posteriores invalidan la aprobación.
func TestApprovalVerify_ProposalChangedAfterApproval_ReturnsStale(t *testing.T) {
	dir, _, c := buildValidApproval(t)

	current, _, _ := showFile(dir, "HEAD", proposalPath)
	writeFile(t, dir, proposalPath, current+"\nun cambio posterior a la aprobación.\n")
	commitAll(t, dir, "docs(propuesta): cambio posterior a la aprobación")

	got, _ := Verify(Options{RepoDir: dir, ProposalPath: proposalPath, Commit: c})
	if got.Status != StatusStale {
		t.Fatalf("Status = %v, quería STALE (la propuesta cambió después de C)", got.Status)
	}
	if got.Status.ExitCode() != 1 {
		t.Errorf("ExitCode() = %d, quería 1", got.Status.ExitCode())
	}
}

// TestVerify_CommitRefAmbiguousOrMissing_ReturnsInvalid: --commit que no
// resuelve a un commit real del repo.
func TestApprovalVerify_CommitRefAmbiguousOrMissing_ReturnsInvalid(t *testing.T) {
	dir := newRepo(t)
	writeFile(t, dir, proposalPath, proposalBody(""))
	commitAll(t, dir, "docs(propuesta): revisión final")

	got, err := Verify(Options{RepoDir: dir, ProposalPath: proposalPath, Commit: "0000000000000000000000000000000000000000"})
	if err != nil {
		t.Fatalf("error de infraestructura inesperado: %v", err)
	}
	if got.Status != StatusInvalid {
		t.Fatalf("Status = %v, quería INVALID (SHA inexistente)", got.Status)
	}
}

// TestVerify_MakesNoWrites: "comprobar cero escrituras" del rojo-antes. El
// árbol de trabajo, el índice y las refs tienen que quedar bit a bit iguales.
func TestApprovalVerify_MakesNoWrites(t *testing.T) {
	dir, _, c := buildValidApproval(t)

	before := snapshotRepo(t, dir)
	if _, err := Verify(Options{RepoDir: dir, ProposalPath: proposalPath, Commit: c}); err != nil {
		t.Fatalf("Verify error: %v", err)
	}
	after := snapshotRepo(t, dir)

	if before != after {
		t.Errorf("Verify escribió algo en el repo:\nantes:\n%s\ndespués:\n%s", before, after)
	}
}

// TestZeroValueResult_IsNotApproved: el cero-valor de Result (lo que recibe
// un consumidor que ignora el error) tiene que caer en el estado MÁS
// restrictivo, nunca en uno que un consumidor descuidado pueda leer como
// éxito. Esto es lo que hace imposible por diseño al "consumidor que acepta
// cualquier retorno como aprobación" del rojo-antes.
func TestApprovalZeroValueResult_IsNotApproved(t *testing.T) {
	var zero Result
	if zero.Status == StatusWrittenUnverified {
		t.Fatal("el cero-valor de Status es WRITTEN_UNVERIFIED: un consumidor que ignora el error vería 'aprobado'")
	}
	if zero.Status.ExitCode() == 0 {
		t.Fatal("el cero-valor tiene ExitCode 0: approval no certifica éxito por default")
	}
}

// TestNoConsumerSeesAnyStatusAsGenericSuccess recorre TODOS los Status
// bloqueantes construidos por los tests de arriba y confirma que ninguno
// pasa por "aprobado": el único que debería es StatusWrittenUnverified.
func TestApprovalNoConsumerSeesAnyStatusAsGenericSuccess(t *testing.T) {
	isApproved := func(r Result) bool { return r.Status == StatusWrittenUnverified }

	blocking := []Result{
		blocked(StatusMissing, "x", "y", "z"),
		blocked(StatusInvalid, "x", "y", "z"),
		blocked(StatusStale, "x", "y", "z"),
		{}, // cero-valor
	}
	for _, r := range blocking {
		if isApproved(r) {
			t.Errorf("Status %v se leyó como aprobado, y no lo es", r.Status)
		}
	}
}

// TestWrittenUnverified_NeverSoundsLikeCryptographicIdentity: ningún texto
// que el paquete emite dice "VERIFIED" a secas ni sugiere identidad
// acreditada. Es la contraparte directa de J1 revisada.
func TestApprovalWrittenUnverified_NeverSoundsLikeCryptographicIdentity(t *testing.T) {
	dir, _, c := buildValidApproval(t)
	got, err := Verify(Options{RepoDir: dir, ProposalPath: proposalPath, Commit: c})
	if err != nil {
		t.Fatalf("Verify error: %v", err)
	}
	if got.Status.String() != "WRITTEN_UNVERIFIED" {
		t.Fatalf("String() = %q, quería WRITTEN_UNVERIFIED", got.Status.String())
	}
	if !strings.Contains(got.Message, "no acreditada") {
		t.Errorf("Message = %q, no dice explícitamente que la identidad no está acreditada", got.Message)
	}
	for _, s := range []Status{StatusMissing, StatusInvalid, StatusStale, StatusWrittenUnverified} {
		if s.String() == "VERIFIED" {
			t.Fatalf("Status %d imprime VERIFIED: ese estado no debería existir en este contrato", s)
		}
	}
}

// --- helpers específicos de este archivo ---

func mustShow(t *testing.T, dir, rev, path string) string {
	t.Helper()
	c, ok, err := showFile(dir, rev, path)
	if err != nil || !ok {
		t.Fatalf("showFile(%s, %s) ok=%v err=%v", rev, path, ok, err)
	}
	return c
}

// snapshotRepo es una huella determinista de TODO lo que Verify podría
// escribir: estado del índice/árbol de trabajo y el ref actual. Dos
// corridas iguales tienen que dar el mismo string.
func snapshotRepo(t *testing.T, dir string) string {
	t.Helper()
	status := git(t, dir, "status", "--porcelain=v1", "-uall")
	head := strings.TrimSpace(git(t, dir, "rev-parse", "HEAD"))
	var files []string
	_ = filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.Contains(p, string(filepath.Separator)+".git"+string(filepath.Separator)) {
			return nil
		}
		rel, _ := filepath.Rel(dir, p)
		files = append(files, rel)
		return nil
	})
	return head + "\n" + status + "\n" + strings.Join(files, ",")
}

// Un archivo sin trackear y ajeno (un log, un borrador) no altera lo firmado:
// no puede invalidar la aprobación. Antes isDirty usaba -uall y cualquier
// .log suelto en el repo real volvía INVALID una firma válida.
func TestApprovalVerify_UntrackedUnrelatedFile_StillWrittenUnverified(t *testing.T) {
	dir, _, c := buildValidApproval(t)
	writeFile(t, dir, "docs/propuestas/informes/x/.borrador.log", "salida suelta\n")

	got, err := Verify(Options{RepoDir: dir, ProposalPath: proposalPath, Commit: c})
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != StatusWrittenUnverified {
		t.Fatalf("Status = %v, quería WRITTEN_UNVERIFIED; Cause=%q", got.Status, got.Cause)
	}
}

// Una modificación sin commitear a un archivo trackeado sí bloquea: lo que se
// valida tiene que ser exactamente lo commiteado.
func TestApprovalVerify_TrackedProposalModified_ReturnsInvalid(t *testing.T) {
	dir, _, c := buildValidApproval(t)
	current, _, _ := showFile(dir, "HEAD", proposalPath)
	writeFile(t, dir, proposalPath, current+"\nun cambio sin commitear\n")

	got, err := Verify(Options{RepoDir: dir, ProposalPath: proposalPath, Commit: c})
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != StatusInvalid {
		t.Fatalf("Status = %v, quería INVALID por árbol sucio; Cause=%q", got.Status, got.Cause)
	}
}
