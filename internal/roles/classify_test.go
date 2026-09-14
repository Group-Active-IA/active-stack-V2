package roles

import (
	"strings"
	"testing"
)

// convocados arma un conjunto con los roles convocados, para asserts legibles.
func convocados(t *testing.T, rep Report) map[string]bool {
	t.Helper()
	out := map[string]bool{}
	for _, r := range rep.Convened() {
		out[r] = true
	}
	return out
}

// exige falla nombrando el rol que faltó y el archivo que lo pedía.
func exige(t *testing.T, rep Report, role string, files ...string) {
	t.Helper()
	d, ok := rep.Decision(role)
	if !ok {
		t.Fatalf("el informe no dice nada sobre %s; los siete roles salen siempre, "+
			"convocados o no: una ausencia es indistinguible de un olvido", role)
	}
	if !d.Convened {
		t.Errorf("%v no convocó a %s. Motivo devuelto: %q", files, role, d.Reason)
	}
	if len(d.Files) == 0 && len(files) > 0 {
		t.Errorf("%s fue convocado sin el hecho del diff que lo justifica: %+v", role, d)
	}
}

// TestMigracionConvocaAlDBA es la primera fila de la tabla de `roles.md`:
// "a migration file, a schema, a query → DBA".
func TestMigracionConvocaAlDBA(t *testing.T) {
	rep := Classify([]string{"migrations/003_tenant.sql"})
	exige(t, rep, RoleDBA, "migrations/003_tenant.sql")
}

// TestEnvActivaElPisoEntero es la segunda fila, la que no se puede errar:
// secretos, variables de entorno, autenticación, pagos o datos de usuario
// convocan el piso entero, sin importar el tamaño del diff.
func TestEnvActivaElPisoEntero(t *testing.T) {
	rep := Classify([]string{".env"})

	if !rep.Floor.Triggered {
		t.Fatalf("un diff que toca .env no activó el piso: %+v", rep.Floor)
	}
	if rep.Floor.File != ".env" {
		t.Errorf("el piso se activó sin decir qué archivo lo disparó: %+v. "+
			"Un piso sin su hecho es una afirmación que nadie puede discutir", rep.Floor)
	}
	got := convocados(t, rep)
	for _, role := range []string{RoleBackend, RoleQA, RoleDBA, RoleArch} {
		if !got[role] {
			t.Errorf("el piso no convocó a %s; el piso es backend, QA, DBA y arquitecto, "+
				"y no es del orquestador la decisión de recortarlo", role)
		}
	}
}

// TestElPisoNoSeRelajaPorTamano es el anti-assert del piso.
//
// Un diff de tres líneas que toca auth.py convoca el piso igual que uno de
// trescientas. Un piso que se salta cuando el cambio parece menor se salta
// exactamente en el cambio donde alguien tenía apuro.
func TestElPisoNoSeRelajaPorTamano(t *testing.T) {
	chico := Classify([]string{"api/auth.py"})
	grande := Classify([]string{
		"api/auth.py", "api/users.go", "api/session.go", "web/login.tsx",
		"migrations/004.sql", "deploy/ci.yml", "docs/readme.md",
	})

	if !chico.Floor.Triggered {
		t.Fatalf("un diff de un solo archivo que toca autenticación no activó el piso")
	}
	for _, role := range []string{RoleBackend, RoleQA, RoleDBA, RoleArch} {
		cd, _ := chico.Decision(role)
		gd, _ := grande.Decision(role)
		if cd.Convened != gd.Convened {
			t.Errorf("%s se convoca en el diff grande (%v) y no en el chico (%v): "+
				"el piso se relajó por tamaño, que es justo lo que no puede hacer",
				role, gd.Convened, cd.Convened)
		}
	}
}

// TestSoloUnReadmeNoConvocaANadieMasQueQA es la última fila de la tabla:
// "none of the above → the small review: QA, or nobody".
func TestSoloUnReadmeNoConvocaANadieMasQueQA(t *testing.T) {
	rep := Classify([]string{"README.md"})

	if rep.Floor.Triggered {
		t.Errorf("un README activó el piso: %+v", rep.Floor)
	}
	for _, role := range []string{RoleDBA, RoleBackend, RoleFrontend, RoleDevOps, RoleArch} {
		if d, _ := rep.Decision(role); d.Convened {
			t.Errorf("un README convocó a %s (%q). Convocar a todos cuesta lo mismo "+
				"que no convocar a nadie: deja de ser señal", role, d.Reason)
		}
	}
	if d, _ := rep.Decision(RoleQA); !d.Convened {
		t.Errorf("la revisión chica no convocó a QA")
	}
}

// TestDiffVacioNoConvocaANadie: un diff vacío significa "nada que revisar", y
// eso es una respuesta, no un fallo.
func TestDiffVacioNoConvocaANadie(t *testing.T) {
	rep := Classify(nil)

	if got := rep.Convened(); len(got) != 0 {
		t.Errorf("un diff vacío convocó a %v; no hay nada que revisar", got)
	}
	if len(rep.Decisions) == 0 {
		t.Errorf("un diff vacío no devolvió ninguna decisión; los siete roles salen " +
			"siempre, también para decir que no hay nada que mirar")
	}
	if rep.Floor.Triggered {
		t.Errorf("un diff vacío activó el piso: %+v", rep.Floor)
	}
}

// TestArchivoAmbiguoNoSeAfirmaComoNoCorresponde es la tercera respuesta.
//
// El comando clasifica por ruta. Un utils.py que resulta tener la validación de
// tokens adentro no lo puede saber, y tiene que decir que no lo puede saber en
// lugar de afirmar que no corresponde.
func TestArchivoAmbiguoNoSeAfirmaComoNoCorresponde(t *testing.T) {
	rep := Classify([]string{"lib/utils.py"})

	if len(rep.Unclear) == 0 {
		t.Fatalf("lib/utils.py salió clasificado sin reservas. El comando clasifica " +
			"por ruta: un nombre genérico no dice qué hay adentro, y afirmar que no " +
			"corresponde es afirmar algo que no se miró")
	}
	if rep.Unclear[0].File != "lib/utils.py" {
		t.Errorf("la reserva no nombra el archivo: %+v", rep.Unclear[0])
	}
}

// TestUnNombreConcretoNoEsAmbiguo es el anti-assert de la tercera respuesta: si
// cada archivo sin rol saliera con reserva, la sección sería ruido en cada
// corrida y se saltearía como se saltea todo lo que aparece siempre.
func TestUnNombreConcretoNoEsAmbiguo(t *testing.T) {
	rep := Classify([]string{"internal/planner/planner.go", "docs/guia.md"})

	if len(rep.Unclear) != 0 {
		t.Errorf("archivos con nombre concreto salieron con reserva: %+v. "+
			"Una reserva que aparece en cada corrida deja de leerse", rep.Unclear)
	}
}

// TestFrontendPorArbolDeUI: "any file under the UI tree → frontend".
func TestFrontendPorArbolDeUI(t *testing.T) {
	exige(t, Classify([]string{"web/components/Boton.tsx"}), RoleFrontend, "web/components/Boton.tsx")
	exige(t, Classify([]string{"internal/tui/model.go"}), RoleFrontend, "internal/tui/model.go")
}

// TestBackendPorEndpoint: "an endpoint, a handler, a validation at the boundary".
func TestBackendPorEndpoint(t *testing.T) {
	exige(t, Classify([]string{"internal/server/handlers/pedidos.go"}), RoleBackend, "handlers/pedidos.go")
}

// TestDevOpsPorCI: "deploy, CI or containers → DevOps".
func TestDevOpsPorCI(t *testing.T) {
	exige(t, Classify([]string{".github/workflows/ci.yml"}), RoleDevOps, ".github/workflows/ci.yml")
	exige(t, Classify([]string{"Dockerfile"}), RoleDevOps, "Dockerfile")
}

// TestLosRolesQueNoCorrespondenSalenConMotivo es la mitad del valor del comando.
func TestLosRolesQueNoCorrespondenSalenConMotivo(t *testing.T) {
	rep := Classify([]string{"migrations/003_tenant.sql"})

	for _, role := range []string{RoleFrontend, RoleDevOps} {
		d, ok := rep.Decision(role)
		if !ok {
			t.Fatalf("%s no aparece en el informe", role)
		}
		if d.Convened {
			t.Fatalf("%s se convocó por una migración", role)
		}
		if strings.TrimSpace(d.Reason) == "" {
			t.Errorf("%s no salió convocado y sin motivo. Una ausencia sin motivo es "+
				"indistinguible de un olvido, que es el defecto que este comando cierra", role)
		}
	}
}

// TestElArquitectoNoLoDecideLaRuta: una frontera no se lee por ruta, así que
// fuera del piso el comando no afirma que el arquitecto no corresponde.
func TestElArquitectoNoLoDecideLaRuta(t *testing.T) {
	d, _ := Classify([]string{"internal/planner/planner.go"}).Decision(RoleArch)
	if d.ReasonCode != "not-decided-by-diff" {
		t.Errorf("el arquitecto salió con %q. El diff no decide sobre una frontera: "+
			"decir 'no corresponde' ahí es afirmar algo que la ruta no muestra", d.ReasonCode)
	}
}

// TestRutasDeWindowsSeNormalizan: el CI cruza Linux y GOOS=windows, y ya cazó
// bugs de rutas. Una barra invertida no puede cambiar a quién se convoca.
func TestRutasDeWindowsSeNormalizan(t *testing.T) {
	rep := Classify([]string{`migrations\003_tenant.sql`, `web\components\Boton.tsx`})
	got := convocados(t, rep)
	if !got[RoleDBA] || !got[RoleFrontend] {
		t.Errorf("rutas con barra invertida clasificaron distinto que con barra normal: %v", got)
	}
}

// TestUnGoldenDeTestdataNoEsDeploy.
//
// Primera corrida sobre el diff real de este repo (2026-09-13): convocó a
// DevOps por `internal/harness/config/testdata/compose_all_toggles_claude.golden`,
// porque "compose" estaba en la lista de DevOps por docker-compose. "Ante la
// duda, convoca" es sobre la duda, no sobre un falso positivo obvio: un revisor
// que aparece en toda corrida deja de ser señal, que es lo mismo que no
// convocarlo.
func TestUnGoldenDeTestdataNoEsDeploy(t *testing.T) {
	rep := Classify([]string{"internal/harness/config/testdata/compose_all_toggles_claude.golden"})
	if d, _ := rep.Decision(RoleDevOps); d.Convened {
		t.Errorf("un golden de testdata convocó a DevOps: %+v", d)
	}
}

// TestUnComposeDeVerdadSigueSiendoDeploy es la otra dirección: angostar el
// término no puede costar el caso que lo justificaba.
func TestUnComposeDeVerdadSigueSiendoDeploy(t *testing.T) {
	exige(t, Classify([]string{"docker-compose.yml"}), RoleDevOps, "docker-compose.yml")
	exige(t, Classify([]string{"deploy/compose.yaml"}), RoleDevOps, "deploy/compose.yaml")
}

// TestMainNoEsUnNombreAmbiguo.
//
// Misma corrida: `cmd/active-stack/main.go` salió con reserva. "main" e "index"
// son nombres POSICIONALES (punto de entrada, barrel), no cajones de sastre, y
// aparecen en casi todo diff de un CLI o de un front. Una reserva que sale
// siempre se saltea igual que un revisor que se convoca siempre.
func TestMainNoEsUnNombreAmbiguo(t *testing.T) {
	rep := Classify([]string{"cmd/active-stack/main.go"})
	if len(rep.Unclear) != 0 {
		t.Errorf("main.go salió con reserva: %+v", rep.Unclear)
	}
}
