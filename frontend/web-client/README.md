# WiChat web client (Next.js)

From **repo root** (pnpm workspace + Husky hooks):

```bash
pnpm install
cp frontend/web-client/.env.example frontend/web-client/.env.local
pnpm --filter web-client dev
```

Or from this directory: `pnpm dev` (after root `pnpm install`).

Dev server: **http://localhost:8080**. **api-gateway** tetap **:3000** (`NEXT_PUBLIC_API_URL`); ms-auth/ms-user **3001/3002** — UI yang pakai port lain, bukan gateway.

**Quality gate before PR:**

```bash
pnpm --filter web-client check   # format + eslint + tsc + build
make check                       # same, from this folder
```

Pre-commit runs **lint-staged** (ESLint + Prettier on staged web-client files).

Conventions: [frontend/README.md](../README.md), `.cursor/rules/wichat-web-*.mdc`.
