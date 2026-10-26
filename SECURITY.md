# Security policy

respawn runs with Kubernetes API permissions and exposes game ports to the
internet, so security reports are taken seriously and handled ahead of feature
work.

## Supported versions

respawn is pre-1.0. Only the latest released minor version receives security
fixes.

| Version         | Supported |
| --------------- | --------- |
| latest `0.x`    | Yes       |
| older `0.x`     | No        |
| unreleased main | Best effort |

## Reporting a vulnerability

**Please do not open a public issue, pull request, or discussion for security
problems.**

Report privately through GitHub's
[private vulnerability reporting](https://github.com/bblaker/respawn/security/advisories/new).
Include:

- the affected component (operator, waker, backup job, Helm chart, manifests)
  and version or commit
- your environment (Kubernetes distribution and version, game, configuration)
- steps to reproduce or a proof of concept
- the impact you believe it has

### What to expect

| Step                                  | Target       |
| ------------------------------------- | ------------ |
| Acknowledgement                       | 3 days       |
| Initial assessment and severity       | 7 days       |
| Fix or mitigation for High/Critical   | 30 days      |
| Fix or mitigation for Medium/Low      | next release |

This is a volunteer project, so these are targets rather than guarantees. You
will be kept informed as the report progresses, and credited in the advisory
unless you would rather not be.

Disclosure is coordinated: a GitHub Security Advisory (with a CVE where
appropriate) is published when the fix is released. Please give us a
reasonable window before disclosing publicly; 90 days is the default.

## Scope

In scope:

- the operator (`cmd/manager`) and its RBAC
- the waker (`cmd/waker`), which parses untrusted network traffic
- backup and restore jobs and their handling of credentials
- generated manifests and the Helm chart, including default security contexts
  and NetworkPolicies
- release artifacts: container images, signatures, SBOMs, provenance

Out of scope:

- vulnerabilities in the upstream game servers or their container images
  (`itzg/minecraft-server`, `factoriotools/factorio`); report those upstream,
  though we would like to hear about them if respawn's defaults make them worse
- denial of service by a player who is legitimately allowed to join the server
- issues that require cluster-admin or equivalent access to exploit
- missing hardening with no demonstrated impact

## Security design

The security model and threat analysis are part of the
[project plan](docs/PROJECT_PLAN.md#7-security). In short:

- the operator and waker run as non-root, distroless containers with a
  read-only root filesystem and all capabilities dropped
- RBAC is least-privilege; each waker can only patch its own `GameServer`
- RCON is never exposed outside the cluster and is restricted by
  NetworkPolicy
- secrets are referenced, never inlined in custom resources, and never logged
- backups are encrypted client-side with restic before they leave the cluster
- release images are signed with cosign and ship with an SBOM and SLSA
  provenance
