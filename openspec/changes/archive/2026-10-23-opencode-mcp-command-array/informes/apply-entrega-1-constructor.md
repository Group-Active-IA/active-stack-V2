# Apply, entrega 1, constructor: opencode-mcp-command-array

Autor: Codex. Proveedor declarado: OpenAI. Modelo/familia exactos sin registro de despacho: SUPUESTO. SHA de entrega: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25.
Procedencia: informe del constructor y salidas textuales de herramientas de esta sesión; copia manual mediada por el constructor, sin captura por hook. No se afirma independencia ni revisión de terceros. Se leyeron las reglas, propuesta y exploración; no hay KB local ni archivo de lecciones. No se recibieron casos dictados por el auditor en esta entrega.

## Paso 1: opencode-mcp-command-array

Cherry-pick sin commit de fae8bbf, disponible desde juancruz (remote local ya agregado). Conflicto de mcp_test.go resuelto conservando el comentario del schema command array con formato gofmt. Para observar el rojo se guardó la versión entrante de mcp.go, se restituyó temporalmente la de HEAD, se corrieron los tests entrantes y recién después se aplicó producción. Se agregó aserción de environment porque el cherry-pick también cambia ese campo. No se tomó ddeb94f ni se tocaron permisos. Se eliminó el trailer de coautoría y se usa el nombre del change en el conventional commit en español. Sin goldens.

### Rojo observado

```text
Hora: 2026-10-23 09:57:29 -0300
SHA: 84fcf0ebc0a6f647c04b5b0eb4dfbbc9939c866d
Comando: go test ./internal/harness/external -run ^TestBuildStdioOverlay_OpenCode$ -count=1
--- FAIL: TestBuildStdioOverlay_OpenCode (0.00s)
    mcp_test.go:521: command = engram, want [engram mcp] (OpenCode schema needs binary+args combined)
    mcp_test.go:524: OpenCode local overlay must not have a separate 'args' key
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/external	0.004s
FAIL
```

```text
Hora: 2026-10-23 09:57:45 -0300
SHA: 84fcf0ebc0a6f647c04b5b0eb4dfbbc9939c866d
Comando: go test ./internal/harness/external -run ^TestBuildStdioOverlay_OpenCode$ -count=1
--- FAIL: TestBuildStdioOverlay_OpenCode (0.00s)
    mcp_test.go:521: command = engram, want [engram mcp] (OpenCode schema needs binary+args combined)
    mcp_test.go:524: OpenCode local overlay must not have a separate 'args' key
    mcp_test.go:527: environment = <nil>, se esperaba LOG=1
    mcp_test.go:529: OpenCode no admite env separado de environment
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/external	0.004s
FAIL
```

### Verde previo al commit

```text
Hora: 2026-10-23 09:57:46 -0300
SHA: 84fcf0ebc0a6f647c04b5b0eb4dfbbc9939c866d
Comando: go test ./...
ok  	github.com/Group-Active-IA/active-stack/assets	(cached)
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.705s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	2.869s
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	0.825s
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
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	0.080s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/install	0.299s
ok  	github.com/Group-Active-IA/active-stack/internal/model	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/planner	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/roles	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/system	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/tui	0.037s
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	(cached)
ok  	github.com/Group-Active-IA/active-stack/internal/verify	(cached)
```

### Cierre del paso 1, después del commit

```text
Hora: 2026-10-23 09:58:10 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
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
```


## Suite posterior al commit de esta entrega

```text
Hora: 2026-10-23 09:58:10 -0300
SHA: 8c6a06be9aca371e7d40b9cd5b4d4b8caa6e8d25
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
```
