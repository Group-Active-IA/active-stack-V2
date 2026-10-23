package verify

import (
	"context"
	"errors"
	"runtime/debug"
	"testing"
)

func TestBinarioInstalado(t *testing.T) {
	originalRead, originalRun := readBinaryInfo, runBinaryHelp
	defer func() { readBinaryInfo, runBinaryHelp = originalRead, originalRun }()
	for _, tc := range []struct {
		name, path      string
		readErr, runErr error
		want            bool
	}{
		{name: "cmd correcto", path: "github.com/Group-Active-IA/active-stack/cmd/active-stack"},
		{name: "test", path: "github.com/Group-Active-IA/active-stack/internal/install.test", want: true},
		{name: "otro comando", path: "example.org/cmd/otro", want: true},
		{name: "ilegible", readErr: errors.New("ilegible"), want: true},
		{name: "roles falla", path: "github.com/Group-Active-IA/active-stack/cmd/active-stack", runErr: errors.New("exit 1"), want: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			readBinaryInfo = func(string) (*debug.BuildInfo, error) { return &debug.BuildInfo{Path: tc.path}, tc.readErr }
			runBinaryHelp = func(_ context.Context, path string, args ...string) error {
				if path != "destino" || len(args) != 2 || args[0] != "roles" || args[1] != "--help" {
					t.Fatalf("invocación inesperada: %s %v", path, args)
				}
				return tc.runErr
			}
			report := BuildReport(RunChecks(context.Background(), BinaryChecks("destino")))
			if report.Ready == tc.want {
				t.Fatalf("Ready=%v, error esperado=%v", report.Ready, tc.want)
			}
		})
	}
}
