package headless

import (
	"fmt"
	"io"

	"github.com/Group-Active-IA/active-stack/internal/pipeline"
)

// writeRollbackReport prints the outcome of the rollback stage.
//
// It exists because "Rollback: succeeded" used to be printed whenever every
// Rollback() call returned nil, which is not the same claim. In the run that
// motivated this, engram's 19 MB binary was still on disk and the user was
// told the machine had been restored. Success now only earns the word
// "succeeded" when nothing was left behind; otherwise the rollback is reported
// as partial, with the list.
func writeRollbackReport(w io.Writer, rb pipeline.StageResult) {
	if rb.Stage != pipeline.StageRollback {
		return
	}

	switch {
	case !rb.Success:
		fmt.Fprintf(w, "Rollback: failed (%v)\n", rb.Err)
	case len(rb.Unreverted) > 0:
		fmt.Fprintln(w, "Rollback: partial — the following were left in place:")
	default:
		fmt.Fprintln(w, "Rollback: succeeded")
	}

	for _, left := range rb.Unreverted {
		fmt.Fprintf(w, "  - %s\n", left)
	}
}
