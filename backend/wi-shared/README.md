# wi-shared

Shared libraries for WiChat backend services (config, **infra/** observability, errors, response writers, gRPC **models**).

Protobuf contracts: [models/README.md](models/README.md) — `make proto` / `make buf-lint` in this module.

Error model and HTTP/gRPC conventions: [docs/architecture/error-handling.md](../../docs/architecture/error-handling.md).

Logging (`log/slog`): [docs/architecture/logging.md](../../docs/architecture/logging.md).

Metrics (Prometheus): [docs/architecture/metrics.md](../../docs/architecture/metrics.md). Use `NewRegistry()` and `NewHTTPMetrics` + `Middleware` (register RED collectors once per registry).

Service/repository/delivery layout in services: `.cursor/rules/wichat-go-layers.mdc`.
