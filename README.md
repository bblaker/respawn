# respawn

> A Kubernetes operator for game servers.

Declare a Minecraft or Factorio server as a Kubernetes custom resource and get
scheduling, persistence, backups, and **scale-to-zero**: the server sleeps when
nobody is online and wakes up when someone tries to join.

> [!WARNING]
> **Status: design phase.** Nothing here is ready to install yet. The API below
> is the *planned* `v1alpha1` shape and will change. See the
> [project plan](docs/PROJECT_PLAN.md) for the roadmap.

## Why

Running a game server for a handful of friends usually means a VM that is on
24/7 and mostly empty, backups that are a cron job somebody wrote once, and a
world that disappears when the disk dies. respawn treats the server as a
declarative resource instead:

- **Declarative.** One `GameServer` object describes the game, version,
  storage, and backup policy. The operator reconciles the rest.
- **Persistent.** World data lives on a PersistentVolumeClaim with an explicit
  reclaim policy, so deleting the server does not delete the world unless you
  ask it to.
- **Backed up.** Scheduled, encrypted, deduplicated backups to any
  S3-compatible bucket, plus one-off backups and restores via
  `GameServerBackup`.
- **Scale-to-zero.** Idle servers save, shut down, and free their CPU and RAM.
  A small "waker" holds the port while the server sleeps and starts it again
  when a player connects.

It is also a deliberately complete example of a production-style operator:
CRDs with CEL validation, an idempotent reconcile loop, owner references,
finalizers, status conditions, events, metrics, and envtest/e2e coverage.

## Planned usage

```yaml
apiVersion: respawn.bblaker.com/v1alpha1
kind: GameServer
metadata:
  name: survival
  namespace: games
spec:
  minecraft:
    eula: true            # you must accept the Minecraft EULA yourself
    version: "1.21.4"
    type: PAPER
    memory: 4Gi
  storage:
    size: 20Gi
    reclaimPolicy: Retain
  power:
    mode: Auto            # AlwaysOn | Auto | Off
    idleTimeout: 15m
  backup:
    schedule: "0 4 * * *"
    repositoryRef:
      name: b2-backups
    retention:
      keepDaily: 7
      keepWeekly: 4
  service:
    type: LoadBalancer
```

```console
$ kubectl get gameservers -n games
NAME       GAME        STATE      PLAYERS   ENDPOINT             AGE
survival   minecraft   Sleeping   0         192.168.1.40:25565   3d
```

## Supported games

| Game      | Image                                                                        | Player detection        | Wake trigger            |
| --------- | ---------------------------------------------------------------------------- | ----------------------- | ----------------------- |
| Minecraft | [`itzg/minecraft-server`](https://github.com/itzg/docker-minecraft-server)   | Server List Ping / RCON | Login handshake         |
| Factorio  | [`factoriotools/factorio`](https://github.com/factoriotools/factorio-docker) | RCON                    | Connection request (UDP) |

## Stack

Go · [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime)
· [Kubebuilder](https://book.kubebuilder.io/) · CRDs · [k3s](https://k3s.io/)
· [restic](https://restic.net/)

## Documentation

- [Project plan](docs/PROJECT_PLAN.md): architecture, API design, milestones,
  testing and security strategy
- [Contributing](CONTRIBUTING.md)
- [Security policy](SECURITY.md)
- [Code of conduct](CODE_OF_CONDUCT.md)
- [Changelog](CHANGELOG.md)

## License

[Apache License 2.0](LICENSE)

Minecraft is a trademark of Mojang AB. Factorio is a trademark of Wube
Software. respawn is not affiliated with or endorsed by either.
