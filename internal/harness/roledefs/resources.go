package roledefs

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Group-Active-IA/active-stack/internal/roles"
)

// Resource is one role-*.md file to materialize on disk.
type Resource struct {
	Name    string // e.g. "role-dba.md" — exactly the filename, no path.
	Content []byte
}

// canonicalRoleIDs combines the seven specialist criteria (classified by
// internal/roles) and the five sequential workflow responsibilities. The
// classifier remains independent: installing a role does not require its attendance.
var canonicalRoleIDs = append([]string{
	roles.RoleDBA, roles.RoleBackend, roles.RoleFrontend, roles.RoleDevOps,
	roles.RoleQA, roles.RoleArch, roles.RoleAnalyst,
}, workflowRoleIDs...)

// workflowRoleIDs are the five sequential responsibilities of the lean flow.
// Unlike the specialists, the flow cannot run without the STACK's version of
// each one: a personal role-auditor.md with another contract would be
// convened as if it were ours. Install reports a conflict on any of them as
// blocking, and Health flags one whose content is not the stack's (5.1).
var workflowRoleIDs = []string{
	"role-planner", "role-plan-reviewer", "role-constructor", "role-breaker", "role-auditor",
}

// IsWorkflowRole reports whether resource name (e.g. "role-auditor.md") is one
// of the five sequential workflow roles.
func IsWorkflowRole(name string) bool {
	id := strings.TrimSuffix(name, ".md")
	for _, w := range workflowRoleIDs {
		if w == id {
			return true
		}
	}
	return false
}

// Senior is a legacy optional procedure, not one of the installed workflow roles.
var forbiddenRoleIDs = map[string]bool{"role-senior": true}

// DefaultResources returns the production resource set, read from the
// embedded assets. It is the only place that touches embeddedAssets; Install
// and Uninstall always go through a []Resource parameter so tests can inject
// a different set without touching the filesystem or go:embed.
func DefaultResources() ([]Resource, error) {
	entries, err := embeddedAssets.ReadDir("assets")
	if err != nil {
		return nil, fmt.Errorf("roledefs: read embedded assets: %w", err)
	}
	out := make([]Resource, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		content, err := embeddedAssets.ReadFile("assets/" + e.Name())
		if err != nil {
			return nil, fmt.Errorf("roledefs: read embedded asset %q: %w", e.Name(), err)
		}
		out = append(out, Resource{Name: e.Name(), Content: content})
	}
	return out, nil
}

// ValidateResources is the gate the propose requires BEFORE bundling:
// exactly the twelve canonical roles, none empty, none forbidden, no stray
// file. It is checked both against the embedded production set (
// TestDefaultResources_PassesValidation) and against whatever Install/
// Uninstall receive at runtime — a malformed catalog or a bad override in a
// test must fail here, not halfway through writing files.
func ValidateResources(res []Resource) error {
	seen := make(map[string]Resource, len(res))
	var dupes []string
	for _, r := range res {
		name := strings.TrimSuffix(r.Name, ".md")
		if forbiddenRoleIDs[name] {
			return fmt.Errorf("roledefs: %q is not one of the twelve canonical roles — "+
				"it must never be bundled by this harness", r.Name)
		}
		if _, dup := seen[r.Name]; dup {
			dupes = append(dupes, r.Name)
		}
		seen[r.Name] = r
	}
	if len(dupes) > 0 {
		sort.Strings(dupes)
		return fmt.Errorf("roledefs: duplicate resource name(s): %s", strings.Join(dupes, ", "))
	}

	var missing, empty []string
	for _, id := range canonicalRoleIDs {
		name := id + ".md"
		r, ok := seen[name]
		if !ok {
			missing = append(missing, name)
			continue
		}
		if len(r.Content) == 0 {
			empty = append(empty, name)
		}
		delete(seen, name)
	}
	if len(missing) > 0 {
		return fmt.Errorf("roledefs: missing canonical role definition(s): %s", strings.Join(missing, ", "))
	}
	if len(empty) > 0 {
		return fmt.Errorf("roledefs: empty role definition(s): %s", strings.Join(empty, ", "))
	}
	if len(seen) > 0 {
		var extra []string
		for name := range seen {
			extra = append(extra, name)
		}
		sort.Strings(extra)
		return fmt.Errorf("roledefs: unexpected resource(s) beyond the twelve canonical roles: %s", strings.Join(extra, ", "))
	}
	return nil
}
