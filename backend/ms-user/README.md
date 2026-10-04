# ms-user

Workspaces, profiles, and membership (gRPC).

```bash
cp .env.example .env   # optional overrides
make dev    # Air live reload
make run    # single run
make check  # fmt, lint, test
```

Default listen: **`PORT=3002`** (gRPC + `/health`, `/ready`, `/metrics`).

From `backend/`: `make dev-ms-user`.
