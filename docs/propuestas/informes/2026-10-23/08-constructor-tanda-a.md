# Constructor, tanda A

Alcance: pasos 0, 1, 2 y 3 del propose firmado en 4748680. Paso 0 limitado a código y diagnóstico por pedido de Juani. Sin reinstalación, instalaciones, push ni lectura de .env. No se tocaron los informes preexistentes sin versionar.

Procedencia: informe escrito por el constructor en esta sesión, con salida textual de herramientas. Proveedor OpenAI; identidad Codex. Modelo/familia exactos sin registro de despacho: SUPUESTO. No se convocaron revisores ni se afirma independencia. Se leyeron CLAUDE.md, propuesta y exploración; no existe knowledge-base/ ni archivo de lecciones local. No se recibieron casos dictados por un auditor en esta entrega.

Las pruebas corren con bwrap: raíz en solo lectura, escritura habilitada en el checkout, /tmp y caché de Go. Los destinos reales ~/.claude, ~/.local/bin y /usr/local/bin quedan en solo lectura. No se hizo go build. Son pruebas locales, no observación de CI remoto.

## Paso 0: instalacion-sana

Diagnóstico estático: internal/install/self_install.go usa os.Executable cuando sourcePath está vacío y copia ese archivo como active-stack. En un test eso puede ser install.test. TestExternalStepCallsExternalInstaller, TestVerifyHookStepRunsAfterApply y TestVerifyHookNilIsNoOp, entre otros, ejecutan planes con el helper buildOptions, que deja NoSelfInstall=false y SelfInstallBinaryPath vacío. No aíslan el resolver del directorio binario. En el código actual pasan HomeDir=t.TempDir(): si /usr/local/bin no es escribible, el destino es ese home temporal, no el home real. Si es escribible, eligen /usr/local/bin. TestBuildPlan_SelfInstallGoosMatchesSnapshotGoos también copia el test, pero con resolver inyectado bajo t.TempDir().

La vía histórica que ignoraba --home está documentada y corregida en self_install.go y TestSelfInstall_WritesBinaryUnderGivenHome. Este último usa un archivo fuente ficticio y dos homes temporales. JUICIO: el auto-install no aislado explica cómo una suite puede copiar un .test; no hay evidencia suficiente para atribuir el archivo real a un test único ni confirmar cuándo ocurrió. No se reprodujo una escritura en destinos reales.

Cambios: version informa versión y revisión disponibles mediante build info; verify valida identidad cmd/active-stack antes de ejecutar roles --help, con timeout. El hook de instalación incorpora el chequeo salvo --no-self-install. Los adapters reales verifican además contenido contra Compose(variante, toggles), incluyendo el prompt primario de OpenCode en settings. Los checks de markers y mode:primary siguen vigentes. No se modificaron backup, filemerge, pipeline ni configuración del usuario. Sin cambios de goldens.

### Rojo observado

```text
Hora: 2026-10-23 09:54:02 -0300
SHA: 474868094a0a69bd8a548b8bff2d17d25f13f400
Comando: go test ./cmd/active-stack -run ^TestVersionSubcommand$ -count=1
--- FAIL: TestVersionSubcommand (0.00s)
    version_test.go:20: version debe salir 0 e identificar el ejecutable: exit status 1, unknown command: "version"
FAIL
FAIL	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.007s
FAIL
```

```text
Hora: 2026-10-23 09:54:19 -0300
SHA: 474868094a0a69bd8a548b8bff2d17d25f13f400
Comando: go test ./internal/verify -run ^TestEntregaComparaCompose$ -count=1
--- FAIL: TestEntregaComparaCompose (0.00s)
    --- FAIL: TestEntregaComparaCompose/claude (0.00s)
        delivery_test.go:38: falla=false, esperada=true: [{ID:config:sdd-orchestrator:claude Description:idempotent marker present exactly once in claude instructions (sdd-orchestrator) Status:passed Error:}]
    --- FAIL: TestEntregaComparaCompose/codex (0.00s)
        delivery_test.go:38: falla=false, esperada=true: [{ID:config:sdd-orchestrator:codex Description:idempotent marker present exactly once in codex instructions (sdd-orchestrator) Status:passed Error:}]
    --- FAIL: TestEntregaComparaCompose/opencode (0.00s)
        delivery_test.go:38: falla=false, esperada=true: [{ID:config:sdd-orchestrator:opencode Description:primary agent "sdd-orchestrator" registered (mode:primary) in opencode settings Status:passed Error:}]
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/verify	0.007s
FAIL
```

```text
Hora: 2026-10-23 09:54:35 -0300
SHA: 474868094a0a69bd8a548b8bff2d17d25f13f400
Comando: go test ./internal/verify -run ^TestBinarioInstalado$ -count=1
FAIL	github.com/Group-Active-IA/active-stack/internal/verify [build failed]
FAIL
```

```text
Hora: 2026-10-23 09:55:46 -0300
SHA: 474868094a0a69bd8a548b8bff2d17d25f13f400
Comando: go test ./cmd/active-stack -run ^TestVerifyInstalacionIncluyeBinario$ -count=1
FAIL	github.com/Group-Active-IA/active-stack/cmd/active-stack [build failed]
FAIL
```

El rojo de binario y wiring es de compilación por API todavía inexistente; version y Compose fallaron en sus aserciones contra el comportamiento viejo.

### Verde previo al commit

```text
Hora: 2026-10-23 09:55:10 -0300
SHA: 474868094a0a69bd8a548b8bff2d17d25f13f400
Comando: go test ./internal/verify ./cmd/active-stack -run TestBinarioInstalado|TestEntregaComparaCompose|TestVersionSubcommand -count=1
ok  	github.com/Group-Active-IA/active-stack/internal/verify	0.011s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.008s
Hora: 2026-10-23 09:56:00 -0300
SHA: 474868094a0a69bd8a548b8bff2d17d25f13f400
Comando: go test ./...
ok  	github.com/Group-Active-IA/active-stack/assets	0.030s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.835s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	2.813s
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	0.904s
ok  	github.com/Group-Active-IA/active-stack/internal/agents	0.012s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	0.009s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	0.013s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	0.007s
ok  	github.com/Group-Active-IA/active-stack/internal/backup	0.114s
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	0.188s
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	0.018s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	0.025s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	2.812s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	0.123s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	0.082s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	0.037s
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	0.013s
ok  	github.com/Group-Active-IA/active-stack/internal/install	0.256s
ok  	github.com/Group-Active-IA/active-stack/internal/model	0.007s
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	0.004s
ok  	github.com/Group-Active-IA/active-stack/internal/planner	0.005s
ok  	github.com/Group-Active-IA/active-stack/internal/roles	0.100s
ok  	github.com/Group-Active-IA/active-stack/internal/system	0.020s
ok  	github.com/Group-Active-IA/active-stack/internal/tui	0.031s
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	0.021s
ok  	github.com/Group-Active-IA/active-stack/internal/verify	0.034s
```
