package roledefs

import "embed"

// embeddedAssets is the bundled set of role-*.md resources. See doc.go for
// why there is no SKILL.md index inside it.
//
//go:embed assets/*.md
var embeddedAssets embed.FS
