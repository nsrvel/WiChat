# Metrics (Prometheus)

WiChat **microservices** expose ops HTTP (`wi-shared/infra/ops`): **ms-auth** shares `/health`, `/ready`, `/metrics` with gRPC on `grpc_addr` (`:3001`). **api-gateway** serves the same routes on public `http_addr`.

## Stack

| Component | Role |
|-----------|------|
| App (`wi-shared/infra/metrics`) | `/metrics` handler, RED via `NewHTTPMetrics` + `Middleware` |
| Prometheus | Scrapes targets every 15s |
| Grafana | Dashboards and exploration (datasource provisioned in `deploy/`) |

Logs use [logging.md](logging.md) (slog → Loki later); metrics are a separate pipeline.

## Ops endpoints (shared)

| Path | Purpose |
|------|---------|
| `GET /health` | **Liveness** — process up; always `{"status":"ok"}` |
| `GET /ready` | **Readiness** — can serve traffic; `503` + `{"status":"not_ready"}` when checks fail |
| `GET /metrics` | Prometheus text exposition |

**api-gateway** (`:3000`): `/health` is liveness-only; `/ready` pings ms-auth (short timeout).

**ms-auth** (`:3001`): gRPC and ops HTTP on one listener (`ops.NewCombinedServer`); `/ready` reflects `grpc.health.v1`.

In production, bind `grpc_addr` to loopback or an internal interface and firewall the port. `ops.NewHTTPServer` is for a separate ops-only listener when a service needs it.

Use `ops.RegisterWithReady(mux, reg, checks...)` on a mux, or `ops.NewCombinedServer` to share a port with gRPC.

## RED HTTP metrics (api-gateway and HTTP services)

Create metrics once at startup:

```go
httpMetrics, err := metrics.NewHTTPMetrics(metrics.HTTPOptions{Service: "api-gateway", Registry: reg})
api := httpMetrics.Middleware(apiMux)
```

`WrapHTTP` is a one-shot helper (registers collectors); do not call it twice on the same registry.

Records:

- `wichat_http_requests_total` — labels: `service`, `method`, `route`, `status_class`
- `wichat_http_request_duration_seconds` — labels: `service`, `method`, `route`

`route` uses normalized paths (`/api/v1/users/{id}`) to limit cardinality.

## ms-auth

- `metrics.NewRegistry()` in `main`
- Ops HTTP: [wi-shared/infra/ops/ops.go](../../backend/wi-shared/infra/ops/ops.go)

## api-gateway

- `metrics.NewRegistry()` + `NewHTTPMetrics` middleware on the API mux
- `/health` and `/metrics` on `http_addr` (default `:3000`)
- Composition: [api-gateway/internal/server](../../backend/api-gateway/internal/server)

## Local observability

1. Start infra: `docker compose -f deploy/docker-compose.yml up -d`
2. Run ms-auth and api-gateway on the host (`:3001` ms-auth, `:3000` api-gateway)
3. Prometheus UI: http://localhost:9090/targets — jobs `ms-auth` and `api-gateway` should be **UP**
4. Grafana: http://localhost:9000 (login `admin` / `admin` on first setup) — Prometheus datasource pre-provisioned

On Linux, Prometheus service includes `extra_hosts` for `host.docker.internal`.

## Label policy

Allowed: `service`, `method`, `route`, `status_class`, gRPC labels (future).

Forbidden: `user_id`, `workspace_id`, raw URLs with unbounded segments.

## Follow-up

- gRPC server metrics (ms-auth unary interceptor)
- Alertmanager rules (5xx rate, target down)
- Loki in compose (logging)
