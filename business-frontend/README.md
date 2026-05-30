# Share Bite — unified frontend

Single Vite app combining guest/community features (from `frontend/`) and business features (boxes, venues, recommendations).

## Local development

1. Start infra: `docker compose -f build/compose.infra.yaml --env-file .env up -d` from repo root.
2. Run migrations: `make migrate-up` from repo root.
3. Start APIs: `make run-guest`, `make run-business`, `make run-auth`.
4. Copy env: `cp .env.example .env.local` and set `VITE_GOOGLE_CLIENT_ID` if using OAuth.
5. Run app:

```bash
npm install
npm run dev
```

Vite proxies `/api/guest`, `/api/auth`, `/api/business` to local services.

## Layout

- **Navbar** (top): search, theme, auth, notifications, profile.
- **Sidebar** (left): community, business, settings, admin links.
- **`/`** home is role-based: `business` → recommendations feed; `user` → guest feed; `admin`/`moderator` → `/admin`.
