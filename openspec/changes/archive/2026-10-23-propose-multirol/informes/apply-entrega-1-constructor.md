# Constructor: tanda B

Alcance: pasos 4, 5, 6, 7 y 7b de la propuesta firmada en 4748680.
Constructor: Codex, proveedor OpenAI; identidad y familia SUPUESTAS sin registro de despacho. No constituye revisión independiente ni aprobación de auditor.
Fuentes: CLAUDE.md, propuesta, fragmentos, tests, documentación y tanda A. KB numerada ausente; lecciones leídas en las referencias L-NNN de la propuesta. Casos dictados por la propuesta auditada; no se afirma convergencia independiente.
Consumidores: Compose de ocho variantes, comandos embebidos, catálogo y documentación. Pruebas textuales, sin acreditar conducta real de agentes. Sin instalación, build de binario ni push. Los informes previos del usuario quedan fuera de estos commits.


# Tareas: propose-multirol

- [x] Contrato observado rojo, implementación y verde registrados abajo.
- [x] Diff de goldens revisado.

## Bucle de verify

No iniciado: esta entrega es apply del CONSTRUCTOR. Vuelta pendiente, sin veredicto ni aprobación de auditor. El siguiente verify registra vuelta 1 antes de convocar revisores. No se simula una vuelta de verify mediante go test.

### rojo

```text
Hora: 2026-10-23T10:13:16.067728-03:00
SHA: 9f2d03b32559f888ef1931dbde221e173921e37b
Comando: go test ./internal/harness/config -run "TestCompose_RolesToggle_HighPropose" -count=1
--- FAIL: TestCompose_RolesToggle_HighProposeIsReadBySeveralRolesBeforeSignature (0.01s)
    roles_toggle_test.go:548: fila propose no exige: where governance is active
    roles_toggle_test.go:548: fila propose no exige: high
    roles_toggle_test.go:548: fila propose no exige: critical
    roles_toggle_test.go:548: fila propose no exige: at least three
    roles_toggle_test.go:548: fila propose no exige: functional analyst
    roles_toggle_test.go:548: fila propose no exige: auditor
    roles_toggle_test.go:548: fila propose no exige: before human approval
    roles_toggle_test.go:555: falta propose multirol: at least three independent specialist reports
    roles_toggle_test.go:555: falta propose multirol: architect, functional analyst and qa
    roles_toggle_test.go:555: falta propose multirol: one additional report for each specialty the scope requires
    roles_toggle_test.go:555: falta propose multirol: one at a time
    roles_toggle_test.go:555: falta propose multirol: exclude other specialists' reports and memories
    roles_toggle_test.go:555: falta propose multirol: auditor reads all n original reports
    roles_toggle_test.go:555: falta propose multirol: n >= 3
    roles_toggle_test.go:555: falta propose multirol: n + 1
    roles_toggle_test.go:555: falta propose multirol: before human approval
    roles_toggle_test.go:555: falta propose multirol: do not replace the human signature
    roles_toggle_test.go:555: falta propose multirol: ## aprobación
    roles_toggle_test.go:555: falta propose multirol: until all required reports exist
    roles_toggle_test.go:555: falta propose multirol: identify which report is missing and why
    roles_toggle_test.go:555: falta propose multirol: informes/
    roles_toggle_test.go:555: falta propose multirol: leeme.md
    roles_toggle_test.go:555: falta propose multirol: model/provider
    roles_toggle_test.go:555: falta propose multirol: provenance
    roles_toggle_test.go:555: falta propose multirol: full return or textual transcript
    roles_toggle_test.go:555: falta propose multirol: contamination
    roles_toggle_test.go:555: falta propose multirol: does not count as independent
    roles_toggle_test.go:555: falta propose multirol: obtain a replacement
    roles_toggle_test.go:555: falta propose multirol: reviewers receive no write
    roles_toggle_test.go:555: falta propose multirol: not proof of isolation
    roles_toggle_test.go:563: regla multirol fuera del lugar aprobado
--- FAIL: TestCompose_RolesToggle_HighProposeIsConditionalOnGovernance (0.00s)
    roles_toggle_test.go:571: falta condición: where governance is active and marks the change high or critical
    roles_toggle_test.go:571: falta condición: if governance is absent, do not invent a high/critical classification
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.009s
FAIL
Salida: 1
```

### verde-dirigido

```text
Hora: 2026-10-23T10:13:48.024232-03:00
SHA: 9f2d03b32559f888ef1931dbde221e173921e37b
Comando: go test ./internal/harness/config -run TestCompose_RolesToggle -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.014s
Salida: 0
```

### goldens

```text
Hora: 2026-10-23T10:13:48.534040-03:00
SHA: 9f2d03b32559f888ef1931dbde221e173921e37b
Comando: go test ./internal/harness/config -run TestCompose_AllToggles -update -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
Salida: 0
```

### verde

```text
Hora: 2026-10-23T10:13:48.784130-03:00
SHA: 9f2d03b32559f888ef1931dbde221e173921e37b
Comando: go test ./...
ok  	github.com/Group-Active-IA/active-stack/assets	(cached)
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.511s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	2.346s
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	0.613s
ok  	github.com/Group-Active-IA/active-stack/internal/agents	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/backup	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	2.631s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/install	0.207s
ok  	github.com/Group-Active-IA/active-stack/internal/model	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/planner	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/roles	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/system	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/tui	0.032s
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/verify	0.034s
Salida: 0
```

### cierre-post-commit

```text
Hora: 2026-10-23T10:13:52.497969-03:00
SHA: 531e1b65c314c56e9e6b3eb43d66096f77cfff7d
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

Commit entregado: 531e1b65c314c56e9e6b3eb43d66096f77cfff7d.

Chequeos adicionales, hora de registro 2026-10-23T10:14:42.093474-03:00: gofmt -l internal/harness/config assets/commands_test.go y go vet ./... terminaron con salida 0, sin salida textual. git diff --check terminó con salida 0.
