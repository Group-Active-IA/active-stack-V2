package external

import (
	"context"
	"os/exec"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/model"
	"github.com/Group-Active-IA/active-stack/internal/system"
)

// ──────────────────────────────────────────────────────────────────────────────
// El instalador npm no miraba si el binario ya estaba
//
// El 2026-09-09, instalando en una maquina que ya tenia openspec 1.3.1 via nvm,
// el harness openspec fallo con "sudo: npm: orden no encontrada" y volteo la
// instalacion entera por rollback. La causa: installNPM corre `npm install -g`
// sin preguntar antes si el binario existe.
//
// Y el sudo lo empeora: nvm vive en el HOME del usuario, asi que `npm` esta en
// el PATH interactivo y NO en el de sudo. Una maquina que YA tiene la
// herramienta funcionando es exactamente donde el instalador se cae.
//
// homebrew.go y el metodo download ya devuelven AlreadyInstalled; npm era el
// unico de los tres que no preguntaba.
// ──────────────────────────────────────────────────────────────────────────────

// TestInstallNPM_SkipsWhenBinaryIsAlreadyOnPath asserts npm is never invoked
// when the package's binary already resolves on PATH.
func TestInstallNPM_SkipsWhenBinaryIsAlreadyOnPath(t *testing.T) {
	origRunner, origLookPath := runner, lookPath
	defer func() { runner, lookPath = origRunner, origLookPath }()

	called := false
	runner = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		called = true
		return nil, nil
	}
	lookPath = func(file string) (string, error) {
		return "/home/user/.nvm/versions/node/v22/bin/" + file, nil
	}

	h := model.Harness{ID: "openspec", External: &model.External{Method: "npm", Pkg: "@fission-ai/openspec"}}
	res, err := installNPM(context.Background(), h, system.PlatformProfile{OS: "linux"})
	if err != nil {
		t.Fatalf("installNPM() con el binario ya presente: %v", err)
	}
	if called {
		t.Error("corrio npm install aunque el binario ya estaba en el PATH: " +
			"en una maquina con la herramienta instalada bajo nvm, ese comando falla bajo sudo " +
			"y voltea la instalacion entera")
	}
	if !res.AlreadyInstalled {
		t.Error("AlreadyInstalled = false: el reporte no distingue 'lo instale' de 'ya estaba'")
	}
	if res.BinaryPath == "" {
		t.Error("BinaryPath vacio: el path resuelto se pierde")
	}
}

// TestInstallNPM_InstallsWhenBinaryIsMissing asserts the skip does not swallow
// the real install path.
func TestInstallNPM_InstallsWhenBinaryIsMissing(t *testing.T) {
	origRunner, origLookPath := runner, lookPath
	defer func() { runner, lookPath = origRunner, origLookPath }()

	called := false
	runner = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		called = true
		return nil, nil
	}
	lookPath = func(file string) (string, error) { return "", exec.ErrNotFound }

	h := model.Harness{ID: "openspec", External: &model.External{Method: "npm", Pkg: "@fission-ai/openspec"}}
	res, err := installNPM(context.Background(), h, system.PlatformProfile{OS: "linux", NpmWritable: true})
	if err != nil {
		t.Fatalf("installNPM(): %v", err)
	}
	if !called {
		t.Error("no corrio npm install con el binario ausente")
	}
	if res.AlreadyInstalled {
		t.Error("AlreadyInstalled = true tras instalarlo de verdad")
	}
}
