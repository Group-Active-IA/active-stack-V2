package roledefs

import "github.com/Group-Active-IA/active-stack/internal/model"

// AgentAdapter is the minimal interface this harness needs per agent,
// following the local-interface pattern of internal/harness/skill.
// AgentsDir(homeDir) returning "" signals the adapter does not support this
// capability (install/uninstall skip it silently) — see doc.go for why
// Codex/OpenCode return "" today (J3/J4, TBD).
type AgentAdapter interface {
	Agent() model.Agent
	AgentsDir(homeDir string) string
}

// Outcome describes what happened to a single resource for a single agent.
type Outcome string

const (
	OutcomeCreated   Outcome = "created"   // absent → written
	OutcomeNoop      Outcome = "noop"      // identical to what we'd write
	OutcomeUpdated   Outcome = "updated"   // ours, intact, backed up + overwritten
	OutcomeConflict  Outcome = "conflict"  // personal/unknown/modified — preserved
	OutcomeDeleted   Outcome = "deleted"   // uninstall: removed (was ours, intact)
	OutcomePreserved Outcome = "preserved" // uninstall: not ours (or modified) — kept
)

// ResourceResult reports the outcome for one (resource, agent) pair.
type ResourceResult struct {
	Agent        model.Agent
	Resource     string // e.g. "role-dba.md"
	Outcome      Outcome
	Path         string // the canonical path this resource lives at
	ConflictPath string // set only when Outcome == OutcomeConflict: where our version landed instead
}

// Result is the aggregate outcome of one Install or Uninstall call.
type Result struct {
	Resources []ResourceResult
	// Skipped lists agents whose AgentsDir capability is empty (unsupported).
	Skipped []model.Agent
}

// Conflicts returns the subset of Resources with Outcome == OutcomeConflict.
func (r Result) Conflicts() []ResourceResult {
	var out []ResourceResult
	for _, rr := range r.Resources {
		if rr.Outcome == OutcomeConflict {
			out = append(out, rr)
		}
	}
	return out
}
