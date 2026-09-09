package pipeline

type Stage string

const (
	StagePrepare  Stage = "prepare"
	StageApply    Stage = "apply"
	StageRollback Stage = "rollback"
)

// Step is the minimal unit of work executed by the Runner.
type Step interface {
	ID() string
	Run() error
}

// RollbackStep extends Step with a Rollback operation.
// Apply steps that write user config MUST implement this interface.
type RollbackStep interface {
	Step
	Rollback() error
}

// UnrevertedReporter is an optional interface for rollback steps that can
// perform side effects outside the snapshot model — installing a global npm
// package, running brew, editing the user PATH. Such a step must declare what
// it left behind so the rollback is reported as partial instead of clean.
//
// It is called after Rollback(); returning nil means "nothing was left".
type UnrevertedReporter interface {
	Unreverted() []string
}

// FailurePolicy controls how the runner behaves when a step fails.
type FailurePolicy int

const (
	// StopOnError stops execution at the first failed step (default).
	StopOnError FailurePolicy = iota
	// ContinueOnError continues executing remaining steps, collecting all errors.
	ContinueOnError
)

// ProgressEvent is emitted by the runner as each step starts and completes.
type ProgressEvent struct {
	StepID string
	Stage  Stage
	Status StepStatus
	Err    error
}

// ProgressFunc is a callback invoked for every step lifecycle event.
type ProgressFunc func(ProgressEvent)

// StagePlan holds the ordered steps for each stage.
type StagePlan struct {
	Prepare []Step
	Apply   []Step
}
