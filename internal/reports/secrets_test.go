package reports

import "testing"

func TestScanSecrets_DetectsKnownPatterns_NeverReturnsTheValue(t *testing.T) {
	cases := map[string]string{
		"github-token":      "token: ghp_abcdefghijklmnopqrstuvwxyz0123456789",
		"slack-token":       "xoxb-1234567890-abcdefgh",
		"aws-access-key-id": "AKIAABCDEFGHIJKLMNOP",
		"private-key-block": "-----BEGIN RSA PRIVATE KEY-----\nMIIBogIBAAJ...",
		"email":             "contacto: alguien@ejemplo.com",
	}
	for wantKind, text := range cases {
		kinds := scanSecrets(text)
		found := false
		for _, k := range kinds {
			if k == wantKind {
				found = true
			}
			if k == text {
				t.Fatalf("scanSecrets devolvió el valor crudo en vez del tipo: %q", k)
			}
		}
		if !found {
			t.Fatalf("texto %q: esperaba detectar %q, encontró %v", text, wantKind, kinds)
		}
	}
}

func TestScanSecrets_PlainTextWithoutSecrets_ReturnsNothing(t *testing.T) {
	kinds := scanSecrets("informe sin nada sensible, sólo prosa de revisión.")
	if len(kinds) != 0 {
		t.Fatalf("falso positivo: %v", kinds)
	}
}

func TestReportsCapture_SecretTakesPriorityOverAmbiguousIdentity(t *testing.T) {
	repo := newRepo(t)
	in := baseInput(repo)
	in.DispatchID = "" // identidad ambigua
	in.Source.Text = "credencial de prueba: ghp_" + "x0000000000000000000000000000000000"

	res, err := Capture(in)
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if res.Status != StatusRetained {
		t.Fatalf("Status = %s, quería RETAINED: un secreto bloquea antes que la identidad ambigua", res.Status)
	}
}
