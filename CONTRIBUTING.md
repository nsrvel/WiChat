# Contributing to WiChat

Thank you for helping make WiChat production-grade. This project is pre-implementation in places; when in doubt, follow [docs/architecture/wichat-system-design.md](docs/architecture/wichat-system-design.md).

## Branches

- **`release`** — active development (default branch; open PRs here).
- **`main`** — stable line; updated from `release` when we cut a shipped version (not day-to-day commits).

## Before you start

1. Open an issue for non-trivial work (feature, breaking change, new dependency).
2. Keep PRs focused — one concern per change when possible.
3. AGPL-3.0 applies to contributions; you agree your contributions are licensed under the same terms.

## Development setup

Documented as each phase lands (Phase 0: `docker compose`, Go toolchain, Node for clients). Until then, architecture docs are the source of truth.

## Code standards

Follow `.cursor/rules/` (same standards for humans and agents). Highlights:

- **Go (core):** see `go-core.mdc` — ~90% unit coverage goal on domain packages.
- **TypeScript (clients):** see `typescript-clients.mdc`.
- **API / WS:** see `api-realtime.mdc` and system design §13.
- **Security:** see `security-trust.mdc`.

## Pull requests

- Describe **why**, not only what.
- Link issues (`Fixes #123`) when applicable.
- Note test evidence (unit, integration, or manual steps).
- Update docs if behavior or ops change.

## Commits

Use clear, imperative subjects (`Add workspace registration modes`). Squash or merge per maintainer preference on the PR.

## Plugin authors (in-house v1)

Plugins run as separate processes with declared permission manifests; see system design §6.8 and §6.2.
