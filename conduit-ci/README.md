# Pilot: yaegi fork-ci via Conduit (1 week)

Goal: prove Conduit replaces `~/actions-runner-yaegi` (deleted 2026-10-26;
its `fork-ci` gate queued forever with no `[self-hosted, windows]` runner
left). Proven 2026-10-26: definition below ran 5/5 green through the
Conduit MCP envelope (`done`, ledger with real outputs).

## What changed in this PR

- `conduit-ci/fork-ci.yaml`: fork-ci steps as a Conduit generic
  definition (host toolchain builds go.mod 1.21 clean).
- `.github/workflows/fork-ci.yml` → `.disabled`: no self-hosted runs
  during the pilot (reversible with `git revert`).
- `main.yml`/`release.yml`/`go-cross.yml` untouched (hosted already).

## How to fire (owner/agent, per push/PR)

`conduit_run` with `kind: "generic"` and this file as `definition`;
poll `conduit_get`, stream `conduit_logs`.

## Success criteria

1 green week of agent-fired Conduit runs → nothing to decommission
(runner already removed).

## Rollback

`git revert` this PR: gate re-enables (needs a `[self-hosted, windows]`
runner to drain).
