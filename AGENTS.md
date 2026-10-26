# respawn: guide for coding agents

respawn is a Kubebuilder v4 operator. Read [CONTRIBUTING.md](CONTRIBUTING.md)
for workflows and code rules, and [docs/PROJECT_PLAN.md](docs/PROJECT_PLAN.md)
for the design. Decisions are recorded in [docs/adr/](docs/adr/README.md).

## Layout

```
cmd/main.go                 operator entrypoint (keep here; the Kubebuilder CLI edits it)
api/v1alpha1/*_types.go     CRD schemas and markers
api/v1alpha1/well_known.go  finalizer, annotation, label, and condition names
internal/controller/        reconcilers and envtest specs
test/e2e/                   e2e suite, run on k3d by `make test-e2e`
config/                     kustomize manifests; samples in config/samples/
```

## Rules

- Never hand-edit generated files: `api/**/zz_generated.*.go`,
  `config/crd/bases/*`, `config/rbac/role.yaml`, `PROJECT`. Run
  `make generate manifests` after changing types or markers.
- Do not delete `// +kubebuilder:scaffold:*` markers.
- Scaffold new APIs and webhooks with the `kubebuilder` CLI, not by hand.
- Images are built with ko (`make image`), not Docker. There is no Dockerfile.
- E2E tests need a dedicated k3d cluster; never point them at a real cluster.
- Never log or surface secrets, RCON passwords, or backup credentials.

## Before finishing a change

```sh
make generate manifests
make lint test
```
