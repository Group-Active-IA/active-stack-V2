AR-1: **CERRADO**. JUICIO del arquitecto, confianza alta.
EVIDENCIA, 2026-10-23 11:34:47 -0300, rama `feat/metodologia-fase-2`, commit `000edd5`.
`git check-ignore -v openspec/changes/metodologia-fase-2/informes/example.md`: sin salida, exit **1**.
`git check-ignore --no-index -v openspec/changes/metodologia-fase-2/informes/example.md`: sin salida, exit **1**.
`git check-ignore --no-index -v openspec/specs/example.md openspec/changes/metodologia-fase-2/proposal.md openspec/changes/otro-change/tasks.md openspec/config.yaml`: exit **0**.
Salida: los cuatro paths siguen ignorados por `.gitignore:36:/openspec/*` o `.gitignore:40:/openspec/changes/*/*`; la excepción queda acotada a `informes/`.