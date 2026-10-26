# 14. Docs as Markdown in the repository, published on Cloudflare Pages at M7

- **Status:** Accepted
- **Date:** 2026-10-26

## Context

Docs could stay as repository Markdown only, or be published as a site on GitHub Pages or Cloudflare Pages.

## Decision

Keep Markdown in `docs/` as the source of truth. At M7, build it with MkDocs Material and publish it to Cloudflare Pages from a `docs.yml` workflow.

## Consequences

- Cloudflare Pages gives a preview deploy for every pull request, which GitHub Pages lacks.
- The deploy token is scoped to Pages, stored in a protected `docs` environment, and never exposed to workflows from forks.
- The API reference is generated from the CRDs during the build, so it cannot drift.
