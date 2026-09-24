# Deployment

## Local infrastructure (Postgres + Redis + observability)

From repo root:

```bash
docker compose -f deploy/docker-compose.yml up -d
```

Stop:

```bash
docker compose -f deploy/docker-compose.yml down
```

### Metrics (Prometheus + Grafana)

- Prometheus: http://localhost:9091 (UI and API)
- Grafana: http://localhost:3000

Scrape config expects **ms-auth** admin metrics on **9090** and **gateway** on **8080** (`host.docker.internal`). Run both binaries locally, then check Prometheus → Status → Targets.

See [docs/architecture/metrics.md](../docs/architecture/metrics.md).

Instance env vars: see `.env.example` at repo root. WiChat services are added to compose as each slice lands.

See `docs/architecture/wichat-system-design.md` §7 and §11.
