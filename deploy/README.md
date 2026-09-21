# Deployment

## Local infrastructure (Postgres + Redis)

From repo root:

```bash
docker compose -f deploy/docker-compose.yml up -d
```

Stop:

```bash
docker compose -f deploy/docker-compose.yml down
```

Instance env vars: see `.env.example` at repo root. WiChat services are added to compose as each slice lands.

See `docs/architecture/wichat-system-design.md` §7 and §11.
