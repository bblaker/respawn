# 11. API group respawn.bblaker.com

- **Status:** Accepted
- **Date:** 2026-10-26

## Context

CRDs need an API group based on a DNS name the project controls. Changing it later breaks every existing resource.

## Decision

Use the group `respawn.bblaker.com`, starting at version `v1alpha1`.

## Consequences

- Conventional: the domain is controlled by the author.
- The group is fixed from M2, when the first real resources ship.
- If the project ever moves to a neutral home, a new group would be a breaking change handled like an API version migration.
