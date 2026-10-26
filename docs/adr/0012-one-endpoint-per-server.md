# 12. One endpoint per server; no hostname routing in v0.1

- **Status:** Accepted
- **Date:** 2026-10-26

## Context

Several Minecraft servers could share one IP and port with hostname-based routing (for example mc-router), or each server could get its own Service.

## Decision

Each GameServer gets its own Service. Hostname routing is deferred until after v0.1.

## Consequences

- The waker and selector flip (ADR 0005) stay simple.
- Works for Factorio too, whose UDP protocol carries no hostname.
- With k3s ServiceLB, servers share the node IPs, so each needs a distinct port (`service.port`). Distinct IPs need MetalLB or kube-vip. The networking guide (M7) covers both.
