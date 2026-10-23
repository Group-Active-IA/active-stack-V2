package config_test

import (
	"github.com/Group-Active-IA/active-stack/internal/harness/config"
	"strings"
	"testing"
)

func TestGovernanceAprobacionHumanaEnProposal(t *testing.T) {
	body, err := config.Compose([]string{"governance"}, "claude")
	if err != nil {
		t.Fatal(err)
	}
	body = strings.ToLower(body)
	i := strings.Index(body, "## agent governance")
	if i < 0 {
		t.Fatal("falta fragmento governance")
	}
	body = body[i:]
	for _, frase := range []string{"## aprobación", "proposal.md", "| elección | persona | fecha | alcance |", "no agent fills", "written, dated and explicit scope", "missing or ambiguous approval stops construction", "new rule after signature", "own row", "before construction", "l-029", "cryptographic verification", "not implemented"} {
		if !strings.Contains(body, frase) {
			t.Errorf("falta contrato de firma: %s", frase)
		}
	}
}
