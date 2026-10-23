package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func approvalGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(cmd.Environ(),
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t",
		"GIT_COMMITTER_EMAIL=t@t", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if _, lookErr := exec.LookPath("git"); lookErr != nil {
			t.Skipf("sin git en el PATH: %v", lookErr)
		}
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

func approvalWriteFile(t *testing.T, dir, rel, body string) {
	t.Helper()
	full := filepath.Join(dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// buildApprovedRepo arma un repo temporal con una aprobación válida completa
// (mismo patrón que internal/approval.buildValidApproval, duplicado a propósito:
// el CLI se prueba de punta a punta, sin importar los internals del test de
// internal/approval).
func buildApprovedRepo(t *testing.T) (dir, c string) {
	t.Helper()
	dir = t.TempDir()
	approvalGit(t, dir, "init", "-q", "-b", "main")
	approvalWriteFile(t, dir, "README.md", "hola\n")
	approvalGit(t, dir, "add", ".")
	approvalGit(t, dir, "commit", "-qm", "inicial")

	const propPath = "docs/propuestas/x.md"
	body := "# Propose\n\n## Aprobación\n\n| Elección | Persona | Fecha | Alcance |\n|---|---|---|---|\n"
	approvalWriteFile(t, dir, propPath, body)
	approvalGit(t, dir, "add", "-A")
	approvalGit(t, dir, "commit", "-qm", "docs(propuesta): revisión final")
	p := strings.TrimSpace(approvalGit(t, dir, "rev-parse", "HEAD"))

	approvalWriteFile(t, dir, propPath, body+"| Pasos 1-4 | Juani | 23/10/2026 | el alcance firmado |\n")
	approvalGit(t, dir, "add", "-A")
	approvalGit(t, dir, "commit", "-qm", "docs(propuestas): aprobación\n\nPropuesta-revisada: "+p)
	c = strings.TrimSpace(approvalGit(t, dir, "rev-parse", "HEAD"))
	return dir, c
}

func TestRunApprovalDispatch_ValidApproval_ExitCode2WithLiteralMessage(t *testing.T) {
	dir, c := buildApprovedRepo(t)
	var out bytes.Buffer

	code := runApprovalDispatch([]string{"--change", "docs/propuestas/x.md", "--commit", c}, dir, &out)
	if code != 2 {
		t.Fatalf("exit = %d, quería 2. Salida:\n%s", code, out.String())
	}
	if !strings.Contains(out.String(), "estructura válida; identidad humana no acreditada") {
		t.Errorf("la salida no trae el mensaje literal de contrato:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "el alcance firmado") {
		t.Errorf("la salida no trae el alcance escrito:\n%s", out.String())
	}
}

func TestRunApprovalDispatch_JSONOutput_ExitCode2(t *testing.T) {
	dir, c := buildApprovedRepo(t)
	var out bytes.Buffer

	code := runApprovalDispatch([]string{"--json", "--change", "docs/propuestas/x.md", "--commit", c}, dir, &out)
	if code != 2 {
		t.Fatalf("exit = %d, quería 2. Salida:\n%s", code, out.String())
	}

	var got struct {
		Status   string `json:"status"`
		ExitCode int    `json:"exit_code"`
		Scope    string `json:"scope"`
	}
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("--json no emitió JSON: %v\n%s", err, out.String())
	}
	if got.Status != "WRITTEN_UNVERIFIED" || got.ExitCode != 2 {
		t.Errorf("JSON = %+v, quería WRITTEN_UNVERIFIED/2", got)
	}
}

func TestRunApprovalDispatch_MissingFlags_ExitCode1(t *testing.T) {
	var out bytes.Buffer
	code := runApprovalDispatch(nil, t.TempDir(), &out)
	if code != 1 {
		t.Fatalf("exit = %d, quería 1 (faltan --change/--commit)", code)
	}
	if !strings.Contains(out.String(), "--change") {
		t.Errorf("el error no explica qué falta:\n%s", out.String())
	}
}

func TestRunApprovalDispatch_Help_ExitCode0(t *testing.T) {
	var out bytes.Buffer
	code := runApprovalDispatch([]string{"--help"}, t.TempDir(), &out)
	if code != 0 {
		t.Fatalf("exit = %d, quería 0 (--help no es un error de uso)", code)
	}
	if !strings.Contains(out.String(), "Usage:") {
		t.Errorf("--help no imprimió el usage:\n%s", out.String())
	}
}

func TestRunApprovalDispatch_NoApprovalTable_ExitCode1(t *testing.T) {
	dir := t.TempDir()
	approvalGit(t, dir, "init", "-q", "-b", "main")
	approvalWriteFile(t, dir, "docs/propuestas/x.md", "# Propose\n\nsin tabla.\n")
	approvalGit(t, dir, "add", "-A")
	approvalGit(t, dir, "commit", "-qm", "sin aprobación")
	c := strings.TrimSpace(approvalGit(t, dir, "rev-parse", "HEAD"))

	var out bytes.Buffer
	code := runApprovalDispatch([]string{"--change", "docs/propuestas/x.md", "--commit", c}, dir, &out)
	if code != 1 {
		t.Fatalf("exit = %d, quería 1 (MISSING). Salida:\n%s", code, out.String())
	}
	if !strings.Contains(out.String(), "MISSING") {
		t.Errorf("la salida no dice MISSING:\n%s", out.String())
	}
}
