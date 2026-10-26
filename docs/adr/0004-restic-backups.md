# 4. restic for backups

- **Status:** Accepted
- **Date:** 2026-10-26

## Context

Backups must be consistent, encrypted, space-efficient for large and slowly changing worlds, and must work on k3s's node-local `local-path` volumes.

## Decision

Back up with restic to S3-compatible storage, running in Jobs the operator creates.

## Consequences

- Client-side authenticated encryption, content-defined deduplication, integrity checks, and retention (`forget --prune`) come from one well-audited tool.
- Works with any storage class, including `local-path`, unlike CSI VolumeSnapshots.
- The restic password becomes a critical secret: losing it makes every backup unrecoverable. The docs must say so loudly.
- CSI snapshots remain a possible later alternative for storage classes that support them.
