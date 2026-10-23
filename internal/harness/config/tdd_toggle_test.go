package config_test

import (
	"github.com/Group-Active-IA/active-stack/internal/harness/config"
	"strings"
	"testing"
)

func TestConstructorContratoApply(t *testing.T) {
	block, err := config.Compose([]string{"tdd"}, "claude")
	if err != nil {
		t.Fatal(err)
	}
	_, body, ok := strings.Cut(block, "# Strict TDD Module")
	if !ok {
		t.Fatal("falta módulo TDD")
	}
	body = strings.Join(strings.Fields(body), " ")
	for _, want := range []string{
		"The constructor receives the full apply and archive instructions in its prompt; phase 1 does not move them out of this block.",
		"Before acting, read the numbered lessons when a knowledge base is available and record the path read or its absence in the apply report.",
		"The brief names consumers outside the assigned area, including other modules, infrastructure and scripts (L-023).",
		"Resolve wording or ordering conflicts and declare the resolution as your own (L-030).",
		"A new rule after signature requires its own human approval before construction; stop rather than infer permission from the brief (L-029).",
		"For each task, paste RED and GREEN under that task in tasks.md: command, date +%T and the last output lines.",
		"Mark [x] only after both executed outputs are pasted; a written test alone is not RED.",
		"If a test was born green, make a temporary 2-4 line mutant that fails at the assertion the test claims to prove (L-019, L-034).",
		"Paste the mutant diff, command, output and failing file:line; a setup or compilation failure does not prove that assertion.",
		"Revert the mutant and rerun green; never commit the mutant.",
		"Without a prior commit, copy the original file before mutating and declare that backup (L-013).",
		"Run closing suites after the last commit and record git rev-parse HEAD beside command, time and output (L-031).",
		"Any later code change invalidates that closing run, even if an earlier run was green.",
		"Green means all required workflow steps, including formatting, lint, types and tests; a test-only run is not complete CI (L-005).",
		"The person adding a file runs the applicable workflow checks; report unavailable or unauthorized checks explicitly.",
		"In verify rounds 2 and 3, the constructor fixes only the auditor's list; new scope goes to debt.",
		"When another provider is available, the constructor may run in another model family, for example codex exec.",
		"If no other provider is available, say so in the report and do not claim independence.",
		"Declare the model/provider and shared family in each report; without a dispatch record, label identity as ASSUMPTION.",
		"The breaker and auditor count the tests themselves; do not accept the constructor's count as their evidence.",
		"Declare which sources were read, whether the auditor dictated exact cases and any contamination; dictated-case agreement is not independent convergence.",
		"With an external identity provider, request a real walkthrough with the person at the keyboard before archive; record unvisited paths as explicit debt and a limitation for the auditor, never as validated.",
		"Preserve the complete apply report in openspec/changes/<change>/informes/ with phase and delivery in its name, separate from verify reports.",
		"Index author, declared model/provider, SHA, provenance and missing reports in LEEME.md; preserve the full return or textual transcript and declare orchestration mediation.",
		"If a complete report cannot be recovered, report it missing and stop closure.",
	} {
		t.Run(want, func(t *testing.T) {
			if !strings.Contains(body, want) {
				t.Errorf("falta contrato del constructor: %s", want)
			}
		})
	}
}
