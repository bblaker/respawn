# 6. The operator polls player counts directly

- **Status:** Accepted
- **Date:** 2026-10-26

## Context

Scale-to-zero needs the player count. Options: the operator queries the game (Server List Ping, RCON), a sidecar exporter reports it, or the operator parses game logs.

## Decision

The operator polls each running server over the pod network, every `power.pollInterval` (default 30s), through the game adapter.

## Consequences

- No changes to upstream game images and no sidecars.
- All power-state decisions live in one place, the reconcile loop.
- The operator needs pod-network access, so `make run` from a laptop cannot poll without port-forwarding.
- Status is written only when the count changes, to keep API writes low.
