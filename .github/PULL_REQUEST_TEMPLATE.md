## What does this change?

<!-- One or two sentences. Link the issue it resolves: "Fixes #123". -->

## Why?

<!-- The problem or use case. Skip for trivial changes. -->

## How was it tested?

<!-- Unit, envtest, e2e, manual against k3d? Paste commands or output where useful. -->

## Checklist

- [ ] Tests added or updated
- [ ] `make generate manifests` run and generated files committed (if `api/` changed)
- [ ] `make lint test` passes locally
- [ ] Docs and `CHANGELOG.md` updated for user-visible changes
- [ ] No secrets, credentials, or RCON passwords in logs, events, status, or error messages
- [ ] Breaking API changes are called out above and marked with `!` in the commit
