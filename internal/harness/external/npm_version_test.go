package external

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/system"
)

// ──────────────────────────────────────────────────────────────────────────────
// Presencia no es versión (review 2026-10-29, hallazgo 8.2)
//
// Después del skip por binario presente, installNPM devolvía AlreadyInstalled
// para cualquier openspec del PATH, también uno viejo al que le faltan los
// subcomandos que el bloque lean usa (`openspec status --change --json`,
// `openspec instructions`). Con min_version en el catálogo se compara; una
// actualización que falla no voltea la instalación (ver npm_already_test.go):
// queda una nota con el comando para correrlo a mano.
// ──────────────────────────────────────────────────────────────────────────────

func fakeNPM(t *testing.T, version string, installErr error) *[]string {
	t.Helper()
	origRunner, origLookPath := runner, lookPath
	t.Cleanup(func() { runner, lookPath = origRunner, origLookPath })
	var calls []string
	lookPath = func(file string) (string, error) { return "/usr/bin/" + file, nil }
	runner = func(_ context.Context, name string, args ...string) ([]byte, error) {
		call := name + " " + strings.Join(args, " ")
		calls = append(calls, call)
		if strings.HasSuffix(call, "--version") {
			return []byte(version), nil
		}
		return []byte("npm output"), installErr
	}
	return &calls
}

func openspecHarness(min string) model.Harness {
	return model.Harness{ID: "openspec", External: &model.External{Method: "npm", Pkg: "@fission-ai/openspec", MinVersion: min}}
}

func installCalls(calls []string) int {
	n := 0
	for _, c := range calls {
		if strings.Contains(c, "install -g") {
			n++
		}
	}
	return n
}

func TestInstallNPM_KeepsBinaryAtOrAboveMinVersion(t *testing.T) {
	for _, v := range []string{"1.14.3", "1.15.0\n", "2.0.0", "openspec version 1.14.3"} {
		calls := fakeNPM(t, v, nil)
		res, err := installNPM(context.Background(), openspecHarness("1.14.3"), system.PlatformProfile{OS: "linux", NpmWritable: true})
		if err != nil {
			t.Fatalf("version %q: %v", v, err)
		}
		if installCalls(*calls) != 0 || !res.AlreadyInstalled {
			t.Errorf("version %q >= min: want skip, got calls %v AlreadyInstalled=%v", v, *calls, res.AlreadyInstalled)
		}
	}
}

func TestInstallNPM_UpgradesBelowMinVersion(t *testing.T) {
	calls := fakeNPM(t, "1.3.1", nil)
	res, err := installNPM(context.Background(), openspecHarness("1.14.3"), system.PlatformProfile{OS: "linux", NpmWritable: true})
	if err != nil {
		t.Fatal(err)
	}
	if installCalls(*calls) != 1 {
		t.Fatalf("1.3.1 < 1.14.3 must run npm install once, calls: %v", *calls)
	}
	if res.AlreadyInstalled {
		t.Error("an upgrade was reported as AlreadyInstalled")
	}
	if !strings.Contains(strings.Join(res.Notes, " "), "1.3.1") {
		t.Errorf("upgrade must be reported with the old version, notes: %v", res.Notes)
	}
}

func TestInstallNPM_FailedUpgradeIsANoteNotARollback(t *testing.T) {
	fakeNPM(t, "1.3.1", errors.New("sudo: npm: command not found"))
	res, err := installNPM(context.Background(), openspecHarness("1.14.3"), system.PlatformProfile{OS: "linux"})
	if err != nil {
		t.Fatalf("a failed upgrade of a working tool must not fail the install: %v", err)
	}
	notes := strings.Join(res.Notes, " ")
	if !strings.Contains(notes, "1.14.3") || !strings.Contains(notes, "npm install -g") {
		t.Errorf("note must name the minimum and the manual command, got: %v", res.Notes)
	}
}

func TestInstallNPM_UnreadableVersionIsReportedNotAssumed(t *testing.T) {
	calls := fakeNPM(t, "garbage", nil)
	res, err := installNPM(context.Background(), openspecHarness("1.14.3"), system.PlatformProfile{OS: "linux", NpmWritable: true})
	if err != nil {
		t.Fatal(err)
	}
	if installCalls(*calls) != 0 {
		t.Errorf("unreadable version must not trigger an install, calls: %v", *calls)
	}
	if !strings.Contains(strings.Join(res.Notes, " "), "version") {
		t.Errorf("unreadable version must be reported, notes: %v", res.Notes)
	}
}

func TestInstallNPM_NoMinVersionKeepsPresenceCheck(t *testing.T) {
	calls := fakeNPM(t, "0.0.1", nil)
	res, err := installNPM(context.Background(), openspecHarness(""), system.PlatformProfile{OS: "linux", NpmWritable: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(*calls) != 0 || !res.AlreadyInstalled {
		t.Errorf("without min_version nothing may run, calls: %v", *calls)
	}
}
