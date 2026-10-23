package approval

import "testing"

func TestParseApprovalTable_LastRowWins(t *testing.T) {
	content := "## Aprobación\n\n" +
		"| Elección | Persona | Fecha | Alcance |\n" +
		"|---|---|---|---|\n" +
		"| vieja | Juani | 01/01/2026 | alcance viejo |\n" +
		"| nueva | Juani | 23/10/2026 | alcance nuevo |\n"

	row, found := parseApprovalTable(content)
	if !found {
		t.Fatal("found = false, la tabla tiene filas")
	}
	if row.Alcance != "alcance nuevo" {
		t.Errorf("Alcance = %q, quería la ÚLTIMA fila (%q)", row.Alcance, "alcance nuevo")
	}
}

func TestParseApprovalTable_NoHeading_NotFound(t *testing.T) {
	_, found := parseApprovalTable("# Propose\n\nsin la sección.\n")
	if found {
		t.Fatal("found = true sin heading '## Aprobación'")
	}
}

func TestParseApprovalTable_StopsAtNextHeading(t *testing.T) {
	content := "## Aprobación\n\n" +
		"| Elección | Persona | Fecha | Alcance |\n" +
		"|---|---|---|---|\n" +
		"| x | Juani | 23/10/2026 | y |\n\n" +
		"## Otra sección\n\n" +
		"| no es tabla de aprobación |\n"

	row, found := parseApprovalTable(content)
	if !found || !row.complete() {
		t.Fatalf("no parseó la fila antes del siguiente heading: found=%v row=%+v", found, row)
	}
}
