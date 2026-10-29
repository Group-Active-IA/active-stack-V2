# Estado del proyecto

Único lugar donde se escribe el estado de Active Stack. README, AGENTS/CLAUDE y ARCHITECTURE
enlazan acá en vez de copiarlo (un test lo verifica). El estado de cada change lo da el CLI:
`openspec list` y `openspec status`. Actualizado: 29/10/2026.

## Resumen

🟢 **Activo.** El núcleo del instalador está implementado: catálogo, modelo, adapters
(Claude, Codex y OpenCode), backup y restore, merge por markers, pipeline, verify y TUI. El
roadmap interno C-01…C-32 está cerrado; el historial está en [CHANGES.md](../CHANGES.md).

| Tema | Estado |
|---|---|
| Perfil predeterminado | `lean-workflow`: cinco roles dentro de las fases OPSX, bloque de ~17 KB contra ~115 KB del legacy. Ahorro de tokens real todavía sin medir en un piloto |
| Rama de trabajo | `stack-actualizado`, sin merge a `main` ni release publicado |
| Plataformas | Windows, Linux y WSL (amd64) |
| Fase 1 metodológica | Completa: nueve changes archivados con OK escrito (informe 22 y cierre del 29/10) |
| `instalacion-sana` | **Archivado el 29/10/2026**: CONF-1 pagado contra la receta §3.6 (instalación real, ensayo aislado en Claude y OpenCode, restore idéntico en contenido y permisos). Evidencia y desvíos en `openspec/changes/archive/2026-10-29-instalacion-sana/informes/cierre-receta-3.6.md` |
| Fase 2 | Paso 2a (`active-stack reports capture`) implementado; 2b (captura por hooks) pendiente. Firmas y el resto: (TBD), sin aprobación implícita |
| Epic post-starters | `uninstall-subcommand` y `tui-menu-hub` archivados; `tui-update-stack`, `tui-configure-models` y `claude-agent-switch-research` pendientes |

## Cómo se verifica

`active-stack version` informa versión y SHA. Los checks posteriores a la instalación validan
que el binario sea `cmd/active-stack`, que `active-stack roles --help` salga 0 y que el bloque
instalado sea igual a `Compose` para la variante y los toggles efectivos.

EVIDENCIA: [base metodológica](metodologia-opsx.md); commits `84fcf0e`, `c325955`; propuesta
firmada `4748680`, decisión `5ab0f4f`; informe 22.
