package config_test

// ──────────────────────────────────────────────────────────────────────────────
// AR-2 (propose docs/propuestas/2026-10-23-fase-2-pasos-1-4.md, "Deuda
// aceptada"): "el enlace posterior a Compose no tiene módulo propietario ni
// dependencias declaradas; asignar propietario y declarar que config no
// depende de install/verify".
//
// El propietario asignado es internal/harness/roledefs (ver su doc.go): es
// el paquete que, DESPUÉS de que Compose devuelve texto puro, materializa en
// disco los role-*.md que ese texto nombra. La otra mitad de AR-2 —que
// config nunca importe install/verify— se declara ACÁ, verificada contra los
// imports reales y no sólo en un comentario: un comentario que dice "esto no
// depende de aquello" se desactualiza en silencio el día que alguien agrega
// el import que resuelve una urgencia; este test no.
// ──────────────────────────────────────────────────────────────────────────────

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// forbiddenDownstreamImports son los paquetes que, si internal/harness/config
// importara, invertirían la dirección de dependencia que AR-2 fija: config es
// upstream (lo componen otros), nunca downstream (compone a otros).
var forbiddenDownstreamImports = []string{
	"github.com/Group-Active-IA/active-stack/internal/install",
	"github.com/Group-Active-IA/active-stack/internal/uninstall",
	"github.com/Group-Active-IA/active-stack/internal/verify",
	"github.com/Group-Active-IA/active-stack/internal/harness/roledefs",
}

// TestConfigPackage_DoesNotImportInstallOrVerify escanea los imports REALES
// de los .go que no son de test en internal/harness/config (no sus
// subpaquetes, como permissions/, que tienen su propio contrato) y falla si
// alguno linkea hacia install, uninstall, verify o roledefs.
//
// No escanea los _test.go del propio paquete a propósito: un test puede
// necesitar ensamblar un escenario con install/verify sin que eso invierta
// la dependencia del PAQUETE (el binario de producción nunca enlaza
// *_test.go). Lo que AR-2 pide blindado es el código que se shipea.
func TestConfigPackage_DoesNotImportInstallOrVerify(t *testing.T) {
	dir := "." // este test corre con cwd = internal/harness/config (go test convention)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir(%q): %v", dir, err)
	}

	scanned := 0
	fset := token.NewFileSet()
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %q: %v", path, err)
		}
		scanned++
		for _, imp := range f.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			for _, forbidden := range forbiddenDownstreamImports {
				if importPath == forbidden {
					t.Errorf("%s imports %q — AR-2 fija que internal/harness/config nunca "+
						"depende de install/verify/roledefs; esa dirección sólo puede ir al "+
						"revés (ellos importan config, nunca al contrario)", path, importPath)
				}
			}
		}
	}
	if scanned == 0 {
		t.Fatal("no se escaneó ningún .go de producción — el test no prueba nada, revisar el filtro")
	}
}
