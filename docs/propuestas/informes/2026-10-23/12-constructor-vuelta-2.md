# Constructor — vuelta 2 de verify, fase 1

Base firmada: `4748680`. Rama: `feat/metodologia-proptech`. Alcance: exclusivamente los cinco ítems del auditor 11; notas RF1-4, RF1-5b y RF1-6 sin cambios. No se ejecuta build, instalación ni push; no se lee .env ni se modifica configuración global. Los logs preexistentes no se incluyen.

## Ítem 1 — contratos de sentido y mutantes aislados

Tests reforzados antes de mutar; cada archivo se restaura en finally, sin regenerar goldens. M4 elimina una regla: su mensaje muestra la negación requerida ausente (no introduce frase nueva). HOME temporal en todas las corridas. GOCACHE/GOPATH reutilizan la caché existente.
### M1

Hora: 2026-10-23T10:22:41.703080-03:00; SHA: `8d05c4eefca7d3fe4ed3c1667806c7c3580402f7`; HOME: `/tmp/active-stack-v2-home-4_ciocdq`.

Comando: `go test ./internal/harness/config -run TestShipEsperaCICompleto -count=1`

```text
--- FAIL: TestShipEsperaCICompleto (0.00s)
    ship_ci_test.go:11: CI contradictorio: A cancelled job counts as green when its rerun passes.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL

exit=1
```

Hora: 2026-10-23T10:22:42.179514-03:00; SHA: `8d05c4eefca7d3fe4ed3c1667806c7c3580402f7`; HOME: `/tmp/active-stack-v2-home-4_ciocdq`.

Comando: `go test ./internal/harness/config -run TestShipEsperaCICompleto -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s

exit=0
```
### M1b

Hora: 2026-10-23T10:22:42.566785-03:00; SHA: `8d05c4eefca7d3fe4ed3c1667806c7c3580402f7`; HOME: `/tmp/active-stack-v2-home-4_ciocdq`.

Comando: `go test ./assets -run TestComandoShipEsperaCICompleto -count=1`

```text
--- FAIL: TestComandoShipEsperaCICompleto (0.00s)
    ship_ci_test.go:17: CI contradictorio: A cancelled job counts as green when its rerun passes.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.003s
FAIL

exit=1
```

Hora: 2026-10-23T10:22:42.785034-03:00; SHA: `8d05c4eefca7d3fe4ed3c1667806c7c3580402f7`; HOME: `/tmp/active-stack-v2-home-4_ciocdq`.

Comando: `go test ./assets -run TestComandoShipEsperaCICompleto -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/assets	0.003s

exit=0
```
### M2

Hora: 2026-10-23T10:22:42.989519-03:00; SHA: `8d05c4eefca7d3fe4ed3c1667806c7c3580402f7`; HOME: `/tmp/active-stack-v2-home-4_ciocdq`.

Comando: `go test ./internal/harness/config -run TestVerifyContratoBucle -count=1`

```text
--- FAIL: TestVerifyContratoBucle (0.00s)
    verify_loop_test.go:13: secuencia contradictoria: if blockers remain, start round 4 automatically.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.005s
FAIL

exit=1
```

Hora: 2026-10-23T10:22:43.456284-03:00; SHA: `8d05c4eefca7d3fe4ed3c1667806c7c3580402f7`; HOME: `/tmp/active-stack-v2-home-4_ciocdq`.

Comando: `go test ./internal/harness/config -run TestVerifyContratoBucle -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s

exit=0
```
### M3

Hora: 2026-10-23T10:22:43.695227-03:00; SHA: `8d05c4eefca7d3fe4ed3c1667806c7c3580402f7`; HOME: `/tmp/active-stack-v2-home-4_ciocdq`.

Comando: `go test ./internal/harness/config -run TestVerifyContratoBucle|TestAuditorDecide -count=1`

```text
--- FAIL: TestVerifyContratoBucle (0.00s)
    verify_loop_test.go:13: secuencia contradictoria: launch the auditor in parallel with the senior.
    verify_loop_test.go:17: falta espera obligatoria; no se permite: Launch the auditor in parallel with the senior.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.008s
FAIL

exit=1
```

Hora: 2026-10-23T10:22:44.186185-03:00; SHA: `8d05c4eefca7d3fe4ed3c1667806c7c3580402f7`; HOME: `/tmp/active-stack-v2-home-4_ciocdq`.

Comando: `go test ./internal/harness/config -run TestVerifyContratoBucle|TestAuditorDecide -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.008s

exit=0
```
### M4

Hora: 2026-10-23T10:22:44.411904-03:00; SHA: `8d05c4eefca7d3fe4ed3c1667806c7c3580402f7`; HOME: `/tmp/active-stack-v2-home-4_ciocdq`.

Comando: `go test ./internal/harness/config -run TestAuditorDecide -count=1`

```text
--- FAIL: TestAuditorDecideSinAutoridadParalela (0.00s)
    verify_loop_test.go:42: falta negación obligatoria: Missing required reports prevent approval.; permitir aprobación sin informes contradice el contrato
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL

exit=1
```

Hora: 2026-10-23T10:22:44.890972-03:00; SHA: `8d05c4eefca7d3fe4ed3c1667806c7c3580402f7`; HOME: `/tmp/active-stack-v2-home-4_ciocdq`.

Comando: `go test ./internal/harness/config -run TestAuditorDecide -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s

exit=0
```
### M5

Hora: 2026-10-23T10:22:45.115274-03:00; SHA: `8d05c4eefca7d3fe4ed3c1667806c7c3580402f7`; HOME: `/tmp/active-stack-v2-home-4_ciocdq`.

Comando: `go test ./assets -run TestVerifyConformidadNoHabilitaArchive -count=1`

```text
--- FAIL: TestVerifyConformidadNoHabilitaArchive (0.00s)
    commands_test.go:116: conformidad anuncia cierre: All checks passed. Proceed to archive.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.002s
FAIL

exit=1
```

Hora: 2026-10-23T10:22:45.330498-03:00; SHA: `8d05c4eefca7d3fe4ed3c1667806c7c3580402f7`; HOME: `/tmp/active-stack-v2-home-4_ciocdq`.

Comando: `go test ./assets -run TestVerifyConformidadNoHabilitaArchive -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/assets	0.003s

exit=0
```
### M6

Hora: 2026-10-23T10:22:45.501840-03:00; SHA: `8d05c4eefca7d3fe4ed3c1667806c7c3580402f7`; HOME: `/tmp/active-stack-v2-home-4_ciocdq`.

Comando: `go test ./internal/harness/config -run TestArchiveBaseCompletaConEvidencia -count=1`

```text
--- FAIL: TestArchiveBaseCompletaConEvidencia (0.00s)
    knowledge_base_toggle_test.go:720: historia contradictoria: Rewrite old lessons freely.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL

exit=1
```

Hora: 2026-10-23T10:22:45.986365-03:00; SHA: `8d05c4eefca7d3fe4ed3c1667806c7c3580402f7`; HOME: `/tmp/active-stack-v2-home-4_ciocdq`.

Comando: `go test ./internal/harness/config -run TestArchiveBaseCompletaConEvidencia -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s

exit=0
```
### M7

Hora: 2026-10-23T10:22:46.213936-03:00; SHA: `8d05c4eefca7d3fe4ed3c1667806c7c3580402f7`; HOME: `/tmp/active-stack-v2-home-4_ciocdq`.

Comando: `go test ./internal/harness/config -run TestGovernance -count=1`

```text
--- FAIL: TestGovernanceAprobacionHumanaEnProposal (0.00s)
    governance_approval_test.go:16: firma contradictoria: An agent may fill it after chat approval.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.003s
FAIL

exit=1
```

Hora: 2026-10-23T10:22:46.672607-03:00; SHA: `8d05c4eefca7d3fe4ed3c1667806c7c3580402f7`; HOME: `/tmp/active-stack-v2-home-4_ciocdq`.

Comando: `go test ./internal/harness/config -run TestGovernance -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.003s

exit=0
```
## Ítem 2 — M8, rechazo de bloques duplicados

Fixture: instalación en HOME de test y duplicación del bloque entero en Claude y Codex. Se invoca específicamente config:compose y se exige error de markers duplicados; no alcanza otro check fallido. OpenCode entrega un prompt JSON sin markers y queda cubierto por TestEntregaComparaCompose.

Hora: 2026-10-23T10:23:20.878013-03:00; SHA: `f3f72af3d2d53c2271b944b45145f4ed1ab4dbea`; HOME: `/tmp/active-stack-v2-home-1lqtgeyh`.

Comando: `go test ./internal/verify -run TestEntregaRechazaBloqueDuplicado -count=1`

```text
--- FAIL: TestEntregaRechazaBloqueDuplicado (0.00s)
    --- FAIL: TestEntregaRechazaBloqueDuplicado/claude (0.00s)
        delivery_test.go:77: bloque duplicado debe fallar por markers duplicados; resultado: <nil>
    --- FAIL: TestEntregaRechazaBloqueDuplicado/codex (0.00s)
        delivery_test.go:77: bloque duplicado debe fallar por markers duplicados; resultado: <nil>
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/verify	0.005s
FAIL

exit=1
```

Hora: 2026-10-23T10:23:21.400470-03:00; SHA: `f3f72af3d2d53c2271b944b45145f4ed1ab4dbea`; HOME: `/tmp/active-stack-v2-home-1lqtgeyh`.

Comando: `go test ./internal/verify -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/internal/verify	0.023s

exit=0
```
## Ítem 3 — RF1-2, condición de auditor ausente

Test con auditor y KB activos: cada línea que ordena marcar la fuente ausente debe anteponer la condición. Rojo contra el texto original, luego verde. Golden: exclusivamente se antepone `If the auditor is absent, ` al mandato del prompt, en las dos variantes.

Hora: 2026-10-23T10:24:13.085325-03:00; SHA: `c3259558cb1cc7677b54c1333542f3578eac2bb2`; HOME: `/tmp/active-stack-v2-home-s8nsqirb`.

Comando: `go test ./internal/harness/config -run TestArchiveBaseCompletaConEvidencia -count=1`

```text
--- FAIL: TestArchiveBaseCompletaConEvidencia (0.00s)
    knowledge_base_toggle_test.go:721: auditor activo: instrucción sin condición: APPEND decisions, lessons L-NNN from auditor LECCIONES and debt from auditor DEUDA, dated, with what each one cost. Only append lessons. Each entry states what changed and typed evidence: SHA, test or decision. Mark auditor source absent and do not invent its lists; closure requires Juani's written decision on degraded scope. Return the two
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL

exit=1
```

Hora: 2026-10-23T10:24:13.476144-03:00; SHA: `c3259558cb1cc7677b54c1333542f3578eac2bb2`; HOME: `/tmp/active-stack-v2-home-s8nsqirb`.

Comando: `go test ./internal/harness/config -run TestArchiveBaseCompletaConEvidencia -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s

exit=0
```

Hora: 2026-10-23T10:24:13.939477-03:00; SHA: `c3259558cb1cc7677b54c1333542f3578eac2bb2`; HOME: `/tmp/active-stack-v2-home-s8nsqirb`.

Comando: `go test ./internal/harness/config -run ^TestCompose_AllToggles -update -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.005s

exit=0
```

Hora: 2026-10-23T10:24:14.158283-03:00; SHA: `c3259558cb1cc7677b54c1333542f3578eac2bb2`; HOME: `/tmp/active-stack-v2-home-s8nsqirb`.

Comando: `go test ./internal/harness/config -run TestArchiveBaseCompletaConEvidencia|^TestCompose_AllToggles -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.008s

exit=0
```
## Ítem 4 — RF1-3, alcance explícito del CI

Los tests de ambas copias prohíben `For this repository` y exigen `For active-stack`. Rojo antes de corregir cada copia. Goldens: sólo cambia el sujeto del requisito de workflow, sin cambiar jobs ni eventos.

Hora: 2026-10-23T10:24:14.418216-03:00; SHA: `13ce7f23fad84c951bf6f01995bab9b0ad3490ad`; HOME: `/tmp/active-stack-v2-home-s8nsqirb`.

Comando: `go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestShipEsperaCICompleto (0.00s)
    ship_ci_test.go:10: CI ambiguo: For this repository
    --- FAIL: TestShipEsperaCICompleto/For_active-stack,_require_.github/workflows/ci.yml_on_push,_with_test_on_Ubuntu_and_Windows. (0.00s)
        ship_ci_test.go:33: falta contrato de CI: For active-stack, require .github/workflows/ci.yml on push, with test on Ubuntu and Windows.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL

exit=1
```

Hora: 2026-10-23T10:24:14.813401-03:00; SHA: `13ce7f23fad84c951bf6f01995bab9b0ad3490ad`; HOME: `/tmp/active-stack-v2-home-s8nsqirb`.

Comando: `go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestComandoShipEsperaCICompleto (0.00s)
    ship_ci_test.go:16: CI ambiguo: For this repository
    --- FAIL: TestComandoShipEsperaCICompleto/For_active-stack,_require_.github/workflows/ci.yml_on_push,_with_test_on_Ubuntu_and_Windows. (0.00s)
        ship_ci_test.go:39: falta contrato de CI: For active-stack, require .github/workflows/ci.yml on push, with test on Ubuntu and Windows.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.003s
FAIL

exit=1
```

Hora: 2026-10-23T10:24:15.016476-03:00; SHA: `13ce7f23fad84c951bf6f01995bab9b0ad3490ad`; HOME: `/tmp/active-stack-v2-home-s8nsqirb`.

Comando: `go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestShipEsperaCICompleto (0.00s)
    ship_ci_test.go:10: CI ambiguo: For this repository
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL

exit=1
```
La primera aserción abarcaba también la frase legítima de autenticación de gh, fuera del contrato CI. Se acotó a `For this repository, require`; no se modifica esa otra instrucción. La salida roja anterior queda registrada para trazabilidad.

Hora: 2026-10-23T10:24:35.951516-03:00; SHA: `13ce7f23fad84c951bf6f01995bab9b0ad3490ad`; HOME: `/tmp/active-stack-v2-home-55xjt12u`.

Comando: `go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.005s

exit=0
```

Hora: 2026-10-23T10:24:36.351599-03:00; SHA: `13ce7f23fad84c951bf6f01995bab9b0ad3490ad`; HOME: `/tmp/active-stack-v2-home-55xjt12u`.

Comando: `go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/assets	0.002s

exit=0
```

Hora: 2026-10-23T10:24:36.568438-03:00; SHA: `13ce7f23fad84c951bf6f01995bab9b0ad3490ad`; HOME: `/tmp/active-stack-v2-home-55xjt12u`.

Comando: `go test ./internal/harness/config -run ^TestCompose_AllToggles -update -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s

exit=0
```
## Ítem 5 — RF1-5a, eliminar autoridad paralela citada

La prohibición existente ahora inspecciona Compose(KnownToggles), no sólo senior.md. El test de ship deja de exigir la cita eliminada. Golden: se elimina la referencia al senior y se conserva la regla local de diagnóstico/consulta de ship; sin tocar senior.md (RF1-5b).

Hora: 2026-10-23T10:24:36.830580-03:00; SHA: `140f5d57ffdcece0a827411beaaa2fb37a108dba`; HOME: `/tmp/active-stack-v2-home-55xjt12u`.

Comando: `go test ./internal/harness/config -run ^TestAuditorDecideSinAutoridadParalela$ -count=1`

```text
--- FAIL: TestAuditorDecideSinAutoridadParalela (0.00s)
    verify_loop_test.go:60: autoridad vieja: the senior labels, you hold
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.008s
FAIL

exit=1
```

Hora: 2026-10-23T10:24:37.245720-03:00; SHA: `140f5d57ffdcece0a827411beaaa2fb37a108dba`; HOME: `/tmp/active-stack-v2-home-55xjt12u`.

Comando: `go test ./internal/harness/config -run TestAuditorDecideSinAutoridadParalela|TestCompose_ShipToggle_TheAnswerBelongsToTheOrchestrator -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.011s

exit=0
```

Hora: 2026-10-23T10:24:37.739424-03:00; SHA: `140f5d57ffdcece0a827411beaaa2fb37a108dba`; HOME: `/tmp/active-stack-v2-home-55xjt12u`.

Comando: `go test ./internal/harness/config -run ^TestCompose_AllToggles -update -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.006s

exit=0
```
## Complemento del ítem 1 — M3 en cada test nombrado

La primera ejecución conjunta dio rojo por TestVerifyContratoBucle. Se refuerza también TestAuditorDecideSinAutoridadParalela y se ejecutan por separado para demostrar que cada uno detecta M3 por su aserción.

Hora: 2026-10-23T10:25:03.983711-03:00; SHA: `61b23ea36f5ac59d68aba1774d65833e56a55d83`; HOME: `/tmp/active-stack-v2-home-4uq57qqy`.

Comando: `go test ./internal/harness/config -run ^TestVerifyContratoBucle$ -count=1`

```text
--- FAIL: TestVerifyContratoBucle (0.00s)
    verify_loop_test.go:13: secuencia contradictoria: launch the auditor in parallel with the senior.
    verify_loop_test.go:17: falta espera obligatoria; no se permite: Launch the auditor in parallel with the senior.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.005s
FAIL

exit=1
```

Hora: 2026-10-23T10:25:04.479388-03:00; SHA: `61b23ea36f5ac59d68aba1774d65833e56a55d83`; HOME: `/tmp/active-stack-v2-home-4uq57qqy`.

Comando: `go test ./internal/harness/config -run ^TestAuditorDecideSinAutoridadParalela$ -count=1`

```text
--- FAIL: TestAuditorDecideSinAutoridadParalela (0.01s)
    verify_loop_test.go:42: secuencia contradictoria: Launch the auditor in parallel with the senior.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.010s
FAIL

exit=1
```

Hora: 2026-10-23T10:25:04.736621-03:00; SHA: `61b23ea36f5ac59d68aba1774d65833e56a55d83`; HOME: `/tmp/active-stack-v2-home-4uq57qqy`.

Comando: `go test ./internal/harness/config -run TestVerifyContratoBucle|TestAuditorDecide -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.010s

exit=0
```
### Regresión adicional de RF1-3 con aserción acotada

Se reintroduce la frase ambigua únicamente en la regla de CI de cada copia. No se exige cambiar la referencia legítima al repositorio al autenticar gh.

Hora: 2026-10-23T10:25:05.183791-03:00; SHA: `61b23ea36f5ac59d68aba1774d65833e56a55d83`; HOME: `/tmp/active-stack-v2-home-4uq57qqy`.

Comando: `go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestShipEsperaCICompleto (0.00s)
    ship_ci_test.go:11: CI ambiguo: For this repository
    --- FAIL: TestShipEsperaCICompleto/For_active-stack,_require_.github/workflows/ci.yml_on_push,_with_test_on_Ubuntu_and_Windows. (0.00s)
        ship_ci_test.go:35: falta contrato de CI: For active-stack, require .github/workflows/ci.yml on push, with test on Ubuntu and Windows.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s
FAIL

exit=1
```

Hora: 2026-10-23T10:25:05.670932-03:00; SHA: `61b23ea36f5ac59d68aba1774d65833e56a55d83`; HOME: `/tmp/active-stack-v2-home-4uq57qqy`.

Comando: `go test ./internal/harness/config -run ^TestShipEsperaCICompleto$ -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/internal/harness/config	0.004s

exit=0
```

Hora: 2026-10-23T10:25:05.893045-03:00; SHA: `61b23ea36f5ac59d68aba1774d65833e56a55d83`; HOME: `/tmp/active-stack-v2-home-4uq57qqy`.

Comando: `go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1`

```text
--- FAIL: TestComandoShipEsperaCICompleto (0.00s)
    ship_ci_test.go:17: CI ambiguo: For this repository
    --- FAIL: TestComandoShipEsperaCICompleto/For_active-stack,_require_.github/workflows/ci.yml_on_push,_with_test_on_Ubuntu_and_Windows. (0.00s)
        ship_ci_test.go:41: falta contrato de CI: For active-stack, require .github/workflows/ci.yml on push, with test on Ubuntu and Windows.
FAIL
FAIL	github.com/Group-Active-IA/active-stack/assets	0.003s
FAIL

exit=1
```

Hora: 2026-10-23T10:25:06.106255-03:00; SHA: `61b23ea36f5ac59d68aba1774d65833e56a55d83`; HOME: `/tmp/active-stack-v2-home-4uq57qqy`.

Comando: `go test ./assets -run ^TestComandoShipEsperaCICompleto$ -count=1`

```text
ok  	github.com/Group-Active-IA/active-stack/assets	0.002s

exit=0
```
