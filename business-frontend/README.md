# Share Bite — Unified Frontend

Single React SPA for guest community features and business tools. Sidebar-first layout (no top navbar on main routes).

## Stack

React 19, Vite, TypeScript, React Router 7, TanStack Query, Axios, RHF + zod, Tailwind v4, shadcn/ui, sonner, lucide.

## Local development

From repo root:

```bash
docker compose -f build/compose.infra.yaml --env-file .env up -d
make migrate-up
make run-guest && make run-business && make run-auth
go run ./cmd/notifications-service
cd business-frontend && npm install && npm run dev
```

Set `VITE_API_BASE_URL=` empty so Vite proxies are used.

OAuth: `VITE_GOOGLE_CLIENT_ID`, redirect `…/oauth/google/callback`.

## Vite proxies

| Service | Port | Proxy |
|---------|------|-------|
| Guest API | 3800 | `/api/guest` → strips `/api/guest` |
| Business API | 3900 | `/api/business` → strips `/api` |
| Auth / admin | 3850 | `/api/auth`, `/api/admin`, `/api/users`, `/api/user` |
| Notifications | 4005 | `/api/notifications` → strips `/api/notifications` |

## Layout

- **Main routes**: left `Sidebar` (brand, avatar + @username + notification bell, theme toggle, post CTA, nav, logout) + full-width content.
- **Auth / OAuth**: minimal SB header only, no sidebar.

## Build

```bash
cd business-frontend && npm run build
```
