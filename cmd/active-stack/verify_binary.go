package main

import (
	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/system"
	"github.com/Group-Active-IA/active-stack/internal/verify"
	"path/filepath"
	"runtime"
)

func buildInstallVerifyHook(h []model.Harness, adapters []verify.Adapter, base, home string, noSelfInstall bool) func() error {
	var extra []verify.Check
	if !noSelfInstall {
		name := "active-stack"
		if runtime.GOOS == "windows" {
			name += ".exe"
		}
		extra = verify.BinaryChecks(filepath.Join(system.BinaryInstallDirForHome(runtime.GOOS, home), name))
	}
	return verify.BuildHook(h, adapters, base, extra...)
}
