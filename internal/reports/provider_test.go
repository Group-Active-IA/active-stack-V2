package reports

import (
	"os"
	"strings"
	"testing"
)

func TestCodexResolveComplete(t *testing.T) {
	raw, err := os.ReadFile("testdata/codex-0.132.0.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	s := ResolveProvider("codex", "0.132.0", "session-fixture", raw)
	if s.Contract != ContractCodexExec132 || s.Text != "Informe inventado: verificación terminada." || s.Session != "session-fixture" || s.Incomplete || s.Raw == "" {
		t.Fatalf("no resolvió informe completo: contract=%q text=%q incomplete=%v", s.Contract, s.Text, s.Incomplete)
	}
	if s.Model != "" || s.Provider != "" || s.Turn != "" {
		t.Fatal("inventó metadatos ausentes")
	}
}

func TestCodexResolvePartial(t *testing.T) {
	raw, _ := os.ReadFile("testdata/codex-0.132.0.jsonl")
	cases := map[string]string{
		"nulo":       `{"type":"thread.started","thread_id":null}`,
		"parcial":    strings.ReplaceAll(string(raw), `"turn.completed"`, `"turn.failed"`),
		"orden":      `{"type":"turn.completed"}` + "\n" + string(raw),
		"ambiguo":    string(raw) + string(raw),
		"sin cierre": strings.Split(string(raw), `{"type":"turn.completed"`)[0],
		"json roto":  "{",
	}
	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			s := ResolveProvider("codex", "0.132.0", "session-fixture", []byte(data))
			if s.Contract != ContractCodexExec132 || !s.Incomplete {
				t.Fatal("evento incompleto no quedó pendiente")
			}
		})
	}
	for _, session := range []string{"", "otra-sesion"} {
		s := ResolveProvider("codex", "0.132.0", session, raw)
		if !s.Incomplete {
			t.Fatal("sesión ajena o ausente aceptada")
		}
	}
}

func TestClaudeOpenCodeUnsupported(t *testing.T) {
	for _, p := range []string{"claude", "opencode", "codex"} {
		s := ResolveProvider(p, "no-confirmada", "session-fixture", []byte(`{"text":"no inventar"}`))
		if knownContract(s.Contract) {
			t.Fatalf("%s habilitado sin contrato confirmado", p)
		}
	}
}

func TestProviderCaptureRawAndPending(t *testing.T) {
	in := baseInput(newRepo(t))
	in.Source.Raw = `{"type":"item.completed","text":"original"}`
	res, err := Capture(in)
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(in.RepoDir + "/" + res.Dir + "/raw.txt")
	if err != nil || string(b) != in.Source.Raw {
		t.Fatal("no conserva evento bruto")
	}
	in = baseInput(newRepo(t))
	in.Source.Incomplete = true
	res, err = Capture(in)
	if err != nil || res.Status != StatusPending {
		t.Fatal("evento parcial acreditado completo")
	}
}

func TestProviderCaptureSecretRawConflict(t *testing.T) {
	in := baseInput(newRepo(t))
	if err := os.Remove(in.RepoDir + "/.gitignore"); err != nil {
		t.Fatal(err)
	}
	first, err := Capture(in)
	if err != nil || first.Status != StatusComplete {
		t.Fatal(err)
	}
	fake := "ghp_" + strings.Repeat("F", 30)
	in.Source.Raw = `{"extra":"` + fake + `"}`
	res, err := Capture(in)
	if err != nil || res.Status != StatusRetained {
		t.Fatal("secreto en bruto/conflicto no retenido")
	}
	testGit(t, in.RepoDir, "add", "-A")
	if strings.Contains(testGit(t, in.RepoDir, "diff", "--cached"), fake) {
		t.Fatal("secreto filtrado a staging")
	}
}

func TestClaudeHandbackBeforeFarewell(t *testing.T) {
	raw, err := os.ReadFile("testdata/claude-2.1.271-hooks.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	s := ResolveProvider("claude", "2.1.271", "session-fixture", raw, "child-fixture")
	if !knownContract(s.Contract) || s.Incomplete || s.Text != "Informe inventado: revisión completa." || s.Agent != "child-fixture" {
		t.Fatal("Handback no resuelto antes de despedida")
	}
	partial := strings.Split(string(raw), "\n")[0]
	if s := ResolveProvider("claude", "2.1.271", "session-fixture", []byte(partial), "child-fixture"); !s.Incomplete {
		t.Fatal("Handback sin Stop completo")
	}
}

func TestProviderCaptureBlankRawSecret(t *testing.T) {
	in := baseInput(newRepo(t))
	in.Source.Text = " "
	in.Source.Raw = `{"text":"ghp_` + strings.Repeat("F", 30) + `"}`
	res, err := Capture(in)
	if err != nil || res.Status != StatusRetained {
		t.Fatal("bruto sensible con texto vacío no retenido")
	}
}

func TestProviderCaptureMetadataSecret(t *testing.T) {
	in := baseInput(newRepo(t))
	fake := "ghp_" + strings.Repeat("F", 30)
	in.Source.Session = fake
	res, err := Capture(in)
	if err != nil || res.Status != StatusRetained {
		t.Fatal("metadatos sensibles no retenidos")
	}
	b, err := os.ReadFile(in.RepoDir + "/" + res.Dir + "/manifest.json")
	if err != nil || strings.Contains(string(b), fake) {
		t.Fatal("metadatos filtrados")
	}
}

func TestClaudeTriangle(t *testing.T) {
	raw, _ := os.ReadFile("testdata/claude-2.1.271-hooks.jsonl")
	for _, data := range []string{
		strings.ReplaceAll(string(raw), "child-fixture", "otra-hija"),
		strings.ReplaceAll(string(raw), `"SubagentStop"`, `"Stop"`),
		strings.ReplaceAll(string(raw), `"message":"Informe inventado: revisión completa."`, `"message":null`),
		string(raw) + string(raw),
		strings.Split(string(raw), "\n")[1],
	} {
		s := ResolveProvider("claude", "2.1.271", "session-fixture", []byte(data), "child-fixture")
		if !s.Incomplete {
			t.Fatal("asociación parcial o errónea completa")
		}
	}
	s := ResolveProvider("claude", "2.1.271", "", raw, "child-fixture")
	if !s.Incomplete {
		t.Fatal("sin sesión completo")
	}
}

func TestCodexTriangleBlankAndRetry(t *testing.T) {
	raw, _ := os.ReadFile("testdata/codex-0.132.0.jsonl")
	for _, text := range []string{"", "   "} {
		data := strings.ReplaceAll(string(raw), "Informe inventado: verificación terminada.", text)
		in := baseInput(newRepo(t))
		in.Source = ResolveProvider("codex", "0.132.0", "session-fixture", []byte(data))
		res, err := Capture(in)
		if err != nil || res.Status != StatusPending {
			t.Fatal("texto vacío completo")
		}
	}
	in := baseInput(newRepo(t))
	in.Source = ResolveProvider("codex", "0.132.0", "session-fixture", raw)
	first, err := Capture(in)
	if err != nil || first.Status != StatusComplete {
		t.Fatal(err)
	}
	second, err := Capture(in)
	if err != nil || second.Status != StatusComplete || second.Dir != first.Dir {
		t.Fatal("reintento no idempotente")
	}
}

func TestClaudePreToolIsNotDelivery(t *testing.T) {
	raw, _ := os.ReadFile("testdata/claude-2.1.271-hooks.jsonl")
	data := strings.ReplaceAll(string(raw), "PostToolUse", "PreToolUse")
	s := ResolveProvider("claude", "2.1.271", "session-fixture", []byte(data), "child-fixture")
	if !s.Incomplete {
		t.Fatal("PreToolUse no acredita entrega exitosa")
	}
}

func TestCodexNullCompletion(t *testing.T) {
	raw, _ := os.ReadFile("testdata/codex-0.132.0.jsonl")
	data := strings.Split(string(raw), `{"type":"turn.completed"`)[0] + `{"type":"turn.completed","usage":null}`
	s := ResolveProvider("codex", "0.132.0", "session-fixture", []byte(data))
	if !s.Incomplete {
		t.Fatal("cierre nulo acreditado")
	}
}

func TestOpenCodeContractNotConfirmed(t *testing.T) {
	raw, err := os.ReadFile("testdata/opencode-health-unconfirmed.json")
	if err != nil {
		t.Fatal(err)
	}
	in := baseInput(newRepo(t))
	in.Source = ResolveProvider("opencode", "fixture-unconfirmed", "session-fixture", raw)
	res, err := Capture(in)
	if err != nil || res.Status != StatusUnsupported {
		t.Fatal("OpenCode sin contrato confirmado no bloqueado")
	}
	if _, err := os.Stat(in.RepoDir + "/" + res.Dir); !os.IsNotExist(err) {
		t.Fatal("UNSUPPORTED escribió artefactos")
	}
}

// La retención también cubre eventos parciales que el parser no puede resolver.
func TestProviderPartialSecretRetained(t *testing.T) {
	for _, provider := range []string{"claude", "codex"} {
		t.Run(provider, func(t *testing.T) {
			fake := "ghp_" + strings.Repeat("F", 30)
			version := "0.132.0"
			if provider == "claude" {
				version = "2.1.271"
			}
			raw := []byte(`{"incomplete":"` + fake + `"}`)
			in := baseInput(newRepo(t))
			if err := os.Remove(in.RepoDir + "/.gitignore"); err != nil {
				t.Fatal(err)
			}
			in.Source = ResolveProvider(provider, version, "session-fixture", raw, "child-fixture")
			res, err := Capture(in)
			if err != nil || res.Status != StatusRetained {
				t.Fatalf("evento parcial sensible: status=%s err=%v", res.Status, err)
			}
			testGit(t, in.RepoDir, "add", "-A")
			if strings.Contains(testGit(t, in.RepoDir, "diff", "--cached"), fake) {
				t.Fatal("secreto en staging")
			}
		})
	}
}

func TestProviderPartialWithoutSecretPending(t *testing.T) {
	for _, p := range []struct{ name, version string }{{"claude", "2.1.271"}, {"codex", "0.132.0"}} {
		in := baseInput(newRepo(t))
		raw := []byte(`{"incomplete":"texto inventado"}`)
		in.Source = ResolveProvider(p.name, p.version, "session-fixture", raw, "child-fixture")
		res, err := Capture(in)
		if err != nil || res.Status != StatusPending {
			t.Fatalf("%s: %s %v", p.name, res.Status, err)
		}
		b, err := os.ReadFile(in.RepoDir + "/" + res.Dir + "/raw.txt")
		if err != nil || string(b) != string(raw) {
			t.Fatal("bruto parcial perdido")
		}
	}
}
