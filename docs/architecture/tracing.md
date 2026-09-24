# Distributed tracing (OpenTelemetry)

Process-wide tracing lives in **`backend/wi-shared/infra/otel`** — separate from [`microservices`](microservices.md) so HTTP/gRPC inter-service kits do not own observability bootstrap.

## Bootstrap

Call once from each service `main`:

```go
shutdown, err := otel.InstallTracer(context.Background(), "api-gateway")
defer func() { _ = shutdown(context.Background()) }()
```

| Variable | Role |
|----------|------|
| `OTEL_EXPORTER_OTLP_ENDPOINT` | When set, enables OTLP HTTP export; when unset, noop (default local dev) |
| `OTEL_SERVICE_NAME` | Fallback service name if the argument to `InstallTracer` is empty |

## gRPC hooks

Compose at dial/server sites (do not hide inside `microservices`):

```go
conn, err := msgrpc.Dial(target, msgrpc.DefaultDialOptions, otel.GRPCClientDialOptions()...)
grpcSrv := msgrpc.NewServer(log, otel.GRPCServerOptions()...)
```

Spans use W3C trace context propagation via the global text map propagator set in `InstallTracer`.

## Follow-up

- `otelhttp` middleware on the api-gateway HTTP stack
- OTLP collector / Tempo in `deploy/`

See also [logging.md](logging.md) for request ID correlation (complements traces, not a replacement).
