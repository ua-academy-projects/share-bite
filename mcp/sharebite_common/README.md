# ShareBite MCP Common

Shared Python helpers for Share Bite MCP servers.

This package supports the standards documented in:

- `docs/mcp/architecture.md`
- `docs/mcp/security.md`

## Included helpers

- `auth`: bearer token extraction and normalization
- `security`: secret redaction, no-secret assertions, confirmation checks
- `errors`: safe HTTP-to-MCP error mapping
- `pagination`: offset pagination clamping
- `request_id`: request ID generation
- `audit`: standard audit event shape with redacted details
- `testing`: fake MCP/context utilities and response security assertions

## Usage direction

New MCP tools should prefer these helpers instead of duplicating local behavior.
Existing servers may migrate incrementally while keeping backward-compatible environment variable aliases.
