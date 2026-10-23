package roledefs

import (
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/roles"
)

// TestDefaultResources_HasExactlySevenCanonicalRoles is RED first: antes de
// escribir ValidateResources, DefaultResources() todavía no existe.
func TestDefaultResources_HasExactlySevenCanonicalRoles(t *testing.T) {
	res, err := DefaultResources()
	if err != nil {
		t.Fatalf("DefaultResources: %v", err)
	}
	if len(res) != 7 {
		t.Fatalf("got %d resources, want 7: %+v", len(res), res)
	}

	want := map[string]bool{
		roles.RoleDBA + ".md":      true,
		roles.RoleBackend + ".md":  true,
		roles.RoleFrontend + ".md": true,
		roles.RoleDevOps + ".md":   true,
		roles.RoleQA + ".md":       true,
		roles.RoleArch + ".md":     true,
		roles.RoleAnalyst + ".md":  true,
	}
	got := map[string]bool{}
	for _, r := range res {
		got[r.Name] = true
		if len(r.Content) == 0 {
			t.Errorf("resource %q has empty content", r.Name)
		}
	}
	for name := range want {
		if !got[name] {
			t.Errorf("missing canonical resource %q", name)
		}
	}
	for name := range got {
		if !want[name] {
			t.Errorf("unexpected resource %q bundled (not one of the seven canonical roles)", name)
		}
	}
}

// TestValidateResources_RejectsSeniorAndAuditor triangula el gate "no importar
// archivos personales ni sumar senior/auditor por accidente": un conjunto de
// recursos que SÍ trae role-senior.md tiene que rechazarse, aunque también
// traiga los siete canónicos.
func TestValidateResources_RejectsSeniorAndAuditor(t *testing.T) {
	base, err := DefaultResources()
	if err != nil {
		t.Fatalf("DefaultResources: %v", err)
	}
	withSenior := append(append([]Resource{}, base...), Resource{
		Name:    "role-senior.md",
		Content: []byte("intruso"),
	})
	if err := ValidateResources(withSenior); err == nil {
		t.Fatal("ValidateResources aceptó un set que incluye role-senior.md, y no debería")
	}

	withAuditor := append(append([]Resource{}, base...), Resource{
		Name:    "role-auditor.md",
		Content: []byte("intruso"),
	})
	if err := ValidateResources(withAuditor); err == nil {
		t.Fatal("ValidateResources aceptó un set que incluye role-auditor.md, y no debería")
	}
}

// TestValidateResources_RejectsMissingRole triangula el otro lado: faltar uno
// de los siete también es inválido, no sólo sobrar uno.
func TestValidateResources_RejectsMissingRole(t *testing.T) {
	base, err := DefaultResources()
	if err != nil {
		t.Fatalf("DefaultResources: %v", err)
	}
	incomplete := base[:len(base)-1]
	if err := ValidateResources(incomplete); err == nil {
		t.Fatal("ValidateResources aceptó un set con sólo 6 roles, y no debería")
	}
}

// TestValidateResources_RejectsEmptyContent cubre el gate de "definición
// faltante/vacía" que el CLI de roles necesita poder distinguir más adelante.
func TestValidateResources_RejectsEmptyContent(t *testing.T) {
	base, err := DefaultResources()
	if err != nil {
		t.Fatalf("DefaultResources: %v", err)
	}
	empty := append([]Resource{}, base...)
	empty[0].Content = nil
	if err := ValidateResources(empty); err == nil {
		t.Fatalf("ValidateResources aceptó %q vacío, y no debería", empty[0].Name)
	}
}

// TestDefaultResources_PassesValidation confirma que el set real embebido
// pasa su propio gate — si algún día se agrega un octavo archivo a assets/
// sin querer, esto rompe antes que cualquier instalación real.
func TestDefaultResources_PassesValidation(t *testing.T) {
	res, err := DefaultResources()
	if err != nil {
		t.Fatalf("DefaultResources: %v", err)
	}
	if err := ValidateResources(res); err != nil {
		t.Fatalf("el set embebido real no pasa su propio gate: %v", err)
	}
}
