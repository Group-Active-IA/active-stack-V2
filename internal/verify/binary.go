package verify

import (
	"context"
	"debug/buildinfo"
	"fmt"
	"os/exec"
	"time"
)

var readBinaryInfo = buildinfo.ReadFile
var runBinaryHelp = func(ctx context.Context, path string, args ...string) error {
	out, err := exec.CommandContext(ctx, path, args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("roles --help: %w: %s", err, out)
	}
	return nil
}

// BinaryChecks valida identidad antes de ejecutar el binario instalado.
func BinaryChecks(path string) []Check {
	return []Check{{ID: "active-stack:binary", Description: "binario de cmd/active-stack y roles --help con salida 0", Run: func(ctx context.Context) error {
		info, err := readBinaryInfo(path)
		if err != nil {
			return fmt.Errorf("leer binario instalado: %w", err)
		}
		if info.Path != "github.com/Group-Active-IA/active-stack/cmd/active-stack" {
			return fmt.Errorf("binario instalado de %q, se esperaba cmd/active-stack", info.Path)
		}
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		return runBinaryHelp(ctx, path, "roles", "--help")
	}}}
}
