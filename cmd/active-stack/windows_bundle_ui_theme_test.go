package main

import (
	"path/filepath"
	"testing"
)

// TestWindowsBundleThemeContainsMinimalInstallerPages asserts the WiX Burn theme
// ships every page the installer navigates to, with its user-facing copy.
//
// The path is built with filepath.Join, NOT with hardcoded backslashes.
// On 2026-09-08 this test read "..\\..\\installer\\windows\\theme\\..." and had
// been failing on Linux with "no such file or directory": there the whole thing
// is one literal filename, so os.ReadFile never found it. It passed on Windows
// and nobody noticed, because the only CI workflow triggers on release tags and
// never runs `go test`.
//
// A test that fails on setup instead of on its assert proves nothing about the
// theme. Its three sibling tests in this package already used filepath.Join;
// this one was the outlier.
func TestWindowsBundleThemeContainsMinimalInstallerPages(t *testing.T) {
	root := filepath.Join("..", "..")
	themeFile := filepath.Join(root, "installer", "windows", "theme", "ActiveStackTheme.xml")

	assertFileContains(t, themeFile, "<Page Name=\"Install\">")
	assertFileContains(t, themeFile, "<Page Name=\"Progress\">")
	assertFileContains(t, themeFile, "<Page Name=\"Success\">")
	assertFileContains(t, themeFile, "<Page Name=\"Modify\">")
	assertFileContains(t, themeFile, "Set up your AI coding workspace")
	assertFileContains(t, themeFile, "Install Active Stack")
	assertFileContains(t, themeFile, "Preparing your setup")
	assertFileContains(t, themeFile, "Active Stack is ready")
}
