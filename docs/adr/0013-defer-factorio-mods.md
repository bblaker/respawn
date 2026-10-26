# 13. Factorio mods deferred past v0.1

- **Status:** Accepted
- **Date:** 2026-10-26

## Context

Downloading Factorio mods needs a factorio.com account token, which is another secret to handle and another failure mode.

## Decision

M6 ships a plain Factorio server only. Mod support is listed under "Later".

## Consequences

- M6 stays focused on proving the game adapter interface with a second game.
- Users who need mods before then can bake them into a custom image via `spec.image`.
