# Share Bite MCP Security Standards

This document defines security requirements for Share Bite MCP servers.

## Security Principles

- MCP servers are adapters, not privileged backdoors.
- Go services remain the source of truth for authorization and business validation.
- Tools must not expose secrets, raw database access, or infrastructure internals.
- Write actions must be explicit about permissions and side effects.
- Admin write actions must produce audit logs.

## Auth Token Forwarding

MCP servers must forward bearer tokens to the owning Go API.

Standard behavior:

1. Extract `Authorization` from MCP request context headers.
2. Normalize tokens to `Bearer <token>` when forwarding.
3. For local stdio development only, optionally use `<SERVICE>_API_AUTH_TOKEN` as a fallback.
4. Never include access tokens, refresh tokens, API keys, passwords, or authorization headers in tool/resource responses.

## Role Checks And Permissions

MCP tools must not invent authorization decisions that conflict with Go APIs. They may enforce pre-checks for defense in depth.

- Read tools should state whether authentication is required.
- Write tools must define required permissions in the description or metadata.
- Admin write tools must require an admin/moderator role as appropriate.
- Permission names should match Go API concepts where possible.

Example permission statement:

```text
Side effects: updates venue hours. Required permission: business:venue:update. Auth: forwards caller bearer token.
```

## Audit Logging

Admin write tools must emit audit logs for both successful and denied attempts.

Audit events should include:

- timestamp;
- tool/action name;
- admin/user ID when known;
- status: `SUCCESS`, `DENIED`, or `ERROR`;
- request ID when available;
- redacted details.

Audit logs must redact secrets before writing.

## Error Mapping

MCP tools must map API errors into safe messages.

| Go/API condition | MCP-safe message |
| --- | --- |
| 400 validation failure | `validation failed` with field details when safe |
| 401 missing/invalid token | `unauthorized` |
| 403 forbidden | `forbidden` |
| 404 not found | `not found` |
| 408/timeout | `upstream request timed out` |
| 409 conflict | `conflict` |
| 429 rate limited | `rate limited` |
| 5xx | `upstream service error` |

Do not return raw stack traces, SQL errors, connection strings, tokens, or complete upstream response bodies.

## Confirmation-Required Actions

Tools that perform high-impact actions must require explicit confirmation.

Examples:

- deleting resources;
- suspending or muting accounts;
- rejecting business verification;
- bulk updates;
- sending notifications to many users;
- operations that cannot be safely retried.

Standard parameter:

```python
confirm: bool = False
```

If `confirm` is false, the tool must return a clear error explaining the required confirmation.

## Rate Limits

MCP servers should rely on deployed gateway/service rate limits and may add local process-level controls for expensive tools.

Minimum standard:

- Document expected rate limit source.
- Avoid unbounded loops and unbounded pagination.
- Clamp `limit` parameters.
- Add confirmation for bulk actions.

## Request IDs

MCP servers must propagate `X-Request-ID` to Go APIs. Generated IDs must be UUID strings. Logs should include request IDs when available.

## Timeout And Retry Security

- Keep timeouts finite.
- Do not retry non-idempotent write actions unless idempotency is implemented.
- Retries must not hide authorization failures.
- Never log tokens during retry handling.

## No Raw Database Access

MCP servers must not open database connections or execute SQL directly. New capabilities must be added to the owning Go API first, then exposed through MCP.

## No Secrets In Responses

MCP responses must be checked for common secret fields before returning data from tools/resources.

Blocked field names include:

- `password`
- `token`
- `access_token`
- `refresh_token`
- `authorization`
- `secret`
- `api_key`
- `jwt`
- `cookie`

When a response needs to mention that a credential exists, return a boolean such as `token_present: true` instead of the secret value.
