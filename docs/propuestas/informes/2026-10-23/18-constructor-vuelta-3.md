# Constructor — vuelta 3 de 3, fase 1

## Alcance y fuentes

Lista cerrada del auditor 17, decisión 16 y mutantes 10/14; propose §3.2–3.4 y registros 11/12 consultados. Un solo constructor, sin delegación: casos dictados por el auditor, sin independencia de revisión. Logs preexistentes excluidos. No build, instalación, push, lectura de .env ni cambios globales.

## SUPUESTO

Los controles de vocabulario delimitan las contradicciones de la lista firmada; no prueban conducta de agentes ni equivalencia semántica universal. Las variantes abiertas que el auditor pasó a deuda siguen siendo deuda. No se supone aprobación: corresponde al auditor o a la decisión escrita de Juani al cerrar esta tercera vuelta. No se encontró base de lecciones separada por nombre en archivos versionados; se leyeron las referencias L del propose.

## EVIDENCIA

Cada mutación se aplica aisladamente y se restaura en finally. SHA identifica HEAD; el diff bajo prueba incluye los tests nuevos y la mutación indicada. Se agrega hash SHA256 del archivo mutado para distinguirlo. HOME temporal; cachés de Go preexistentes reutilizadas. No se regeneran goldens durante mutaciones.

### Ítem 1 archivo nombrado ./internal/harness/config

Hora: 2026-10-23T10:39:09.189447-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `internal/harness/config/assets/review.md`; 'Conformity passed. Ready to archive.'; SHA256: `b6b1569aa015e58f1ae11baef55125bf1e619fe6d51276754bae603a585675c8`.

Comando: `go test ./internal/harness/config -run ^TestVerifyVocabularioCompuesto$ -count=1`

```text
--- FAIL: TestVerifyVocabularioCompuesto (0.13s)
    --- FAIL: TestVerifyVocabularioCompuesto/claude (0.04s)
        verify_loop_test.go:103: contrato contradicho (ready (to|for) archive[.!]|proceed to archive): ready to archive.
    --- FAIL: TestVerifyVocabularioCompuesto/codex (0.04s)
        verify_loop_test.go:103: contrato contradicho (ready (to|for) archive[.!]|proceed to archive): ready to archive.
    --- FAIL: TestVerifyVocabularioCompuesto/opencode (0.04s)
        verify_loop_test.go:103: contrato contradicho (ready (to|for) archive[.!]|proceed to archive): ready to archive.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.132s
FAIL

exit=1
```

### Ítem 1 archivo nombrado ./assets

Hora: 2026-10-23T10:39:09.706220-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `assets/commands/claude/opsx/verify.md`; 'Conformity passed. Ready to archive.'; SHA256: `7c68f07c6e53a40dfcc69d89f4704ebc42c44171cf638c65ce93184bcfc8793d`.

Comando: `go test ./assets -run ^TestVerifyConformidadNoHabilitaArchive$ -count=1`

```text
--- FAIL: TestVerifyConformidadNoHabilitaArchive (0.00s)
    commands_test.go:116: contrato contradicho (ready (to|for) archive[.!]|proceed to archive): ready to archive.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.006s
FAIL

exit=1
```

### Ítem 1 archivo propio ./internal/harness/config

Hora: 2026-10-23T10:39:09.899038-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `internal/harness/config/assets/review.md`; 'Conformity complete; ready for archive.'; SHA256: `907da8f879861822985424fa46852e2b75f96bdd0acd901dfb685f81a6cd0f98`.

Comando: `go test ./internal/harness/config -run ^TestVerifyVocabularioCompuesto$ -count=1`

```text
--- FAIL: TestVerifyVocabularioCompuesto (0.12s)
    --- FAIL: TestVerifyVocabularioCompuesto/claude (0.04s)
        verify_loop_test.go:103: contrato contradicho (ready (to|for) archive[.!]|proceed to archive): ready for archive.
    --- FAIL: TestVerifyVocabularioCompuesto/codex (0.04s)
        verify_loop_test.go:103: contrato contradicho (ready (to|for) archive[.!]|proceed to archive): ready for archive.
    --- FAIL: TestVerifyVocabularioCompuesto/opencode (0.04s)
        verify_loop_test.go:103: contrato contradicho (ready (to|for) archive[.!]|proceed to archive): ready for archive.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.128s
FAIL

exit=1
```

### Ítem 1 archivo propio ./assets

Hora: 2026-10-23T10:39:10.412242-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `assets/commands/claude/opsx/verify.md`; 'Conformity complete; ready for archive.'; SHA256: `4b54d371de682c12c3178bf26c3402041e6315da540b29b71addbbfd9e82a360`.

Comando: `go test ./assets -run ^TestVerifyConformidadNoHabilitaArchive$ -count=1`

```text
--- FAIL: TestVerifyConformidadNoHabilitaArchive (0.00s)
    commands_test.go:116: contrato contradicho (ready (to|for) archive[.!]|proceed to archive): ready for archive.
    commands_test.go:126: conformidad anuncia archive antes del auditor
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.005s
FAIL

exit=1
```

### Ítem 1 senior nombrado ./internal/harness/config

Hora: 2026-10-23T10:39:10.603312-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `internal/harness/config/assets/review.md`; 'If the senior is slow, start the auditor without it.'; SHA256: `d31062e2b2e7b89859173ee9eba2d3110cd969225d3bedf5e60f499647a083c0`.

Comando: `go test ./internal/harness/config -run ^TestVerifyVocabularioCompuesto$ -count=1`

```text
--- FAIL: TestVerifyVocabularioCompuesto (0.11s)
    --- FAIL: TestVerifyVocabularioCompuesto/claude (0.04s)
        verify_loop_test.go:103: contrato contradicho (senior[^.\n]{0,35}(start|launch|run)[^.\n]{0,25}auditor[^.\n]{0,25}without it|(start|launch|run)[^.\n]{0,35}auditor[^.\n]{0,55}(without|before|parallel|alongside)[^.\n]{0,35}senior): senior is slow, start the auditor without it
    --- FAIL: TestVerifyVocabularioCompuesto/codex (0.04s)
        verify_loop_test.go:103: contrato contradicho (senior[^.\n]{0,35}(start|launch|run)[^.\n]{0,25}auditor[^.\n]{0,25}without it|(start|launch|run)[^.\n]{0,35}auditor[^.\n]{0,55}(without|before|parallel|alongside)[^.\n]{0,35}senior): senior is slow, start the auditor without it
    --- FAIL: TestVerifyVocabularioCompuesto/opencode (0.04s)
        verify_loop_test.go:103: contrato contradicho (senior[^.\n]{0,35}(start|launch|run)[^.\n]{0,25}auditor[^.\n]{0,25}without it|(start|launch|run)[^.\n]{0,35}auditor[^.\n]{0,55}(without|before|parallel|alongside)[^.\n]{0,35}senior): senior is slow, start the auditor without it
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.118s
FAIL

exit=1
```

### Ítem 1 senior nombrado ./assets

Hora: 2026-10-23T10:39:11.108959-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `assets/commands/claude/opsx/verify.md`; 'If the senior is slow, start the auditor without it.'; SHA256: `4c78678aa3216952075578615505c7eab15779836aea22ca0c9d2cae9c5ebcd0`.

Comando: `go test ./assets -run ^TestVerifyConformidadNoHabilitaArchive$ -count=1`

```text
--- FAIL: TestVerifyConformidadNoHabilitaArchive (0.00s)
    commands_test.go:116: contrato contradicho (senior[^.\n]{0,35}(start|launch|run)[^.\n]{0,25}auditor[^.\n]{0,25}without it|(start|launch|run)[^.\n]{0,35}auditor[^.\n]{0,55}(without|before|parallel|alongside)[^.\n]{0,35}senior): senior is slow, start the auditor without it
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.006s
FAIL

exit=1
```

### Ítem 1 senior propio ./internal/harness/config

Hora: 2026-10-23T10:39:11.323694-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `internal/harness/config/assets/review.md`; 'Run the auditor before the senior returns.'; SHA256: `03d6cc3f1ba4950eeac73f6df7904cc6025d4ed52f087d084b9592fc88eb168b`.

Comando: `go test ./internal/harness/config -run ^TestVerifyVocabularioCompuesto$ -count=1`

```text
--- FAIL: TestVerifyVocabularioCompuesto (0.13s)
    --- FAIL: TestVerifyVocabularioCompuesto/claude (0.04s)
        verify_loop_test.go:103: contrato contradicho (senior[^.\n]{0,35}(start|launch|run)[^.\n]{0,25}auditor[^.\n]{0,25}without it|(start|launch|run)[^.\n]{0,35}auditor[^.\n]{0,55}(without|before|parallel|alongside)[^.\n]{0,35}senior): run the auditor before the senior
    --- FAIL: TestVerifyVocabularioCompuesto/codex (0.04s)
        verify_loop_test.go:103: contrato contradicho (senior[^.\n]{0,35}(start|launch|run)[^.\n]{0,25}auditor[^.\n]{0,25}without it|(start|launch|run)[^.\n]{0,35}auditor[^.\n]{0,55}(without|before|parallel|alongside)[^.\n]{0,35}senior): run the auditor before the senior
    --- FAIL: TestVerifyVocabularioCompuesto/opencode (0.04s)
        verify_loop_test.go:103: contrato contradicho (senior[^.\n]{0,35}(start|launch|run)[^.\n]{0,25}auditor[^.\n]{0,25}without it|(start|launch|run)[^.\n]{0,35}auditor[^.\n]{0,55}(without|before|parallel|alongside)[^.\n]{0,35}senior): run the auditor before the senior
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.131s
FAIL

exit=1
```

### Ítem 1 senior propio ./assets

Hora: 2026-10-23T10:39:11.949435-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `assets/commands/claude/opsx/verify.md`; 'Run the auditor before the senior returns.'; SHA256: `d73c3207c107293651b744a3dc3296d29183fc88d2cfb843f34dc493ad986e7f`.

Comando: `go test ./assets -run ^TestVerifyConformidadNoHabilitaArchive$ -count=1`

```text
--- FAIL: TestVerifyConformidadNoHabilitaArchive (0.00s)
    commands_test.go:116: contrato contradicho (senior[^.\n]{0,35}(start|launch|run)[^.\n]{0,25}auditor[^.\n]{0,25}without it|(start|launch|run)[^.\n]{0,35}auditor[^.\n]{0,55}(without|before|parallel|alongside)[^.\n]{0,35}senior): run the auditor before the senior
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.006s
FAIL

exit=1
```

### Ítem 1 lentes nombrado ./internal/harness/config

Hora: 2026-10-23T10:39:12.160440-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `internal/harness/config/assets/review.md`; 'Wiring: omitted.'; SHA256: `25b406fe6c720a15c383560b9cad99e3d7349d097cd77f9be2c45fb3c7538723`.

Comando: `go test ./internal/harness/config -run ^TestVerifyVocabularioCompuesto$ -count=1`

```text
--- FAIL: TestVerifyVocabularioCompuesto (0.11s)
    --- FAIL: TestVerifyVocabularioCompuesto/claude (0.04s)
        verify_loop_test.go:103: contrato contradicho ((skip|omit|remove|without|disable)[^.\n]{0,45}(wiring|comments/documentation)|(wiring|comments/documentation)[^.\n]{0,30}(omitted|optional|unnecessary)): wiring: omitted
    --- FAIL: TestVerifyVocabularioCompuesto/codex (0.04s)
        verify_loop_test.go:103: contrato contradicho ((skip|omit|remove|without|disable)[^.\n]{0,45}(wiring|comments/documentation)|(wiring|comments/documentation)[^.\n]{0,30}(omitted|optional|unnecessary)): wiring: omitted
    --- FAIL: TestVerifyVocabularioCompuesto/opencode (0.04s)
        verify_loop_test.go:103: contrato contradicho ((skip|omit|remove|without|disable)[^.\n]{0,45}(wiring|comments/documentation)|(wiring|comments/documentation)[^.\n]{0,30}(omitted|optional|unnecessary)): wiring: omitted
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.116s
FAIL

exit=1
```

### Ítem 1 lentes nombrado ./assets

Hora: 2026-10-23T10:39:12.719984-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `assets/commands/claude/opsx/verify.md`; 'Wiring: omitted.'; SHA256: `a33c2e630efa488079a0ef270e14265b92d87e8ce5f500382c73c69f98b6bd0a`.

Comando: `go test ./assets -run ^TestVerifyConformidadNoHabilitaArchive$ -count=1`

```text
--- FAIL: TestVerifyConformidadNoHabilitaArchive (0.00s)
    commands_test.go:116: contrato contradicho ((skip|omit|remove|without|disable)[^.\n]{0,45}(wiring|comments/documentation)|(wiring|comments/documentation)[^.\n]{0,30}(omitted|optional|unnecessary)): wiring: omitted
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.006s
FAIL

exit=1
```

### Ítem 1 lentes propio ./internal/harness/config

Hora: 2026-10-23T10:39:12.943033-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `internal/harness/config/assets/review.md`; 'Skip comments/documentation for this round.'; SHA256: `26658d3bb32f78da51a460b953a174cdf35d6cacb0bbf3a1fe89f438f645087a`.

Comando: `go test ./internal/harness/config -run ^TestVerifyVocabularioCompuesto$ -count=1`

```text
--- FAIL: TestVerifyVocabularioCompuesto (0.11s)
    --- FAIL: TestVerifyVocabularioCompuesto/claude (0.04s)
        verify_loop_test.go:103: contrato contradicho ((skip|omit|remove|without|disable)[^.\n]{0,45}(wiring|comments/documentation)|(wiring|comments/documentation)[^.\n]{0,30}(omitted|optional|unnecessary)): skip comments/documentation
    --- FAIL: TestVerifyVocabularioCompuesto/codex (0.03s)
        verify_loop_test.go:103: contrato contradicho ((skip|omit|remove|without|disable)[^.\n]{0,45}(wiring|comments/documentation)|(wiring|comments/documentation)[^.\n]{0,30}(omitted|optional|unnecessary)): skip comments/documentation
    --- FAIL: TestVerifyVocabularioCompuesto/opencode (0.04s)
        verify_loop_test.go:103: contrato contradicho ((skip|omit|remove|without|disable)[^.\n]{0,45}(wiring|comments/documentation)|(wiring|comments/documentation)[^.\n]{0,30}(omitted|optional|unnecessary)): skip comments/documentation
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.116s
FAIL

exit=1
```

### Ítem 1 lentes propio ./assets

Hora: 2026-10-23T10:39:13.540864-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `assets/commands/claude/opsx/verify.md`; 'Skip comments/documentation for this round.'; SHA256: `25f1edcebc2a1efc93d58968ecdbc108eb6b69aacd98f684c37d99baa29b1a01`.

Comando: `go test ./assets -run ^TestVerifyConformidadNoHabilitaArchive$ -count=1`

```text
--- FAIL: TestVerifyConformidadNoHabilitaArchive (0.00s)
    commands_test.go:116: contrato contradicho ((skip|omit|remove|without|disable)[^.\n]{0,45}(wiring|comments/documentation)|(wiring|comments/documentation)[^.\n]{0,30}(omitted|optional|unnecessary)): skip comments/documentation
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.006s
FAIL

exit=1
```

### Ítem 1 critical nombrado ./internal/harness/config

Hora: 2026-10-23T10:39:13.744433-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `internal/harness/config/assets/review.md`; 'Approval with unresolved CRITICAL is allowed.'; SHA256: `3c1da0f18d81ed3f8b7de035cabc6c52b4def0a56360c551b5c102f29b51d067`.

Comando: `go test ./internal/harness/config -run ^TestVerifyVocabularioCompuesto$ -count=1`

```text
--- FAIL: TestVerifyVocabularioCompuesto (0.13s)
    --- FAIL: TestVerifyVocabularioCompuesto/claude (0.04s)
        verify_loop_test.go:103: contrato contradicho ((approve|approval|aprobado)[^.\n]{0,45}(despite|with|even if)[^.\n]{0,25}(unresolved|pending|open) critical): approval with unresolved critical
    --- FAIL: TestVerifyVocabularioCompuesto/codex (0.04s)
        verify_loop_test.go:103: contrato contradicho ((approve|approval|aprobado)[^.\n]{0,45}(despite|with|even if)[^.\n]{0,25}(unresolved|pending|open) critical): approval with unresolved critical
    --- FAIL: TestVerifyVocabularioCompuesto/opencode (0.04s)
        verify_loop_test.go:103: contrato contradicho ((approve|approval|aprobado)[^.\n]{0,45}(despite|with|even if)[^.\n]{0,25}(unresolved|pending|open) critical): approval with unresolved critical
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.132s
FAIL

exit=1
```

### Ítem 1 critical nombrado ./assets

Hora: 2026-10-23T10:39:14.340432-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `assets/commands/claude/opsx/verify.md`; 'Approval with unresolved CRITICAL is allowed.'; SHA256: `3f6c71290b4c3ad72d63b9f3a0c93a4a2c3d8bcfad8550fb2b270a85d61a07da`.

Comando: `go test ./assets -run ^TestVerifyConformidadNoHabilitaArchive$ -count=1`

```text
--- FAIL: TestVerifyConformidadNoHabilitaArchive (0.00s)
    commands_test.go:116: contrato contradicho ((approve|approval|aprobado)[^.\n]{0,45}(despite|with|even if)[^.\n]{0,25}(unresolved|pending|open) critical): approval with unresolved critical
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.006s
FAIL

exit=1
```

### Ítem 1 critical propio ./internal/harness/config

Hora: 2026-10-23T10:39:14.549910-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `internal/harness/config/assets/review.md`; 'APROBADO even if pending CRITICAL remains.'; SHA256: `b50f4724c01da53c9ec08ff623e85556bc0fe8fc403c8e43c7b9cd3442541ee1`.

Comando: `go test ./internal/harness/config -run ^TestVerifyVocabularioCompuesto$ -count=1`

```text
--- FAIL: TestVerifyVocabularioCompuesto (0.13s)
    --- FAIL: TestVerifyVocabularioCompuesto/claude (0.05s)
        verify_loop_test.go:103: contrato contradicho ((approve|approval|aprobado)[^.\n]{0,45}(despite|with|even if)[^.\n]{0,25}(unresolved|pending|open) critical): aprobado even if pending critical
    --- FAIL: TestVerifyVocabularioCompuesto/codex (0.04s)
        verify_loop_test.go:103: contrato contradicho ((approve|approval|aprobado)[^.\n]{0,45}(despite|with|even if)[^.\n]{0,25}(unresolved|pending|open) critical): aprobado even if pending critical
    --- FAIL: TestVerifyVocabularioCompuesto/opencode (0.04s)
        verify_loop_test.go:103: contrato contradicho ((approve|approval|aprobado)[^.\n]{0,45}(despite|with|even if)[^.\n]{0,25}(unresolved|pending|open) critical): aprobado even if pending critical
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.134s
FAIL

exit=1
```

### Ítem 1 critical propio ./assets

Hora: 2026-10-23T10:39:15.162432-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `assets/commands/claude/opsx/verify.md`; 'APROBADO even if pending CRITICAL remains.'; SHA256: `844ac8659713e31b78f22ffda04221d18a3969fb21684c161c5323b7c6440c01`.

Comando: `go test ./assets -run ^TestVerifyConformidadNoHabilitaArchive$ -count=1`

```text
--- FAIL: TestVerifyConformidadNoHabilitaArchive (0.00s)
    commands_test.go:116: contrato contradicho ((approve|approval|aprobado)[^.\n]{0,45}(despite|with|even if)[^.\n]{0,25}(unresolved|pending|open) critical): aprobado even if pending critical
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.006s
FAIL

exit=1
```

### Ítem 1 contador nombrado ./internal/harness/config

Hora: 2026-10-23T10:39:15.377038-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `internal/harness/config/assets/review.md`; 'On resume reset the counter.'; SHA256: `796b77288d7f4cb9eb75d97a4e0359a80d21e3ac6cf0e538e4915e6202213491`.

Comando: `go test ./internal/harness/config -run ^TestVerifyVocabularioCompuesto$ -count=1`

```text
--- FAIL: TestVerifyVocabularioCompuesto (0.12s)
    --- FAIL: TestVerifyVocabularioCompuesto/claude (0.04s)
        verify_loop_test.go:103: contrato contradicho ((always|may|must|should|resume|recovery)[^.\n]{0,30}(reset|restart)[^.\n]{0,30}(counter|round count)): resume reset the counter
    --- FAIL: TestVerifyVocabularioCompuesto/codex (0.04s)
        verify_loop_test.go:103: contrato contradicho ((always|may|must|should|resume|recovery)[^.\n]{0,30}(reset|restart)[^.\n]{0,30}(counter|round count)): resume reset the counter
    --- FAIL: TestVerifyVocabularioCompuesto/opencode (0.04s)
        verify_loop_test.go:103: contrato contradicho ((always|may|must|should|resume|recovery)[^.\n]{0,30}(reset|restart)[^.\n]{0,30}(counter|round count)): resume reset the counter
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.119s
FAIL

exit=1
```

### Ítem 1 contador nombrado ./assets

Hora: 2026-10-23T10:39:15.954006-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `assets/commands/claude/opsx/verify.md`; 'On resume reset the counter.'; SHA256: `478aa77a17d20a328c6f5bcc84893d664e96eb9640d2b2c6af8ec003fed0b9b9`.

Comando: `go test ./assets -run ^TestVerifyConformidadNoHabilitaArchive$ -count=1`

```text
--- FAIL: TestVerifyConformidadNoHabilitaArchive (0.00s)
    commands_test.go:116: contrato contradicho ((always|may|must|should|resume|recovery)[^.\n]{0,30}(reset|restart)[^.\n]{0,30}(counter|round count)): resume reset the counter
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.006s
FAIL

exit=1
```

### Ítem 1 contador propio ./internal/harness/config

Hora: 2026-10-23T10:39:16.166847-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `internal/harness/config/assets/review.md`; 'Recovery must restart the round count.'; SHA256: `6467994d8179cdcab68a8db9db30098743a79d4fb6d39293f9ac2e2719bfc728`.

Comando: `go test ./internal/harness/config -run ^TestVerifyVocabularioCompuesto$ -count=1`

```text
--- FAIL: TestVerifyVocabularioCompuesto (0.12s)
    --- FAIL: TestVerifyVocabularioCompuesto/claude (0.04s)
        verify_loop_test.go:103: contrato contradicho ((always|may|must|should|resume|recovery)[^.\n]{0,30}(reset|restart)[^.\n]{0,30}(counter|round count)): recovery must restart the round count
    --- FAIL: TestVerifyVocabularioCompuesto/codex (0.04s)
        verify_loop_test.go:103: contrato contradicho ((always|may|must|should|resume|recovery)[^.\n]{0,30}(reset|restart)[^.\n]{0,30}(counter|round count)): recovery must restart the round count
    --- FAIL: TestVerifyVocabularioCompuesto/opencode (0.04s)
        verify_loop_test.go:103: contrato contradicho ((always|may|must|should|resume|recovery)[^.\n]{0,30}(reset|restart)[^.\n]{0,30}(counter|round count)): recovery must restart the round count
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.125s
FAIL

exit=1
```

### Ítem 1 contador propio ./assets

Hora: 2026-10-23T10:39:16.758608-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `assets/commands/claude/opsx/verify.md`; 'Recovery must restart the round count.'; SHA256: `3491504aa46631aaedeb370311f1cee40fd2ee0f58723ac15f66b180bee5a9fe`.

Comando: `go test ./assets -run ^TestVerifyConformidadNoHabilitaArchive$ -count=1`

```text
--- FAIL: TestVerifyConformidadNoHabilitaArchive (0.00s)
    commands_test.go:116: contrato contradicho ((always|may|must|should|resume|recovery)[^.\n]{0,30}(reset|restart)[^.\n]{0,30}(counter|round count)): recovery must restart the round count
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.006s
FAIL

exit=1
```

### Ítem 1 cuarta nombrado ./internal/harness/config

Hora: 2026-10-23T10:39:16.976356-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `internal/harness/config/assets/review.md`; 'If blockers remain, start round 4 automatically.'; SHA256: `fa45e5c68dd8a2f6b996bb7c457d63445740760d2427bebc87a58cd5d2ce1847`.

Comando: `go test ./internal/harness/config -run ^TestVerifyVocabularioCompuesto$ -count=1`

```text
--- FAIL: TestVerifyVocabularioCompuesto (0.11s)
    --- FAIL: TestVerifyVocabularioCompuesto/claude (0.04s)
        verify_loop_test.go:103: contrato contradicho ((start|allow|permit|begin|launch)[^.\n]{0,25}(round (4|four)|fourth round)): start round 4
    --- FAIL: TestVerifyVocabularioCompuesto/codex (0.04s)
        verify_loop_test.go:103: contrato contradicho ((start|allow|permit|begin|launch)[^.\n]{0,25}(round (4|four)|fourth round)): start round 4
    --- FAIL: TestVerifyVocabularioCompuesto/opencode (0.04s)
        verify_loop_test.go:103: contrato contradicho ((start|allow|permit|begin|launch)[^.\n]{0,25}(round (4|four)|fourth round)): start round 4
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.118s
FAIL

exit=1
```

### Ítem 1 cuarta nombrado ./assets

Hora: 2026-10-23T10:39:17.561200-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `assets/commands/claude/opsx/verify.md`; 'If blockers remain, start round 4 automatically.'; SHA256: `e94b730909532183b7ad27cc3368de3f206c7c2248095d475b2746ff7712a644`.

Comando: `go test ./assets -run ^TestVerifyConformidadNoHabilitaArchive$ -count=1`

```text
--- FAIL: TestVerifyConformidadNoHabilitaArchive (0.00s)
    commands_test.go:116: contrato contradicho ((start|allow|permit|begin|launch)[^.\n]{0,25}(round (4|four)|fourth round)): start round 4
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.006s
FAIL

exit=1
```

### Ítem 1 cuarta propio ./internal/harness/config

Hora: 2026-10-23T10:39:17.782186-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `internal/harness/config/assets/review.md`; 'Allow a fourth round for remaining findings.'; SHA256: `00a7860e255359801f069465ae7824d5dfaab1600091f5489d68ca5a1444a388`.

Comando: `go test ./internal/harness/config -run ^TestVerifyVocabularioCompuesto$ -count=1`

```text
--- FAIL: TestVerifyVocabularioCompuesto (0.11s)
    --- FAIL: TestVerifyVocabularioCompuesto/claude (0.04s)
        verify_loop_test.go:103: contrato contradicho ((start|allow|permit|begin|launch)[^.\n]{0,25}(round (4|four)|fourth round)): allow a fourth round
    --- FAIL: TestVerifyVocabularioCompuesto/codex (0.04s)
        verify_loop_test.go:103: contrato contradicho ((start|allow|permit|begin|launch)[^.\n]{0,25}(round (4|four)|fourth round)): allow a fourth round
    --- FAIL: TestVerifyVocabularioCompuesto/opencode (0.04s)
        verify_loop_test.go:103: contrato contradicho ((start|allow|permit|begin|launch)[^.\n]{0,25}(round (4|four)|fourth round)): allow a fourth round
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.118s
FAIL

exit=1
```

### Ítem 1 cuarta propio ./assets

Hora: 2026-10-23T10:39:18.379067-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Mutación: `assets/commands/claude/opsx/verify.md`; 'Allow a fourth round for remaining findings.'; SHA256: `6f7199e88f9c5c8b06c6c0445f4bceca827418fc05743560c185b957cbec3eb1`.

Comando: `go test ./assets -run ^TestVerifyConformidadNoHabilitaArchive$ -count=1`

```text
--- FAIL: TestVerifyConformidadNoHabilitaArchive (0.00s)
    commands_test.go:116: contrato contradicho ((start|allow|permit|begin|launch)[^.\n]{0,25}(round (4|four)|fourth round)): allow a fourth round
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.006s
FAIL

exit=1
```

### Ítem 1 verde

Hora: 2026-10-23T10:39:18.593863-03:00; SHA: `248173fbd0df45d0c35baa62c838ab31ccbb87c1`; HOME: `/tmp/active-stack-v3-home-rdcenx5a`.

Comando: `go test ./internal/harness/config ./assets -run "TestVerify(VocabularioCompuesto|ConformidadNoHabilitaArchive)" -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.139s
ok  	github.com/Group-Active-IA/active-stack/assets	0.005s

exit=0
```

### Ítem 2 pendiente nombrado ./internal/harness/config

Hora: 2026-10-23T10:40:07.301128-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Mutación: `internal/harness/config/assets/ship.md`; 'A pending job is green.'; SHA256: `0d347aa5a8b73dc2407fe69d1cc50bc6b4b5a7129c7de1d316b2f633ab59d018`.

Comando: `go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestShipEsperaCICompleto (0.05s)
    ship_ci_test.go:20: contrato contradicho ((pending|queued|in_progress|cancelled|skipped|omitted|neutral|timed_out|another sha|different sha|zero runs|no runs|failed query|query failure)[^.\n]{0,75}(is|are|counts as|considered|treat as|means) (green|success)\b): pending job is green
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.052s
FAIL

exit=1
```

### Ítem 2 pendiente propio ./internal/harness/config

Hora: 2026-10-23T10:40:07.771476-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Mutación: `internal/harness/config/assets/ship.md`; 'Treat CI as green even with queued jobs.'; SHA256: `ee7ab6c26eef6a112507fb6069033221d2257ac91a3393f296835e7f83cc0c40`.

Comando: `go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestShipEsperaCICompleto (0.05s)
    ship_ci_test.go:20: contrato contradicho ((green|success)\b[^.\n]{0,45}(despite|with|even with)[^.\n]{0,35}(pending|queued|skipped|omitted|another sha|different sha|zero runs|no runs|failed query|query failure)): green even with queued
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.052s
FAIL

exit=1
```

### Ítem 2 omitido nombrado ./internal/harness/config

Hora: 2026-10-23T10:40:08.235271-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Mutación: `internal/harness/config/assets/ship.md`; 'A skipped job counts as green.'; SHA256: `de37cdd2852ce24a8e5937da65229a13d7f2149a2c8fc81724e352d8fc7b18c9`.

Comando: `go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestShipEsperaCICompleto (0.04s)
    ship_ci_test.go:20: contrato contradicho ((pending|queued|in_progress|cancelled|skipped|omitted|neutral|timed_out|another sha|different sha|zero runs|no runs|failed query|query failure)[^.\n]{0,75}(is|are|counts as|considered|treat as|means) (green|success)\b): skipped job counts as green
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.047s
FAIL

exit=1
```

### Ítem 2 omitido propio ./internal/harness/config

Hora: 2026-10-23T10:40:08.687757-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Mutación: `internal/harness/config/assets/ship.md`; 'CI is green despite omitted jobs.'; SHA256: `8cb9d8519b911baeee72dc725749672ddebe566c77631670dbae5db675de7603`.

Comando: `go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestShipEsperaCICompleto (0.05s)
    ship_ci_test.go:20: contrato contradicho ((green|success)\b[^.\n]{0,45}(despite|with|even with)[^.\n]{0,35}(pending|queued|skipped|omitted|another sha|different sha|zero runs|no runs|failed query|query failure)): green despite omitted
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.051s
FAIL

exit=1
```

### Ítem 2 SHA ajeno nombrado ./internal/harness/config

Hora: 2026-10-23T10:40:09.148151-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Mutación: `internal/harness/config/assets/ship.md`; 'A result for another SHA is green.'; SHA256: `bfd7b9a88f14de6fa5c4f78b5896e00225d2c37f0ccda6c5d86cd9503ba306a7`.

Comando: `go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestShipEsperaCICompleto (0.05s)
    ship_ci_test.go:20: contrato contradicho ((pending|queued|in_progress|cancelled|skipped|omitted|neutral|timed_out|another sha|different sha|zero runs|no runs|failed query|query failure)[^.\n]{0,75}(is|are|counts as|considered|treat as|means) (green|success)\b): another sha is green
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.049s
FAIL

exit=1
```

### Ítem 2 SHA ajeno propio ./internal/harness/config

Hora: 2026-10-23T10:40:09.593921-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Mutación: `internal/harness/config/assets/ship.md`; 'A different SHA means success.'; SHA256: `8aae0c8984471996bb9016a54103e52a0a7728404559b2b5e8df5b0d55dc94fa`.

Comando: `go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestShipEsperaCICompleto (0.05s)
    ship_ci_test.go:20: contrato contradicho ((pending|queued|in_progress|cancelled|skipped|omitted|neutral|timed_out|another sha|different sha|zero runs|no runs|failed query|query failure)[^.\n]{0,75}(is|are|counts as|considered|treat as|means) (green|success)\b): different sha means success
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.053s
FAIL

exit=1
```

### Ítem 2 cero runs nombrado ./internal/harness/config

Hora: 2026-10-23T10:40:10.068577-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Mutación: `internal/harness/config/assets/ship.md`; 'Zero runs are green.'; SHA256: `b185297b06899ecb4eaa92c2e9ba54ed12ef0a2f02f3d4ced568f49d62673dcf`.

Comando: `go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestShipEsperaCICompleto (0.05s)
    ship_ci_test.go:20: contrato contradicho ((pending|queued|in_progress|cancelled|skipped|omitted|neutral|timed_out|another sha|different sha|zero runs|no runs|failed query|query failure)[^.\n]{0,75}(is|are|counts as|considered|treat as|means) (green|success)\b): zero runs are green
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.051s
FAIL

exit=1
```

### Ítem 2 cero runs propio ./internal/harness/config

Hora: 2026-10-23T10:40:10.533727-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Mutación: `internal/harness/config/assets/ship.md`; 'No runs counts as success.'; SHA256: `a5a8d0a57f852f317edbe77b9afcd6bed044e93b725f206172295c4b9afc2593`.

Comando: `go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestShipEsperaCICompleto (0.04s)
    ship_ci_test.go:20: contrato contradicho ((pending|queued|in_progress|cancelled|skipped|omitted|neutral|timed_out|another sha|different sha|zero runs|no runs|failed query|query failure)[^.\n]{0,75}(is|are|counts as|considered|treat as|means) (green|success)\b): no runs counts as success
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.049s
FAIL

exit=1
```

### Ítem 2 consulta fallida nombrado ./internal/harness/config

Hora: 2026-10-23T10:40:10.989218-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Mutación: `internal/harness/config/assets/ship.md`; 'A failed query is green.'; SHA256: `3954fa8084cb0a248ddfc74fdfce97c67b754d1601d9e3d1ce99f9a6eb866193`.

Comando: `go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestShipEsperaCICompleto (0.05s)
    ship_ci_test.go:20: contrato contradicho ((pending|queued|in_progress|cancelled|skipped|omitted|neutral|timed_out|another sha|different sha|zero runs|no runs|failed query|query failure)[^.\n]{0,75}(is|are|counts as|considered|treat as|means) (green|success)\b): failed query is green
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.056s
FAIL

exit=1
```

### Ítem 2 consulta fallida propio ./internal/harness/config

Hora: 2026-10-23T10:40:11.450460-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Mutación: `internal/harness/config/assets/ship.md`; 'A query failure means success.'; SHA256: `e71c1c494db0cea1ec8a138fe980e4c3775b6c561bc1b97f3cf564b800bf85da`.

Comando: `go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestShipEsperaCICompleto (0.05s)
    ship_ci_test.go:20: contrato contradicho ((pending|queued|in_progress|cancelled|skipped|omitted|neutral|timed_out|another sha|different sha|zero runs|no runs|failed query|query failure)[^.\n]{0,75}(is|are|counts as|considered|treat as|means) (green|success)\b): query failure means success
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.052s
FAIL

exit=1
```

### Ítem 2 pendiente nombrado ./assets

Hora: 2026-10-23T10:40:11.913113-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Mutación: `assets/commands/claude/opsx/ship.md`; 'A pending job is green.'; SHA256: `85a1649b41a57ad1787a954b8c695b760d5bd812a664bdd3754ed78f57f9435d`.

Comando: `go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestComandoShipEsperaCICompleto (0.01s)
    ship_ci_test.go:23: contrato contradicho ((pending|queued|in_progress|cancelled|skipped|omitted|neutral|timed_out|another sha|different sha|zero runs|no runs|failed query|query failure)[^.\n]{0,75}(is|are|counts as|considered|treat as|means) (green|success)\b): pending job is green
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.015s
FAIL

exit=1
```

### Ítem 2 pendiente propio ./assets

Hora: 2026-10-23T10:40:12.136672-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Mutación: `assets/commands/claude/opsx/ship.md`; 'Treat CI as green even with queued jobs.'; SHA256: `f9cbbea39dcee5d8e32031205500d508c222bd821607aaab0fee29843b456b6b`.

Comando: `go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestComandoShipEsperaCICompleto (0.01s)
    ship_ci_test.go:23: contrato contradicho ((green|success)\b[^.\n]{0,45}(despite|with|even with)[^.\n]{0,35}(pending|queued|skipped|omitted|another sha|different sha|zero runs|no runs|failed query|query failure)): green even with queued
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.013s
FAIL

exit=1
```

### Ítem 2 omitido nombrado ./assets

Hora: 2026-10-23T10:40:12.352373-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Mutación: `assets/commands/claude/opsx/ship.md`; 'A skipped job counts as green.'; SHA256: `43eb8eb5f91ed7123be962630fa28f00837f432419640a8df5bebfb2cc6eba6c`.

Comando: `go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestComandoShipEsperaCICompleto (0.01s)
    ship_ci_test.go:23: contrato contradicho ((pending|queued|in_progress|cancelled|skipped|omitted|neutral|timed_out|another sha|different sha|zero runs|no runs|failed query|query failure)[^.\n]{0,75}(is|are|counts as|considered|treat as|means) (green|success)\b): skipped job counts as green
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.013s
FAIL

exit=1
```

### Ítem 2 omitido propio ./assets

Hora: 2026-10-23T10:40:12.566314-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Mutación: `assets/commands/claude/opsx/ship.md`; 'CI is green despite omitted jobs.'; SHA256: `87fb18eb7a92437b69b1768907b0f52193fe66bd9f5b1d07147012a61f66d9a3`.

Comando: `go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestComandoShipEsperaCICompleto (0.01s)
    ship_ci_test.go:23: contrato contradicho ((green|success)\b[^.\n]{0,45}(despite|with|even with)[^.\n]{0,35}(pending|queued|skipped|omitted|another sha|different sha|zero runs|no runs|failed query|query failure)): green despite omitted
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.013s
FAIL

exit=1
```

### Ítem 2 SHA ajeno nombrado ./assets

Hora: 2026-10-23T10:40:12.781785-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Mutación: `assets/commands/claude/opsx/ship.md`; 'A result for another SHA is green.'; SHA256: `6ab404cc7975f5f3a60955f5c032eac3f8ea6c0dae29b4e30555081392d1f17e`.

Comando: `go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestComandoShipEsperaCICompleto (0.01s)
    ship_ci_test.go:23: contrato contradicho ((pending|queued|in_progress|cancelled|skipped|omitted|neutral|timed_out|another sha|different sha|zero runs|no runs|failed query|query failure)[^.\n]{0,75}(is|are|counts as|considered|treat as|means) (green|success)\b): another sha is green
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.013s
FAIL

exit=1
```

### Ítem 2 SHA ajeno propio ./assets

Hora: 2026-10-23T10:40:12.998565-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Mutación: `assets/commands/claude/opsx/ship.md`; 'A different SHA means success.'; SHA256: `d69f54cb85485df43be421f46db2c1183b993dc6a62fc10558d6d95bba79ad42`.

Comando: `go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestComandoShipEsperaCICompleto (0.01s)
    ship_ci_test.go:23: contrato contradicho ((pending|queued|in_progress|cancelled|skipped|omitted|neutral|timed_out|another sha|different sha|zero runs|no runs|failed query|query failure)[^.\n]{0,75}(is|are|counts as|considered|treat as|means) (green|success)\b): different sha means success
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.013s
FAIL

exit=1
```

### Ítem 2 cero runs nombrado ./assets

Hora: 2026-10-23T10:40:13.219357-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Mutación: `assets/commands/claude/opsx/ship.md`; 'Zero runs are green.'; SHA256: `84a731c92dd3c26979bf94df5b56b22bce6afb0e8bef620101d5ff5ba3dc5217`.

Comando: `go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestComandoShipEsperaCICompleto (0.01s)
    ship_ci_test.go:23: contrato contradicho ((pending|queued|in_progress|cancelled|skipped|omitted|neutral|timed_out|another sha|different sha|zero runs|no runs|failed query|query failure)[^.\n]{0,75}(is|are|counts as|considered|treat as|means) (green|success)\b): zero runs are green
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.015s
FAIL

exit=1
```

### Ítem 2 cero runs propio ./assets

Hora: 2026-10-23T10:40:13.433619-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Mutación: `assets/commands/claude/opsx/ship.md`; 'No runs counts as success.'; SHA256: `a9d499c8cd983985b0ba8bc7d7e713e30baad2b5d26968dfae07ec7a81630b24`.

Comando: `go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestComandoShipEsperaCICompleto (0.01s)
    ship_ci_test.go:23: contrato contradicho ((pending|queued|in_progress|cancelled|skipped|omitted|neutral|timed_out|another sha|different sha|zero runs|no runs|failed query|query failure)[^.\n]{0,75}(is|are|counts as|considered|treat as|means) (green|success)\b): no runs counts as success
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.014s
FAIL

exit=1
```

### Ítem 2 consulta fallida nombrado ./assets

Hora: 2026-10-23T10:40:13.653015-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Mutación: `assets/commands/claude/opsx/ship.md`; 'A failed query is green.'; SHA256: `d547e0029d529ab72bcd49b7444fe46643a2077a1961113a82ba5990c31372a3`.

Comando: `go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestComandoShipEsperaCICompleto (0.01s)
    ship_ci_test.go:23: contrato contradicho ((pending|queued|in_progress|cancelled|skipped|omitted|neutral|timed_out|another sha|different sha|zero runs|no runs|failed query|query failure)[^.\n]{0,75}(is|are|counts as|considered|treat as|means) (green|success)\b): failed query is green
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.012s
FAIL

exit=1
```

### Ítem 2 consulta fallida propio ./assets

Hora: 2026-10-23T10:40:13.871649-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Mutación: `assets/commands/claude/opsx/ship.md`; 'A query failure means success.'; SHA256: `a562d2300403dc8d102d0443920e085a5704a76470ca9155cee2cfa93f730f38`.

Comando: `go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestComandoShipEsperaCICompleto (0.01s)
    ship_ci_test.go:23: contrato contradicho ((pending|queued|in_progress|cancelled|skipped|omitted|neutral|timed_out|another sha|different sha|zero runs|no runs|failed query|query failure)[^.\n]{0,75}(is|are|counts as|considered|treat as|means) (green|success)\b): query failure means success
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.014s
FAIL

exit=1
```

### Ítem 2 verde

Hora: 2026-10-23T10:40:14.096064-03:00; SHA: `9831a41d7ff532e2ee0c8f5768185f8a95e4387b`; HOME: `/tmp/active-stack-v3-home-losd9pr7`.

Comando: `go test ./internal/harness/config ./assets -run "Test(ShipEsperaCICompleto|ComandoShipEsperaCICompleto)" -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.054s
ok  	github.com/Group-Active-IA/active-stack/assets	0.013s

exit=0
```
