# 7. The operator creates scheduled backups instead of using CronJobs

- **Status:** Accepted
- **Date:** 2026-10-26

## Context

Backups of a running server must pause world saves (quiesce), run, then resume, with crash recovery in between. A Kubernetes CronJob can start a Job on schedule but cannot coordinate with the live server.

## Decision

The GameServer reconciler parses `backup.schedule` and creates a `GameServerBackup` named after the scheduled time. The GameServerBackup reconciler does the quiesce, Job, and resume sequence.

## Consequences

- Deterministic names make a retried reconcile a no-op instead of a duplicate backup.
- Missed schedules (operator down) can be handled explicitly: at most one catch-up backup.
- The operator owns cron parsing and time-based requeues, tested with an injected clock.
