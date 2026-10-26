# 8. Server-side apply for owned objects

- **Status:** Accepted
- **Date:** 2026-10-26

## Context

The operator creates and updates StatefulSets, Services, Secrets, NetworkPolicies, and more. It could use get, mutate, and update, or server-side apply.

## Decision

Apply owned objects with server-side apply using the field manager `respawn`.

## Consequences

- Field ownership is explicit, so other controllers and users can own other fields without fighting.
- Idempotency is simpler: applying the same desired object twice makes no change.
- Desired objects must be built completely on every reconcile, which suits the pure builders in `internal/resources`.
