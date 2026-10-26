# 5. Per-server waker and Service selector flip

- **Status:** Accepted
- **Date:** 2026-10-26

## Context

While a server sleeps, something must hold its port, answer server-list pings, and notice a real join attempt. That could be a shared gateway proxy that is always in the data path, or a small component per server that only takes traffic while the server sleeps.

## Decision

Run one small waker Deployment per GameServer. The operator flips the game Service's selector between the game pod (awake) and the waker pod (asleep).

## Consequences

- Awake servers have no extra network hop, and a waker bug cannot affect a running server.
- Each waker gets a narrowly scoped Role (`resourceNames` on its own GameServer).
- Costs one extra small pod (about 10 MiB) per server.
- Traffic switches as fast as Endpoints propagate, typically a few seconds.
