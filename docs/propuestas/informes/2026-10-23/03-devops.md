**JUICIO DevOps, confianza alta:** recomiendo resolver OPS-1 a OPS-4 antes de firmar. “BLOQUEANTE” indica esa recomendación, no un veto.

**EVIDENCIA, 23/10/2026 09:39 ART:** `git status --short --branch` → `feat/metodologia-proptech`; sólo aparece `?? docs/propuestas/informes/`. No leí esos informes ni modifiqué archivos.

**OPS-1 · BLOQUEANTE · Reinstalación sin criterio observable de recuperación.**
**EVIDENCIA:** `go version -m ~/.local/bin/active-stack` → `cmd/active-stack/headless.test`; `stat -c '%y' ~/.claude/CLAUDE.md` → `2026-09-09 16:09:51`.
**EVIDENCIA:** `nl -ba docs/propuestas/2026-10-23-metodologia-proptech-en-opsx.md` → línea 122 pide reinstalar “con backup”, sin procedimiento ni ensayo de restore.
**JUICIO DevOps, alta:** antes de firmar, exigir receta desde clon limpio, snapshot identificado, reinstalación repetida sin duplicaciones y fallo inducido con restauración comparada por contenido y permisos. El rollback debe incluir binario, instrucciones y comandos; registrar errores y ubicación del backup.
**SUPUESTO:** no ejecuté instalación, tests ni rollback por la restricción de sólo lectura; mi criterio de terminado queda sin demostrar.

**OPS-2 · BLOQUEANTE · “CI completo” todavía puede confundirse con un run verde.**
**EVIDENCIA:** `nl -ba docs/propuestas/2026-10-23-metodologia-proptech-en-opsx.md` → línea 89 propone `gh run list --commit <sha>` y `gh run watch`; línea 103 condiciona los tests del comando a cambiar su ruta.
**EVIDENCIA:** `cat .github/workflows/ci.yml` → eventos push y pull_request, matriz Ubuntu/Windows; `rg -n 'cycle ends' assets/commands/claude/opsx/ship.md` → línea 381 termina en PR abierto.
**JUICIO DevOps, alta:** definir workflows/eventos exigidos, espera por aparición del run, identidad del SHA y tratamiento de cancelado, omitido y timeout. Un run exitoso no demuestra el conjunto completo. Probar ambas copias aunque sus rutas sigan iguales, incluyendo ausencia de CI y consulta fallida.

**OPS-3 · BLOQUEANTE · El hook necesita contrato de persistencia y falla.**
**EVIDENCIA:** `nl -ba docs/propuestas/2026-10-23-metodologia-proptech-en-opsx.md` → línea 132 guarda el último mensaje en `.git/informes-<change>/`; `nl -ba docs/exploraciones/2026-10-23-metodologia-proptech.md` → líneas 173-175 declaran supuestos sobre payload y transcript.
**JUICIO DevOps, alta:** exigir spike verificable antes de comprometer implementación: identificación por sesión/agente/vuelta, escritura atómica sin sobrescribir informes, resolución de gitdir para worktrees y señal visible si falta el original. Instalar/desinstalar el hook debe preservar hooks existentes con backup y rollback probado.
**SUPUESTO:** no validé el contrato externo actual de `SubagentStop`; no lo doy por estable ni disponible.

**OPS-4 · BLOQUEANTE · La compatibilidad requiere aceptación por agente y destino.**
**EVIDENCIA:** `rg --files assets/commands` → ship/verify propios sólo bajo `claude`; `rg -n 'CommandsDir|ConfigDelivery' internal/agents/{codex,opencode}/adapter.go` → Codex sin directorio de comandos; OpenCode entrega como agente primario.
**EVIDENCIA:** `nl -ba internal/verify/harness_checks.go` → líneas 99-102 verifican OpenCode en settings; `nl -ba internal/harness/config/primary_agent.go` → líneas 50 y 71 respaldan instrucciones y settings mediante el mismo snapshotDir.
**JUICIO DevOps, alta:** explicitar fallback sin hook para Codex/OpenCode y comparar `Compose` contra el destino real, con variante y toggles efectivos. Incluir restauración conjunta de instrucciones/settings de OpenCode y comprobar que ambos respaldos sobrevivan. Los goldens Claude/OpenCode no prueban esa entrega.

**OPS-5 · NOTA · El aislamiento de self-install no termina en `--home`.**
**EVIDENCIA:** `nl -ba internal/system/binary_install_dir.go` → líneas 61-62 priorizan `/usr/local/bin` si resulta escribible, aun con home explícito.
**JUICIO DevOps, alta:** el diagnóstico del paso 0 debe revisar también ese destino; los ensayos necesitan resolver binario y PATH hacia un entorno aislado. No atribuir todavía el incidente a un test concreto.

**JUICIO DevOps sobre las cinco recomendaciones, confianza media:**

1. **Un rompedor, tres lentes:** aceptable; declarar pérdida de independencia y exigir evidencia separada por lente. Vueltas posteriores deben cubrir también consumidores afectados por la corrección.
2. **Senior antes del auditor:** de acuerdo; conserva los informes originales y permite que el auditor contraste memoria con reproducciones actuales.
3. **Firma en fase 2, configurada por Juani:** de acuerdo con postergar; firma válida no prueba aprobación del alcance. Falta definir identidad confiable y vínculo entre commit firmado, propuesta y revisión aprobada.
4. **Modelo declarado:** alcanza como declaración trazable, tipada SUPUESTO cuando no haya registro verificable. No venderla como independencia controlada.
5. **Archive a sonnet:** razonable como preferencia con sustitución declarada; el cierre sigue dependiendo del diff documental y su evidencia, no del nombre del modelo.