// Package contracttest provides vocabulary checks for the textual harness contracts.
// These checks delimit known contradictions; they do not prove agent behaviour.
package contracttest

import (
	"regexp"
	"strings"
	"testing"
)

func Reject(t *testing.T, body string, patterns ...string) {
	t.Helper()
	body = strings.ToLower(strings.Join(strings.Fields(strings.NewReplacer("*", "", "`", "").Replace(body)), " "))
	for _, pattern := range patterns {
		if found := regexp.MustCompile(pattern).FindString(body); found != "" {
			t.Errorf("contrato contradicho (%s): %s", pattern, found)
		}
	}
}

var Verify = []string{
	`ready (to|for) archive[.!]|proceed to archive`,
	`(skip|omit|remove|without|disable)[^.\n]{0,45}(wiring|comments/documentation)|(wiring|comments/documentation)[^.\n]{0,30}(omitted|optional|unnecessary)`,
	`(approve|approval|aprobado)[^.\n]{0,45}(despite|with|even if)[^.\n]{0,25}(unresolved|pending|open) critical`,
	`senior[^.\n]{0,35}(start|launch|run)[^.\n]{0,25}auditor[^.\n]{0,25}without it|(start|launch|run)[^.\n]{0,35}auditor[^.\n]{0,55}(without|before|parallel|alongside)[^.\n]{0,35}senior`,
	`(always|may|must|should|resume|recovery)[^.\n]{0,30}(reset|restart)[^.\n]{0,30}(counter|round count)`,
	`(start|allow|permit|begin|launch)[^.\n]{0,25}(round (4|four)|fourth round)`,
}
