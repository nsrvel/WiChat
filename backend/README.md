# Backend (Go workspace)

Monorepo **microservices** under one `go.work` — each module is self-contained (future split-friendly).

## Modules

| Directory | Role |
|-----------|------|
| `wi-shared/` | Config, errors, response, gRPC models, infra helpers |
| `ms-auth/` | Auth service |
| `ms-user/` | Users / workspaces (skeleton) |
| `api-gateway/` | Public HTTP `/api/v1` |

## Commands (from `backend/`)

```bash
make check              # all modules
make dev-ms-auth        # Air reload
make dev-api-gateway
```

Per module: `cd ms-auth && make dev` — copy `.env.example` → `.env` in **that** directory.

## Config

- **No** shared backend `.env` at repo root.
- Env templates: each module’s `.env.example`.
- Infra defaults: [deploy/README.md](../deploy/README.md).

Layout: [docs/architecture/monorepo-layout.md](../docs/architecture/monorepo-layout.md).
