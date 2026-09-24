# Application logging

WiChat services use **`wi-shared/logger`** (`log/slog`) for structured application logs.

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

`gateway/internal/middleware` sets `X-Request-ID` on every response (generated when the client omits it).

## Follow-up (not v1)

- Propagate `request_id` into slog context and downstream gRPC metadata
- gRPC unary interceptor for internal errors (ms-auth has recovery + internal log today)
- Metrics (Prometheus) and tracing (OpenTelemetry) — separate slices
