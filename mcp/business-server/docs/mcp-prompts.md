# MCP Prompt Guide — Business Server

## 1. Overview

The Business MCP server exposes the **Business API** (`localhost:3900`) as MCP tools and resources. It covers:

- **Discovery** — venue search, nearby venues, personalized feed, food boxes (read-only, mostly anonymous).
- **Business Profile** — read and update brand-level profile (name, avatar, banner, description).
- **Venue Management** — list, read, update venues and their working hours (business-owner only).
- **System Health** — connectivity and operational status.

**Auth model:**
- **Anonymous** — no token required for `search_venues`, `get_recommended_venues`, `search_boxes`.
- **Bearer required** — `get_feed_items` and all mutation tools (`update_business_profile`, `update_venue_details`, `update_venue_hours`).

**Return contract:** All tools return a `dict` with a strict envelope:
```json
{
  "ok": true,
  "error": null,
  "validation_errors": [],
  "changed_fields": ["name"],
  "result": {...}
}
```
On failure:
```json
{
  "ok": false,
  "error": "human-readable message",
  "validation_errors": [{"field": "latitude", "message": "..."}],
  "changed_fields": [],
  "result": null
}
```

**Error field mapping:** Business API returns `{"error": "..."}`. The MCP server normalizes this into the envelope above.

---

## 2. Tool Prompts by Category

### 2.1 Health & Status

#### `business_health_check`

**When to use:** Verify that the Business API container is reachable.

**Example prompts:**
- "Is the business API up?"
- "Check business service health."
- "Can you reach the business backend?"

**Parameters:** None.

**Auth:** Bearer token forwarded if present (used for `/swagger/doc.json` which may require auth in some environments).

**Edge cases:**
- Raises `RuntimeError` on 5xx or connection failure — surface as "Business API is currently unavailable."

---

#### `get_business_api_status`

**When to use:** Deep diagnostics — OpenAPI metadata, path count, version.

**Example prompts:**
- "What's the business API version?"
- "How many endpoints does the business API have?"
- "Business API status check."

**Parameters:** None.

**Auth:** Bearer token required (forwarded from MCP context or `BUSINESS_API_AUTH_TOKEN` env).

**Edge cases:**
- `401` → "Unauthorized. Please check your authentication token."

---

### 2.2 Discovery — Venues

#### `search_venues`

**When to use:** User wants to find venues by keyword or location tags.

**Example prompts:**
- "Find coffee shops."
- "Search for vegan restaurants."
- "Any romantic places nearby?"
- "Look for venues with tag 'breakfast'."

**Parameters:**

| Parameter | Type | Default | Source / Example |
|-----------|------|---------|-----------------|
| `q` | `str` | `None` | Keyword: "coffee", "pizza" |
| `tags` | `str` | `None` | Comma-separated slugs: "vegan,romantic" |
| `skip` | `int` | `0` | Pagination offset |
| `limit` | `int` | `10` | Page size (max `100`) |
| `request_id` | `str` | `None` | Optional `X-Request-ID` |

**Validation:** At least one of `q` or `tags` must be provided. Otherwise returns `validation_errors`.

**Auth:** Anonymous. Bearer forwarded if present.

**Edge cases:**
- Empty filters → `{"ok": false, "validation_errors": [...]}` or `"at least one search filter is required"`.
- Empty result → "No venues found. Try different keywords or tags."
- Downstream failure → `{"ok": false, "error": "business-api returned 500: ..."}`

---

#### `get_recommended_venues`

**When to use:** List venues closest to the user's coordinates, sorted by distance.

**Example prompts:**
- "What restaurants are near me?"
- "Show venues close to 50.45, 30.52."
- "Nearest coffee shops."
- "Places around here."

**Parameters:**

| Parameter | Type | Default | Example |
|-----------|------|---------|---------|
| `lat` | `float` | — | `50.45` |
| `lon` | `float` | — | `30.52` |
| `skip` | `int` | `0` | Pagination |
| `limit` | `int` | `10` | Max `100` |
| `request_id` | `str` | `None` | Trace ID |

**Validation:** `lat` must be `-90..90`, `lon` must be `-180..180`.

**Auth:** Anonymous. Bearer forwarded if present.

**Edge cases:**
- Invalid coordinates → `validation_errors` with `"latitude must be a number between -90 and 90"`.
- Empty result → "No venues found near this location."

---

### 2.3 Discovery — Feed

#### `get_feed_items`

**When to use:** Show personalized post recommendations based on user behavior and location. Uses the 5-3-2-1-1 weighted tag quota algorithm.

**Example prompts:**
- "Show me recommended posts."
- "What should I see near my location?"
- "Personalized feed for me."
- "Recommendations near 50.45, 30.52."

**Parameters:**

| Parameter | Type | Default | Example |
|-----------|------|---------|---------|
| `lat` | `float` | — | `50.45` |
| `lon` | `float` | — | `30.52` |
| `skip` | `int` | `0` | Pagination |
| `limit` | `int` | `10` | Max `100` (default `24` in Business API) |
| `request_id` | `str` | `None` | Trace ID |

**Validation:** Same coordinate and pagination rules as `get_recommended_venues`.

**Auth:** **Bearer required.** Returns `{"ok": false, "error": "Unauthorized: Missing authentication token"}` if no token.

**Edge cases:**
- Missing token → "Authentication required for personalized recommendations."
- Invalid coordinates → `validation_errors`.
- Downstream failure → `{"ok": false, "error": "..."}`.

---

### 2.4 Discovery — Boxes

#### `search_boxes`

**When to use:** Find available food boxes (surprise bags, discount meals) sorted by distance.

**Example prompts:**
- "Find surprise boxes near me."
- "Available food boxes around here."
- "Discount meals near 50.45, 30.52."
- "Show dessert boxes from Org 5."

**Parameters:**

| Parameter | Type | Default | Example |
|-----------|------|---------|---------|
| `lat` | `float` | — | `50.45` |
| `lon` | `float` | — | `30.52` |
| `skip` | `int` | `0` | Pagination |
| `limit` | `int` | `10` | Max `100` |
| `org_id` | `int` | `None` | Filter by organization ID |
| `category_id` | `int` | `None` | Filter by category ID |
| `request_id` | `str` | `None` | Trace ID |

**Validation:** Coordinates and pagination validated. `org_id` / `category_id` are optional and added to params only when not `None`.

**Auth:** Anonymous. Bearer forwarded if present.

**Edge cases:**
- Invalid coordinates → `validation_errors`.
- Empty result → "No available boxes near this location."

---

### 2.5 Business Profile (Mutations)

#### `get_business_profile`

**When to use:** Retrieve a brand profile by ID.

**Example prompts:**
- "Show me business profile 10."
- "What's the profile for brand ID 42?"
- "Get business details for ID 7."

**Parameters:**

| Parameter | Type | Example |
|-----------|------|---------|
| `business_id` | `int` | `10` |
| `auth_token` | `str` | Forwarded automatically |
| `request_id` | `str` | `None` |

**Auth:** Bearer forwarded from context or env.

**Return:** `{"ok": true, "business_id": 10, "result": {"id": 10, "name": "..."}}`

---

#### `update_business_profile`

**When to use:** Update mutable brand fields (name, avatar, banner, description).

**Example prompts:**
- "Update my business name to 'New Name'."
- "Change the business description."
- "Set a new avatar for brand 10."

**Parameters:**

| Parameter | Type | Example |
|-----------|------|---------|
| `business_id` | `int` | `10` |
| `payload` | `dict` | `{"name": "New Name"}` |
| `auth_token` | `str` | Forwarded automatically |
| `request_id` | `str` | `None` |

**Validation:** Only allowed fields: `name`, `avatar`, `banner`, `description`. Unknown fields → `validation_errors`. Empty payload → `validation_errors`.

**Return:** `{"ok": true, "business_id": 10, "changed_fields": ["name"], "result": {...}}`

**Edge cases:**
- Validation failure → `{"ok": false, "validation_errors": [{"field": "badField", "message": "unknown field"}]}`
- Forbidden / not found → `{"ok": false, "error": "..."}`

---

### 2.6 Venue Management (Mutations)

#### `list_business_venues`

**When to use:** List all venues belonging to a brand.

**Example prompts:**
- "Show me all venues for business 10."
- "List my brand locations."
- "What venues does brand 42 have?"

**Parameters:**

| Parameter | Type | Default | Example |
|-----------|------|---------|---------|
| `business_id` | `int` | — | `10` |
| `skip` | `int` | `0` | Pagination |
| `limit` | `int` | `10` | Max `100` |
| `auth_token` | `str` | `None` | Forwarded if present |
| `request_id` | `str` | `None` | Trace ID |

**Auth:** Optional for reading; mutations require ownership.

**Return:** `{"ok": true, "business_id": 10, "result": {"items": [...], "total": 5}}`

---

#### `get_venue_details`

**When to use:** Retrieve detailed information about a specific venue, with ownership verification.

**Example prompts:**
- "Show me venue 7."
- "Get details for venue ID 15."
- "What do we know about venue 3?"

**Parameters:**

| Parameter | Type | Example |
|-----------|------|---------|
| `business_id` | `int` | `10` (for ownership check) |
| `venue_id` | `int` | `7` |
| `auth_token` | `str` | Forwarded automatically |
| `request_id` | `str` | `None` |

**Ownership check:** The venue must belong to `business_id`. If `venue.brand.id != business_id`, returns `{"ok": false, "error": "unauthorized access to another business venue"}`.

---

#### `update_venue_details`

**When to use:** Update venue fields (name, avatar, banner, description, latitude, longitude, tagIds).

**Example prompts:**
- "Update venue 7 name to 'Downtown Branch'."
- "Change venue 15 description."
- "Set new coordinates for venue 3."

**Parameters:**

| Parameter | Type | Example |
|-----------|------|---------|
| `business_id` | `int` | `10` |
| `venue_id` | `int` | `7` |
| `payload` | `dict` | `{"name": "New Name", "latitude": 50.45}` |
| `auth_token` | `str` | Forwarded automatically |
| `request_id` | `str` | `None` |

**Validation:**
- Allowed fields: `name`, `avatar`, `banner`, `description`, `latitude`, `longitude`, `tagIds`.
- `latitude`: `-90..90`
- `longitude`: `-180..180`
- `tagIds`: list of integers, max 5 items.

**Ownership check:** Performed before update.

**Return:** `{"ok": true, "business_id": 10, "venue_id": 7, "changed_fields": ["name"], "result": {...}}`

**Edge cases:**
- Validation failure → `validation_errors` with field-level messages.
- Foreign venue → `{"ok": false, "error": "unauthorized access to another business venue"}`.

---

#### `update_venue_hours`

**When to use:** Update weekly working hours for a venue.

**Example prompts:**
- "Set venue 7 hours to 9-18 Mon-Fri."
- "Update working hours for venue 15."
- "Make venue 3 closed on Sundays."

**Parameters:**

| Parameter | Type | Example |
|-----------|------|---------|
| `business_id` | `int` | `10` |
| `venue_id` | `int` | `7` |
| `payload` | `dict` | `{"days": [{"weekday": 1, "openTime": "09:00", "closeTime": "18:00"}]}` |
| `auth_token` | `str` | Forwarded automatically |
| `request_id` | `str` | `None` |

**Validation (payload):**
- `days` must be a non-empty array (1–7 items).
- Each day: `weekday` (integer `1..7`), `openTime` / `closeTime` (`HH:MM` string or `null` for closed).
- `openTime` must be before `closeTime`.
- No duplicate weekdays.
- Both `openTime` and `closeTime` must be provided together (or both `null`).

**Ownership check:** Performed before update.

**Return:** `{"ok": true, "business_id": 10, "venue_id": 7, "changed_fields": ["days"], "result": {...}}`

**Edge cases:**
- Partial pair (`openTime` without `closeTime`) → `validation_errors`: `"both openTime and closeTime must be provided together"`.
- Invalid time format → `"openTime must be HH:MM"`.
- `openTime >= closeTime` → `"openTime must be before closeTime"`.
- Duplicate weekday → `"duplicate weekday"`.

---

## 3. Resource Prompts

Resources are read-only contextual data. They help the AI understand API structure and constraints.

### 3.1 `sharebite://business/api-info`

**Content:** Service name, base URL, timeout, auth forwarding policy, request ID propagation.

**When to surface:** User asks "How is this server configured?" or "What API does this connect to?"

---

### 3.2 `sharebite://business/openapi-summary`

**Content:** Live OpenAPI metadata from `/swagger/doc.json` — title, version, path count, sorted path list.

**When to surface:** User asks "What endpoints are available?" or "How many API paths exist?"

**Note:** This resource makes an actual HTTP call to the Business API. It may fail if the API is down.

---

### 3.3 `sharebite://business/profile-schema`

**Content:** JSON Schema for `update_business_profile` payload.

**Allowed fields:**
- `name`: string, `minLength: 3`, `maxLength: 40`
- `avatar`: string or null
- `banner`: string or null
- `description`: string or null

**When to surface:** Before calling `update_business_profile`, to validate the user's intended payload.

---

### 3.4 `sharebite://business/venue-schema`

**Content:** JSON Schema for `update_venue_details` payload.

**Allowed fields:**
- `name`: string or null
- `avatar`: string or null
- `banner`: string or null
- `description`: string or null
- `latitude`: number or null, `-90..90`
- `longitude`: number or null, `-180..180`
- `tagIds`: array of integers or null, max 5 items

**When to surface:** Before calling `update_venue_details`, to validate the user's intended payload.

---

### 3.5 `sharebite://business/venue-hours-format`

**Content:** JSON example of valid venue hours payload.

**Example structure:**
```json
{
  "days": [
    {"weekday": 1, "openTime": "09:00", "closeTime": "18:00"},
    {"weekday": 7, "openTime": null, "closeTime": null}
  ]
}
```

**When to surface:** Before calling `update_venue_hours`, to show the user the correct format.

**Key rules to relay:**
- `weekday`: `1` (Monday) through `7` (Sunday).
- `openTime` / `closeTime`: `HH:MM` strings or `null` for closed days.
- Both must be provided together (or both `null`).
- `openTime` must be strictly before `closeTime`.
- No duplicate weekdays in the array.

---

## 4. Cross-Tool Patterns

### 4.1 "Find venues near me"

1. `get_recommended_venues(lat=..., lon=..., limit=10)` → present results with distance.

**If empty:** "No venues found near this location. Try expanding your search radius or different coordinates."

---

### 4.2 "Show me recommended posts"

1. `get_feed_items(lat=..., lon=..., limit=24)` → present personalized posts.

**If 401:** "Authentication required for personalized recommendations. Please provide a Bearer token."

---

### 4.3 "Search for food boxes"

1. `search_boxes(lat=..., lon=..., limit=10)` → list available boxes.
2. If user specifies category: `search_boxes(..., category_id=2)`.
3. If user specifies organization: `search_boxes(..., org_id=5)`.

---

### 4.4 "Update my business profile"

1. `get_business_profile(business_id=10)` → show current state.
2. `update_business_profile(business_id=10, payload={"name": "New Name"})` → confirm `changed_fields`.

**If validation fails:** Present `validation_errors` as a bullet list of what to fix.

---

### 4.5 "Update venue hours"

1. `get_venue_details(business_id=10, venue_id=7)` → show current hours (if any).
2. `update_venue_hours(business_id=10, venue_id=7, payload={"days": [...]})` → confirm `changed_fields: ["days"]`.

**Validation helper:** Before calling, read `sharebite://business/venue-hours-format` to show the correct structure.

**If validation fails:**
- `"both openTime and closeTime must be provided together"` → explain that partial pairs are not allowed.
- `"openTime must be before closeTime"` → suggest swapping times.
- `"duplicate weekday"` → remove the duplicate entry.

---

### 4.6 "Show me all venues for my business"

1. `list_business_venues(business_id=10, skip=0, limit=100)` → paginated list.
2. If user picks one: `get_venue_details(business_id=10, venue_id=<id>)` → full details.

---

## 5. Guardrails & Warnings

### 5.1 Authentication boundaries

| Tool | Token required? | What happens without token |
|------|-----------------|---------------------------|
| `business_health_check` | Optional | Works; token forwarded if present. |
| `get_business_api_status` | Yes | `RuntimeError` (401). |
| `search_venues` | No | Works. |
| `get_recommended_venues` | No | Works. |
| `get_feed_items` | **Yes** | `{"ok": false, "error": "Unauthorized: Missing authentication token"}` |
| `search_boxes` | No | Works. |
| `get_business_profile` | Yes | `RuntimeError` or `{"ok": false, "error": "..."}` |
| `update_business_profile` | Yes | `RuntimeError` or `{"ok": false, "error": "..."}` |
| `list_business_venues` | Optional | Works for reading. |
| `get_venue_details` | Yes | Ownership check + auth required. |
| `update_venue_details` | Yes | Ownership check + auth required. |
| `update_venue_hours` | Yes | Ownership check + auth required. |

### 5.2 Mutation safety

- **Business ID must never be guessed.** It must come from authenticated context or explicit trusted input.
- **Ownership verification** is enforced for all venue mutations. If the venue does not belong to the provided `business_id`, the tool returns `{"ok": false, "error": "unauthorized access to another business venue"}`.
- **Validation first.** All mutation payloads are validated before reaching the Business API. Invalid fields return `validation_errors` without making upstream calls.
- **Changed fields tracking.** Successful mutations return `changed_fields` so the user knows exactly what was modified.

### 5.3 Common Integration Pitfalls

When building integrations or debugging the Business MCP, watch out for these recurring issues:

- **Schema drift in mutations.** `update_business_profile` silently ignores unknown fields. If you pass `logo` instead of `avatar`, nothing breaks — but nothing updates either. Always check `sharebite://business/profile-schema` before building a payload.
- **tagIds type confusion.** The API expects `tagIds: [1, 2, 3]` (list of integers). Passing `"1,2,3"` or `["1", "2", "3"]` will fail validation.
- **Partial time pairs in hours.** `update_venue_hours` rejects partial `openTime` / `closeTime` pairs. Both must be `HH:MM` strings, or both must be `null`. A single `null` in a pair returns a validation error.
- **Guessing IDs.** `business_id` and `venue_id` must come from authenticated context or previous tool results. Hardcoding IDs leads to `404` or ownership violations.
- **Anonymous feed calls.** `get_feed_items` requires a Bearer token. Calling it without auth always returns `401`. Ensure your integration forwards the token from MCP context or sets `BUSINESS_API_AUTH_TOKEN`.

### 5.4 Error handling cheat sheet

| Envelope state | Assistant message |
|----------------|-------------------|
| `ok: false` + `validation_errors` | Present each error as: `"<field>": <message>`. Ask user to fix and retry. |
| `ok: false` + `error: "Unauthorized..."` | "Authentication required. Please provide a valid Bearer token." |
| `ok: false` + `error: "unauthorized access..."` | "This venue does not belong to the specified business. Check your business_id and venue_id." |
| `ok: false` + `error: "not found"` | "The requested item was not found. Verify the ID is correct." |
| `ok: false` + `error: "business-api returned 5xx..."` | "Business API is temporarily unavailable. Please try again later." |
| `ok: true` + `changed_fields: [...]` | "Successfully updated. Changed: <fields>." |
