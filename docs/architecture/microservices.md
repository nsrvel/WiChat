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
| Liveness | `GET /health` on `http_addr` | `GET /health` on `grpc_addr` (same port as gRPC) |
| Readiness | `GET /ready` (deps, e.g. auth Ping) | `GET /ready` + `grpc.health.v1` on `grpc_addr` |

Use **liveness** for “restart the pod”; **readiness** for “send user traffic / register in load balancer”. Do not fail gateway liveness when a downstream is down.

## CallUnary

Wrap generated stub calls (package `microservices`):

```go
cfg := microservices.CallConfigFromSettings(settings, microservices.RetryIdempotent)
err := microservices.CallUnary(ctx, cfg, log, authv1.AuthService_Ping_FullMethodName, func(callCtx context.Context) error {
    out, err = client.Ping(callCtx, &authv1.PingRequest{})
    return err
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
| `RetryIdempotent` | Ping, Get*, safe reads |
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

## Gateway dev check

Non-production: `GET /api/v1/dev/auth-ping` calls `AuthService.Ping` on ms-auth (requires `auth_grpc_addr`, default `localhost:3001`).
