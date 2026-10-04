# Frontend

Primary web client: **`frontend/web-client/`** (Next.js App Router, PWA-capable). Desktop (Electron) and mobile (React Native) attach later and reuse the same API.

Work **one vertical slice at a time**. UI-only slices with MSW/fixtures are OK if Zod shapes match planned **api-gateway** `/api/v1/` contracts.

## Layout (`web-client/src/`)

```
app/              # routing + layouts only (thin pages)
features/         # domain screens & feature-local components
components/ui/    # shadcn
components/layout/# shared shell
services/         # *.key | *.service | *.query | *.mutation
schemas/          # Zod
integrations/     # api-client, tanstack-query, i18n, msw, realtime (later)
hooks/            # generic UI hooks
lib/
```

Agent rules: `.cursor/rules/wichat-web-layout.mdc`, `.cursor/rules/wichat-web-data.mdc`. Skill: `.cursor/skills/wichat-web-client/`.

Env: copy `web-client/.env.example` → `web-client/.env.local`. From repo root: `pnpm install` then `pnpm --filter web-client dev` (UI **8080**, API **3000**). Monorepo layout: [docs/architecture/monorepo-layout.md](../docs/architecture/monorepo-layout.md).

Pattern reference (not runtime): `references/react-tanstack/`.
