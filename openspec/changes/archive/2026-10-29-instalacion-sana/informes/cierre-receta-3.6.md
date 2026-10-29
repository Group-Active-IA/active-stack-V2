# Cierre de instalacion-sana contra la receta §3.6 (CONF-1)

Fecha: 2026-10-29. Autor: orquestador (sesión principal Claude), por pedido de Juani.
Rama `stack-actualizado`; código verificado en `7910d54`. Receta:
`docs/propuestas/2026-10-23-metodologia-proptech-en-opsx.md` §3.6.

## Tarea del change: el `.test` en `~/.local/bin`

EVIDENCIA: `active-stack version` imprimía la ayuda de un test y "PASS"; el binario era
`headless.test` (strings: `cmd/active-stack/headless.test`). Se aisló corriendo cada test con un
binario precompilado y HOME falso: `TestRunWindowsStartersInstall_Success` instala un starter sin
HomeDir y el self-install copiaba el binario de tests al home real. Arreglo `9683781`: `TestMain`
del paquete `headless` aísla HOME/USERPROFILE/LOCALAPPDATA/APPDATA. Suite completa con HOME falso:
ningún paquete escribe el binario.

## Punto 1. Registro desde el SHA aprobado

EVIDENCIA (17:36 -0300, máquina de Juani): ejecutable efectivo `~/.local/bin/active-stack`,
`active-stack version` → SHA `d634565`; sin rastro de `headless.test`; `/usr/local/bin/active-stack`
no existe. Destinos: Claude `~/.claude/CLAUDE.md` (variante claude, marcadores 1 par), Codex
`~/.codex/AGENTS.md` (variante sin `subagent_type`, sin directorio de comandos inventado).
Toggles efectivos: los del catálogo (lean-workflow y los 12 históricos).

## Punto 2. Snapshot antes de escribir

EVIDENCIA: el instalador hace snapshot vía `internal/backup` con manifiesto (destino, existencia,
hash, permisos). Se agregó un tar completo previo de `~/.claude`, `~/.claude.json`, `~/.codex` y
`~/.agents` (`~/backups-claude/`, 29/10).

DESVÍO: el snapshot de instalación iba siempre a `backups/install` y la reinstalación lo pisaba
(el ensayo del 23/10 lo esquivó copiándolo a mano). En la máquina de Juani el snapshot del stack
quedó en el de la última reinstalación; el estado original está en el tar. Corregido en `7910d54`:
un `backups/install-<fecha>` por instalación.

## Punto 3. Ensayo aislado: instalar, reinstalar, inducir falla, restaurar

EVIDENCIA (HOME aislados, binario de la rama):

- Claude: home vacío, reinstalación idempotente (mismos hashes, marcador único), home con roles,
  comandos y skill personales: restore deja el home idéntico salvo el binario engram descargado.
- OpenCode (17:44): config previa con `opencode.json` 0600 y `AGENTS.md` propio; instalación con
  permiso conservado (600), reinstalación idempotente, restore desde el snapshot de la primera
  instalación → **idéntico al original en contenido y permisos**.
- Claude (17:45): `settings.json` y `.claude.json` 0600 → conservados; restore idéntico.
- Falla inducida: `--toggles` con verificación desalineada produjo rollback completo (hallado y
  corregido en `6fbb065`); conflictos de roles reportados como BLOQUEANTE sin abortar.

Bugs hallados en este punto y corregidos con test: backups de comandos pisados en el mismo
segundo, snapshot sin contenido de carpetas existentes (restore incompleto), mensaje "verified"
con 0 checks (`fab5966`); snapshot único pisado y permisos 0600→0644 (`7910d54`).

## Punto 4. Reinstalación en la máquina de Juani

EVIDENCIA: instalación real 29/10 13:50 con OK escrito de Juani (sesión), actualizaciones 16:01
(Claude y Codex). Aceptación: binario de `cmd/active-stack`, `active-stack roles --help` → exit 0
(17:36), contenido igual a `Compose` (check `config:compose` en verde), texto del usuario fuera
de markers preservado (sección `# Juani`, mismo hash antes y después).

DESVÍO: la instalación dejó `~/.codex/config.toml` en 0644 (era 0600 según el tar); se devolvió a
0600 a mano (sin secretos en el archivo) y el bug quedó corregido en `7910d54`.

## Aceptación por destino

- Claude: bloque y comandos `verify`/`ship` instalados. ✅
- Codex: instrucciones por su adapter, sin directorio de comandos. ✅
- OpenCode: agente primario `sdd-orchestrator` con el bloque lean como prompt, referencia en
  `opencode.json`, instrucciones y settings recuperados juntos por restore. ✅ (sólo en HOME
  aislado: Juani no usa OpenCode)

## Desvíos declarados

- La receta pedía no modificar backup, filemerge ni rollback sin propuesta y review ALTO propios.
  Se modificaron (`fab5966`, `7910d54`) con OK escrito de Juani en la sesión, sin review ALTO
  formal.
- Sin rompedor ni auditor: Juani pidió trabajar directo en este repo. No hay veredicto APROBADO.
- Restore no se ejecutó en la máquina real (sólo en HOME aislados).

## Veredicto

JUICIO (orquestador, confianza alta en lo ejecutado): los cuatro puntos de §3.6 tienen evidencia
y los defectos que encontró la receta están corregidos con test. CONF-1 queda pagado con los
desvíos de arriba declarados. Se archiva por decisión de Juani.
