package roledefs

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/Group-Active-IA/active-stack/internal/model"
)

// HealthIssue reports one canonical role definition that is missing or
// empty for one agent — the "distingue definición faltante/vacía" the
// propose asks the roles CLI to surface, factored out here so both the
// verify health check ("reconstruye Compose+enlace y compara referencias y
// recursos") and the CLI can call the same logic instead of drifting apart.
type HealthIssue struct {
	Agent    model.Agent
	Resource string // e.g. "role-dba.md"
	Problem  string // "missing" | "empty" | "foreign" (workflow role only)
	Path     string
}

// Health re-derives what should exist (resources — in production always
// DefaultResources(), the twelve canonical roles) and reports every one that
// is not actually there, or is there with zero bytes, per adapter.
// A WORKFLOW role (IsWorkflowRole) is also checked for content: present but
// not byte-identical to the stack's version is "foreign" — a personal or
// edited file the lean flow would convene with someone else's contract (5.1).
// Specialists are not content-checked: a personal one is legitimate.
// Adapters with an empty AgentsDir are skipped — not an issue, "not
// applicable".
func Health(adapters []AgentAdapter, homeDir string, resources []Resource) ([]HealthIssue, error) {
	var issues []HealthIssue
	for _, adapter := range adapters {
		dir := adapter.AgentsDir(homeDir)
		if dir == "" {
			continue
		}
		for _, res := range resources {
			path := filepath.Join(dir, res.Name)
			info, err := os.Stat(path)
			switch {
			case os.IsNotExist(err):
				issues = append(issues, HealthIssue{Agent: adapter.Agent(), Resource: res.Name, Problem: "missing", Path: path})
			case err != nil:
				return issues, err
			case info.Size() == 0:
				issues = append(issues, HealthIssue{Agent: adapter.Agent(), Resource: res.Name, Problem: "empty", Path: path})
			case IsWorkflowRole(res.Name):
				data, rerr := os.ReadFile(path)
				if rerr != nil {
					return issues, rerr
				}
				if !bytes.Equal(data, res.Content) {
					issues = append(issues, HealthIssue{Agent: adapter.Agent(), Resource: res.Name, Problem: "foreign", Path: path})
				}
			}
		}
	}
	return issues, nil
}
