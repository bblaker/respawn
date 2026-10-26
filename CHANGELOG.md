# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
While the version is `0.x`, minor releases may contain breaking changes; they
are always called out under **Changed** or **Removed**.

## [Unreleased]

### Added

- Kubebuilder v4 scaffold with `GameServer`, `GameServerBackup`, and
  `BackupRepository` v1alpha1 types (validation and defaults; reconcilers are
  still no-ops).
- Makefile with pinned tools, ko image builds (amd64 and arm64), k3d
  development and e2e clusters, and a `verify` target for generated files.
- envtest specs for CRD defaults and validation; e2e suite on k3d.
- CI workflows: lint, verify, tests with coverage, multi-arch build, e2e,
  CodeQL, govulncheck, dependency review, and OpenSSF Scorecard.
- Architecture decision records in `docs/adr/`.
- Repository health files: README, license, contributing guide, security
  policy, code of conduct, issue and pull request templates.
- Project plan covering architecture, API design, milestones, testing, and
  security.
