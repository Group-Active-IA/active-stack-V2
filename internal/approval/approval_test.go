package approval

import (
	"strings"
	"testing"
)

const proposalPath = "docs/propuestas/x.md"
const informePath = "docs/propuestas/informes/x/01-arquitecto.md"

// buildValidApproval arma el escenario feliz: P = propuesta revisada con su
// informe referenciado, C = commit documental separado que agrega la fila de
// Aprobación completa y cita a P en el mensaje. Devuelve (dir, P, C).
func buildValidApproval(t *testing.T) (dir, p, c string) {
	t.Helper()
	dir = newRepo(t)

	body := proposalBody("ver [informe](informes/x/01-arquitecto.md).\n\n")
	writeFile(t, dir, proposalPath, body)
	writeFile(t, dir, informePath, "hallazgos del arquitecto\n")
	p = commitAll(t, dir, "docs(propuesta): revisión final")

	current, ok, err := showFile(dir, "HEAD", proposalPath)
	if err != nil || !ok {
		t.Fatalf("no pude leer la propuesta recién commiteada: ok=%v err=%v", ok, err)
	}
	writeFile(t, dir, proposalPath, current+approvalRowLine(
		"Pasos 1 a 4, modo rápido", "Juani", "23/10/2026", "tal como está el propose"))
	git(t, dir, "add", proposalPath)
	git(t, dir, "commit", "-qm", "docs(propuestas): aprobación de Juani\n\nPropuesta-revisada: "+p)
	c = strings.TrimSpace(git(t, dir, "rev-parse", "HEAD"))
	return dir, p, c
}

// TestVerify_ValidWrittenApproval_ReturnsWrittenUnverified es el caso feliz
// del contrato: estructura completa, coincidente y vigente → código 2, y
// NUNCA un estado que suene a identidad acreditada.
func TestVerify_ValidWrittenApproval_ReturnsWrittenUnverified(t *testing.T) {
	dir, p, c := buildValidApproval(t)

	got, err := Verify(Options{RepoDir: dir, ProposalPath: proposalPath, Commit: c})
	if err != nil {
		t.Fatalf("Verify devolvió error de infraestructura en un escenario válido: %v", err)
	}
	if got.Status != StatusWrittenUnverified {
		t.Fatalf("Status = %v, quería WRITTEN_UNVERIFIED (código 2). Cause=%q Field=%q",
			got.Status, got.Cause, got.Field)
	}
	if got.Status.ExitCode() != 2 {
		t.Errorf("ExitCode() = %d, quería 2", got.Status.ExitCode())
	}
	if got.Message != writtenUnverifiedMessage {
		t.Errorf("Message = %q, quería el texto literal de contrato %q", got.Message, writtenUnverifiedMessage)
	}
	if got.Scope != "tal como está el propose" {
		t.Errorf("Scope = %q, quería el texto de la columna Alcance", got.Scope)
	}
	if got.ProposalSHA != p {
		t.Errorf("ProposalSHA = %q, quería P = %q", got.ProposalSHA, p)
	}
	if got.CommitSHA != c {
		t.Errorf("CommitSHA = %q, quería C = %q", got.CommitSHA, c)
	}
	if len(got.Manifest) != 2 {
		t.Fatalf("Manifest tiene %d entradas, quería 2 (propuesta + informe): %+v", len(got.Manifest), got.Manifest)
	}
}
