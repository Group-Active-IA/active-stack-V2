package reports

import "testing"

func TestBuildKey_MatchesProposeLiteralFormForVerify(t *testing.T) {
	got, err := BuildKey("fase-2-pasos-1-4", "verify", 2, "backend")
	if err != nil {
		t.Fatalf("BuildKey: %v", err)
	}
	want := "opsx/fase-2-pasos-1-4/verify/round-2/backend"
	if got != want {
		t.Fatalf("got %q, want %q (forma literal del propose §2a)", got, want)
	}
}

func TestBuildKey_RejectsNonPositiveRound(t *testing.T) {
	if _, err := BuildKey("c", "verify", 0, "backend"); err == nil {
		t.Fatalf("round 0 tenía que rechazarse")
	}
	if _, err := BuildKey("c", "verify", -1, "backend"); err == nil {
		t.Fatalf("round negativo tenía que rechazarse")
	}
}

func TestBuildKey_RejectsTraversalInAnySegment(t *testing.T) {
	cases := map[string]struct{ change, phase, role string }{
		"change":    {"../escape", "verify", "backend"},
		"phase":     {"c", "../escape", "backend"},
		"role":      {"c", "verify", "../escape"},
		"separator": {"c", "verify", "back/end"},
	}
	for name, c := range cases {
		if _, err := BuildKey(c.change, c.phase, 1, c.role); err == nil {
			t.Fatalf("%s: tenía que rechazarse", name)
		}
	}
}

func TestRegisterContract_EnablesAFutureAdapterContract(t *testing.T) {
	const id = "test-adapter-v1"
	if knownContract(id) {
		t.Fatalf("%q ya estaba registrado antes del test", id)
	}
	RegisterContract(id)
	if !knownContract(id) {
		t.Fatalf("RegisterContract no habilitó %q", id)
	}
}
