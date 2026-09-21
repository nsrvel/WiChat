# AGENTS.md — WiChat

Guidance for humans and coding agents working in this repository.

## Product

WiChat is a **self-hosted**, **AGPL-3.0** team communication platform with a **plugin architecture** — a **major product** built for real daily use. Operability, security, and deployability are part of the application, not optional polish.

**Canonical spec:** [docs/architecture/wichat-system-design.md](docs/architecture/wichat-system-design.md)

## Architecture constraints (do not drift)

- Core backend: **Go monorepo microservices** under `services/` (gRPC + Kafka between services; not one monolith binary).
- Separate processes: **LiveKit**, **plugins (gRPC)**, not core domains.
- Data: PostgreSQL, Redis, Elasticsearch, MinIO, Kafka.
- Clients: Next.js (web + iOS PWA), Electron, React Native (Android).
- **Non-goals:** multi-tenant public SaaS, federation, E2E encryption v1, public plugin marketplace v1.

## Repository layout (target)

```
apps/          # web, desktop, mobile clients
services/      # Go microservices (gateway, identity, platform, chat, …)
pkg/           # protobuf + shared types only
plugins/       # optional plugin services + bundles
deploy/        # docker compose, env templates
docs/          # architecture & ops (system design lives here)
```

Phases in system design §15 — implement in order unless explicitly reprioritized.

## Quality bar (production-ready)

- Security: TLS, rate limits, no secrets in repo, bcrypt/argon2, guest/session scoping per spec.
- Tests: Go unit (~90% goal on core), testcontainers for integration, Playwright for critical E2E paths.
- Observability: Grafana + Loki (when services exist).
- API: REST `/api/v1/`, WebSocket for realtime; UTC timestamps; i18n keys from backend.
- Keep core **light** — heavy features belong in plugins.

## Legal

License: [LICENSE](LICENSE). Preserve AGPL headers on new source files. Do not add proprietary-licensed dependencies without explicit approval.

## Cursor rules

Only **always-on** rules for now (keep context lean). Add file-specific `.mdc` rules when that code exists (Go, web, deploy, etc.).

| Rule | Scope |
|------|--------|
| `wichat-program.mdc` | How we work: `release`, micro-steps, vertical slices |
| `wichat-core.mdc` | Product + architecture guardrails |

Detail standards: [system design](docs/architecture/wichat-system-design.md) and this file until narrower rules return.
