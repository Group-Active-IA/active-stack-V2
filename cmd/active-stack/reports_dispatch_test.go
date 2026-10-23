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

func reportsGit(t *testing.T, dir string, args ...string) string {
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

func reportsRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	reportsGit(t, dir, "init", "-q", "-b", "main")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	reportsGit(t, dir, "add", ".")
	reportsGit(t, dir, "commit", "-qm", "inicial")
	return dir
}

func TestRunReportsCaptureDispatch_HappyPath_ReturnsZero(t *testing.T) {
	dir := reportsRepo(t)
	var buf bytes.Buffer
	args := []string{"capture",
		"--change", "fase-2-pasos-1-4",
		"--phase", "verify",
		"--round", "1",
		"--role", "backend",
		"--dispatch-id", "d-1",
		"--sha", "abc123",
		"--report-text", "informe de prueba",
	}
	code := runReportsDispatch(args, dir, &buf)
	if code != 0 {
		t.Fatalf("code = %d, salida:\n%s", code, buf.String())
	}
	if !strings.Contains(buf.String(), "COMPLETE") {
		t.Fatalf("salida no menciona COMPLETE: %s", buf.String())
	}
}

func TestRunReportsCaptureDispatch_JSON_ReportsCompleteStatusAsString(t *testing.T) {
	dir := reportsRepo(t)
	var buf bytes.Buffer
	args := []string{"capture",
		"--change", "fase-2-pasos-1-4",
		"--phase", "verify",
		"--round", "1",
		"--role", "backend",
		"--dispatch-id", "d-1",
		"--sha", "abc123",
		"--report-text", "informe de prueba",
		"--json",
	}
	code := runReportsDispatch(args, dir, &buf)
	if code != 0 {
		t.Fatalf("code = %d, salida:\n%s", code, buf.String())
	}
	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("salida no es JSON válido: %v\n%s", err, buf.String())
	}
	if out["status"] != "COMPLETE" {
		t.Fatalf("status = %v, quería \"COMPLETE\" (string, no entero)", out["status"])
	}
}

func TestRunReportsCaptureDispatch_SecretInReport_ReturnsRetainedExitCodeAndNeverLeaksIt(t *testing.T) {
	dir := reportsRepo(t)
	const fakeToken = "ghp_FAKEFAKEFAKEFAKEFAKEFAKEFAKEFAKEFAKE"
	var buf bytes.Buffer
	args := []string{"capture",
		"--change", "fase-2-pasos-1-4",
		"--phase", "verify",
		"--round", "1",
		"--role", "backend",
		"--dispatch-id", "d-1",
		"--sha", "abc123",
		"--report-text", "credencial de prueba: " + fakeToken,
	}
	code := runReportsDispatch(args, dir, &buf)
	if code != 4 {
		t.Fatalf("code = %d, quería 4 (RETAINED). Salida:\n%s", code, buf.String())
	}
	if strings.Contains(buf.String(), fakeToken) {
		t.Fatalf("el token apareció en la salida del CLI: %s", buf.String())
	}
}

func TestRunReportsCaptureDispatch_EmptyDirUsesWorkingDirectory(t *testing.T) {
	// Mismo defecto que ya mordió a approval (ver informe del paso 1): un
	// dir vacío tiene que resolver os.Getwd(), no reventar con RepoDir
	// obligatorio sin completar. Probamos moviendo el proceso de test a un
	// repo temporal real.
	dir := reportsRepo(t)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(wd)

	var buf bytes.Buffer
	args := []string{"capture",
		"--change", "fase-2-pasos-1-4",
		"--phase", "verify",
		"--round", "1",
		"--role", "backend",
		"--dispatch-id", "d-1",
		"--sha", "abc123",
		"--report-text", "informe de prueba",
	}
	code := runReportsDispatch(args, "", &buf)
	if code != 0 {
		t.Fatalf("code = %d, salida:\n%s", code, buf.String())
	}
}

func TestRunReportsCaptureDispatch_MissingRequiredFlags_ReturnsOne(t *testing.T) {
	var buf bytes.Buffer
	code := runReportsDispatch([]string{"capture"}, "repo-no-usado", &buf)
	if code != 1 {
		t.Fatalf("code = %d, quería 1 por flags faltantes", code)
	}
}

func TestRunReportsDispatch_UnknownSubcommand_ReturnsOne(t *testing.T) {
	var buf bytes.Buffer
	code := runReportsDispatch([]string{"bogus"}, "repo-no-usado", &buf)
	if code != 1 {
		t.Fatalf("code = %d, quería 1 por subcomando desconocido", code)
	}
}

func TestRunReportsCaptureDispatch_Help_ReturnsZero(t *testing.T) {
	var buf bytes.Buffer
	code := runReportsDispatch([]string{"capture", "--help"}, "repo-no-usado", &buf)
	if code != 0 {
		t.Fatalf("code = %d, quería 0 para --help", code)
	}
	if !strings.Contains(buf.String(), "Usage:") {
		t.Fatalf("--help no mostró usage: %s", buf.String())
	}
}
