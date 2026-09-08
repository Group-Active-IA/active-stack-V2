package verify

import "testing"

// TestGlobsIntersect pins the predicate that decides whether a later
// permission.bash entry can override a floor deny rule.
//
// The question is NOT "does the later pattern match the literal text of the
// floor rule". Both sides are globs, so the question is whether the two
// patterns have any command in common — an intersection of two languages.
// Until 2026-09-08 this was a pattern-vs-text match, and it answered "no" for
// the most likely way to break the floor by hand: a narrower rule appended at
// the end of the file. globsIntersect("rm -rf /home/*", "rm -rf /*") was false,
// even though `rm -rf /home/proyectos` matches both.
//
// Every case is asserted in BOTH directions: intersection is symmetric, and a
// predicate that only works with the floor on one side would be a trap for the
// next caller.
func TestGlobsIntersect(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want bool
	}{
		// ── The regression that motivated this: partial overlap ──
		{"narrower suffix overlaps broader", "rm -rf /home/*", "rm -rf /*", true},
		{"broader wildcard overlaps everything", "*", "rm -rf /*", true},
		{"identical patterns", "rm -rf /*", "rm -rf /*", true},

		// ── No overlap: these must NOT be reported as shadowing ──
		{"disjoint literals", "git status", "rm -rf /*", false},
		{"common prefix then divergence", "rm -rf /home/*", "rm -rf /etc/*", false},
		{"literal is not a prefix of the glob", "rm", "rm -rf /*", false},

		// ── Wildcards on both sides ──
		{"same suffix after star", "*.env", "*.env", true},
		{"different suffix after star", "*x", "*y", false},
		{"star matches the empty string", "*", "", true},
		{"empty vs literal", "", "git status", false},
		{"two empties", "", "", true},
		{"star absorbs a whole literal", "git *", "git status", true},
		{"star cannot bridge a mismatch", "git *", "svn status", false},
		{"multiple stars", "rm*-rf*", "rm -rf /home", true},
		{"trailing star vs exact prefix", "git status*", "git status", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := globsIntersect(tt.a, tt.b); got != tt.want {
				t.Errorf("globsIntersect(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
			if got := globsIntersect(tt.b, tt.a); got != tt.want {
				t.Errorf("globsIntersect(%q, %q) = %v, want %v (predicate must be symmetric)", tt.b, tt.a, got, tt.want)
			}
		})
	}
}
