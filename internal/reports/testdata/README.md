# Contratos de captura (BE-4)

Las fixtures reproducen estructuras públicas con texto e identificadores inventados.
No se leyó ni copió contenido de sesiones personales. Son ensayos de formato, no
respuestas obtenidas invocando a un modelo. El CLI importa un archivo explícito
por despacho, fuera del repo destino; no descubre directorios del proveedor.

| Proveedor | Fuente y acceso | Formato admitido | Fixture |
|---|---|---|---|
| Codex | stdout de `codex exec --json`, archivo por invocación | CLI 0.132.0; JSONL `thread.started → turn.started → item.completed(agent_message) → turn.completed` | `codex-0.132.0.jsonl` |
| Claude | payloads de hooks PostToolUse/SubagentStop, exportados como JSONL por despacho | Claude Code 2.1.271, `SubagentHandback.tool_input.message` seguido del cierre de la misma hija/sesión | `claude-2.1.271-hooks.jsonl` |
| OpenCode | candidato: API HTTP `GET /session/:id/message/:messageID`, `/global/health` para versión | **Ninguno habilitado**; versión y asociación padre/hija/mensaje (TBD). Siempre UNSUPPORTED | `opencode-health-unconfirmed.json`: sólo forma pública de health, NO un esquema de informe inventado |

Fuentes consultadas el 2026-10-23:

- [Codex no interactivo](https://learn.chatgpt.com/docs/non-interactive-mode).
- [Esquema Rust fijado a rust-v0.132.0](https://github.com/openai/codex/blob/rust-v0.132.0/codex-rs/exec/src/exec_events.rs).
- [Claude SubagentStop/Handback](https://code.claude.com/docs/en/hooks#subagentstop): Handback desde 2.1.271. La documentación es móvil; la whitelist admite sólo la versión declarada, no versiones futuras por semejanza.
- [OpenCode Server](https://opencode.ai/docs/server/): el esquema detallado se obtiene de `/doc` de una instancia concreta. No se obtuvo esa instancia versionada ni se ensayó asociación hija/mensaje; no se habilita API por inferencia.

Límites y procedencia:

- `--source-session` y, en Claude, `--source-agent` son la identidad esperada
  registrada por el llamador antes del despacho, no extraída de la prosa.
  El CLI no acredita autenticidad del llamador ni la versión del ejecutable que
  produjo el archivo. Modelo/proveedor ausentes permanecen vacíos. Codex no
  proporciona turn ID en este formato; no se fabrica uno.
- Codex exige una sesión, un turno, exactamente un mensaje completo y cierre;
  múltiples mensajes/turnos, error, nulos o desorden quedan PENDING. Otros items
  conocidos se omiten del bruto persistido para no volcar comandos/transcript.
  `raw.txt` contiene los eventos necesarios, `informe.md` el texto íntegro,
  con hashes y tamaños separados. Formato de manifiesto 1: campos aditivos
  `raw_hash` y `raw_size_bytes`; lectura de manifiestos previos compatible.
- Claude exige Handback exitoso (PostToolUse) y SubagentStop de la misma hija;
  Stop aislado, PreToolUse y cierre ausente/duplicado son PENDING, nunca se usa
  la despedida. No se siguen rutas de transcript ni se abre el HOME del proveedor.
- Importación explícita, sin integración instalada: **no se implementó aún el
  ciclo automático de instalación/update/uninstall de hooks de 2b**. No se
  anuncia paridad con captura automática. Ownership y backup deberán ensayarse
  antes de habilitar esa integración; el presente comando no escribe settings.
- Límite de entrada CLI/bruto/texto: 8 MiB, rechazo visible sin truncar. No hay
  conexiones API ni espera de red; timeout/retención del núcleo siguen como deuda
  heredada. RETAINED (RETENIDO) prevalece incluso sobre vacío, evento parcial o conflicto;
  el bruto de un evento parcial se conserva para escanearlo antes de persistirlo;
  patrones propios heredados no equivalen a un detector universal de secretos.
