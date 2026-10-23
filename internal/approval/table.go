package approval

import "strings"

// approvalRow es la última fila de datos de la tabla "## Aprobación":
// Elección | Persona | Fecha | Alcance.
type approvalRow struct {
	Eleccion string
	Persona  string
	Fecha    string
	Alcance  string
}

// complete dice si las cuatro columnas tienen contenido no vacío.
func (r approvalRow) complete() bool {
	return r.Eleccion != "" && r.Persona != "" && r.Fecha != "" && r.Alcance != ""
}

// missingFields devuelve los nombres de columna vacíos, en orden fijo.
func (r approvalRow) missingFields() []string {
	var out []string
	if r.Eleccion == "" {
		out = append(out, "Elección")
	}
	if r.Persona == "" {
		out = append(out, "Persona")
	}
	if r.Fecha == "" {
		out = append(out, "Fecha")
	}
	if r.Alcance == "" {
		out = append(out, "Alcance")
	}
	return out
}

// parseApprovalTable busca el encabezado "## Aprobación" en content y parsea
// su tabla markdown. Devuelve found=false si no hay sección o la tabla no
// tiene ninguna fila de datos (tabla vacía: MISSING, no INVALID).
//
// Sólo se usa la ÚLTIMA fila de datos: es la decisión vigente si alguna vez se
// agrega una corrección a la tabla en el mismo commit documental.
func parseApprovalTable(content string) (row approvalRow, found bool) {
	lines := strings.Split(content, "\n")
	headingIdx := -1
	for i, l := range lines {
		if strings.TrimSpace(l) == "## Aprobación" {
			headingIdx = i
			break
		}
	}
	if headingIdx == -1 {
		return approvalRow{}, false
	}

	var headerCols []string
	var rows [][]string
	for i := headingIdx + 1; i < len(lines); i++ {
		l := strings.TrimSpace(lines[i])
		if strings.HasPrefix(l, "## ") {
			break // siguiente sección
		}
		if !strings.HasPrefix(l, "|") {
			continue
		}
		cells := splitTableRow(l)
		if headerCols == nil {
			headerCols = cells
			continue
		}
		if isSeparatorRow(cells) {
			continue
		}
		rows = append(rows, cells)
	}

	if len(rows) == 0 {
		return approvalRow{}, false
	}

	last := rows[len(rows)-1]
	idx := columnIndex(headerCols)
	get := func(name string) string {
		i, ok := idx[name]
		if !ok || i >= len(last) {
			return ""
		}
		return strings.TrimSpace(last[i])
	}
	row = approvalRow{
		Eleccion: get("elección"),
		Persona:  get("persona"),
		Fecha:    get("fecha"),
		Alcance:  get("alcance"),
	}
	return row, true
}

func columnIndex(header []string) map[string]int {
	idx := make(map[string]int, len(header))
	for i, h := range header {
		idx[strings.ToLower(strings.TrimSpace(h))] = i
	}
	return idx
}

func splitTableRow(line string) []string {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "|")
	line = strings.TrimSuffix(line, "|")
	parts := strings.Split(line, "|")
	out := make([]string, len(parts))
	for i, p := range parts {
		out[i] = strings.TrimSpace(p)
	}
	return out
}

// isSeparatorRow detecta la fila "|---|---|---|" que separa encabezado de datos.
func isSeparatorRow(cells []string) bool {
	for _, c := range cells {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		if strings.Trim(c, "-: ") != "" {
			return false
		}
	}
	return true
}
