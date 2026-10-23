package main

import (
	"testing"
)

func TestVerifyInstalacionIncluyeBinario(t *testing.T) {
	home := t.TempDir()
	if err := buildInstallVerifyHook(nil, nil, home, home, true)(); err != nil {
		t.Fatalf("sin self-install: %v", err)
	}
	if err := buildInstallVerifyHook(nil, nil, home, home, false)(); err == nil {
		t.Fatal("binario ausente aceptado")
	}
}
