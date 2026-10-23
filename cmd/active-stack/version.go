package main

import (
	"fmt"
	"runtime/debug"
)

func versionText() string {
	version, revision := "devel", "desconocido"
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" {
			version = info.Main.Version
		}
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				revision = setting.Value
			}
		}
	}
	return fmt.Sprintf("active-stack %s (SHA %s)", version, revision)
}
