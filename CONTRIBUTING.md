# Contributing to WiChat

Thank you for helping make WiChat production-grade. This project is pre-implementation in places; when in doubt, follow [docs/architecture/wichat-system-design.md](docs/architecture/wichat-system-design.md).

## Before you start

1. Open an issue for non-trivial work (feature, breaking change, new dependency).
2. Keep PRs focused — one concern per change when possible.
3. AGPL-3.0 applies to contributions; you agree your contributions are licensed under the same terms.

## Development setup

Documented as each phase lands (Phase 0: `docker compose`, Go toolchain, Node for clients). Until then, architecture docs are the source of truth.

## Code standards

- **Go (core):** idiomatic Go, table-driven tests, `testify` where it helps readability; target ~90% unit coverage on core modules over time.
- **TypeScript (clients):** strict typing, shadcn/Radix patterns for UI; match existing folder layout.
- **API:** REST under `/api/v1/`; realtime via WebSocket; errors use standard HTTP status + human-readable message (see system design §13).
- **i18n:** backend emits message keys; user-facing strings live in translation files.
- **Security:** no secrets in git; rate limits at the gateway; validate at trust boundaries.

## Pull requests

- Describe **why**, not only what.
- Link issues (`Fixes #123`) when applicable.
- Note test evidence (unit, integration, or manual steps).
- Update docs if behavior or ops change.

## Commits

Use clear, imperative subjects (`Add workspace registration modes`). Squash or merge per maintainer preference on the PR.

## Plugin authors (in-house v1)

Plugins run as separate processes with declared permission manifests; see system design §6.8 and §6.2.
