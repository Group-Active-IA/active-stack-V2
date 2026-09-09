package pipeline

import "time"

// StepStatus describes the outcome of a single step execution.
type StepStatus string

const (
	StepStatusPending    StepStatus = "pending"
	StepStatusRunning    StepStatus = "running"
	StepStatusSucceeded  StepStatus = "succeeded"
	StepStatusFailed     StepStatus = "failed"
	StepStatusRolledBack StepStatus = "rolled-back"
	StepStatusSkipped    StepStatus = "skipped"
	// StepStatusDegraded is emitted by a best-effort harness whose install failed
	// (C-32). It is distinct from StepStatusFailed (which aborts the pipeline) and
	// from StepStatusSucceeded (which is clean). A degraded step returns nil from
	// Run() so the pipeline continues, but it is NOT reported as a clean success.
	StepStatusDegraded StepStatus = "degraded"
)

// StepResult captures the outcome and timing of a single step.
type StepResult struct {
	StepID     string
	Status     StepStatus
	StartedAt  time.Time
	FinishedAt time.Time
	Err        error
}

// StageResult captures the outcome of running all steps in one stage.
type StageResult struct {
	Stage   Stage
	Steps   []StepResult
	Success bool
	Err     error
	// Unreverted lists side effects that were applied and that rollback could
	// not undo — a globally installed npm package, a brew formula, a PATH
	// entry. It is only ever populated for StageRollback.
	//
	// Success answers "did every Rollback() call return nil"; it does NOT
	// answer "is the machine back to how it was". Without this list a rollback
	// that reverted nothing at all was indistinguishable from a real one, and
	// the installer reported "Rollback: succeeded" while leaving a 19 MB binary
	// behind. A non-empty Unreverted with Success true means partial.
	Unreverted []string
}

// ExecutionResult is the top-level result returned by Orchestrator.Execute.
type ExecutionResult struct {
	Prepare  StageResult
	Apply    StageResult
	Rollback StageResult
	Err      error
}
