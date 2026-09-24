# models

gRPC/protobuf API contracts shared by WiChat services.

- **Sources:** `auth/v1/*.proto` (add more domains under `models/<domain>/v1/`).
- **Generated:** `gen/` — run `make proto` from `wi-shared/` (requires [buf](https://buf.build/docs/installation)). Do not hand-edit `gen/`.

Import generated types from service code:

```go
import authv1 "github.com/wichat/wichat/backend/wi-shared/models/gen/auth/v1"
```

When services live in separate repos, depend on a tagged **wi-shared** module release that includes `models/gen`.
