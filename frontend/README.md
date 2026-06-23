# Share Bite — Unified SPA

Single deployable React app for guest/community and business features. Legacy `frontend/` is reference-only.

## Stack

React 19, Vite, TypeScript, React Router 7, TanStack Query, Axios, RHF + Zod, Tailwind v4, shadcn/ui, sonner, lucide.

## Local development

Start backend services from repo root:

```bash
docker compose -f build/compose.infra.yaml --env-file .env up -d
make migrate-up
make run-guest && make run-business && make run-auth
go run ./cmd/notifications-service
```

Run the SPA:

```bash
cd frontend
npm install
npm run dev
```

Leave `VITE_API_BASE_URL` empty so Vite proxies API calls:

| Proxy path | Service | Port |
|------------|---------|------|
| `/api/guest` | Guest API | 3800 |
| `/api/business` | Business API | 3900 |
| `/api/auth`, `/api/admin`, `/api/users`, `/api/user` | Auth / admin | 3850 |
| `/api/notifications` | Notifications | 4005 |

OAuth (optional): set `VITE_GOOGLE_CLIENT_ID` and `VITE_GOOGLE_REDIRECT_URI` (default `…/oauth/google/callback`).

## Layout

Sidebar-first chrome (no top Navbar): brand, user row with notifications, theme toggle, gold **+ Share a Bite** CTA, nav sections, logout dialog. Minimal header only on `/auth` and `/oauth/*`.

## Build

```bash
npm run build
```
