# 1. One GameServer kind with a per-game union block

- **Status:** Accepted
- **Date:** 2026-10-26

## Context

respawn supports several games (Minecraft and Factorio first). Each could get its own kind (`MinecraftServer`, `FactorioServer`), or there could be one `GameServer` kind with a block per game.

## Decision

Use one `GameServer` kind. Game-specific settings live in a union of optional blocks (`spec.minecraft`, `spec.factorio`), and exactly one must be set. Game behavior sits behind the `internal/game` adapter interface.

## Consequences

- One controller, one set of status conditions, one printer-column layout, and one set of docs and tooling.
- Adding a game means one adapter and one spec block, not a new CRD and controller.
- The "exactly one block" and "game is immutable" rules have to be enforced with CEL (M2), because OpenAPI cannot express them.
- Per-game fields are namespaced under their block, which makes `kubectl explain` output slightly deeper.
