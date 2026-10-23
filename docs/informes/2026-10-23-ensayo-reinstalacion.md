# Ensayo aislado de reinstalación de Active Stack — fase 1

Fecha: 23 de octubre de 2026. Zona: America/Argentina/Buenos_Aires (UTC−03).
Final de las comprobaciones: 2026-10-23T12:27:07.778275-03:00.
Directorio exclusivo del ensayo: `/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/as-review/ensayo`.

## Resultado

Instalación headless **Full / Claude**, reinstalación idempotente, falla inducida con recuperación del estado comprobado y restauración del CLAUDE.md original: comprobados. Verify explícito: **17 aprobados, 0 fallidos, 0 advertencias, 1 omitido**. El omitido es conectividad de Context7, cuyo check no está implementado. Esto acredita el ensayo; no autoriza ni acredita reinstalación en la configuración real.

La falla devolvió código 1 y el CLI anunció rollback **partial** por una limitación declarada sobre herramientas externas. La comparación independiente de **116 archivos**, directorios y permisos del HOME (excepto el almacén de backups) fue idéntica antes/después. Engram ya existía por la primera instalación; no se instaló de nuevo en la corrida fallida.

## Aislamiento y procedencia

Se leyó README y AGENTS.md del proyecto, las opciones de CLI y el código necesario para resolver instalación, backups y verify. El checkout del repositorio estaba en `fe0294d9e1a8f818c03e0a26a5bee4afaaaa9424`, con archivos ajenos sin seguimiento; no se modificó. La rama **main** resolvía a **020402249c4d7363d62d025bea8db778fed1ce6d**, exactamente el commit solicitado. Para no compilar el checkout distinto ni escribir en el repositorio, se exportó main a `src/` mediante `git archive main` y se compiló esa copia.

Del HOME real se copiaron exclusivamente `.claude/CLAUDE.md` y 20 archivos Markdown de `.claude/agents/`; listado en [logs/copia.json](logs/copia.json). La copia excluye enlaces simbólicos, `.credentials.json`, `*.key`, nombres que comienzan por `.env` y nombres que contienen token. No se copió settings, skills, historial, ni credenciales del usuario. La lista de nombres versionados tampoco contiene `.credentials.json`, `*.key` ni archivos que comiencen por `.env`.

La prohibición general de escribir en `.claude` se interpretó con la excepción explícita de este directorio de ensayo, que está dentro de ella. No se escribieron rutas de configuración reales fuera de esa excepción.

**HOME falso solo no alcanza:** `BinaryInstallDirForHome` prueba `/usr/local/bin` creando un temporal, y npm puede intentar instalación global. Se utilizó Bubblewrap: raíz de solo lectura, HOME real oculto por tmpfs, sólo el ensayo montado con escritura; `/dev`, `/proc` y `/opt` son montajes privados. El runtime Node existente se expuso de solo lectura en `/opt/node`. No se instalaron paquetes globales reales ni se concedió sudo. El entorno de ejecución se limpia (`--clearenv`), sin heredar tokens; git no usa configuración global ni pide autenticación. El archivo ejecutable [run-sandbox](run-sandbox) contiene el aislamiento exacto.

Variables de ejecución: HOME=`/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/as-review/ensayo/home`, TMPDIR=`/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/as-review/ensayo/tmp`, npm_config_prefix=`/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/as-review/ensayo/home/.local`, npm_config_cache=`/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/as-review/ensayo/cache/npm`; PATH sólo incluye el bin falso, runtime Node y `/usr/bin:/bin`. Los binarios Engram y Active Stack quedan bajo el HOME falso. OpenSpec se detecta en el runtime Node expuesto, de solo lectura; no se reinstala allí.

El SHA256 del CLAUDE.md real al terminar sigue siendo el del original copiado: `25d2ecaf5a3f369e44c6e4aecb6d1f5fd5f646c5d97c378910d9d91955e6dfb7`. Evidencia: [logs/real-sin-cambios.txt](logs/real-sin-cambios.txt).

## Compilación

Comandos desde el directorio del ensayo:

```bash
mkdir -p src bin cache/go cache/mod logs home/.claude tmp
git -C /home/juanisarmiento/Proyectos/active-stack-V2 rev-parse main
git -C /home/juanisarmiento/Proyectos/active-stack-V2 archive main | tar -x -C src --exclude='.env' --exclude='*.key' --exclude='*token*' --exclude='.credentials.json'
HOME="$PWD/home" TMPDIR="$PWD/tmp" GOCACHE="$PWD/cache/go" GOMODCACHE="$PWD/cache/mod" GOPATH="$PWD/cache/gopath" bash -c 'cd src && go build -o ../bin/active-stack ./cmd/active-stack'
python3 setup.py
```

Build exitoso; Go descargó el toolchain 1.26.0 y dependencias en las cachés del ensayo. Repetición registrada en `logs/build.time`, `logs/build.txt`, `logs/build.exit`: código 0 y sin salida. SHA256 binario: `e878eb5240fdaf5ad0df00b92578e1d403dc74114db89ed43cd37c6bfd1d2fda`.

Incidencias de preparación: primer intento de build desde la raíz del ensayo falló por falta de go.mod; se corrigió ejecutando desde src. El primer armado del sandbox no pudo crear `/opt/node` sobre raíz de solo lectura; se añadió `/opt` privado. Una primera ejecución de instalación no pudo abrir `/dev/null`, falló en la primera skill y anunció rollback succeeded; se corrigió con `/dev` privado. Su snapshot se conserva en `backup-intento-inicial/`. Ninguna de estas incidencias se presenta como instalación exitosa.

## Qué cambiaría en el CLAUDE.md real de Juani

**Tu texto personal queda intacto, byte a byte fuera del bloque administrado.** Se mantienen identidad, idioma, reglas core, stack, referencias a skills (incluida Impeccable), memoria Engram, Cerebro Digital/Obsidian y su estructura, política de informes, workflows pesados y autorización personal de delegación. Los 20 agentes copiados también conservaron contenido en la comparación del ensayo.

Se **reemplaza entero el contenido del bloque existente** entre `<!-- active-stack:sdd-orchestrator -->` y `<!-- /active-stack:sdd-orchestrator -->`. No se agrega una segunda copia del orquestador. Aunque el diff muestra altas y bajas, la unidad que administra el instalador es ese bloque completo, no cada párrafo por separado. Los fragmentos se componen con los 12 toggles de Full: delegation, model-routing, governance, evidence, auditor, roles, review, knowledge-base, senior, ship, engram y tdd.

Dentro de ese bloque:

- Se agrega una sección de coordinación (“Dispatch”): alcance claro de cada agente, dependencias, especialistas según el diff, controles mecánicos una sola vez y registro de lo delegado. Se agregan reglas para agrupar correcciones por zona.
- Se incorpora aprobación humana escrita para propose de riesgo alto/crítico, con tabla de aprobación, y se exige reproducir hallazgos antes de llamarlos bloqueantes.
- Se actualiza el trabajo del auditor: decisión explícita y umbral de cierre; se amplía cómo elegir y contar revisores según el diff y qué informes deben existir antes de firmar.
- Se reemplaza la revisión de tres pases independientes por **un rompedor con tres lentes obligatorias**, bucle de hasta tres vueltas, decisión del auditor y tratamiento explícito de operación degradada e informes completos.
- Se amplía archive: reescribir y agregar a la base de conocimiento, registrar pendientes y deuda, e informe de knowledge base como condición de cierre.
- Se reemplaza el alcance del senior “por fase” por participación **sólo en verify**, con sus límites; se ajustan partes de su texto previo.
- Se agrega exigencia de esperar todo el CI del SHA publicado en ship.
- Se actualiza TDD/apply: brief del constructor, consumidores y límites, evidencia rojo/verde por tarea, cierre sobre el commit entregado, correcciones por lista del auditor y familia del modelo. Se reemplaza la vieja plantilla de resumen TDD.
- Hay ajustes menores en referencias de modelos y redacción; el diff completo permite revisarlos sin omisiones.

**¿Se pierde algo suyo?** No se perdió texto personal fuera del bloque: comprobación exacta del prefijo y sufijo, ambos iguales. Sí desaparecen/reformulan **123 líneas del bloque administrado** y aparecen **653 líneas**. Si Juani había personalizado frases dentro de ese bloque, esas frases están sujetas a reemplazo; el instalador no distingue autoría dentro de los markers y este ensayo no permite atribuir cada línea retirada. El original completo y el diff quedan conservados para revisar esas bajas. El archivo pasa de 1.702 a 2.232 líneas (86.988 a 124.505 bytes).

Diff solicitado: [diff-CLAUDE.md.diff](diff-CLAUDE.md.diff), 1.042 líneas. Copias: [CLAUDE.original.md](CLAUDE.original.md) y [CLAUDE.instalado.md](CLAUDE.instalado.md). Análisis mecánico: [logs/analisis-diff.json](logs/analisis-diff.json).

## Instalación inicial, snapshot y manifiesto

Comandos:

```bash
./run-sandbox ./bin/active-stack install --help
./run-sandbox ./bin/active-stack install --mode full --agent claude --yes --dry-run
./run-sandbox ./bin/active-stack install --mode full --agent claude --yes
cp -a home/.active-stack/backups/install backup-primera-instalacion
cp home/.claude/CLAUDE.md CLAUDE.instalado.md
diff -u --label CLAUDE.original.md --label CLAUDE.instalado.md CLAUDE.original.md CLAUDE.instalado.md > diff-CLAUDE.md.diff
```

`--mode` y `--agent` activan headless. Se deja self-install activo para probar también la copia del binario, bajo aislamiento. Tier omitido: balanceado por defecto; no se leyó el tier ni settings reales de Juani. `--help` muestra las opciones pero termina con `error: flag: help requested`; es comportamiento del CLI.

El snapshot es **install**, creado a las **2026-10-23T15:24:33.505324542Z** (12:24:33 local), en `/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/as-review/ensayo/home/.active-stack/backups/install`. Es comprimido (`archive.tar.gz` y `manifest.json`); sólo había un archivo existente a respaldar, CLAUDE.md. También registra destinos ausentes para poder eliminarlos en restore.

**El slot no es histórico:** las siguientes instalaciones reemplazan `install/`. Por eso se preservó su primera versión en [backup-primera-instalacion/manifest.json](backup-primera-instalacion/manifest.json), junto con su archive. Su `root_dir` sigue apuntando al slot original: para restaurarlo mediante CLI se repuso la copia en ese slot, sin editar el manifiesto. También se preservó `backup-fallo/` antes del restore. El slot actual corresponde a la instalación final.

Manifiesto original completo:

```json
{
  "id": "install",
  "created_at": "2026-10-23T15:24:33.505324542Z",
  "root_dir": "/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/as-review/ensayo/home/.active-stack/backups/install",
  "entries": [
    {
      "original_path": "/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/as-review/ensayo/home/.claude/skills",
      "snapshot_path": "",
      "existed": false,
      "is_dir": true
    },
    {
      "original_path": "/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/as-review/ensayo/home/.claude.json",
      "snapshot_path": "",
      "existed": false
    },
    {
      "original_path": "/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/as-review/ensayo/home/.claude/CLAUDE.md",
      "snapshot_path": "files/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/as-review/ensayo/home/.claude/CLAUDE.md",
      "existed": true,
      "mode": 436
    },
    {
      "original_path": "/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/as-review/ensayo/home/.claude/settings.json",
      "snapshot_path": "",
      "existed": false
    },
    {
      "original_path": "/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/as-review/ensayo/home/.claude/commands",
      "snapshot_path": "",
      "existed": false,
      "is_dir": true
    },
    {
      "original_path": "/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/as-review/ensayo/home/.local/bin/active-stack",
      "snapshot_path": "",
      "existed": false
    }
  ],
  "file_count": 1,
  "compressed": true,
  "checksum": "7f54593b74d6b5487b5bee42e8464c7aac1b90edc2b612747d0125d18eb2bd49"
}
```

## Reinstalación y ausencia de duplicados

```bash
./run-sandbox ./bin/active-stack install --mode full --agent claude --yes
sha256sum CLAUDE.instalado.md home/.claude/CLAUDE.md
```

Código 0. Ambos SHA256: `ad19611ebd537a0bdca62fa0616c8bc10d6e4cff931fe5014e23d0ce982c2de0`. Un marker de apertura y uno de cierre del orquestador; contenido idéntico en la segunda instalación. El check de composición verifica además igualdad con Compose. El recuento “0 checks” del resumen headless no describe los checks ejecutados por el hook: el hook se ejecutó exitosamente, pero su informe no se propaga a ese resumen. Para evitar depender de esa salida vacía se ejecutó también verify explícito.

## Falla a mitad y rollback

Se partió de la instalación exitosa, se repuso CLAUDE.md original para obligar una actualización observable y se dejó `commands/active/init.md` con un texto centinela previo. Luego se montó **sólo commands/** como de solo lectura dentro de otra ejecución Bubblewrap; el resto del HOME siguió escribible. No se modificó el instalador. Se prefirió un montaje de solo lectura porque `WriteFileAtomic` hace chmod del directorio y puede neutralizar una prueba basada sólo en permisos Unix.

```bash
cp run-sandbox run-sandbox-fallo
# run-sandbox-fallo agrega antes de --clearenv:
# --ro-bind "$BASE/home/.claude/commands" "$BASE/home/.claude/commands"
cp CLAUDE.original.md home/.claude/CLAUDE.md
printf 'Contenido previo para probar rollback\n' > home/.claude/commands/active/init.md
python3 estado.py logs/estado-antes-fallo.json
./run-sandbox-fallo ./bin/active-stack install --mode full --agent claude --yes
python3 estado.py logs/estado-despues-fallo.json
diff -u logs/estado-antes-fallo.json logs/estado-despues-fallo.json > logs/rollback.diff
```

Snapshot y pasos previos, incluido `config:sdd-orchestrator`, terminaron; falló `command:starter-add-command` con `read-only file system`. Código 1 esperado. El snapshot restaura CLAUDE.md; el centinela anterior queda igual. `rollback.diff` está vacío: mismo conjunto de archivos y directorios, hashes y permisos. Se excluyeron `.active-stack/` y sus nuevos backups, que son evidencia conservada deliberadamente. No se comparan mtime ni procesos externos; el alcance es estado de archivos y permisos, no toda la máquina. El mensaje “brew formula stays installed” es una advertencia genérica del rollback del harness Engram, no evidencia de que se haya ejecutado brew en esta corrida.

## Restore byte a byte

Se repuso `backup-primera-instalacion/` en `home/.active-stack/backups/install/`, tras preservar el snapshot de la falla. Se usó el comando público de backup que existe en este binario bajo el namespace windows; funciona en Linux y llama a `backup.RestoreService`. No se construyó una restauración manual de CLAUDE.md para acreditar esta prueba.

```bash
./run-sandbox ./bin/active-stack windows backups list
./run-sandbox ./bin/active-stack windows backups restore --id install
sha256sum CLAUDE.original.md home/.claude/CLAUDE.md
cmp CLAUDE.original.md home/.claude/CLAUDE.md
```

Restore código 0, respuesta success=true; cmp código 0. Original y restaurado: `25d2ecaf5a3f369e44c6e4aecb6d1f5fd5f646c5d97c378910d9d91955e6dfb7`. Así se verificó igualdad **byte a byte**, no sólo igualdad visual. Después se instaló de nuevo Full para dejar una instalación sana y poder ejecutar el verify final. El HOME falso termina instalado; el original restaurado queda acreditado en los logs y conservado aparte.

## Verify de instalación sana y versión

El binario no tiene subcomando `verify`. Su instalación ejecuta el hook de `internal/verify`. Para registrar cada resultado se agregó únicamente a la copia de fuentes del ensayo un ejecutable auxiliar en [src/cmd/ensayo-verify/main.go](src/cmd/ensayo-verify/main.go): carga el catálogo Full, el adapter Claude, `ChecksForHarness`, `BinaryChecks`, `RunChecks`, `BuildReport` y `RenderReport` de producción. No cambia paquetes productivos ni suplanta checks con mocks.

```bash
HOME="$PWD/home" TMPDIR="$PWD/tmp" GOCACHE="$PWD/cache/go" GOMODCACHE="$PWD/cache/mod" GOPATH="$PWD/cache/gopath" bash -c 'cd src && go build -o ../bin/ensayo-verify ./cmd/ensayo-verify'
./run-sandbox ./bin/active-stack install --mode full --agent claude --yes
./run-sandbox ./bin/ensayo-verify
./run-sandbox ./home/.local/bin/active-stack roles --help
./run-sandbox ./home/.local/bin/active-stack version
```

Todos código 0. Verify verifica identidad Go del binario (`cmd/active-stack`), `roles --help`, skills, permisos, configuración MCP y bloque igual a Compose. No prueba interacción real con Claude ni consulta MCP remota: Context7 está omitido por falta de implementación.

Versión: `active-stack (devel) (SHA desconocido)`. Es una limitación real de compilar la exportación sin `.git`: el runtime no tiene metadata VCS. No se inyectó una versión ficticia; la procedencia se acredita con main resuelto, exportación del commit y hash del binario. El ejecutable instalado pasó el check de identidad, aunque el texto de version no permite inferir el SHA.

## Comandos reproducibles y evidencia íntegra

Los scripts [setup.py](setup.py), [run-sandbox](run-sandbox), [run-sandbox-fallo](run-sandbox-fallo), [estado.py](estado.py) y el auxiliar Go dejan disponibles las operaciones exactas. Las horas siguientes se capturaron justo antes de cada etapa; las salidas y códigos están en logs. Para compilación se registró también una repetición final exitosa; no se inventó una hora exacta del primer build.

### build

Hora: 2026-10-23T12:27:07-03:00. Código: 0.

Comando: `go build -o ../bin/active-stack ./cmd/active-stack (desde src, variables locales detalladas arriba)`

```text
(sin salida)
```

### install-1

Hora: 2026-10-23T12:24:33-03:00. Código: 0.

Comando: `./run-sandbox ./bin/active-stack install --mode full --agent claude --yes`

```text
→ snapshot running...
  ✓ snapshot
  → skill:agent-instruction running...
  ✓ skill:agent-instruction
  → external:context7 running...
  ✓ external:context7
  → external:engram running...
  ✓ external:engram
  → skill:find-skill running...
  ✓ skill:find-skill
  → skill:kb-creator running...
  ✓ skill:kb-creator
  → external:openspec running...
  ✓ external:openspec
  → skill:active-orchestrator running...
  ✓ skill:active-orchestrator
  → permissions:permissions running...
  ✓ permissions:permissions
  → skill:project-scaffold running...
  ✓ skill:project-scaffold
  → skill:roadmap-generator running...
  ✓ skill:roadmap-generator
  → config:sdd-orchestrator running...
  ✓ config:sdd-orchestrator
  → skill:skill-creator running...
  ✓ skill:skill-creator
  → skill:skill-registry running...
  ✓ skill:skill-registry
  → command:starter-add-command running...
  ✓ command:starter-add-command
  → skill:web-scraper running...
  ✓ skill:web-scraper
  → skill:discovery-research running...
  ✓ skill:discovery-research
  → self-install running...
  ✓ self-install
  → verify-hook running...
  ✓ verify-hook
Verification checks: 0 passed, 0 failed, 0 warnings, 0 skipped
Installation verified. You're ready to start building.
```

### install-2

Hora: 2026-10-23T12:25:04-03:00. Código: 0.

Comando: `./run-sandbox ./bin/active-stack install --mode full --agent claude --yes`

```text
→ snapshot running...
  ✓ snapshot
  → skill:agent-instruction running...
  ✓ skill:agent-instruction
  → external:context7 running...
  ✓ external:context7
  → external:engram running...
  ✓ external:engram
  → skill:find-skill running...
  ✓ skill:find-skill
  → skill:kb-creator running...
  ✓ skill:kb-creator
  → external:openspec running...
  ✓ external:openspec
  → skill:active-orchestrator running...
  ✓ skill:active-orchestrator
  → permissions:permissions running...
  ✓ permissions:permissions
  → skill:project-scaffold running...
  ✓ skill:project-scaffold
  → skill:roadmap-generator running...
  ✓ skill:roadmap-generator
  → config:sdd-orchestrator running...
  ✓ config:sdd-orchestrator
  → skill:skill-creator running...
  ✓ skill:skill-creator
  → skill:skill-registry running...
  ✓ skill:skill-registry
  → command:starter-add-command running...
  ✓ command:starter-add-command
  → skill:web-scraper running...
  ✓ skill:web-scraper
  → skill:discovery-research running...
  ✓ skill:discovery-research
  → self-install running...
  ✓ self-install
  → verify-hook running...
  ✓ verify-hook
Verification checks: 0 passed, 0 failed, 0 warnings, 0 skipped
Installation verified. You're ready to start building.
```

### fallo

Hora: 2026-10-23T12:25:54-03:00. Código: 1.

Comando: `./run-sandbox-fallo ./bin/active-stack install --mode full --agent claude --yes`

```text
→ snapshot running...
  ✓ snapshot
  → skill:agent-instruction running...
  ✓ skill:agent-instruction
  → external:context7 running...
  ✓ external:context7
  → external:engram running...
  ✓ external:engram
  → skill:find-skill running...
  ✓ skill:find-skill
  → skill:kb-creator running...
  ✓ skill:kb-creator
  → external:openspec running...
  ✓ external:openspec
  → skill:active-orchestrator running...
  ✓ skill:active-orchestrator
  → permissions:permissions running...
  ✓ permissions:permissions
  → skill:project-scaffold running...
  ✓ skill:project-scaffold
  → skill:roadmap-generator running...
  ✓ skill:roadmap-generator
  → config:sdd-orchestrator running...
  ✓ config:sdd-orchestrator
  → skill:skill-creator running...
  ✓ skill:skill-creator
  → skill:skill-registry running...
  ✓ skill:skill-registry
  → command:starter-add-command running...
  ✗ command:starter-add-command: agent "claude": command: write "/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/as-review/ensayo/home/.claude/commands/active/init.md": open /home/juanisarmiento/.claude/jobs/f3e370f0/tmp/as-review/ensayo/home/.claude/commands/active/init.md: read-only file system

Installation failed: agent "claude": command: write "/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/as-review/ensayo/home/.claude/commands/active/init.md": open /home/juanisarmiento/.claude/jobs/f3e370f0/tmp/as-review/ensayo/home/.claude/commands/active/init.md: read-only file system
Rollback: partial — the following were left in place:
  - brew formula "gentleman-programming/tap/engram" stays installed (undo with: brew uninstall gentleman-programming/tap/engram)
```

### restore

Hora: 2026-10-23T12:26:25-03:00. Código: 0.

Comando: `./run-sandbox ./bin/active-stack windows backups restore --id install`

```text
{"success":true,"message":"Backup restored successfully.","id":"install"}
```

### install-final

Hora: 2026-10-23T12:26:25-03:00. Código: 0.

Comando: `./run-sandbox ./bin/active-stack install --mode full --agent claude --yes`

```text
→ snapshot running...
  ✓ snapshot
  → skill:agent-instruction running...
  ✓ skill:agent-instruction
  → external:context7 running...
  ✓ external:context7
  → external:engram running...
  ✓ external:engram
  → skill:find-skill running...
  ✓ skill:find-skill
  → skill:kb-creator running...
  ✓ skill:kb-creator
  → external:openspec running...
  ✓ external:openspec
  → skill:active-orchestrator running...
  ✓ skill:active-orchestrator
  → permissions:permissions running...
  ✓ permissions:permissions
  → skill:project-scaffold running...
  ✓ skill:project-scaffold
  → skill:roadmap-generator running...
  ✓ skill:roadmap-generator
  → config:sdd-orchestrator running...
  ✓ config:sdd-orchestrator
  → skill:skill-creator running...
  ✓ skill:skill-creator
  → skill:skill-registry running...
  ✓ skill:skill-registry
  → command:starter-add-command running...
  ✓ command:starter-add-command
  → skill:web-scraper running...
  ✓ skill:web-scraper
  → skill:discovery-research running...
  ✓ skill:discovery-research
  → self-install running...
  ✓ self-install
  → verify-hook running...
  ✓ verify-hook
Verification checks: 0 passed, 0 failed, 0 warnings, 0 skipped
Installation verified. You're ready to start building.
```

### verify-final

Hora: 2026-10-23T12:26:36-03:00. Código: 0.

Comando: `./run-sandbox ./bin/ensayo-verify`

```text
Verification checks: 17 passed, 0 failed, 0 warnings, 1 skipped
[ok] active-stack:binary - binario de cmd/active-stack y roles --help con salida 0
[ok] external:binary:openspec - binary "openspec" in PATH after install
[ok] external:binary:engram - binary "engram" in PATH after install
[ok] external:mcp-config:context7:claude - MCP config exists and is valid JSON at /home/juanisarmiento/.claude/jobs/f3e370f0/tmp/as-review/ensayo/home/.claude.json for claude
[--] external:mcp-endpoint:context7 - MCP endpoint https://mcp.context7.com is reachable (check not implemented)
[ok] config:sdd-orchestrator:claude - idempotent marker present exactly once in claude instructions (sdd-orchestrator)
[ok] config:compose:claude - contenido instalado igual a Compose para variante y toggles efectivos
[ok] permissions:permissions:claude - security floor applied in claude settings
[ok] skill:active-orchestrator:claude - SKILL.md present and non-empty for claude/active-orchestrator
[ok] skill:kb-creator:claude - SKILL.md present and non-empty for claude/kb-creator
[ok] skill:project-scaffold:claude - SKILL.md present and non-empty for claude/project-scaffold
[ok] skill:roadmap-generator:claude - SKILL.md present and non-empty for claude/roadmap-generator
[ok] skill:agent-instruction:claude - SKILL.md present and non-empty for claude/agent-instruction
[ok] skill:skill-registry:claude - SKILL.md present and non-empty for claude/skill-registry
[ok] skill:web-scraper:claude - SKILL.md present and non-empty for claude/web-scraper
[ok] skill:discovery-research:claude - SKILL.md present and non-empty for claude/discovery-research
[ok] skill:find-skill:claude - SKILL.md present and non-empty for claude/find-skill
[ok] skill:skill-creator:claude - SKILL.md present and non-empty for claude/skill-creator
Installation verified. You're ready to start building.
```

### Evidencia help

```text
Usage of install:
  -agent string
    	comma-separated list of agents (e.g. claude,opencode)
  -custom string
    	comma-separated harness IDs (only with --mode custom)
  -dry-run
    	print plan steps; do not execute
  -headless
    	non-interactive install; implied by --mode or --agent
  -home string
    	override home directory (default: os.UserHomeDir())
  -lang string
    	output language: en|es (default: en)
  -mode string
    	install mode: lite|full|custom
  -no-self-install
    	skip copying the running binary into the PATH bin dir (CI / reproducible builds)
  -tier string
    	permission tier: estricto|balanceado|bypass
  -windows-events-file string
    	internal: write windows install JSON events to this file
  -y	alias for --yes
  -yes
    	confirm without prompt
error: flag: help requested
```

### Evidencia dry-run

```text
Dry-run: plan steps (not executed):
  [prepare] snapshot
  [apply]   skill:agent-instruction
  [apply]   external:context7
  [apply]   external:engram
  [apply]   skill:find-skill
  [apply]   skill:kb-creator
  [apply]   external:openspec
  [apply]   skill:active-orchestrator
  [apply]   permissions:permissions
  [apply]   skill:project-scaffold
  [apply]   skill:roadmap-generator
  [apply]   config:sdd-orchestrator
  [apply]   skill:skill-creator
  [apply]   skill:skill-registry
  [apply]   command:starter-add-command
  [apply]   skill:web-scraper
  [apply]   skill:discovery-research
  [apply]   self-install
  [apply]   verify-hook
```

### Evidencia idempotencia

```text
ad19611ebd537a0bdca62fa0616c8bc10d6e4cff931fe5014e23d0ce982c2de0  CLAUDE.instalado.md
ad19611ebd537a0bdca62fa0616c8bc10d6e4cff931fe5014e23d0ce982c2de0  home/.claude/CLAUDE.md
```

### Evidencia backups-list

```text
{"backups":[{"id":"install","created_at":"2026-10-23T15:24:33Z","source":"unknown source","description":"","file_count":1,"pinned":false,"compressed":true,"display_label":"unknown source — 2026-10-23 12:24 (1 files)","manifest_path":"/home/juanisarmiento/.claude/jobs/f3e370f0/tmp/as-review/ensayo/home/.active-stack/backups/install/manifest.json"}]}
```

### Evidencia restore-sha256

```text
25d2ecaf5a3f369e44c6e4aecb6d1f5fd5f646c5d97c378910d9d91955e6dfb7  CLAUDE.original.md
25d2ecaf5a3f369e44c6e4aecb6d1f5fd5f646c5d97c378910d9d91955e6dfb7  home/.claude/CLAUDE.md
```

### Evidencia version

```text
active-stack (devel) (SHA desconocido)
```

### Evidencia roles-help

```text
Usage: active-stack roles [flags]

Lee un diff y devuelve qué revisores corresponden. No convoca a nadie: la
convocatoria sigue siendo del orquestador.

Sin flags se usan los archivos cambiados y no commiteados.

Flags:
  --base <ref>         comparar contra esa referencia (suma los sin trackear)
  --files a.go b.go    lista explícita, sin tocar git
  --json               salida para consumir sin parsear texto

Los roles que NO corresponden también salen, con el motivo: una ausencia sin
motivo es indistinguible de un olvido.
```

### Evidencia real-sin-cambios

```text
25d2ecaf5a3f369e44c6e4aecb6d1f5fd5f646c5d97c378910d9d91955e6dfb7  /home/juanisarmiento/.claude/CLAUDE.md
25d2ecaf5a3f369e44c6e4aecb6d1f5fd5f646c5d97c378910d9d91955e6dfb7  CLAUDE.original.md
```

## Límites de la conclusión

El ensayo demuestra comportamiento con CLAUDE.md y agents copiados, no con settings/MCP/skills reales, que deliberadamente no se leyeron. Dependencias remotas se obtuvieron durante la ejecución; sus versiones no están fijadas por este ensayo. El rollback comprobado incluye todos los archivos del HOME falso excepto sus backups, pero no garantiza todas las fallas posibles. No se ejecutó la suite general del repositorio: se verificó la instalación real aislada solicitada. No se tocó main ni se hicieron commits.

> Copia del informe original del ensayo, que vivía en una carpeta temporal de la sesión. Los enlaces relativos a `logs/`, `run-sandbox` y los diffs apuntan a archivos que no se versionaron: el diff del `CLAUDE.md` contiene datos personales de Juani y no va al repo.
