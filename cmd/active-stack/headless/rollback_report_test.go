package headless

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/pipeline"
)

func TestWriteRollbackReport(t *testing.T) {
	tests := []struct {
		name        string
		rb          pipeline.StageResult
		wantContain []string
		wantAbsent  []string
	}{
		{
			name:        "clean rollback says succeeded",
			rb:          pipeline.StageResult{Stage: pipeline.StageRollback, Success: true},
			wantContain: []string{"Rollback: succeeded"},
			wantAbsent:  []string{"partial", "failed"},
		},
		{
			// The reported bug: every Rollback() returned nil, so the installer
			// said "succeeded" while a 19 MB binary was still on disk.
			name: "success with leftovers says partial and lists them",
			rb: pipeline.StageResult{
				Stage:      pipeline.StageRollback,
				Success:    true,
				Unreverted: []string{`npm package "openspec" stays installed globally`},
			},
			wantContain: []string{"Rollback: partial", `npm package "openspec" stays installed globally`},
			wantAbsent:  []string{"Rollback: succeeded"},
		},
		{
			name: "failed rollback still lists what was left",
			rb: pipeline.StageResult{
				Stage:      pipeline.StageRollback,
				Success:    false,
				Err:        errors.New("permission denied"),
				Unreverted: []string{"brew formula \"engram\" stays installed"},
			},
			wantContain: []string{"Rollback: failed", "permission denied", `brew formula "engram" stays installed`},
			wantAbsent:  []string{"Rollback: succeeded"},
		},
		{
			name:       "no rollback stage prints nothing",
			rb:         pipeline.StageResult{Stage: pipeline.StageApply, Success: true},
			wantAbsent: []string{"Rollback"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writeRollbackReport(&buf, tt.rb)
			got := buf.String()

			for _, want := range tt.wantContain {
				if !strings.Contains(got, want) {
					t.Errorf("output missing %q; got:\n%s", want, got)
				}
			}
			for _, absent := range tt.wantAbsent {
				if strings.Contains(got, absent) {
					t.Errorf("output must not contain %q; got:\n%s", absent, got)
				}
			}
		})
	}
}
