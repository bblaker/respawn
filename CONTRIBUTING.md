# Contributing to respawn

Thanks for your interest. respawn is early, so the most useful contributions
right now are design feedback on the [project plan](docs/PROJECT_PLAN.md),
bug reports, and small focused pull requests.

By participating you agree to follow the [code of conduct](CODE_OF_CONDUCT.md).
For security issues, follow [SECURITY.md](SECURITY.md) instead of opening an
issue.

## Before you start

- **Bugs:** open an issue using the bug report template.
- **Small fixes** (typos, docs, obvious bugs): open a pull request directly.
- **Features, API changes, new games:** open an issue first so we can agree on
  the design before you write code. Changes to the CRD API are hard to undo.

## Development environment

You will need:

| Tool                                   | Used for                           |
| -------------------------------------- | ---------------------------------- |
| Go (version in `go.mod`)               | building and testing               |
| Docker                                 | the local k3d cluster and images   |
| kubectl                                | talking to the cluster             |
| make                                   | all workflows below                |

Everything else (controller-gen, envtest binaries, golangci-lint, kustomize,
ko, and k3d) is installed into `./bin` by the Makefile at pinned versions, so
you do not need to install it globally.

### Common tasks

```sh
make help            # list all targets
make generate        # regenerate deepcopy code
make manifests       # regenerate CRDs and RBAC from markers
make verify          # fail if generated files or go.mod are out of date
make lint            # golangci-lint
make test            # unit + envtest integration tests, with -race
make image           # build the operator image into local Docker (ko.local/respawn:dev)
make dev-cluster     # create the local k3d cluster (respawn-dev)
make dev-deploy      # build, import, and deploy the operator into respawn-dev
make run             # run the operator locally against the current kubeconfig
make test-e2e        # end-to-end tests on a throwaway k3d cluster (respawn-e2e)
```

`make dev-cluster` and `make test-e2e` add their cluster to your kubeconfig
and switch the current context to it, as k3d does by default. `make
test-e2e` deletes its cluster when the tests pass.

## Making changes

1. Fork the repository and create a branch from `main`.
2. Make your change, with tests. See the
   [testing strategy](docs/PROJECT_PLAN.md#6-testing-strategy) for what kind of
   test belongs where.
3. If you changed anything under `api/`, run `make generate manifests` and
   commit the generated files. CI fails if they are out of date.
4. Run `make lint test` locally.
5. Update docs and add an entry under `Unreleased` in
   [CHANGELOG.md](CHANGELOG.md) for user-visible changes.
6. Open a pull request and fill in the template.

### Commit messages

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(controller): scale StatefulSet to zero after idle timeout
fix(waker): reject handshakes with oversized packet length
docs: describe backup retention fields
```

Common types: `feat`, `fix`, `docs`, `test`, `refactor`, `perf`, `build`,
`ci`, `chore`. Mark breaking changes with `!` (`feat(api)!: ...`) and explain
them in the commit body.

### Sign-off

Sign off your commits to certify the
[Developer Certificate of Origin](https://developercertificate.org/):

```sh
git commit -s
```

## Code guidelines

- Follow standard Go style; `gofmt` and `golangci-lint` are authoritative.
- Reconcilers must be **idempotent**: running reconcile twice with the same
  state must produce the same result and no extra API writes.
- Do not block in reconcile. Use `RequeueAfter` instead of sleeping.
- Never log secrets, RCON passwords, or backup credentials, and never put them
  in status, events, or error messages.
- Keep game-specific logic behind the game adapter interface in
  `internal/game/`. The controller should not know what Minecraft is.
- Code that parses network input (the waker, RCON, ping clients) needs fuzz
  tests.
- Exported identifiers and every CRD field need doc comments. Those comments
  become the API reference and `kubectl explain` output.

## Review

Pull requests need a passing CI run and one approving review from a
maintainer. Expect feedback; it is about the code, not you. Maintainers may
squash commits on merge.

## Licensing

By contributing you agree that your contributions are licensed under the
[Apache License 2.0](LICENSE).
