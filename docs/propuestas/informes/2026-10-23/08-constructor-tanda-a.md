Tanda A completada. `go test ./...` pasó después de cada commit, el 23/10/2026, hora Argentina:

| Paso | Commit | Rojo observado | Suite final |
|---|---|---|---|
| 0, instalacion-sana | `84fcf0e` | `version` desconocido, Compose alterado aceptado y APIs ausentes | OK 09:57:03 |
| 1, opencode-mcp-command-array | `8c6a06b` | `command` incorrecto y `environment` ausente | OK 09:58:10 |
| 2, ship-espera-ci | `60b5e83` | Contrato CI ausente en ambas copias; 14 mutantes detectados | OK 10:00:14 |
| 3, apply-constructor | `f3fc02a` | Contrato del constructor ausente | OK 10:03:44 |

El `.test` puede copiarse mediante `os.Executable()` y auto-install en tests como `TestVerifyHookNilIsNoOp`. Hoy usan homes temporales, pero el resolver puede elegir `/usr/local/bin`. No quedó probado qué test originó el archivo real.

[Informe con salidas rojas y verdes](/home/juanisarmiento/Proyectos/active-stack-V2/docs/propuestas/informes/2026-10-23/08-constructor-tanda-a.md). El cierre del paso 3 quedó en su `tasks.md`.

Sin reinstalación, instalaciones, `go build`, push ni trailers de coautoría.