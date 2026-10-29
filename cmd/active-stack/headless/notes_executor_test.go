package headless_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/cmd/active-stack/headless"
	"github.com/Group-Active-IA/active-stack/internal/backup"
	extinstaller "github.com/Group-Active-IA/active-stack/internal/harness/external"
	"github.com/Group-Active-IA/active-stack/internal/install"
	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/system"
)

// TestHeadlessExecutor_PrintsSkippedNotes (lean-agent-workflow 5.3): el
// salto del MCP engram porque el plugin ya lo provee se imprime en la
// salida del instalador, inline y en el resumen final, con exit 0.
func TestHeadlessExecutor_PrintsSkippedNotes(t *testing.T) {
	restoreSnap := install.SetSnapshotCreate(func(string, []string) (backup.Manifest, error) { return backup.Manifest{}, nil })
	defer restoreSnap()
	note := `MCP "engram" no registrado en /h/.claude.json: el plugin engram@engram ya está habilitado y lo provee`
	restoreExt := install.SetExternalInstallFn(func(context.Context, model.Harness, system.PlatformProfile, []extinstaller.AgentAdapter, string) (extinstaller.Result, error) {
		return extinstaller.Result{AlreadyInstalled: true, Notes: []string{note}}, nil
	})
	defer restoreExt()

	h := model.Harness{
		ID: "engram", Type: model.HarnessExternal,
		External:     &model.External{Method: "homebrew", MCP: &model.MCP{Name: "engram", Command: "engram"}},
		InstallModes: []model.InstallMode{model.ModeFull},
	}
	cat := &fakeExecCatalog{harnesses: []model.Harness{h}}
	reg := &fakeExecRegistry{adapters: map[model.Agent]install.AgentAdapter{model.AgentClaude: fakeExecAdapter{agent: model.AgentClaude}}}

	var out bytes.Buffer
	code := headless.RunHeadless(headless.ParsedFlags{
		Yes: true, HomeDir: t.TempDir(), NoSelfInstall: true,
		Intent: install.Intent{Agents: []model.Agent{model.AgentClaude}, Mode: model.ModeFull},
	}, cat, reg, &out)
	if code != 0 {
		t.Fatalf("exit %d, want 0:\n%s", code, out.String())
	}
	if n := strings.Count(out.String(), "engram@engram"); n < 2 {
		t.Fatalf("note must appear inline and in the end summary (seen %d times):\n%s", n, out.String())
	}
}
