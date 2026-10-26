# 3. CEL validation instead of admission webhooks in v1alpha1

- **Status:** Accepted
- **Date:** 2026-10-26

## Context

Cross-field and transition rules (one game block, immutable fields, size only grows) need more than OpenAPI schema validation. The options are validating and defaulting webhooks, or CEL rules in the CRD (`x-kubernetes-validations`).

## Decision

Use CEL rules and OpenAPI defaults. Ship no admission webhooks in v1alpha1.

## Consequences

- No cert-manager dependency and no webhook TLS to manage, which matters on small k3s clusters.
- No availability risk: an unavailable webhook can block every write to the resource.
- Rules that need to look at other objects (for example, "does this Secret exist?") cannot be CEL. They become status conditions instead.
- A conversion webhook will still be needed when v1beta1 arrives (post-0.1).
