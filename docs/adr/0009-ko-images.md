# 9. ko for container images

- **Status:** Accepted
- **Date:** 2026-10-26

## Context

The operator and waker are pure Go binaries. Images must be small, non-root, multi-arch (amd64 and arm64), and reproducible, with SBOMs.

## Decision

Build images with ko on Chainguard's static base image (`cgr.dev/chainguard/static`), pinned by digest in `.ko.yaml`. There is no Dockerfile.

## Consequences

- Fast multi-arch builds with no Docker build step, and SBOMs built in.
- Chainguard was chosen over Google's distroless (gcr.io) because ko's built-in Google keychain sends the local `gcloud` login to gcr.io when one is installed, and gcr.io rejects the request if that login is stale. Chainguard avoids that.
- Chainguard's free tier only publishes the `latest` tag, so the digest pin must be bumped deliberately; Dependabot does not track it.
- Kubebuilder's `docker-build` targets are replaced by `make image` and `make image-all`.
