package roles

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestSalidaLegibleNombraLosDosLados: los convocados con su hecho del diff, y
// los que no corresponden con su motivo. Lo segundo es la mitad del valor del
// comando: hoy un rol no convocado no aparece en ningún lado, y esa ausencia es
// indistinguible de un olvido.
func TestSalidaLegibleNombraLosDosLados(t *testing.T) {
	out := RenderText(Classify([]string{"migrations/003_tenant.sql", "README.md"}))

	if !strings.Contains(out, "role-dba") || !strings.Contains(out, "migrations/003_tenant.sql") {
		t.Errorf("la salida no convoca al DBA con el archivo que lo justifica:\n%s", out)
	}
	if !strings.Contains(out, "role-frontend") ||
		!strings.Contains(out, "ningún archivo bajo el árbol de UI") {
		t.Errorf("la salida no nombra al frontend con su motivo de ausencia:\n%s", out)
	}
	if !strings.Contains(out, "role-devops") ||
		!strings.Contains(out, "deploy y configuración sin cambios") {
		t.Errorf("la salida no nombra a DevOps con su motivo de ausencia:\n%s", out)
	}
	if !strings.Contains(out, "role-qa") {
		t.Errorf("la salida no nombra a QA:\n%s", out)
	}
}

// TestLosConvocadosSalenPrimero: el orden es fijo y los convocados van arriba.
func TestLosConvocadosSalenPrimero(t *testing.T) {
	out := RenderText(Classify([]string{"migrations/003_tenant.sql"}))
	iDBA := strings.Index(out, "role-dba")
	iFront := strings.Index(out, "role-frontend")
	if iDBA < 0 || iFront < 0 || iDBA > iFront {
		t.Errorf("los convocados no salen antes que los que no corresponden:\n%s", out)
	}
}

// TestElPisoSeDiceExplicitamenteConSuArchivo.
func TestElPisoSeDiceExplicitamenteConSuArchivo(t *testing.T) {
	out := RenderText(Classify([]string{"api/auth.py"}))

	if !strings.Contains(strings.ToLower(out), "piso") {
		t.Fatalf("un diff que toca autenticación no dice que activó el piso:\n%s", out)
	}
	if !strings.Contains(out, "api/auth.py") {
		t.Errorf("el piso se anuncia sin el archivo que lo disparó:\n%s", out)
	}
	if !strings.Contains(out, "tamaño") {
		t.Errorf("la salida no dice que el piso no se relaja por tamaño del diff. "+
			"Un piso que se salta cuando el cambio parece menor se salta justo en el "+
			"cambio donde alguien tenía apuro:\n%s", out)
	}
}

// TestElDiffVacioEsUnaRespuesta: no es un fallo, y no se inventa un diff.
func TestElDiffVacioEsUnaRespuesta(t *testing.T) {
	out := RenderText(Classify(nil))
	if !strings.Contains(out, "nada que revisar") {
		t.Errorf("un diff vacío no dice 'nada que revisar':\n%s", out)
	}
	if strings.Contains(out, "role-dba") {
		t.Errorf("un diff vacío igual lista roles; no hay nada sobre qué decidir:\n%s", out)
	}
}

// TestLaReservaSaleSoloCuandoHayReserva es el anti-ruido de la tercera
// respuesta: si el bloque apareciera en toda corrida, se saltearía.
func TestLaReservaSaleSoloCuandoHayReserva(t *testing.T) {
	con := RenderText(Classify([]string{"lib/utils.py"}))
	if !strings.Contains(con, "lib/utils.py") || !strings.Contains(con, "no lo sabe") {
		t.Errorf("un archivo ambiguo no sale con su reserva:\n%s", con)
	}

	sin := RenderText(Classify([]string{"internal/planner/planner.go"}))
	if strings.Contains(sin, "no lo sabe") {
		t.Errorf("la reserva salió en una corrida sin nada ambiguo:\n%s", sin)
	}
}

// TestJSONLlevaLosDosLadosYElPiso: el orquestador lo consume sin parsear texto.
func TestJSONLlevaLosDosLadosYElPiso(t *testing.T) {
	raw, err := RenderJSON(Classify([]string{"api/auth.py", "README.md"}))
	if err != nil {
		t.Fatalf("RenderJSON error = %v", err)
	}

	var got Report
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("el JSON no vuelve a leerse como Report: %v\n%s", err, raw)
	}
	if !got.Floor.Triggered || got.Floor.File != "api/auth.py" {
		t.Errorf("el JSON perdió el piso: %+v", got.Floor)
	}
	if len(got.Decisions) != len(ordenCanonico) {
		t.Errorf("el JSON trae %d roles y son %d: los que no corresponden también van, "+
			"o el consumidor no puede distinguir ausencia de olvido",
			len(got.Decisions), len(ordenCanonico))
	}
	for _, d := range got.Decisions {
		if strings.TrimSpace(d.ReasonCode) == "" {
			t.Errorf("%s viaja sin reason_code; el consumidor tendría que parsear "+
				"la prosa en español para saber por qué", d.Role)
		}
	}
}

// TestElPisoNoSeAnunciaDosVecesEnLaMismaLinea: la etiqueta de la fila no repite
// la palabra que ya trae el motivo. Cazado mirando la salida real, no un test.
func TestElPisoNoSeAnunciaDosVecesEnLaMismaLinea(t *testing.T) {
	out := RenderText(Classify([]string{"api/auth.py"}))
	if strings.Contains(out, "piso: piso:") {
		t.Errorf("la fila anuncia el piso dos veces:\n%s", out)
	}
}
