# 2. StatefulSet with replicas 0 or 1 and an operator-managed PVC

- **Status:** Accepted
- **Date:** 2026-10-26

## Context

A game world must never be written by two processes at once. The workload also has to scale to zero without losing data, and the operator needs full control over what happens to the volume on deletion.

## Decision

Run each server as a StatefulSet scaled between 0 and 1. The PersistentVolumeClaim is created and owned by the operator, not by `volumeClaimTemplates`.

## Consequences

- StatefulSet guarantees at most one pod per ordinal, even during rollouts and node failures. A Deployment with `Recreate` gives weaker guarantees.
- The operator decides the PVC's fate through `storage.reclaimPolicy`: it removes the owner reference for `Retain` and deletes the claim for `Delete`.
- Restoring into a fresh PVC before the first start (`restoreFrom`) is straightforward, because the operator creates the claim itself.
