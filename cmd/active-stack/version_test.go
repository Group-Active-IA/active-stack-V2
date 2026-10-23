package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestVersionSubcommand(t *testing.T) {
	if os.Getenv("ACTIVE_STACK_VERSION_HELPER") == "1" {
		os.Args = []string{"active-stack", "version"}
		main()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestVersionSubcommand$")
	cmd.Env = append(os.Environ(), "ACTIVE_STACK_VERSION_HELPER=1")
	out, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(out), "active-stack ") {
		t.Fatalf("version debe salir 0 e identificar el ejecutable: %v, %s", err, out)
	}
}
