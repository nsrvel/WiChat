# Metrics (Prometheus)

WiChat **microservices** expose metrics on a dedicated **ops HTTP port** (`wi-shared/infra/ops`, ms-auth `metrics_addr`). The **gateway** serves the same routes on public HTTP (`http_addr`).

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
| `GET /health` | Liveness JSON `{"status":"ok"}` |
| `GET /metrics` | Prometheus text exposition |

Configure per service, e.g. ms-auth `metrics_addr` (default `:9090`). In production prefer loopback or an internal bind (e.g. `127.0.0.1:9090`) and firewall the port. `ops.NewHTTPServer` sets read/write timeouts.

Use `ops.Register(mux, reg)` to mount `/health` and `/metrics` on an existing mux (gateway), or `ops.NewHTTPServer(addr, reg)` for a separate listener (gRPC microservices).

## RED HTTP metrics (gateway and HTTP services)

Create metrics once at startup:

```go
httpMetrics, err := metrics.NewHTTPMetrics(metrics.HTTPOptions{Service: "gateway", Registry: reg})
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

## gateway

- `metrics.NewRegistry()` + `NewHTTPMetrics` middleware on the API mux
- `/health` and `/metrics` on `http_addr` (default `:8080`)
- Composition: [gateway/internal/server](../../backend/gateway/internal/server)

## Local observability

1. Start infra: `docker compose -f deploy/docker-compose.yml up -d`
2. Run ms-auth and gateway on the host (`:9090` admin, `:8080` gateway)
3. Prometheus UI: http://localhost:9091/targets — jobs `ms-auth` and `gateway` should be **UP**
4. Grafana: http://localhost:3000 (default login `admin` / `admin` on first setup) — Prometheus datasource pre-provisioned

On Linux, Prometheus service includes `extra_hosts` for `host.docker.internal`.

## Label policy

Allowed: `service`, `method`, `route`, `status_class`, gRPC labels (future).

Forbidden: `user_id`, `workspace_id`, raw URLs with unbounded segments.

## Follow-up

- gRPC server metrics (ms-auth unary interceptor)
- Alertmanager rules (5xx rate, target down)
- Loki in compose (logging)
