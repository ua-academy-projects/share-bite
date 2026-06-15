# MCP Prompt Guide — Guest Server

## 1. Overview

The Guest MCP server exposes the **Guest API** (`localhost:3800`) as MCP tools and resources. It covers:

- **UGC posts** — published reviews, photos, and mentions.
- **Customer profiles** — public profiles, followers, and following lists.
- **Collections** — user-curated venue lists with collaborators.
- **System health** — connectivity and operational status.

**Auth model:**
- **Public endpoints** — no token required (`get_customer_by_username`, `search_posts`, `get_post`, `get_post_authors`).
- **Optional auth** — token improves access but is not mandatory (`get_collection`, `get_collection_venues`).
- **Required auth** — Bearer token mandatory (`list_my_collections`, `get_guest_api_status`).

**Error contract:** All tools return a JSON string. On failure the payload contains:
```json
{"error": "<code>", "message": "<human-readable>"}
```
Common codes: `unauthorized`, `forbidden`, `not_found`, `downstream_failure`.

**Edge cases:**
- Raises `RuntimeError` on any API error (4xx, 5xx, connection failure) — surface as "Guest API is currently unavailable."

---

## 2. Tool Prompts by Category

### 2.1 Health & Status

#### `guest_health_check`

**When to use:** Verify that the Guest API container is reachable before making substantive calls.

**Example prompts:**
- "Is the guest API up?"
- "Check service health."
- "Can you reach the backend?"

**Parameters:** None.

**Edge cases:**
- Raises `RuntimeError` on 5xx or connection failure — surface as "Guest API is currently unavailable."

---

#### `get_guest_api_status`

**When to use:** Deep diagnostics — check PostgreSQL, Redis, and internal component health.

**Example prompts:**
- "What's the system status?"
- "Is the database connected?"
- "Any backend issues right now?"

**Parameters:** None.

**Auth:** Bearer token required (forwarded from MCP context or `GUEST_API_AUTH_TOKEN` env).

**Edge cases:**
- `401` → "Unauthorized. Please check your authentication token."
- `503` + `database: disconnected` → "Database connectivity issue detected."

---

### 2.2 Discovery — Posts

#### `search_posts`

**When to use:** User wants to browse UGC reviews, find posts about a topic, or see recent activity.

**Example prompts:**
- "Find recent posts about coffee."
- "Show me reviews of vegan restaurants."
- "What did people post last week?"
- "List posts by user alice."

**Parameters:**

| Parameter | Type | Default | Source / Example |
|-----------|------|---------|-----------------|
| `limit` | `int` | `20` | User says "show 50" → `50` (max `100`) |
| `offset` | `int` | `0` | Pagination cursor |
| `author_id` | `str` (UUID) | `None` | From `get_customer_by_username` result; **not** a raw username |

**Edge cases:**
- Empty result → "No posts found. Try different keywords or a broader search."
- `401` → "Authentication issue while searching posts."
- `downstream_failure` → "Guest API is temporarily unavailable. Please try again later."

**Cross-tool note:** To search by *username*, chain: `get_customer_by_username` → `search_posts(author_id=<uuid>)`.

---

#### `get_post`

**When to use:** User references a specific post ID and wants details (text, images, likes, venue linkage).

**Example prompts:**
- "Show me post 42."
- "What was in that review?"
- "Details of post #123."

**Parameters:**

| Parameter | Type | Example |
|-----------|------|---------|
| `id` | `int` | `42` |

**Edge cases:**
- `404` → "Post not found. It may have been deleted or the ID is incorrect."

---

#### `get_post_authors`

**When to use:** Explain who contributed to a collaborative post (owner + accepted collaborators).

**Example prompts:**
- "Who wrote this post?"
- "Who are the authors of post 42?"
- "Did multiple people review this place?"

**Parameters:**

| Parameter | Type | Example |
|-----------|------|---------|
| `id` | `int` | `42` |

**Edge cases:**
- `404` → "Post not found or has no recorded authors."

---

### 2.3 Discovery — Customers

#### `get_customer_by_username`

**When to use:** Resolve a human-readable username to a customer profile and UUID.

**Example prompts:**
- "Find user alice."
- "Who is @bob?"
- "Show me the profile for charlie."

**Parameters:**

| Parameter | Type | Example |
|-----------|------|---------|
| `username` | `str` | `"alice"` |

**Edge cases:**
- `404` → "No customer found with that username."

**Cross-tool note:** The returned `customer.id` (UUID) is the input for `author_id` in `search_posts`, `get_customer_followers`, and `get_customer_following`.

---

#### `get_customer_followers`

**When to use:** Show who follows a specific customer.

**Example prompts:**
- "Who follows alice?"
- "Show me bob's followers."
- "How many followers does charlie have?"

**Parameters:**

| Parameter | Type | Default | Example |
|-----------|------|---------|---------|
| `id` | `str` (UUID) | — | From `get_customer_by_username` |
| `page_size` | `int` | `20` | "Show 50 followers" → `50` |
| `page_token` | `str` | `None` | Pagination cursor from previous response |

**Edge cases:**
- `403` → "This user's follower list is private."
- `404` → "Customer not found."

---

#### `get_customer_following`

**When to use:** Show who a specific customer follows.

**Example prompts:**
- "Who does alice follow?"
- "Show me who bob is following."
- "What accounts does charlie subscribe to?"

**Parameters:**

| Parameter | Type | Default | Example |
|-----------|------|---------|---------|
| `id` | `str` (UUID) | — | From `get_customer_by_username` |
| `page_size` | `int` | `20` | — |
| `page_token` | `str` | `None` | Pagination cursor |

**Edge cases:**
- `403` → "This user's following list is private."
- `404` → "Customer not found."

---

### 2.4 Discovery — Collections

#### `list_my_collections`

**When to use:** User wants to see their own saved collections.

**Example prompts:**
- "Show my collections."
- "What lists have I created?"
- "My saved places."

**Parameters:**

| Parameter | Type | Default | Example |
|-----------|------|---------|---------|
| `page_size` | `int` | `20` | — |
| `page_token` | `str` | `None` | Pagination cursor |

**Auth:** **Required.** Bearer token must be forwarded via MCP context or `GUEST_API_AUTH_TOKEN` env.

**Edge cases:**
- Missing token → `{"error": "unauthorized", "message": "Authentication required to list your collections."}`
- `downstream_failure` → "Unable to load collections right now."

---

#### `get_collection`

**When to use:** Retrieve a specific collection by ID.

**Example prompts:**
- "Show me collection col-abc-123."
- "What's in my 'Favorites' list?"
- "Details of collection XYZ."

**Parameters:**

| Parameter | Type | Example |
|-----------|------|---------|
| `collection_id` | `str` (UUID) | `"col-abc-123"` |

**Auth:** Optional. Public collections accessible without auth; private collections require ownership.

**Edge cases:**
- `404` → "Collection not found, is private, or does not belong to you."

---

#### `get_collection_venues`

**When to use:** List venues inside a collection, ordered by sort order.

**Example prompts:**
- "What venues are in my Favorites collection?"
- "Show places from collection col-abc-123."
- "List restaurants in that collection."

**Parameters:**

| Parameter | Type | Example |
|-----------|------|---------|
| `collection_id` | `str` (UUID) | `"col-abc-123"` |

**Auth:** Optional for public collections.

**Edge cases:**
- `404` → "Collection not found or inaccessible."
- Empty venues → "This collection has no venues yet."

---

## 3. Resource Prompts

Resources are read-only contextual data exposed to the LLM client. They are **not** invoked as tools; instead, the client may read them to understand API capabilities.

### 3.1 `sharebite://guest/api-info`

**Content:** API version, commit, build time, environment.

**When to surface:** User asks "What version is running?" or "What API is this?"

---

### 3.2 `sharebite://guest/openapi-summary`

**Content:** Full Swagger/OpenAPI JSON specification.

**When to surface:** User asks "What endpoints exist?" or "How is the API structured?"

---

### 3.3 `sharebite://guest/search-filters`

**Content:** Reference JSON describing available filters for posts, customers, and collections.

**When to surface:** User asks "How do I search?" or "What filters are available?"

**Key facts to relay:**
- Posts: `limit` (1–100), `offset` (0–1000), optional `author_id` (UUID).
- Customers: `page_size` (1–100), optional `page_token`.
- Collections: `page_size` (1–100), optional `page_token`.

---

### 3.4 `sharebite://guest/recommendation-signals`

**Content:** Explanation of how behavioral recommendations are generated.

**When to surface:** User asks "Why am I seeing this?" or "How do recommendations work?"

**Key facts to relay:**
- This resource returns **static metadata** about the recommendation algorithm.
- Algorithm: weighted tag quotas (5-3-2-1-1) with H3 spatial indexing.
- The algorithm itself requires `lat`, `lon`, and Bearer token for personalization — but that applies to the **Business MCP** `get_feed_items` tool, not this static resource.
- Fallback: random nearby posts when no like history exists.

**Note:** This resource is **static**. For personalized recommendations, the user must use the **Business MCP** `get_feed_items` tool (requires `lat`, `lon`, and Bearer token).

---

## 4. Cross-Tool Patterns

These are common multi-step workflows the assistant should handle natively.

### 4.1 "Show me posts by user X"

1. `get_customer_by_username(username="X")` → extract `customer.id` (UUID).
2. `search_posts(author_id=<uuid>, limit=20)` → present results.

**If step 1 fails (404):** "User X not found."

---

### 4.2 "Show me who follows user X"

1. `get_customer_by_username(username="X")` → extract `customer.id`.
2. `get_customer_followers(id=<uuid>, page_size=20)` → present results.

**If step 1 fails:** "User X not found."
**If step 2 returns 403:** "X's follower list is private."

---

### 4.3 "What's in my collections?"

1. `list_my_collections(page_size=20)` → list collection names and IDs.
2. If user picks one: `get_collection_venues(collection_id=<id>)` → list venues.

**If step 1 returns 401:** "Please authenticate to view your collections."

---

### 4.4 "Tell me about post X"

1. `get_post(id=X)` → basic post data.
2. `get_post_authors(id=X)` → list of contributors.

Present combined result: content + authors + venue linkage.

---

## 5. Guardrails & Warnings

### 5.1 Authentication boundaries

| Tool | Token required? | What happens without token |
|------|-----------------|---------------------------|
| `guest_health_check` | No | Works always. |
| `get_guest_api_status` | Yes | `RuntimeError` (401). |
| `search_posts` | No | Works; some posts may be hidden. |
| `get_post` | No | Works for public posts. |
| `get_customer_by_username` | No | Works for public profiles. |
| `get_customer_followers` | No | May return 403 for private profiles. |
| `list_my_collections` | **Yes** | `{"error": "unauthorized"}`. |
| `get_collection` | Optional | Public collections accessible. |
| `get_collection_venues` | Optional | Public collections accessible. |


### 5.2 Common Integration Pitfalls

When building integrations or debugging the Guest MCP, watch out for these recurring issues:

- **Username vs UUID mismatch.** `search_posts` expects `author_id` as a UUID, not a raw username. Always resolve via `get_customer_by_username` first, or the filter will silently fail.
- **Private profile assumptions.** `get_customer_followers` and `get_customer_following` return `403` for private profiles. Don't assume every user has a public follower list.
- **Missing auth on collections.** `list_my_collections` requires a Bearer token. If your integration calls it without auth, you'll get a clear `401` — but it's better to check for the token upfront.
- **Post source confusion.** Guest posts (`/posts/`) are UGC reviews written by customers. Business posts (`/business/posts/`) are brand-owned content. They are different APIs, different schemas, and different MCP servers.

### 5.3 Error handling cheat sheet

| Downstream code | Assistant message |
|-----------------|-------------------|
| `401` | "Authentication required. Please provide a valid token." |
| `403` | "Access denied. This content is private or restricted." |
| `404` | "Not found. The requested item may not exist or the ID is incorrect." |
| `downstream_failure` | "Service temporarily unavailable. Please try again later." |
| `api_error` (other 4xx) | Relay the exact `message` from the payload. |
