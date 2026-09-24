# Service-to-service (`wi-shared/microservices`)

The package name is **`microservices`**, not `rpc`, because WiChat will call **other runtimes and transports** (Python/JS plugins, HTTP APIs) — not only Go gRPC. v1 ships **gRPC** under `microservices/grpc`; a future `microservices/http` can share correlation and retry patterns for outbound HTTP.

Contracts for Go gRPC live in `wi-shared/models`. Cross-cutting **request ID** (`RequestIDHeader`, `WithRequestID`) stays at the `microservices` root for gateway HTTP and gRPC metadata.

Tracing bootstrap is **`wi-shared/infra/otel`** — see [tracing.md](tracing.md).

## gRPC dial and server

```go
conn, err := msgrpc.Dial("localhost:50051", msgrpc.DefaultDialOptions, otel.GRPCClientDialOptions()...)
grpcSrv := msgrpc.NewServer(appLog, otel.GRPCServerOptions()...)
```

- Outbound: `x-request-id` metadata from context.
- Inbound: read `x-request-id` into context, panic recovery.

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

Non-production: `GET /api/v1/dev/auth-ping` calls `AuthService.Ping` on ms-auth (requires `auth_grpc_addr`, default `localhost:50051`).
