# Architecture decision records

Short records of decisions that shaped respawn, in the order they were made. They expand on the decision log in [the project plan](../PROJECT_PLAN.md#12-decision-log).

To propose a change, add a new record that supersedes the old one instead of editing it, and update the status of both.

| # | Decision | Status |
| --- | --- | --- |
| [1](0001-one-gameserver-kind.md) | One GameServer kind with a per-game union block | Accepted |
| [2](0002-statefulset-and-managed-pvc.md) | StatefulSet with replicas 0 or 1 and an operator-managed PVC | Accepted |
| [3](0003-cel-not-webhooks.md) | CEL validation instead of admission webhooks in v1alpha1 | Accepted |
| [4](0004-restic-backups.md) | restic for backups | Accepted |
| [5](0005-per-server-waker.md) | Per-server waker and Service selector flip | Accepted |
| [6](0006-operator-polls-players.md) | The operator polls player counts directly | Accepted |
| [7](0007-operator-creates-scheduled-backups.md) | The operator creates scheduled backups instead of using CronJobs | Accepted |
| [8](0008-server-side-apply.md) | Server-side apply for owned objects | Accepted |
| [9](0009-ko-images.md) | ko for container images | Accepted |
| [10](0010-apache-2-license.md) | Apache-2.0 license | Accepted |
| [11](0011-api-group.md) | API group respawn.bblaker.com | Accepted |
| [12](0012-one-endpoint-per-server.md) | One endpoint per server; no hostname routing in v0.1 | Accepted |
| [13](0013-defer-factorio-mods.md) | Factorio mods deferred past v0.1 | Accepted |
| [14](0014-docs-on-cloudflare-pages.md) | Docs as Markdown in the repository, published on Cloudflare Pages at M7 | Accepted |
