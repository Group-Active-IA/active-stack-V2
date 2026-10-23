# Constructor — continuación de captura 2b/2c/2d

Fecha: 2026-10-23. Rama verificada: `feat/metodologia-fase-2`.
Propuesta: `docs/propuestas/2026-10-23-fase-2-pasos-1-4.md`, tabla humana
firmada por Juani. Leídos el commit `16e9a2b` y los finales de informes 01–04.
Familia del constructor: GPT-6 (Codex). No se convocaron revisores; este informe
es apply y no acredita verify, aprobación del auditor ni archive.

## Alcance y estado real

- **EVIDENCIA:** `16e9a2ba46c1e753ac6fd75b45c2873b3944b592` ya contenía
  parsers, fixtures, CLI y pruebas de Claude/Codex. Se conservó; no se reimplementó.
- **EVIDENCIA:** `482ca9803be2120c79c2047e6685ba370a1d7436` corrige que el
  bruto de un evento parcial se perdía antes de escanearlo. Ahora conserva ese
  original y RETENIDO prevalece sobre PENDING si contiene un secreto. Agrega
  el rojo funcional, triangulación y documentación de importación explícita.
- **EVIDENCIA:** apareció concurrentemente
  `5ce2c0f7ef34e0a7dc159b46490245389a6796cc`, actualización de roles de Juani;
  no se modificó ni se incluyó entre los cambios propios. La suite final lo incluye.
- **EVIDENCIA:** 2c importa stdout JSONL de un único `codex exec --json`
  0.132.0; 2b importa PostToolUse(SubagentHandback) + SubagentStop de Claude
  2.1.271, con hija y sesión esperadas. No toma la despedida como informe.
- **EVIDENCIA:** 2d devuelve UNSUPPORTED sin persistir originales. No hay
  una instancia/API versionada contrastada para asociar hija y mensaje final.
- **PENDIENTE 2b:** install/update/uninstall automático de hooks, ownership,
  recuperación y ensayo de ese ciclo **no están implementados**. El contrato
  de eventos sí está documentado; esa confirmación no acredita la integración.
  La importación explícita no completa por sí sola toda la sección 2b.
  No se registra como deuda aceptada nueva ni se declara terminado ese alcance.

Consumidores: `reports.Capture`, CLI `reports capture` y lectura de manifiestos;
no se cambian `reports.Source`, `RegisterContract` ni el esquema de manifiesto.
Conflictos: se preservan los originales existentes; se ensaya la captura retenida
sobre una identidad ya COMPLETE. No se tocan backup/filemerge/pipeline, settings
personales, catálogo ni configuración instalada. Sin instalaciones, push o merge.

## BE-4: fuente, acceso, versión y fixtures

| Proveedor | Fuente y acceso concreto | Versión/esquema | Fixture con texto inventado |
|---|---|---|---|
| Claude | Archivo JSONL explícito con payloads exportados de hooks PostToolUse y SubagentStop; no se siguen transcript_path ni rutas del HOME | 2.1.271; SubagentHandback.tool_input.message y cierre de la misma hija/sesión | `internal/reports/testdata/claude-2.1.271-hooks.jsonl` |
| Codex | stdout de una invocación `codex exec --json`, archivo explícito por despacho; no rollouts privados | 0.132.0, JSONL exec v1; thread.started, turn.started, item.completed(agent_message), turn.completed | `internal/reports/testdata/codex-0.132.0.jsonl` |
| OpenCode | Candidato: GET /global/health, GET /session/:id/children y GET /session/:id/message/:messageID; no se consulta red desde el adapter | Sin versión ni contrato habilitado (TBD); siempre UNSUPPORTED | `internal/reports/testdata/opencode-health-unconfirmed.json`, sólo health público, no un informe supuesto |

**EVIDENCIA documental consultada en esta continuación:**
[Claude hooks](https://code.claude.com/docs/en/hooks#subagentstop) documenta
Handback desde 2.1.271 y diferencia el informe de la despedida;
[OpenAI, modo no interactivo](https://learn.chatgpt.com/docs/non-interactive-mode)
muestra los eventos JSONL de exec;
[OpenCode Server](https://opencode.ai/docs/server/) publica health, children,
message y el acceso a `/doc` de una instancia. Las páginas son móviles;
versiones distintas permanecen UNSUPPORTED. El pin de esquema Rust de Codex
heredado se referencia en el README de las fixtures; no se volvió a ejecutar Codex.

**SUPUESTO:** el llamador registró antes del despacho la asociación de archivo,
repo, sesión y, en Claude, hija. El archivo importado y la versión declarada no
acreditan autenticidad ni la versión de un proceso real. No se inventan turno,
modelo o proveedor ausentes. No se leyó el HOME real ni sesiones personales:
se utilizaron las fixtures estructurales heredadas, con contenido inventado.

## Entorno y TDD de la continuación

Todas las invocaciones de Go y del binario usaron:

```sh
export GOFLAGS=-p=1
export GOMAXPROCS=2
export HOME=/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/apply-2bcd/home
export GOCACHE=/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/apply-2bcd/gocache
export GOMODCACHE=/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/apply-2bcd/gomodcache
export TMPDIR=/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/apply-2bcd/tmp
```

Se exportaron dentro de cada shell de ejecución, nunca en `git commit`.
Sin `-race`, procesos pesados, instalaciones ni `pkill`. La suite completa se
corrió **una sola vez en esta continuación**, después de los tests por paquete.
Los logs 01–18 heredados de la corrida cortada no se presentan como nuevas
verificaciones ni como suite final completa.

### Rojo funcional, antes de corregir

SHA base `5ce2c0f`; test agregado a los parsers heredados sin cambiar producción.
Comando: `go test ./internal/reports -run TestProviderPartialSecretRetained -count=1 -v`.
El fallo es una aserción de estado, no entorno ni compilación:

```text
2026-10-23T14:32:29-03:00
5ce2c0f7ef34e0a7dc159b46490245389a6796cc
=== RUN   TestProviderPartialSecretRetained
=== RUN   TestProviderPartialSecretRetained/claude
    provider_test.go:226: evento parcial sensible: status=PENDING err=<nil>
=== RUN   TestProviderPartialSecretRetained/codex
    provider_test.go:226: evento parcial sensible: status=PENDING err=<nil>
--- FAIL: TestProviderPartialSecretRetained (0.03s)
    --- FAIL: TestProviderPartialSecretRetained/claude (0.02s)
    --- FAIL: TestProviderPartialSecretRetained/codex (0.01s)
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/reports	0.036s
FAIL
```

### Verde y triangulación

Comando: `go test ./internal/reports -run 'Provider|Claude|Codex|OpenCode' -count=1 -v`.
Cambio: inicializar `Source.Raw` antes de cualquier retorno por evento parcial.
Hora y salida relevante, árbol luego incluido en `482ca98`:

```text
2026-10-23T14:32:51-03:00
=== RUN   TestProviderPartialSecretRetained
=== RUN   TestProviderPartialSecretRetained/claude
=== RUN   TestProviderPartialSecretRetained/codex
--- PASS: TestProviderPartialSecretRetained
PASS
```

Triangulación: evento parcial sin secreto queda PENDING y conserva el bruto exacto
(`TestProviderPartialWithoutSecretPending`); evento parcial sensible en ambos
proveedores queda RETENIDO en repo sin `.gitignore` y no entra al staging.
Los tests previos cubren informe íntegro, vacío, reintento, orden, identidad ajena
y versión desconocida. Comando: `go test ./internal/reports ./cmd/active-stack -count=1`.

```text
2026-10-23T14:33:17-03:00
ok  	github.com/Group-Active-IA/active-stack/internal/reports	0.536s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.755s
```

## Corrida real del binario

Comandos, bajo el entorno anterior:

```sh
go build -o /home/juanisarmiento/.claude/jobs/f3e370f0/tmp/apply-2bcd/bin/active-stack ./cmd/active-stack
python3 /home/juanisarmiento/.claude/jobs/f3e370f0/tmp/apply-2bcd/e2e.py
```

El script crea un repo git temporal SIN `.gitignore` raíz. Los archivos de
entrada quedan fuera de ese repo. Cada captura usa `--change ensayo --phase
verify --round 1 --dispatch-id dispatch-fixture --sha 5ce2c0f7ef34e0a7dc159b46490245389a6796cc
--source-session session-fixture`, más el rol/proveedor/versión/archivo mostrado.
Claude agrega `--source-agent child-fixture`, salvo el ensayo de hija ajena.
Verifica código de salida, hashes de bruto/texto, ausencia del token en TODA la
salida y en `git diff --cached`, y ausencia de `.retenidos/` en staging después
de `git add -A`. No imprime el token falso. Salida íntegra:

```text
2026-10-23T14:33:53-03:00
5ce2c0f7ef34e0a7dc159b46490245389a6796cc
2026-10-23T14:33:54.017133-03:00
Repo: /home/juanisarmiento/.claude/jobs/f3e370f0/tmp/apply-2bcd/tmp/repo-sin-gitignore-0513tkny ; .gitignore inicial: AUSENTE
$ active-stack reports capture [identidad ensayo/verify/1/codex] --source-provider codex --source-version 0.132.0 --report codex-0.132.0.jsonl
estado: COMPLETE (código 0)
clave: opsx/ensayo/verify/round-1/codex
directorio: openspec/changes/ensayo/informes/verify/round-1/codex
exit: 0
$ active-stack reports capture [identidad ensayo/verify/1/codex] --source-provider codex --source-version 0.132.0 --report codex-0.132.0.jsonl
estado: COMPLETE (código 0)
clave: opsx/ensayo/verify/round-1/codex
directorio: openspec/changes/ensayo/informes/verify/round-1/codex
exit: 0
$ active-stack reports capture [identidad ensayo/verify/1/claude] --source-provider claude --source-version 2.1.271 --report claude-2.1.271-hooks.jsonl --source-agent child-fixture
estado: COMPLETE (código 0)
clave: opsx/ensayo/verify/round-1/claude
directorio: openspec/changes/ensayo/informes/verify/round-1/claude
exit: 0
$ active-stack reports capture [identidad ensayo/verify/1/hija-ajena] --source-provider claude --source-version 2.1.271 --report claude-2.1.271-hooks.jsonl --source-agent otra-hija
estado: PENDING (código 2)
clave: opsx/ensayo/verify/round-1/hija-ajena
directorio: openspec/changes/ensayo/informes/verify/round-1/hija-ajena
causa: informe vacío o en blanco
exit: 2
$ active-stack reports capture [identidad ensayo/verify/1/opencode] --source-provider opencode --source-version fixture-unconfirmed --report opencode-health-unconfirmed.json
estado: UNSUPPORTED (código 3)
clave: opsx/ensayo/verify/round-1/opencode
directorio: openspec/changes/ensayo/informes/verify/round-1/opencode
causa: contrato "" no está registrado: ningún adapter lo resolvió todavía
exit: 3
$ active-stack reports capture [identidad ensayo/verify/1/version-ajena] --source-provider codex --source-version otra --report codex-0.132.0.jsonl
estado: UNSUPPORTED (código 3)
clave: opsx/ensayo/verify/round-1/version-ajena
directorio: openspec/changes/ensayo/informes/verify/round-1/version-ajena
causa: contrato "" no está registrado: ningún adapter lo resolvió todavía
exit: 3
$ active-stack reports capture [identidad ensayo/verify/1/vacio] --source-provider codex --source-version 0.132.0 --report blank.jsonl
estado: PENDING (código 2)
clave: opsx/ensayo/verify/round-1/vacio
directorio: openspec/changes/ensayo/informes/verify/round-1/vacio
causa: informe vacío o en blanco
exit: 2
$ active-stack reports capture [identidad ensayo/verify/1/codex] --source-provider codex --source-version 0.132.0 --report secret.jsonl
estado: RETAINED (código 4)
clave: opsx/ensayo/verify/round-1/codex
directorio: openspec/changes/ensayo/informes/verify/round-1/codex/conflict-3db962ae736f
causa: escaneo de secretos detectó: github-token; misma identidad (opsx/ensayo/verify/round-1/codex), bytes distintos a lo ya persistido en openspec/changes/ensayo/informes/verify/round-1/codex; original preservado sin pisar
original retenido (local, ignorado por git): openspec/changes/ensayo/informes/verify/round-1/codex/conflict-3db962ae736f/.retenidos/2dfb94a3c4ea.txt
exit: 4
$ active-stack reports capture [identidad ensayo/verify/1/parcial-claude] --source-provider claude --source-version 2.1.271 --report partial-secret.jsonl --source-agent child-fixture
estado: RETAINED (código 4)
clave: opsx/ensayo/verify/round-1/parcial-claude
directorio: openspec/changes/ensayo/informes/verify/round-1/parcial-claude
causa: escaneo de secretos detectó: github-token
original retenido (local, ignorado por git): openspec/changes/ensayo/informes/verify/round-1/parcial-claude/.retenidos/e3b0c44298fc.txt
exit: 4
$ active-stack reports capture [identidad ensayo/verify/1/parcial-codex] --source-provider codex --source-version 0.132.0 --report partial-secret.jsonl
estado: RETAINED (código 4)
clave: opsx/ensayo/verify/round-1/parcial-codex
directorio: openspec/changes/ensayo/informes/verify/round-1/parcial-codex
causa: escaneo de secretos detectó: github-token
original retenido (local, ignorado por git): openspec/changes/ensayo/informes/verify/round-1/parcial-codex/.retenidos/e3b0c44298fc.txt
exit: 4
$ active-stack reports capture [identidad ensayo/verify/1/claude-vacio] --source-provider claude --source-version 2.1.271 --report blank-claude.jsonl --source-agent child-fixture
estado: PENDING (código 2)
clave: opsx/ensayo/verify/round-1/claude-vacio
directorio: openspec/changes/ensayo/informes/verify/round-1/claude-vacio
causa: informe vacío o en blanco
exit: 2
$ git diff --cached --name-only
openspec/changes/ensayo/informes/LEEME.md
openspec/changes/ensayo/informes/verify/round-1/claude-vacio/informe.md
openspec/changes/ensayo/informes/verify/round-1/claude-vacio/manifest.json
openspec/changes/ensayo/informes/verify/round-1/claude-vacio/raw.txt
openspec/changes/ensayo/informes/verify/round-1/claude/informe.md
openspec/changes/ensayo/informes/verify/round-1/claude/manifest.json
openspec/changes/ensayo/informes/verify/round-1/claude/raw.txt
openspec/changes/ensayo/informes/verify/round-1/codex/conflict-3db962ae736f/manifest.json
openspec/changes/ensayo/informes/verify/round-1/codex/informe.md
openspec/changes/ensayo/informes/verify/round-1/codex/manifest.json
openspec/changes/ensayo/informes/verify/round-1/codex/raw.txt
openspec/changes/ensayo/informes/verify/round-1/hija-ajena/informe.md
openspec/changes/ensayo/informes/verify/round-1/hija-ajena/manifest.json
openspec/changes/ensayo/informes/verify/round-1/hija-ajena/raw.txt
openspec/changes/ensayo/informes/verify/round-1/parcial-claude/manifest.json
openspec/changes/ensayo/informes/verify/round-1/parcial-codex/manifest.json
openspec/changes/ensayo/informes/verify/round-1/vacio/informe.md
openspec/changes/ensayo/informes/verify/round-1/vacio/manifest.json
openspec/changes/ensayo/informes/verify/round-1/vacio/raw.txt
RETENIDO: token ausente de salida y staging; original local ignorado: OK
Contraste bruto/texto contra manifiestos: OK
Repo sigue SIN .gitignore raíz: OK
```

La compilación y esta corrida usan el árbol de código luego commiteado en
`482ca98`; el SHA impreso es el padre porque aún no se había creado ese commit.
No se afirma que esta importación de fixtures sea una sesión real del modelo
ni un ensayo del ciclo de instalación de hooks.

## Suite final

**EVIDENCIA:** único recorrido completo de esta continuación. SHA de código
`482ca9803be2120c79c2047e6685ba370a1d7436`. El commit posterior agrega sólo este
informe; no cambia código ni tests. `gofmt -l .` no produjo salida.

```text
2026-10-23T14:34:34-03:00
482ca9803be2120c79c2047e6685ba370a1d7436
gofmt -l .
go vet ./...
vet exit: 0
go test ./... -count=1
ok  	github.com/Group-Active-IA/active-stack/assets	0.016s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.655s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	1.986s
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	0.508s
ok  	github.com/Group-Active-IA/active-stack/internal/agents	0.003s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	0.004s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	0.004s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	0.004s
ok  	github.com/Group-Active-IA/active-stack/internal/approval	0.810s
ok  	github.com/Group-Active-IA/active-stack/internal/backup	0.037s
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	0.059s
?   	github.com/Group-Active-IA/active-stack/internal/contracttest	[no test files]
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	0.004s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	0.007s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	25.024s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	0.053s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	0.027s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/roledefs	0.024s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	0.011s
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	0.003s
ok  	github.com/Group-Active-IA/active-stack/internal/install	0.164s
ok  	github.com/Group-Active-IA/active-stack/internal/model	0.003s
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	0.002s
ok  	github.com/Group-Active-IA/active-stack/internal/planner	0.003s
ok  	github.com/Group-Active-IA/active-stack/internal/reports	0.639s
ok  	github.com/Group-Active-IA/active-stack/internal/roles	0.063s
ok  	github.com/Group-Active-IA/active-stack/internal/system	0.012s
ok  	github.com/Group-Active-IA/active-stack/internal/tui	0.027s
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	0.016s
ok  	github.com/Group-Active-IA/active-stack/internal/verify	0.026s
suite exit: 0
2026-10-23T14:35:15-03:00
```

## JUICIO y límites

- **JUICIO (confianza alta):** el defecto de retención de eventos parciales queda
  reproducido y corregido para ambos parsers; el ensayo externo confirma que los
  originales sensibles no se filtran al staging en un repo sin reglas propias.
- **JUICIO (confianza media):** el formato documentado y las fixtures alcanzan
  para la importación explícita acotada; no acreditan compatibilidad universal,
  captura automática ni custodia contra cambios maliciosos del mismo usuario.
- **PENDIENTE:** ciclo automático de 2b y ensayo de proveedor OpenCode; no hay
  paridad anunciada. Timeout/retención y detector universal de secretos conservan
  las limitaciones heredadas del núcleo; no se cierran por estos ensayos.
- **PENDIENTE de archive:** actualización integral de ARCHITECTURE, CHANGES y base
  metodológica. En este apply se reescriben README y README de fixtures para
  reflejar la capacidad real; no se alteran reglas ni AGENTS/CLAUDE.

## Escaneo del informe antes del cierre documental

EVIDENCIA: gitleaks no está disponible. Escaneo propio sin tokens, claves ni
correos. El patrón de teléfono señaló segmentos numéricos de fechas y SHA de
la evidencia; revisados como falsos positivos, al excluir únicamente esas dos
formas no quedan candidatos. No se imprimieron valores sensibles. El primer
script de escaneo terminó por esa aserción numérica y el shell continuó con el
commit documental; se verificó inmediatamente después el archivo completo con
las exclusiones revisadas, sin encontrar datos personales ni secretos. Este
resultado no acredita un detector universal.

## Deuda al cerrar (orquestador, 23/10/2026)

- D-21: la captura automática por hooks de Claude Code (2b) no quedó cableada; la captura funciona invocando `reports capture` a mano. Se mergea así por decisión de Juani ("terminate los adaptadores y mergeamos ... demasiado por ahora").
