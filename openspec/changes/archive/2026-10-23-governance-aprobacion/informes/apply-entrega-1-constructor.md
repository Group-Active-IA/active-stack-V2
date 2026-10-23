# Constructor: tanda B

Alcance: pasos 4, 5, 6, 7 y 7b de la propuesta firmada en 4748680.
Constructor: Codex, proveedor OpenAI; identidad y familia SUPUESTAS sin registro de despacho. No constituye revisión independiente ni aprobación de auditor.
Fuentes: CLAUDE.md, propuesta, fragmentos, tests, documentación y tanda A. KB numerada ausente; lecciones leídas en las referencias L-NNN de la propuesta. Casos dictados por la propuesta auditada; no se afirma convergencia independiente.
Consumidores: Compose de ocho variantes, comandos embebidos, catálogo y documentación. Pruebas textuales, sin acreditar conducta real de agentes. Sin instalación, build de binario ni push. Los informes previos del usuario quedan fuera de estos commits.


# Tareas: governance-aprobacion

- [x] Contrato observado rojo, implementación y verde registrados abajo.
- [x] Diff de goldens revisado.

## Bucle de verify

No iniciado: esta entrega es apply del CONSTRUCTOR. Vuelta pendiente, sin veredicto ni aprobación de auditor. El siguiente verify registra vuelta 1 antes de convocar revisores. No se simula una vuelta de verify mediante go test.

### rojo

```text
Hora: 2026-10-23T10:12:17.443422-03:00
SHA: 84c570d86cb5622c659a975a4e1c02b10d81db8c
Comando: go test ./internal/harness/config -run ^TestGovernanceAprobacionHumanaEnProposal$ -count=1
--- FAIL: TestGovernanceAprobacionHumanaEnProposal (0.00s)
    governance_approval_test.go:17: falta contrato de firma: ## aprobación
    governance_approval_test.go:17: falta contrato de firma: proposal.md
    governance_approval_test.go:17: falta contrato de firma: | elección | persona | fecha | alcance |
    governance_approval_test.go:17: falta contrato de firma: no agent fills
    governance_approval_test.go:17: falta contrato de firma: written, dated and explicit scope
    governance_approval_test.go:17: falta contrato de firma: missing or ambiguous approval stops construction
    governance_approval_test.go:17: falta contrato de firma: new rule after signature
    governance_approval_test.go:17: falta contrato de firma: own row
    governance_approval_test.go:17: falta contrato de firma: before construction
    governance_approval_test.go:17: falta contrato de firma: l-029
    governance_approval_test.go:17: falta contrato de firma: cryptographic verification
    governance_approval_test.go:17: falta contrato de firma: not implemented
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.003s
FAIL
Salida: 1
```

### verde-dirigido

```text
Hora: 2026-10-23T10:12:17.842384-03:00
SHA: 84c570d86cb5622c659a975a4e1c02b10d81db8c
Comando: go test ./internal/harness/config -run ^TestGovernanceAprobacionHumanaEnProposal$ -count=1
--- FAIL: TestGovernanceAprobacionHumanaEnProposal (0.00s)
    governance_approval_test.go:17: falta contrato de firma: | elección | persona | fecha | alcance |
    governance_approval_test.go:17: falta contrato de firma: no agent fills
    governance_approval_test.go:17: falta contrato de firma: missing or ambiguous approval stops construction
    governance_approval_test.go:17: falta contrato de firma: new rule after signature
    governance_approval_test.go:17: falta contrato de firma: own row
    governance_approval_test.go:17: falta contrato de firma: before construction
    governance_approval_test.go:17: falta contrato de firma: l-029
    governance_approval_test.go:17: falta contrato de firma: cryptographic verification
    governance_approval_test.go:17: falta contrato de firma: not implemented
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.003s
FAIL
Salida: 1
```

### goldens

```text
Hora: 2026-10-23T10:12:18.326285-03:00
SHA: 84c570d86cb5622c659a975a4e1c02b10d81db8c
Comando: go test ./internal/harness/config -run TestCompose_AllToggles -update -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
Salida: 0
```

### verde

```text
Hora: 2026-10-23T10:12:18.582404-03:00
SHA: 84c570d86cb5622c659a975a4e1c02b10d81db8c
Comando: go test ./...
ok  	github.com/Group-Active-IA/active-stack/assets	(cached)
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.517s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	2.333s
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	0.599s
ok  	github.com/Group-Active-IA/active-stack/internal/agents	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/backup	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	(cached)
--- FAIL: TestGovernanceAprobacionHumanaEnProposal (0.00s)
    governance_approval_test.go:17: falta contrato de firma: | elección | persona | fecha | alcance |
    governance_approval_test.go:17: falta contrato de firma: no agent fills
    governance_approval_test.go:17: falta contrato de firma: missing or ambiguous approval stops construction
    governance_approval_test.go:17: falta contrato de firma: new rule after signature
    governance_approval_test.go:17: falta contrato de firma: own row
    governance_approval_test.go:17: falta contrato de firma: before construction
    governance_approval_test.go:17: falta contrato de firma: l-029
    governance_approval_test.go:17: falta contrato de firma: cryptographic verification
    governance_approval_test.go:17: falta contrato de firma: not implemented
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	2.612s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/install	0.210s
ok  	github.com/Group-Active-IA/active-stack/internal/model	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/planner	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/roles	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/system	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/tui	0.032s
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/verify	0.034s
FAIL
Salida: 1
```

### cierre-post-commit

```text
Hora: 2026-10-23T10:12:22.226524-03:00
SHA: 7f467d534485318778bc17db4efcf98e182ded8f
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
--- FAIL: TestGovernanceAprobacionHumanaEnProposal (0.00s)
    governance_approval_test.go:17: falta contrato de firma: | elección | persona | fecha | alcance |
    governance_approval_test.go:17: falta contrato de firma: no agent fills
    governance_approval_test.go:17: falta contrato de firma: missing or ambiguous approval stops construction
    governance_approval_test.go:17: falta contrato de firma: new rule after signature
    governance_approval_test.go:17: falta contrato de firma: own row
    governance_approval_test.go:17: falta contrato de firma: before construction
    governance_approval_test.go:17: falta contrato de firma: l-029
    governance_approval_test.go:17: falta contrato de firma: cryptographic verification
    governance_approval_test.go:17: falta contrato de firma: not implemented
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	2.350s
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
FAIL
Salida: 1
```

### rojo-aserciones-helper-corregido

```text
Hora: 2026-10-23T10:12:43.792366-03:00
SHA: 7f467d534485318778bc17db4efcf98e182ded8f
Comando: go test ./internal/harness/config -run ^TestGovernanceAprobacionHumanaEnProposal$ -count=1
--- FAIL: TestGovernanceAprobacionHumanaEnProposal (0.00s)
    governance_approval_test.go:22: falta contrato de firma: ## aprobación
    governance_approval_test.go:22: falta contrato de firma: proposal.md
    governance_approval_test.go:22: falta contrato de firma: | elección | persona | fecha | alcance |
    governance_approval_test.go:22: falta contrato de firma: no agent fills
    governance_approval_test.go:22: falta contrato de firma: written, dated and explicit scope
    governance_approval_test.go:22: falta contrato de firma: missing or ambiguous approval stops construction
    governance_approval_test.go:22: falta contrato de firma: new rule after signature
    governance_approval_test.go:22: falta contrato de firma: own row
    governance_approval_test.go:22: falta contrato de firma: before construction
    governance_approval_test.go:22: falta contrato de firma: l-029
    governance_approval_test.go:22: falta contrato de firma: cryptographic verification
    governance_approval_test.go:22: falta contrato de firma: not implemented
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.003s
FAIL
Salida: 1
```

### verde-corregido

```text
Hora: 2026-10-23T10:12:44.233810-03:00
SHA: 7f467d534485318778bc17db4efcf98e182ded8f
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
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	2.370s
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

### cierre-post-commit-corregido

```text
Hora: 2026-10-23T10:12:47.289621-03:00
SHA: 9f2d03b32559f888ef1931dbde221e173921e37b
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

Commit entregado: 9f2d03b32559f888ef1931dbde221e173921e37b.
