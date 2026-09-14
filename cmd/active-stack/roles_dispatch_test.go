package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// TestRolesConArchivosExplicitosNombraLosDosLados.
func TestRolesConArchivosExplicitosNombraLosDosLados(t *testing.T) {
	var out bytes.Buffer
	code := runRolesDispatch([]string{"--files", "migrations/003_tenant.sql"}, t.TempDir(), &out)
	if code != 0 {
		t.Fatalf("exit = %d, salida:\n%s", code, out.String())
	}

	s := out.String()
	if !strings.Contains(s, "role-dba") || !strings.Contains(s, "migrations/003_tenant.sql") {
		t.Errorf("no convoca al DBA con el archivo que lo justifica:\n%s", s)
	}
	if !strings.Contains(s, "ningún archivo bajo el árbol de UI") {
		t.Errorf("el rol que NO corresponde salió sin motivo, y esa es la mitad del "+
			"valor del comando:\n%s", s)
	}
}

// TestRolesJSONEsConsumibleSinParsearTexto.
func TestRolesJSONEsConsumibleSinParsearTexto(t *testing.T) {
	var out bytes.Buffer
	code := runRolesDispatch([]string{"--json", "--files", "api/auth.py"}, t.TempDir(), &out)
	if code != 0 {
		t.Fatalf("exit = %d, salida:\n%s", code, out.String())
	}

	var got struct {
		Floor struct {
			Triggered bool   `json:"triggered"`
			File      string `json:"file"`
		} `json:"floor"`
		Roles []struct {
			Role     string `json:"role"`
			Convened bool   `json:"convened"`
		} `json:"roles"`
	}
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("--json no emitió JSON: %v\n%s", err, out.String())
	}
	if !got.Floor.Triggered || got.Floor.File != "api/auth.py" {
		t.Errorf("el JSON no trae el piso con su archivo: %+v", got.Floor)
	}
	if len(got.Roles) != 7 {
		t.Errorf("el JSON trae %d roles y son 7", len(got.Roles))
	}
}

// TestRolesSinRepoYSinArchivosFalla: no se inventa un diff vacío.
func TestRolesSinRepoYSinArchivosFalla(t *testing.T) {
	var out bytes.Buffer
	code := runRolesDispatch(nil, t.TempDir(), &out)
	if code == 0 {
		t.Fatalf("fuera de un repo y sin --files salió 0. Un diff vacío inventado se "+
			"lee igual que 'nada que revisar':\n%s", out.String())
	}
	if !strings.Contains(out.String(), "--files") {
		t.Errorf("el error no dice cómo salir de la situación:\n%s", out.String())
	}
}

// TestRolesFlagDesconocidaFalla.
func TestRolesFlagDesconocidaFalla(t *testing.T) {
	var out bytes.Buffer
	code := runRolesDispatch([]string{"--todos"}, t.TempDir(), &out)
	if code == 0 {
		t.Fatal("una flag desconocida salió 0")
	}
	for _, flag := range []string{"--base", "--files", "--json"} {
		if !strings.Contains(out.String(), flag) {
			t.Errorf("el usage no nombra %s:\n%s", flag, out.String())
		}
	}
}

// TestRolesBaseSinValorFalla.
func TestRolesBaseSinValorFalla(t *testing.T) {
	var out bytes.Buffer
	if code := runRolesDispatch([]string{"--base"}, t.TempDir(), &out); code == 0 {
		t.Fatalf("--base sin referencia salió 0:\n%s", out.String())
	}
}

// TestRolesHelpNoEsUnError: pedir ayuda no es equivocarse.
func TestRolesHelpNoEsUnError(t *testing.T) {
	var out bytes.Buffer
	code := runRolesDispatch([]string{"--help"}, t.TempDir(), &out)
	if code != 0 {
		t.Errorf("--help salió %d; pedir ayuda no es un error de uso", code)
	}
	if strings.Contains(out.String(), "error:") {
		t.Errorf("--help se reporta como error:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "--files") {
		t.Errorf("--help no muestra el usage:\n%s", out.String())
	}
}
