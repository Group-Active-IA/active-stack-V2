# Constructor: tanda B

Alcance: pasos 4, 5, 6, 7 y 7b de la propuesta firmada en 4748680.
Constructor: Codex, proveedor OpenAI; identidad y familia SUPUESTAS sin registro de despacho. No constituye revisión independiente ni aprobación de auditor.
Fuentes: CLAUDE.md, propuesta, fragmentos, tests, documentación y tanda A. KB numerada ausente; lecciones leídas en las referencias L-NNN de la propuesta. Casos dictados por la propuesta auditada; no se afirma convergencia independiente.
Consumidores: Compose de ocho variantes, comandos embebidos, catálogo y documentación. Pruebas textuales, sin acreditar conducta real de agentes. Sin instalación, build de binario ni push. Los informes previos del usuario quedan fuera de estos commits.


# Tareas: archive-base-completa

- [x] Contrato observado rojo, implementación y verde registrados abajo.
- [x] Diff de goldens revisado.

## Bucle de verify

No iniciado: esta entrega es apply del CONSTRUCTOR. Vuelta pendiente, sin veredicto ni aprobación de auditor. El siguiente verify registra vuelta 1 antes de convocar revisores. No se simula una vuelta de verify mediante go test.

### rojo

```text
Hora: 2026-10-23T10:11:13.207316-03:00
SHA: 66ad39ec9b61364a69e518b749c14c023f9b7f14
Comando: go test ./internal/harness/config -run "TestArchive|TestCompose_KnowledgeBaseToggle_TheDetector" -count=1
--- FAIL: TestCompose_KnowledgeBaseToggle_TheDetectorIsTheRewriteList (0.00s)
    knowledge_base_toggle_test.go:396: the detector section is missing "suspicious: append-only" — the detector itself
--- FAIL: TestArchiveBaseCompletaConEvidencia (0.00s)
    knowledge_base_toggle_test.go:720: falta contrato de archive: union of knowledge base impact and the actual diff
    knowledge_base_toggle_test.go:720: falta contrato de archive: even when the impact section exists
    knowledge_base_toggle_test.go:720: falta contrato de archive: lessons l-nnn
    knowledge_base_toggle_test.go:720: falta contrato de archive: only append lessons
    knowledge_base_toggle_test.go:720: falta contrato de archive: lecciones
    knowledge_base_toggle_test.go:720: falta contrato de archive: deuda
    knowledge_base_toggle_test.go:720: falta contrato de archive: each entry
    knowledge_base_toggle_test.go:720: falta contrato de archive: what changed
    knowledge_base_toggle_test.go:720: falta contrato de archive: sha, test or decision
    knowledge_base_toggle_test.go:720: falta contrato de archive: typed evidence
    knowledge_base_toggle_test.go:720: falta contrato de archive: auditor source absent
    knowledge_base_toggle_test.go:720: falta contrato de archive: written decision
    knowledge_base_toggle_test.go:720: falta contrato de archive: do not invent
    knowledge_base_toggle_test.go:736: prompt no entrega: actual diff
    knowledge_base_toggle_test.go:736: prompt no entrega: lessons l-nnn
    knowledge_base_toggle_test.go:736: prompt no entrega: lecciones
    knowledge_base_toggle_test.go:736: prompt no entrega: deuda
    knowledge_base_toggle_test.go:736: prompt no entrega: sha, test or decision
    knowledge_base_toggle_test.go:736: prompt no entrega: typed evidence
    knowledge_base_toggle_test.go:736: prompt no entrega: suspicious: append-only
--- FAIL: TestArchiveSonnetEnOchoVariantes (0.00s)
    --- FAIL: TestArchiveSonnetEnOchoVariantes/claude (0.00s)
        verify_phase_test.go:357: archive usa haiku, se esperaba sonnet
        verify_phase_test.go:361: motivo archive no declara evidence
        verify_phase_test.go:361: motivo archive no declara substitution
        verify_phase_test.go:361: motivo archive no declara unavailable
        verify_phase_test.go:366: instrucción contradictoria: file operations
    --- FAIL: TestArchiveSonnetEnOchoVariantes/opencode (0.00s)
        verify_phase_test.go:357: archive usa haiku, se esperaba sonnet
        verify_phase_test.go:361: motivo archive no declara evidence
        verify_phase_test.go:361: motivo archive no declara substitution
        verify_phase_test.go:361: motivo archive no declara unavailable
        verify_phase_test.go:366: instrucción contradictoria: file operations
    --- FAIL: TestArchiveSonnetEnOchoVariantes/codex (0.00s)
        verify_phase_test.go:357: archive usa haiku, se esperaba sonnet
        verify_phase_test.go:361: motivo archive no declara evidence
        verify_phase_test.go:361: motivo archive no declara substitution
        verify_phase_test.go:361: motivo archive no declara unavailable
        verify_phase_test.go:366: instrucción contradictoria: file operations
    --- FAIL: TestArchiveSonnetEnOchoVariantes/cursor (0.00s)
        verify_phase_test.go:357: archive usa haiku, se esperaba sonnet
        verify_phase_test.go:361: motivo archive no declara evidence
        verify_phase_test.go:361: motivo archive no declara substitution
        verify_phase_test.go:361: motivo archive no declara unavailable
        verify_phase_test.go:366: instrucción contradictoria: file operations
    --- FAIL: TestArchiveSonnetEnOchoVariantes/gemini (0.00s)
        verify_phase_test.go:357: archive usa haiku, se esperaba sonnet
        verify_phase_test.go:361: motivo archive no declara evidence
        verify_phase_test.go:361: motivo archive no declara substitution
        verify_phase_test.go:361: motivo archive no declara unavailable
        verify_phase_test.go:366: instrucción contradictoria: file operations
    --- FAIL: TestArchiveSonnetEnOchoVariantes/windsurf (0.00s)
        verify_phase_test.go:357: archive usa haiku, se esperaba sonnet
        verify_phase_test.go:361: motivo archive no declara evidence
        verify_phase_test.go:361: motivo archive no declara substitution
        verify_phase_test.go:361: motivo archive no declara unavailable
        verify_phase_test.go:366: instrucción contradictoria: less on archive
        verify_phase_test.go:366: instrucción contradictoria: file operations
    --- FAIL: TestArchiveSonnetEnOchoVariantes/antigravity (0.00s)
        verify_phase_test.go:357: archive usa haiku, se esperaba sonnet
        verify_phase_test.go:361: motivo archive no declara evidence
        verify_phase_test.go:361: motivo archive no declara substitution
        verify_phase_test.go:361: motivo archive no declara unavailable
        verify_phase_test.go:366: instrucción contradictoria: less on archive
        verify_phase_test.go:366: instrucción contradictoria: file operations
    --- FAIL: TestArchiveSonnetEnOchoVariantes/generic (0.00s)
        verify_phase_test.go:357: archive usa haiku, se esperaba sonnet
        verify_phase_test.go:361: motivo archive no declara evidence
        verify_phase_test.go:361: motivo archive no declara substitution
        verify_phase_test.go:361: motivo archive no declara unavailable
        verify_phase_test.go:366: instrucción contradictoria: file operations
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.007s
FAIL
Salida: 1
```

### verde-dirigido

```text
Hora: 2026-10-23T10:11:47.574315-03:00
SHA: 66ad39ec9b61364a69e518b749c14c023f9b7f14
Comando: go test ./internal/harness/config -run "TestArchive|TestCompose_KnowledgeBase|TestOrchestrator" -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.023s
Salida: 0
```

### goldens

```text
Hora: 2026-10-23T10:11:48.104631-03:00
SHA: 66ad39ec9b61364a69e518b749c14c023f9b7f14
Comando: go test ./internal/harness/config -run TestCompose_AllToggles -update -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.006s
Salida: 0
```

### verde

```text
Hora: 2026-10-23T10:11:48.387198-03:00
SHA: 66ad39ec9b61364a69e518b749c14c023f9b7f14
Comando: go test ./...
ok  	github.com/Group-Active-IA/active-stack/assets	(cached)
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.564s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	2.595s
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	0.654s
ok  	github.com/Group-Active-IA/active-stack/internal/agents	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/backup	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	2.862s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/install	0.226s
ok  	github.com/Group-Active-IA/active-stack/internal/model	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/planner	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/roles	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/system	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/tui	0.033s
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/verify	0.034s
Salida: 0
```

### cierre-post-commit

```text
Hora: 2026-10-23T10:11:52.445183-03:00
SHA: 84c570d86cb5622c659a975a4e1c02b10d81db8c
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

Commit entregado: 84c570d86cb5622c659a975a4e1c02b10d81db8c.
