package roledefs

import (
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/roles"
)

// Default resources preserve specialists and add the five workflow roles.
func TestDefaultResources_HasExactlyTwelveCanonicalRoles(t *testing.T) {
	res, err := DefaultResources()
	if err != nil {
		t.Fatalf("DefaultResources: %v", err)
	}
	if len(res) != 12 {
		t.Fatalf("got %d resources, want 12", len(res))
	}

	want := map[string]bool{
		"role-planner.md":          true,
		"role-plan-reviewer.md":    true,
		"role-constructor.md":      true,
		"role-breaker.md":          true,
		"role-auditor.md":          true,
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
			t.Errorf("unexpected resource %q bundled (not one of the twelve canonical roles)", name)
		}
	}
}

// Senior remains outside the installed workflow; duplicate auditor is invalid.
func TestValidateResources_RejectsSeniorAndDuplicateAuditor(t *testing.T) {
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
		t.Fatal("ValidateResources aceptó un set que duplica role-auditor.md, y no debería")
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
		t.Fatal("ValidateResources aceptó un set con un rol faltante, y no debería")
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
// pasa su propio gate — si algún día se agrega un archivo inesperado a assets/
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

func TestDefaultResources_NoFixedModel(t *testing.T) {
	res, err := DefaultResources()
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range res {
		for _, line := range strings.Split(string(r.Content), "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), "model:") {
				t.Errorf("%s fixes model instead of caller selection", r.Name)
			}
		}
	}
}

func TestWorkflowResources_Contracts(t *testing.T) {
	res, err := DefaultResources()
	if err != nil {
		t.Fatal(err)
	}
	contracts := map[string][]string{
		"role-planner.md":       {"proposal.md", "design.md", "tasks.md", "LOW", "CRITICAL"},
		"role-plan-reviewer.md": {"escenarios", "viabilidad", "código inexistente", "No editás"},
		"role-constructor.md":   {"TDD", "RED", "GREEN", "antes", "No hay código"},
		"role-breaker.md":       {"Correctitud", "Cableado", "Documentación", "reproduc", "No editás"},
		"role-auditor.md":       {"APROBADO", "BLOQUEANTES", "código", "pruebas", "informes", "No editás", "Reproducí", "Ejecutá comprobaciones"},
	}
	found := map[string]bool{}
	for _, r := range res {
		expected, ok := contracts[r.Name]
		if !ok {
			continue
		}
		found[r.Name] = true
		for _, text := range expected {
			if !strings.Contains(string(r.Content), text) {
				t.Errorf("%s missing contract %q", r.Name, text)
			}
		}
		for _, text := range []string{"EVIDENCIA", "JUICIO", "SUPOSICIÓN", "Fuentes", "contaminación"} {
			if !strings.Contains(string(r.Content), text) {
				t.Errorf("%s missing reporting contract %q", r.Name, text)
			}
		}
	}
	for name := range contracts {
		if !found[name] {
			t.Errorf("missing workflow role %s", name)
		}
	}
}

// roleDescription returns the frontmatter description: Claude picks subagents
// by it, so it carries the routing contract.
func roleDescription(content string) string {
	start := strings.Index(content, "description:")
	end := strings.Index(content, "\ntools:")
	if start < 0 || end < start {
		return ""
	}
	return content[start:end]
}

func TestWorkflowResources_CurrentLeanContracts(t *testing.T) {
	res, err := DefaultResources()
	if err != nil {
		t.Fatal(err)
	}
	contracts := map[string][]string{
		"role-planner.md":       {"En LOW el orquestador puede"},
		"role-plan-reviewer.md": {"En LOW el orquestador puede"},
		"role-constructor.md":   {"lecciones numeradas (L-NNN)", "sólo la lista del auditor", "agrupadas por archivo"},
		"role-breaker.md":       {"lecciones numeradas (L-NNN)", "En LOW, una sola ronda"},
		"role-auditor.md":       {"lecciones numeradas (L-NNN)", "openspec/changes/archive/", "Engram", "qué fuentes usaste", "reescribe algo que ya existe", "recorrido real"},
	}
	phases := map[string]string{
		"role-planner.md": "propose", "role-plan-reviewer.md": "propose",
		"role-constructor.md": "apply", "role-breaker.md": "verify", "role-auditor.md": "verify",
	}
	for _, r := range res {
		content := string(r.Content)
		for _, text := range contracts[r.Name] {
			if !strings.Contains(content, text) {
				t.Errorf("%s missing lean contract %q", r.Name, text)
			}
		}
		if phase, ok := phases[r.Name]; ok && !strings.Contains(roleDescription(content), phase) {
			t.Errorf("%s description does not name its phase %q", r.Name, phase)
		}
	}
}

func TestSpecialistResources_OptionalNotMandatory(t *testing.T) {
	res, err := DefaultResources()
	if err != nil {
		t.Fatal(err)
	}
	specialists := map[string]bool{
		"role-analyst.md": true, "role-architect.md": true, "role-backend.md": true, "role-dba.md": true,
		"role-devops.md": true, "role-frontend.md": true, "role-qa.md": true,
	}
	seen := 0
	for _, r := range res {
		if !specialists[r.Name] {
			continue
		}
		seen++
		desc := roleDescription(string(r.Content))
		if !strings.Contains(desc, "Opcional") {
			t.Errorf("%s description does not say it is optional", r.Name)
		}
		for _, mandatory := range []string{"siempre", "Obligatorio", "obligatorio", "sin excepción"} {
			if strings.Contains(desc, mandatory) {
				t.Errorf("%s description still forces attendance: %q", r.Name, mandatory)
			}
		}
	}
	if seen != len(specialists) {
		t.Fatalf("saw %d specialists, want %d", seen, len(specialists))
	}
}
