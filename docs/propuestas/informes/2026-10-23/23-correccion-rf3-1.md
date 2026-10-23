# Corrección puntual RF3-1

Decisión: informe 22, opción B del informe 21. Alcance permanente: una expresión en `internal/contracttest/text.go`, este informe y la respuesta separada solicitada. Sin invocar go build ni push. Los pasos se ejecutaron con `python3 /tmp/rf3_1.py`; las inserciones y restauraciones muestran su operación Python reproducible, y los comandos Go son los argumentos exactos de subprocess. Las horas corresponden al registro al terminar cada paso.

## Estado inicial

Hora: 2026-10-23T10:53:57.227575-03:00

Comando:
```sh
git status --short
```

Salida:
```text
?? docs/propuestas/informes/2026-10-23/.01-analista.log
?? docs/propuestas/informes/2026-10-23/.02-arquitecto.log
?? docs/propuestas/informes/2026-10-23/.03-devops.log
?? docs/propuestas/informes/2026-10-23/.04-qa.log
?? docs/propuestas/informes/2026-10-23/.05-auditor.log
?? docs/propuestas/informes/2026-10-23/.06-correccion.log
?? docs/propuestas/informes/2026-10-23/.07-auditor.log
?? docs/propuestas/informes/2026-10-23/.08-a.log
?? docs/propuestas/informes/2026-10-23/.09-b.log
?? docs/propuestas/informes/2026-10-23/.12.log
?? docs/propuestas/informes/2026-10-23/.13.log
?? docs/propuestas/informes/2026-10-23/.18.log
?? docs/propuestas/informes/2026-10-23/.19.log
?? docs/propuestas/informes/2026-10-23/23-correccion-rf3-1.md
```

Exit: 0

## Rama y SHA iniciales

Hora: 2026-10-23T10:53:57.232365-03:00

Comando:
```sh
git log -1 '--format=%H %D'
```

Salida:
```text
e46294a0ee08a41a11946c526f61266d4773ad4e HEAD -> feat/metodologia-proptech
```

Exit: 0

## VIEJO: insertar mutante

Hora: 2026-10-23T10:53:57.232850-03:00

Comando:
```sh
python3 - <<'PY'
from pathlib import Path
p = Path('internal/harness/config/assets/claude/sdd-orchestrator.md')
p.write_bytes(p.read_bytes() + b'\nConformity passed. Ready to archive\n')
PY
```

Salida:
```text
internal/harness/config/assets/claude/sdd-orchestrator.md: Conformity passed. Ready to archive
```

Exit: 0

## VIEJO: TestVerifyVocabularioCompuesto — Conformity passed. Ready to archive

Hora: 2026-10-23T10:53:57.953937-03:00

Comando:
```sh
go test ./internal/harness/config -run '^TestVerifyVocabularioCompuesto$' -count=1 -v
```

Salida:
```text
=== RUN   TestVerifyVocabularioCompuesto
=== RUN   TestVerifyVocabularioCompuesto/claude
=== RUN   TestVerifyVocabularioCompuesto/codex
=== RUN   TestVerifyVocabularioCompuesto/opencode
--- PASS: TestVerifyVocabularioCompuesto (0.17s)
    --- PASS: TestVerifyVocabularioCompuesto/claude (0.06s)
    --- PASS: TestVerifyVocabularioCompuesto/codex (0.05s)
    --- PASS: TestVerifyVocabularioCompuesto/opencode (0.05s)
PASS
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.171s
```

Exit: 0

## Revertir mutante

Hora: 2026-10-23T10:53:57.954264-03:00

Comando:
```sh
python3 - <<'PY'
from pathlib import Path
p = Path('internal/harness/config/assets/claude/sdd-orchestrator.md')
suffix = b'\nConformity passed. Ready to archive\n'
data = p.read_bytes()
assert data.endswith(suffix)
p.write_bytes(data[:-len(suffix)])
PY
```

Salida:
```text
Restaurado byte a byte al contenido previo.
```

Exit: 0

## VIEJO: insertar mutante

Hora: 2026-10-23T10:53:57.954462-03:00

Comando:
```sh
python3 - <<'PY'
from pathlib import Path
p = Path('assets/commands/claude/opsx/verify.md')
p.write_bytes(p.read_bytes() + b'\nConformity passed. Ready to archive\n')
PY
```

Salida:
```text
assets/commands/claude/opsx/verify.md: Conformity passed. Ready to archive
```

Exit: 0

## VIEJO: TestVerifyConformidadNoHabilitaArchive — Conformity passed. Ready to archive

Hora: 2026-10-23T10:53:58.185861-03:00

Comando:
```sh
go test ./assets -run '^TestVerifyConformidadNoHabilitaArchive$' -count=1 -v
```

Salida:
```text
=== RUN   TestVerifyConformidadNoHabilitaArchive
--- PASS: TestVerifyConformidadNoHabilitaArchive (0.00s)
PASS
ok  	github.com/Group-Active-IA/active-stack/assets	0.005s
```

Exit: 0

## Revertir mutante

Hora: 2026-10-23T10:53:58.186111-03:00

Comando:
```sh
python3 - <<'PY'
from pathlib import Path
p = Path('assets/commands/claude/opsx/verify.md')
suffix = b'\nConformity passed. Ready to archive\n'
data = p.read_bytes()
assert data.endswith(suffix)
p.write_bytes(data[:-len(suffix)])
PY
```

Salida:
```text
Restaurado byte a byte al contenido previo.
```

Exit: 0

## VIEJO: insertar mutante

Hora: 2026-10-23T10:53:58.186295-03:00

Comando:
```sh
python3 - <<'PY'
from pathlib import Path
p = Path('internal/harness/config/assets/claude/sdd-orchestrator.md')
p.write_bytes(p.read_bytes() + b'\nReady for archive after conformity.\n')
PY
```

Salida:
```text
internal/harness/config/assets/claude/sdd-orchestrator.md: Ready for archive after conformity.
```

Exit: 0

## VIEJO: TestVerifyVocabularioCompuesto — Ready for archive after conformity.

Hora: 2026-10-23T10:53:58.811152-03:00

Comando:
```sh
go test ./internal/harness/config -run '^TestVerifyVocabularioCompuesto$' -count=1 -v
```

Salida:
```text
=== RUN   TestVerifyVocabularioCompuesto
=== RUN   TestVerifyVocabularioCompuesto/claude
=== RUN   TestVerifyVocabularioCompuesto/codex
=== RUN   TestVerifyVocabularioCompuesto/opencode
--- PASS: TestVerifyVocabularioCompuesto (0.13s)
    --- PASS: TestVerifyVocabularioCompuesto/claude (0.05s)
    --- PASS: TestVerifyVocabularioCompuesto/codex (0.04s)
    --- PASS: TestVerifyVocabularioCompuesto/opencode (0.04s)
PASS
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.133s
```

Exit: 0

## Revertir mutante

Hora: 2026-10-23T10:53:58.811399-03:00

Comando:
```sh
python3 - <<'PY'
from pathlib import Path
p = Path('internal/harness/config/assets/claude/sdd-orchestrator.md')
suffix = b'\nReady for archive after conformity.\n'
data = p.read_bytes()
assert data.endswith(suffix)
p.write_bytes(data[:-len(suffix)])
PY
```

Salida:
```text
Restaurado byte a byte al contenido previo.
```

Exit: 0

## VIEJO: insertar mutante

Hora: 2026-10-23T10:53:58.811548-03:00

Comando:
```sh
python3 - <<'PY'
from pathlib import Path
p = Path('assets/commands/claude/opsx/verify.md')
p.write_bytes(p.read_bytes() + b'\nReady for archive after conformity.\n')
PY
```

Salida:
```text
assets/commands/claude/opsx/verify.md: Ready for archive after conformity.
```

Exit: 0

## VIEJO: TestVerifyConformidadNoHabilitaArchive — Ready for archive after conformity.

Hora: 2026-10-23T10:53:59.025924-03:00

Comando:
```sh
go test ./assets -run '^TestVerifyConformidadNoHabilitaArchive$' -count=1 -v
```

Salida:
```text
=== RUN   TestVerifyConformidadNoHabilitaArchive
    commands_test.go:126: conformidad anuncia archive antes del auditor
--- FAIL: TestVerifyConformidadNoHabilitaArchive (0.00s)
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.006s
FAIL
```

Exit: 1

## Revertir mutante

Hora: 2026-10-23T10:53:59.026149-03:00

Comando:
```sh
python3 - <<'PY'
from pathlib import Path
p = Path('assets/commands/claude/opsx/verify.md')
suffix = b'\nReady for archive after conformity.\n'
data = p.read_bytes()
assert data.endswith(suffix)
p.write_bytes(data[:-len(suffix)])
PY
```

Salida:
```text
Restaurado byte a byte al contenido previo.
```

Exit: 0

## Corrección de la expresión

Hora: 2026-10-23T10:53:59.026295-03:00

Comando:
```sh
python3 - <<'PY'
from pathlib import Path
p = Path('internal/contracttest/text.go')
p.write_text(p.read_text().replace('ready (to|for) archive[.!]', 'ready (to|for) archive([^?]|$)'))
PY
```

Salida:
```text
Expresión reemplazada; se conserva |proceed to archive.
```

Exit: 0

## CORREGIDO: insertar mutante

Hora: 2026-10-23T10:53:59.026463-03:00

Comando:
```sh
python3 - <<'PY'
from pathlib import Path
p = Path('internal/harness/config/assets/claude/sdd-orchestrator.md')
p.write_bytes(p.read_bytes() + b'\nConformity passed. Ready to archive\n')
PY
```

Salida:
```text
internal/harness/config/assets/claude/sdd-orchestrator.md: Conformity passed. Ready to archive
```

Exit: 0

## CORREGIDO: TestVerifyVocabularioCompuesto — Conformity passed. Ready to archive

Hora: 2026-10-23T10:53:59.566586-03:00

Comando:
```sh
go test ./internal/harness/config -run '^TestVerifyVocabularioCompuesto$' -count=1 -v
```

Salida:
```text
=== RUN   TestVerifyVocabularioCompuesto
=== RUN   TestVerifyVocabularioCompuesto/claude
    verify_loop_test.go:105: contrato contradicho (ready (to|for) archive([^?]|$)|proceed to archive): ready to archive 
=== RUN   TestVerifyVocabularioCompuesto/codex
=== RUN   TestVerifyVocabularioCompuesto/opencode
--- FAIL: TestVerifyVocabularioCompuesto (0.13s)
    --- FAIL: TestVerifyVocabularioCompuesto/claude (0.05s)
    --- PASS: TestVerifyVocabularioCompuesto/codex (0.04s)
    --- PASS: TestVerifyVocabularioCompuesto/opencode (0.04s)
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.139s
FAIL
```

Exit: 1

## Revertir mutante

Hora: 2026-10-23T10:53:59.566817-03:00

Comando:
```sh
python3 - <<'PY'
from pathlib import Path
p = Path('internal/harness/config/assets/claude/sdd-orchestrator.md')
suffix = b'\nConformity passed. Ready to archive\n'
data = p.read_bytes()
assert data.endswith(suffix)
p.write_bytes(data[:-len(suffix)])
PY
```

Salida:
```text
Restaurado byte a byte al contenido previo.
```

Exit: 0

## CORREGIDO: insertar mutante

Hora: 2026-10-23T10:53:59.566935-03:00

Comando:
```sh
python3 - <<'PY'
from pathlib import Path
p = Path('assets/commands/claude/opsx/verify.md')
p.write_bytes(p.read_bytes() + b'\nConformity passed. Ready to archive\n')
PY
```

Salida:
```text
assets/commands/claude/opsx/verify.md: Conformity passed. Ready to archive
```

Exit: 0

## CORREGIDO: TestVerifyConformidadNoHabilitaArchive — Conformity passed. Ready to archive

Hora: 2026-10-23T10:53:59.761091-03:00

Comando:
```sh
go test ./assets -run '^TestVerifyConformidadNoHabilitaArchive$' -count=1 -v
```

Salida:
```text
=== RUN   TestVerifyConformidadNoHabilitaArchive
    commands_test.go:116: contrato contradicho (ready (to|for) archive([^?]|$)|proceed to archive): ready to archive
--- FAIL: TestVerifyConformidadNoHabilitaArchive (0.00s)
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.006s
FAIL
```

Exit: 1

## Revertir mutante

Hora: 2026-10-23T10:53:59.761385-03:00

Comando:
```sh
python3 - <<'PY'
from pathlib import Path
p = Path('assets/commands/claude/opsx/verify.md')
suffix = b'\nConformity passed. Ready to archive\n'
data = p.read_bytes()
assert data.endswith(suffix)
p.write_bytes(data[:-len(suffix)])
PY
```

Salida:
```text
Restaurado byte a byte al contenido previo.
```

Exit: 0

## CORREGIDO: insertar mutante

Hora: 2026-10-23T10:53:59.761611-03:00

Comando:
```sh
python3 - <<'PY'
from pathlib import Path
p = Path('internal/harness/config/assets/claude/sdd-orchestrator.md')
p.write_bytes(p.read_bytes() + b'\nReady for archive after conformity.\n')
PY
```

Salida:
```text
internal/harness/config/assets/claude/sdd-orchestrator.md: Ready for archive after conformity.
```

Exit: 0

## CORREGIDO: TestVerifyVocabularioCompuesto — Ready for archive after conformity.

Hora: 2026-10-23T10:54:00.312314-03:00

Comando:
```sh
go test ./internal/harness/config -run '^TestVerifyVocabularioCompuesto$' -count=1 -v
```

Salida:
```text
=== RUN   TestVerifyVocabularioCompuesto
=== RUN   TestVerifyVocabularioCompuesto/claude
    verify_loop_test.go:105: contrato contradicho (ready (to|for) archive([^?]|$)|proceed to archive): ready for archive 
=== RUN   TestVerifyVocabularioCompuesto/codex
=== RUN   TestVerifyVocabularioCompuesto/opencode
--- FAIL: TestVerifyVocabularioCompuesto (0.14s)
    --- FAIL: TestVerifyVocabularioCompuesto/claude (0.04s)
    --- PASS: TestVerifyVocabularioCompuesto/codex (0.04s)
    --- PASS: TestVerifyVocabularioCompuesto/opencode (0.05s)
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.139s
FAIL
```

Exit: 1

## Revertir mutante

Hora: 2026-10-23T10:54:00.312598-03:00

Comando:
```sh
python3 - <<'PY'
from pathlib import Path
p = Path('internal/harness/config/assets/claude/sdd-orchestrator.md')
suffix = b'\nReady for archive after conformity.\n'
data = p.read_bytes()
assert data.endswith(suffix)
p.write_bytes(data[:-len(suffix)])
PY
```

Salida:
```text
Restaurado byte a byte al contenido previo.
```

Exit: 0

## CORREGIDO: insertar mutante

Hora: 2026-10-23T10:54:00.312735-03:00

Comando:
```sh
python3 - <<'PY'
from pathlib import Path
p = Path('assets/commands/claude/opsx/verify.md')
p.write_bytes(p.read_bytes() + b'\nReady for archive after conformity.\n')
PY
```

Salida:
```text
assets/commands/claude/opsx/verify.md: Ready for archive after conformity.
```

Exit: 0

## CORREGIDO: TestVerifyConformidadNoHabilitaArchive — Ready for archive after conformity.

Hora: 2026-10-23T10:54:00.518550-03:00

Comando:
```sh
go test ./assets -run '^TestVerifyConformidadNoHabilitaArchive$' -count=1 -v
```

Salida:
```text
=== RUN   TestVerifyConformidadNoHabilitaArchive
    commands_test.go:116: contrato contradicho (ready (to|for) archive([^?]|$)|proceed to archive): ready for archive 
    commands_test.go:126: conformidad anuncia archive antes del auditor
--- FAIL: TestVerifyConformidadNoHabilitaArchive (0.00s)
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.006s
FAIL
```

Exit: 1

## Revertir mutante

Hora: 2026-10-23T10:54:00.518803-03:00

Comando:
```sh
python3 - <<'PY'
from pathlib import Path
p = Path('assets/commands/claude/opsx/verify.md')
suffix = b'\nReady for archive after conformity.\n'
data = p.read_bytes()
assert data.endswith(suffix)
p.write_bytes(data[:-len(suffix)])
PY
```

Salida:
```text
Restaurado byte a byte al contenido previo.
```

Exit: 0

## Tabla de ruteo legítima sin falso positivo

Hora: 2026-10-23T10:54:01.083898-03:00

Comando:
```sh
go test ./internal/harness/config ./assets -run '^(TestVerifyVocabularioCompuesto|TestVerifyConformidadNoHabilitaArchive)$' -count=1 -v
```

Salida:
```text
=== RUN   TestVerifyVocabularioCompuesto
=== RUN   TestVerifyVocabularioCompuesto/claude
=== RUN   TestVerifyVocabularioCompuesto/codex
=== RUN   TestVerifyVocabularioCompuesto/opencode
--- PASS: TestVerifyVocabularioCompuesto (0.14s)
    --- PASS: TestVerifyVocabularioCompuesto/claude (0.05s)
    --- PASS: TestVerifyVocabularioCompuesto/codex (0.04s)
    --- PASS: TestVerifyVocabularioCompuesto/opencode (0.04s)
PASS
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.141s
=== RUN   TestVerifyConformidadNoHabilitaArchive
--- PASS: TestVerifyConformidadNoHabilitaArchive (0.00s)
PASS
ok  	github.com/Group-Active-IA/active-stack/assets	0.006s
```

Exit: 0

## Pregunta presente en los tres bloques evaluados

Hora: 2026-10-23T10:54:01.088443-03:00

Comando:
```sh
rg -n 'ready to archive\?' internal/harness/config/assets/claude/sdd-orchestrator.md internal/harness/config/assets/codex/sdd-orchestrator.md internal/harness/config/assets/opencode/sdd-orchestrator.md
```

Salida:
```text
internal/harness/config/assets/opencode/sdd-orchestrator.md:109:| "verify", "check it matches", "ready to archive?" | `openspec-verify-change` |
internal/harness/config/assets/codex/sdd-orchestrator.md:99:| "verify", "check it matches", "ready to archive?" | `openspec-verify-change` |
internal/harness/config/assets/claude/sdd-orchestrator.md:109:| "verify", "check it matches", "ready to archive?" | `openspec-verify-change` |
```

Exit: 0

## Suite final con HOME temporal

Hora: 2026-10-23T10:54:29.238568-03:00

Comando:
```sh
python3 - <<'PY'
import os, subprocess, tempfile
with tempfile.TemporaryDirectory(prefix="rf3-1-home-") as home:
    env = os.environ.copy()
    env["HOME"] = home
    r = subprocess.run(["go", "test", "./..."], env=env, text=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    print(r.stdout, end="")
    print("Exit:", r.returncode)
PY
```

Salida:
```text
HOME temporal usado: /tmp/rf3-1-home-6k3revq4
go: downloading go1.26.0 (linux/amd64)
go: downloading gopkg.in/yaml.v3 v3.0.1
go: downloading github.com/pelletier/go-toml/v2 v2.4.2
go: downloading github.com/charmbracelet/bubbletea v1.3.5
go: downloading github.com/charmbracelet/lipgloss v1.1.0
go: downloading github.com/charmbracelet/x/exp/teatest v0.0.0-20260527151214-009e6338d40d
go: downloading github.com/charmbracelet/x/exp/golden v0.0.0-20260527151214-009e6338d40d
go: downloading github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6
go: downloading golang.org/x/sync v0.16.0
go: downloading github.com/charmbracelet/x/ansi v0.10.1
go: downloading github.com/charmbracelet/x/term v0.2.1
go: downloading github.com/muesli/cancelreader v0.2.2
go: downloading github.com/rivo/uniseg v0.4.7
go: downloading github.com/charmbracelet/x/cellbuf v0.0.13-0.20250311204145-2c3ea96c31dd
go: downloading github.com/muesli/termenv v0.16.0
go: downloading github.com/aymanbagabas/go-udiff v0.4.1
go: downloading golang.org/x/sys v0.35.0
go: downloading github.com/mattn/go-runewidth v0.0.16
go: downloading github.com/charmbracelet/colorprofile v0.3.2
go: downloading github.com/lucasb-eyer/go-colorful v1.2.0
go: downloading github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e
go: downloading github.com/aymanbagabas/go-osc52/v2 v2.0.1
go: downloading github.com/mattn/go-isatty v0.0.20
ok  	github.com/Group-Active-IA/active-stack/assets	0.075s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack	0.593s
ok  	github.com/Group-Active-IA/active-stack/cmd/active-stack/headless	2.462s
ok  	github.com/Group-Active-IA/active-stack/e2e/starter	0.747s
ok  	github.com/Group-Active-IA/active-stack/internal/agents	0.009s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/claude	0.010s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/codex	0.009s
ok  	github.com/Group-Active-IA/active-stack/internal/agents/opencode	0.009s
ok  	github.com/Group-Active-IA/active-stack/internal/backup	0.107s
ok  	github.com/Group-Active-IA/active-stack/internal/catalog	0.171s
?   	github.com/Group-Active-IA/active-stack/internal/contracttest	[no test files]
ok  	github.com/Group-Active-IA/active-stack/internal/filemerge	0.008s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/command	0.016s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	2.931s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config/permissions	0.115s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/external	0.076s
ok  	github.com/Group-Active-IA/active-stack/internal/harness/skill	0.031s
ok  	github.com/Group-Active-IA/active-stack/internal/i18n	0.012s
ok  	github.com/Group-Active-IA/active-stack/internal/install	0.298s
ok  	github.com/Group-Active-IA/active-stack/internal/model	0.009s
ok  	github.com/Group-Active-IA/active-stack/internal/pipeline	0.006s
ok  	github.com/Group-Active-IA/active-stack/internal/planner	0.016s
ok  	github.com/Group-Active-IA/active-stack/internal/roles	0.125s
ok  	github.com/Group-Active-IA/active-stack/internal/system	0.037s
ok  	github.com/Group-Active-IA/active-stack/internal/tui	0.057s
ok  	github.com/Group-Active-IA/active-stack/internal/uninstall	0.045s
ok  	github.com/Group-Active-IA/active-stack/internal/verify	0.065s
```

Exit: 0

## Lectura de los resultados

EVIDENCIA: con la expresión vieja, ambos mutantes sobreviven en el bloque compuesto; en assets sobrevive el anuncio sin punto. El segundo mutante ya falla en assets por el guard independiente `strings.Contains(body, "ready for archive")`: no sería veraz informar PASS allí. Con la expresión corregida ambos mutantes fallan por `contrato contradicho` en ambos lugares. Sin mutantes, la pregunta de ruteo `ready to archive?` permanece y los tests pasan. Mutantes revertidos byte a byte.

JUICIO: RF3-1 queda corregido en el alcance autorizado. No se archiva ningún change ni se sustituye el OK escrito de Juani requerido para archive.
