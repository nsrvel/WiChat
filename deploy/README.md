# Deployment

## Local infrastructure (all services)

From repo root:

```bash
docker compose -f deploy/docker-compose.yml up -d
```

Stop:

```bash
docker compose -f deploy/docker-compose.yml down
```

## Per-component compose (infra / partial install)

Files live in [`deploy/compose/`](compose/). Config trees: `compose/prometheus/`, `compose/grafana/`. Use project name **`wichat`** when starting stacks separately.

From repo root:

```bash
# Data stores
docker compose -p wichat -f deploy/compose/postgres.yml up -d
docker compose -p wichat -f deploy/compose/redis.yml up -d

# Observability (both files, or one Kuma service per file — start Prometheus first)
docker compose -p wichat -f deploy/compose/prometheus.yml -f deploy/compose/grafana.yml up -d
```

Kuma (Shell), working directory = repo root:

```bash
# Prometheus
docker compose -p wichat -f deploy/compose/prometheus.yml up

# Grafana (after Prometheus is up, same -p wichat)
docker compose -p wichat -f deploy/compose/grafana.yml up
```

Stop one component (example):

```bash
docker compose -p wichat -f deploy/compose/redis.yml down
docker compose -p wichat -f deploy/compose/prometheus.yml down
docker compose -p wichat -f deploy/compose/grafana.yml down
```

Postgres data volume: `wichat_wichat-postgres` (Compose prefixes volume names with the project name).

### Local ports (defaults)

| Service | Host port | Notes |
|---------|-----------|--------|
| Postgres | **5432** | `postgres://wichat:wichat@localhost:5432/wichat` |
| Redis | **6379** | `redis://localhost:6379/0` |
| Prometheus | **9090** | UI + API |
| Grafana | **9000** | Host maps to container **3000** |
| api-gateway | **3000** | `/health`, `/ready`, `/metrics` on public HTTP |
| ms-auth | **3001** | gRPC + `/health`, `/ready`, `/metrics` on one port |
| ms-user | **3002** | gRPC + `/health`, `/ready`, `/metrics` on one port |

### Metrics (Prometheus + Grafana)

- Prometheus: http://localhost:9090
- Grafana: http://localhost:9000

Scrape config expects **ms-auth** on **3001**, **ms-user** on **3002**, and **api-gateway** on **3000** (`host.docker.internal`). Run binaries locally, then check Prometheus → Status → Targets.

See [docs/architecture/metrics.md](../docs/architecture/metrics.md).

Instance env: [`deploy/.env.example`](.env.example) (Compose overrides). Go services: each module’s `backend/<service>/.env.example`. See [docs/architecture/monorepo-layout.md](../docs/architecture/monorepo-layout.md).

See `docs/architecture/wichat-system-design.md` §7 and §11.

### Go services (local, not Docker)

From `backend/`: `make dev-ms-auth` then `make dev-api-gateway` (Air live reload), or `make run-*` without reload. Kuma Shell: working directory `backend/ms-auth` or `backend/api-gateway`, command `make dev`.
