# Constructor: tanda B

Alcance: pasos 4, 5, 6, 7 y 7b de la propuesta firmada en 4748680.
Constructor: Codex, proveedor OpenAI; identidad y familia SUPUESTAS sin registro de despacho. No constituye revisión independiente ni aprobación de auditor.
Fuentes: CLAUDE.md, propuesta, fragmentos, tests, documentación y tanda A. KB numerada ausente; lecciones leídas en las referencias L-NNN de la propuesta. Casos dictados por la propuesta auditada; no se afirma convergencia independiente.
Consumidores: Compose de ocho variantes, comandos embebidos, catálogo y documentación. Pruebas textuales, sin acreditar conducta real de agentes. Sin instalación, build de binario ni push. Los informes previos del usuario quedan fuera de estos commits.


# Tareas: senior-solo-en-verify

- [x] Contrato observado rojo, implementación y verde registrados abajo.
- [x] Diff de goldens revisado.

## Bucle de verify

No iniciado: esta entrega es apply del CONSTRUCTOR. Vuelta pendiente, sin veredicto ni aprobación de auditor. El siguiente verify registra vuelta 1 antes de convocar revisores. No se simula una vuelta de verify mediante go test.

### rojo

```text
Hora: 2026-10-23T10:09:54.390356-03:00
SHA: db3de543858475a1cd8e574ce77b6f60933679f1
Comando: go test ./internal/harness/config -run "TestCompose_SeniorToggle_(Preamble|Asks|Runs|IsNotAnEighth|Declares)|TestDocsSenior|TestSeniorCatalogo" -count=1
--- FAIL: TestCompose_SeniorToggle_PreambleSaysItIsEveryPhase (0.00s)
    senior_toggle_test.go:197: falta alcance de senior: only in **verify**
--- FAIL: TestCompose_SeniorToggle_AsksADifferentQuestionInEachPhase (0.00s)
    senior_toggle_test.go:341: falta alcance de senior: dated decisions
--- FAIL: TestCompose_SeniorToggle_RunsAtTheCloseOfEachPhase (0.00s)
    senior_toggle_test.go:356: falta alcance de senior: after the breaker returns and before the auditor
    senior_toggle_test.go:357: falta alcance de senior: quiet diff
--- FAIL: TestCompose_SeniorToggle_IsNotAnEighthRole (0.00s)
    senior_toggle_test.go:374: falta alcance de senior: once per verify round when active
--- FAIL: TestCompose_SeniorToggle_DeclaresItsCost (0.00s)
    senior_toggle_test.go:495: falta alcance de senior: one more sub-agent per verify round
    senior_toggle_test.go:496: falta alcance de senior: at most three per change
--- FAIL: TestSeniorCatalogoDescribeVerify (0.00s)
    senior_toggle_test.go:581: catálogo no describe senior sólo en verify
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.014s
FAIL
Salida: 1
```

### rojo-docs

```text
Hora: 2026-10-23T10:10:25.765304-03:00
SHA: db3de543858475a1cd8e574ce77b6f60933679f1
Comando: go test ./internal/harness/config -run ^TestDocsSeniorSoloVerify$ -count=1
--- FAIL: TestDocsSeniorSoloVerify (0.00s)
    docs_toggles_test.go:268: README.md:197: falta senior sólo en verify
    docs_toggles_test.go:268: ARCHITECTURE.md:84: falta senior sólo en verify
    docs_toggles_test.go:268: AGENTS.md:86: falta senior sólo en verify
    docs_toggles_test.go:268: CLAUDE.md:92: falta senior sólo en verify
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL
Salida: 1
```

### verde-dirigido

```text
Hora: 2026-10-23T10:10:26.224105-03:00
SHA: db3de543858475a1cd8e574ce77b6f60933679f1
Comando: go test ./internal/harness/config ./internal/catalog -run "TestCompose_Senior|TestDocsSenior|TestSeniorCatalogo|TestLoad" -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.015s
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	0.008s
Salida: 0
```

### goldens

```text
Hora: 2026-10-23T10:10:38.590568-03:00
SHA: db3de543858475a1cd8e574ce77b6f60933679f1
Comando: go test ./internal/harness/config -run TestCompose_AllToggles -update -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.005s
Salida: 0
```

### verde

```text
Hora: 2026-10-23T10:10:38.850914-03:00
SHA: db3de543858475a1cd8e574ce77b6f60933679f1
Comando: go test ./...
ok  	github.com/Group-Active-IA/active-stack/assets	(cached)
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.533s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	2.422s
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	0.630s
ok  	github.com/Group-Active-IA/active-stack/internal/agents	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/backup	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	0.090s
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	2.659s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/install	0.225s
ok  	github.com/Group-Active-IA/active-stack/internal/model	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/planner	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/roles	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/system	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/tui	0.034s
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/verify	0.044s
Salida: 0
```

### cierre-post-commit

```text
Hora: 2026-10-23T10:10:43.032451-03:00
SHA: 66ad39ec9b61364a69e518b749c14c023f9b7f14
Comando: go test ./...
ok  	github.com/Group-Active-IA/active-stack/assets	(cached)
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	(cached)
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	(cached)
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/backup	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/install	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/model	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/planner	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/roles	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/system	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/tui	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/verify	(cached)
Salida: 0
```

Commit entregado: 66ad39ec9b61364a69e518b749c14c023f9b7f14.
