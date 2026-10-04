# Service-to-service (`wi-shared/microservices`)

The package name is **`microservices`**, not `rpc`, because WiChat will call **other runtimes and transports** (Python/JS plugins, HTTP APIs) — not only Go gRPC. v1 ships **gRPC** under `microservices/grpc`; a future `microservices/http` can share correlation and retry patterns for outbound HTTP.

Contracts for Go gRPC live in `wi-shared/models`. Cross-cutting **request ID** (`RequestIDHeader`, `WithRequestID`) stays at the `microservices` root for api-gateway HTTP and gRPC metadata.

Tracing bootstrap is **`wi-shared/infra/otel`** — see [tracing.md](tracing.md).

## gRPC dial and server

```go
conn, err := msgrpc.Dial("localhost:3001", msgrpc.DefaultDialOptions, otel.GRPCClientDialOptions()...)
grpcSrv := msgrpc.NewServer(appLog, otel.GRPCServerOptions()...)
```

- Outbound: `x-request-id` metadata from context; client **keepalive** (30s) on all dials.
- Inbound: read `x-request-id` into context, panic recovery, matching server keepalive.

### Health

| Check | api-gateway | gRPC microservices (e.g. ms-auth) |
|-------|-------------|-----------------------------------|
| Liveness | `GET /health` on `PORT` | `GET /health` on `ms_*_host` listen addr (same port as gRPC) |
| Readiness | `GET /ready` (deps, e.g. HTTP `GET` ms-auth `/ready`) | `GET /ready` + `grpc.health.v1` on listen addr |

Use **liveness** for “restart the pod”; **readiness** for “send user traffic / register in load balancer”. Do not fail gateway liveness when a downstream is down.

## CallUnary

Wrap generated stub calls (package `microservices`):

```go
cfg := microservices.CallConfigFromSettings(settings, microservices.RetryIdempotent)
// Example: domain RPC via CallUnary (when protos exist)
err := microservices.CallUnary(ctx, cfg, log, someMethod, func(callCtx context.Context) error {
    return stub.DoSomething(callCtx, req)
})
```

| Setting (env) | Default |
|---------------|---------|
| `ms_attempt_timeout_ms` | 5000 |
| `ms_overall_timeout_ms` | 15000 |
| `ms_max_retries` | 3 |
| `ms_retry_backoff_ms` | 100 |
| `ms_breaker_failure_threshold` | 5 (0 disables) |
| `ms_breaker_open_ms` | 30000 |

Outbound gRPC uses a small **circuit breaker** on consecutive transport/`Unavailable` failures (`CallConfig.Breaker`). Gateway handlers should use `response/http.WriteOutboundError` for downstream errors.

### Retry policy

| Policy | Use |
|--------|-----|
| `RetryNone` | Creates/updates without idempotency key |
| `RetryIdempotent` | Get*, safe reads |
| `RetrySafe` | Reserved for idempotency-key writes |

No retry on gRPC codes such as `InvalidArgument`, `NotFound`, `PermissionDenied`. Retries `Unavailable`, `ResourceExhausted`, and transport-style failures (with faster retry on obvious disconnects).

## Facades (per caller, not wi-shared)

Put typed methods next to the service that calls them:

```text
ms-auth/internal/client/user/
  client.go       # msgrpc.Dial + UserServiceClient + CallConfig
  create_user.go  # CreateUser → CallUnary(..., RetryNone)
```

Map domain errors in the **service** layer (`exception`, `response/grpc`).

## Ops HTTP health

Liveness/readiness use **`wi-shared/infra/ops`**: `GET /health`, `GET /ready`, `GET /metrics` on each service listen port.

**api-gateway** `/ready` may probe upstreams with `ops.HTTPGetReadyCheck` (e.g. `http://MS_AUTH_HOST/ready`). Microservices combine gRPC + ops on one port; `/ready` uses `grpc.health.v1` internally.
