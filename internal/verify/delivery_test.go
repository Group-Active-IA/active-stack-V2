package verify_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/agents/claude"
	"github.com/Group-Active-IA/active-stack/internal/agents/codex"
	"github.com/Group-Active-IA/active-stack/internal/agents/opencode"
	"github.com/Group-Active-IA/active-stack/internal/harness/config"
	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/verify"
)

func TestEntregaComparaCompose(t *testing.T) {
	for _, a := range []config.AgentAdapter{claude.NewAdapter(), codex.NewAdapter(), opencode.NewAdapter()} {
		t.Run(string(a.Agent()), func(t *testing.T) {
			home := t.TempDir()
			h := model.Harness{ID: "sdd-orchestrator", Type: model.HarnessConfig, Toggles: []string{"tdd"}}
			if _, err := config.Install(h, []config.AgentAdapter{a}, home); err != nil {
				t.Fatal(err)
			}
			va := a.(verify.Adapter)
			check := func(wantFailure bool) {
				t.Helper()
				results := verify.RunChecks(context.Background(), verify.ChecksForHarness(h, []verify.Adapter{va}, home))
				failed := false
				for _, r := range results {
					if r.Status == verify.CheckStatusFailed {
						failed = true
					}
				}
				if failed != wantFailure {
					t.Fatalf("falla=%v, esperada=%v: %+v", failed, wantFailure, results)
				}
			}
			check(false)
			path := a.InstructionsPath(home)
			if a.ConfigDelivery() == model.ConfigDeliveryPrimaryAgent {
				path = a.SettingsPath(home)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			altered := strings.Replace(string(data), "Strict TDD Mode: enabled", "Strict TDD Mode: disabled", 1)
			if altered == string(data) {
				t.Fatal("fixture sin texto para alterar")
			}
			if err := os.WriteFile(path, []byte(altered), 0600); err != nil {
				t.Fatal(err)
			}
			check(true)
		})
	}
}

func TestEntregaRechazaBloqueDuplicado(t *testing.T) {
	for _, a := range []config.AgentAdapter{claude.NewAdapter(), codex.NewAdapter()} {
		t.Run(string(a.Agent()), func(t *testing.T) {
			home := t.TempDir()
			h := model.Harness{ID: "sdd-orchestrator", Type: model.HarnessConfig, Toggles: []string{"tdd"}}
			if _, err := config.Install(h, []config.AgentAdapter{a}, home); err != nil {
				t.Fatal(err)
			}
			path := a.InstructionsPath(home)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, append(data, data...), 0600); err != nil {
				t.Fatal(err)
			}
			checks := verify.ChecksForHarness(h, []verify.Adapter{a.(verify.Adapter)}, home)
			for _, c := range checks {
				if c.ID != "config:compose:"+string(a.Agent()) {
					continue
				}
				err := c.Run(context.Background())
				if err == nil || !strings.Contains(err.Error(), "duplicados") {
					t.Fatalf("bloque duplicado debe fallar por markers duplicados; resultado: %v", err)
				}
				return
			}
			t.Fatal("falta check config:compose para bloque duplicado")
		})
	}
}
