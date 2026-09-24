# Application logging

WiChat services use **`wi-shared/infra/logger`** (`log/slog`) for structured application logs.

## Output

| Environment | Format | Destination |
|-------------|--------|-------------|
| `production` | JSON | stdout |
| Other (e.g. `development`) | text | stdout |

Default fields on every line: `service`, `env`. Additional fields via `log.Info("msg", "key", value)`.

Container platforms and **Promtail/Alloy → Loki → Grafana** collect stdout; the app does not push logs directly in v1.

## Configuration

Per-service config (e.g. `log_level` in ms-auth) maps to slog levels: `debug`, `info`, `warn`, `error`.

```go
appLog := logger.New(logger.Options{
    Service: "ms-auth",
    Env:     cfg.Env,
    Level:   cfg.LogLevel,
})
```

Inject `logger.Logger` into repository, service, and delivery layers ([wichat-go-layers.mdc](../../.cursor/rules/wichat-go-layers.mdc)).

## HTTP errors and logging

Gateway handlers should pass:

```go
responsehttp.WriteError(w, r, err, responsehttp.WithLogger(logger.HTTPWriteErrorLogger(appLog)))
```

| Situation | Logged? |
|-----------|---------|
| Raw error → 500 | Yes |
| `exception.Internal(cause)` → 500 | Yes (`cause`) |
| Expected 4xx (`NotFound`, `InvalidInput`, …) | No |
| gRPC downstream internal | Yes (`WriteErrorFromGRPC`) |

Client responses stay i18n codes only — see [error-handling.md](error-handling.md).

## Do not log

- Passwords, tokens, refresh cookies, API secrets
- Full message bodies or PII unless explicitly required for audit (audit uses Postgres, not app logs)

## Gateway request ID

`gateway/internal/middleware` sets `microservices.RequestIDHeader` (`X-Request-ID`) on every response and stores the value on `r.Context()` via `microservices.WithRequestID`.

Outbound gRPC from the gateway uses `wi-shared/microservices` client interceptors to send `x-request-id` metadata; ms-auth server interceptors read it back into context (recovery logs include `request_id` when present).

## Follow-up

- Add `request_id` to slog lines from context (automatic on every log call)
- HTTP OpenTelemetry (`otelhttp`) on the gateway
- Deploy OTLP collector / Tempo in `deploy/` when ops slice lands

See [microservices.md](microservices.md) for CallUnary and retries; [tracing.md](tracing.md) for `wi-shared/infra/otel`.
