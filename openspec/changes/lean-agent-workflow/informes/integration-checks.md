# Integrated checks

EVIDENCE — integration checker, 2026-10-29 11:00:14–11:00:54 -03. Base HEAD `61887a442ecf8abed1f5106915cf2d56e6bbb278`; tested working diff is uncommitted, not a delivered-commit certification. No build or installation performed.

- Single suite run: `GOCACHE=/tmp/active-stack-go-cache go test -json ./... -count=1` → exit 1. Independently counted JSON events: 2507 passing named tests/subtests, 1 failing named test (`internal/harness/external: TestDownloadBinary_EmptyProfileOS`). Package verdicts and exact failure are available in `/tmp/lean-checks.tests.json`; stderr at `/tmp/lean-checks.tests.stderr`. Counts include nested subtests and are not a claim of 2508 distinct top-level tests.
- `GOCACHE=/tmp/active-stack-go-cache go vet ./...` → exit 0, empty `/tmp/lean-checks.vet.log`.
- `gofmt -l .` → exit 0, empty `/tmp/lean-checks.gofmt.log`.
- `git diff --check` → empty `/tmp/lean-checks.diff.log`.
- `cmp AGENTS.md CLAUDE.md` → empty `/tmp/lean-checks.mirror.log` (identical).
- `openspec validate lean-agent-workflow --strict` → `Change 'lean-agent-workflow' is valid` in `/tmp/lean-checks.openspec.log`.

JUDGMENT — confidence high in the executed results; full suite is not green. The failing zone is external binary download, outside this change's implementation zones. No fix or rerun was attempted; this report does not classify the failure as pre-existing without a baseline reproduction.

ASSUMPTION / LIMITS — tests cannot certify runtime obedience, plugin/provider behavior, installation integrity in the user's dotfiles, Windows CI or token/cost savings. No commit/SHA change was created during this check.
