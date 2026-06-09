# Admin-Auth MCP Server

Model Context Protocol (MCP) server for admin/auth operations.

The server exposes administrative functionality as MCP **tools** and **resources**, enforcing strict role-based access
control (RBAC) and audit logging for LLM clients interacting with the ShareBite infrastructure.

## Project Structure

```text
admin-auth-server/
├── .env.example
├── audit.log
├── pyproject.toml
├── README.md
├── main.py
├── tests/
│   └── test_server.py
└── mcp_app/
    ├── audit.py
    ├── auth.py
    ├── config.py
    ├── constants.py
    ├── http_client.py
    ├── resources.py
    ├── server.py
    └── tools.py
```

## Architecture

- *Tools* — executable administrative operations exposed to LLM clients, protected by strict role requirements

- *Resources* — read-only contextual information (e.g., schemas and permission maps)

- *HTTP Client* — shared async API client with connection pooling, automatic silent token rotation, and graceful
  shutdown

- *Security & RBAC* — strict admin or moderator clearance validation applied via decorators before execution

- *Audit Logger* — immutable security tracking that automatically redacts sensitive data (tokens, passwords) using regex
  filters

## Prerequisites

- Python 3.12+
- uv
- Running Admin API instance

## Configuration

Create a local environment file:

```bash
cp .env.example .env
```

Example configuration:

```env
ADMIN_AUTH_API_BASE_URL=http://localhost:3850
ADMIN_API_AUTH_REFRESH_TOKEN=
REQUEST_TIMEOUT=10
```

### Environment Variables

| Variable                       | Required | Description                                              |
|--------------------------------|----------|----------------------------------------------------------|
| `ADMIN_AUTH_API_BASE_URL`      | Yes      | Base URL of the Admin Auth GO API                        |
| `TIMEOUT_SECONDS`              | No       | HTTP timeout in seconds (default: 10)                    |
| `ADMIN_API_AUTH_REFRESH_TOKEN` | No       | Long-lived token used to silently refresh session tokens |

## Installation

Install dependencies:

```bash
uv sync
```

## Running

### Stdio Transport

Recommended for local LLM desktop clients such as Claude Desktop.

```bash
uv run main.py --transport stdio
```

### HTTP Transport

Recommended for remote deployments.

```bash
uv run main.py --transport http
```

## Testing

Use MCP Inspector to debug and test tools/resources locally without consuming LLM API credits.

```bash
npx @modelcontextprotocol/inspector uv run main.py --transport stdio
```

Once started:

1. Open the generated localhost URL in your browser
2. Add `ADMIN_AUTH_API_BASE_URL` in **Environment Variables**
3. Click **Connect**

Optional authentication can be provided via:

```env
ADMIN_API_AUTH_REFRESH_TOKEN=your_token_here
```

## Claude Desktop Integration

Open:

**File → Settings → Developer → Edit Config**

(or **Claude → Settings...** on macOS)

Example configuration:

```json
{
  "mcpServers": {
    "admin-auth-server": {
      "command": "uv",
      "args": [
        "--directory",
        "/absolute/path/to/share-bite/mcp/admin-auth-server",
        "run",
        "python",
        "-m",
        "main",
        "--transport",
        "stdio"
      ],
      "env": {
        "ADMIN_AUTH_API_BASE_URL": "http://localhost:3850",
        "REQUEST_TIMEOUT": "10",
        "ADMIN_API_AUTH_REFRESH_TOKEN": ""
      }
    }
  }
}
```

### Windows + WSL

If Claude Desktop runs on Windows while the project lives inside WSL, use `wsl.exe` as a bridge:

```json
{
  "mcpServers": {
    "admin-auth-server": {
      "command": "wsl.exe",
      "args": [
        "bash",
        "-c",
        "cd /absolute/path/to/share-bite/mcp/admin-auth-server && uv run python -m main --transport stdio"
      ],
      "env": {
        "ADMIN_AUTH_API_BASE_URL": "http://localhost:3850",
        "REQUEST_TIMEOUT": "10",
        "ADMIN_API_AUTH_REFRESH_TOKEN": "",
        "WSLENV": "ADMIN_AUTH_API_BASE_URL:REQUEST_TIMEOUT:ADMIN_API_AUTH_REFRESH_TOKEN"
      }
    }
  }
}
```

## Security Model & Limitations

### RBAC Matrix
- **Admin**: Full access to all diagnostic, context, and permission tools.
- **Moderator**: Access limited to health checks and self-context retrieval.

### Local Development Limitations
- **Token Persistence**: When a refresh token is rotated, the server rewrites the local `.env` file directly. This design is strictly intended for single-user local development and local stdio/WSL environments. It is not safe for shared multi-tenant multi-user remote production environments.

Restart Claude Desktop after updating the configuration.
