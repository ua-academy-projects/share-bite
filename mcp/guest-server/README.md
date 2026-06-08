# Guest API MCP Server

Model Context Protocol (MCP) server for the Guest API built with **FastMCP**.

The server exposes selected Guest API functionality as MCP **tools** and **resources**, allowing LLM clients to interact with the Guest ecosystem through a structured interface.

## Project Structure

```text
guest-server/
├── .env.example
├── pyproject.toml
├── README.md
└── app/
    ├── auth.py
    ├── config.py
    ├── constants.py
    ├── http_client.py
    ├── main.py
    ├── server.py
    ├── resources/
    │   ├── __init__.py
    │   └── api.py
    └── tools/
        ├── __init__.py
        └── health.py
```

## Architecture

- **Tools** — executable operations exposed to LLM clients
- **Resources** — read-only contextual information
- **HTTP Client** — shared async Guest API client with connection pooling, centralized auth handling, and graceful shutdown

## Prerequisites

- Python 3.12+
- uv
- Running Guest API instance

## Configuration

Create a local environment file:

```bash
cp .env.example .env
```

Example configuration:

```env
GUEST_API_BASE_URL=http://localhost:3800
TIMEOUT_SECONDS=10
GUEST_API_AUTH_TOKEN=
```

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `GUEST_API_BASE_URL` | Yes | Base URL of the Guest API |
| `TIMEOUT_SECONDS` | No | HTTP timeout in seconds |
| `GUEST_API_AUTH_TOKEN` | No | JWT token for protected endpoints |

`GUEST_API_AUTH_TOKEN` is optional. Public endpoints can be accessed without authentication.

## Installation

Install dependencies:

```bash
uv sync
```

## Running

### Stdio Transport

Recommended for local LLM desktop clients such as Claude Desktop.

```bash
uv run python -m app.main --transport stdio
```

### HTTP Transport

Recommended for remote deployments.

```bash
uv run python -m app.main --transport http
```

## Testing

Use MCP Inspector to debug and test tools/resources locally without consuming LLM API credits.

```bash
npx @modelcontextprotocol/inspector \
uv run python -m app.main --transport stdio
```

Once started:

1. Open the generated localhost URL in your browser
2. Add `GUEST_API_BASE_URL` in **Environment Variables**
3. Click **Connect**

Optional authentication can be provided via:

```env
GUEST_API_AUTH_TOKEN=
```

### Unit Tests

Run all tests:

```bash
uv run pytest tests/ -v
```

Run a specific test file:

```bash
uv run pytest tests/unit/tools/test_health.py -v
```

## Claude Desktop Integration

Open:

**File → Settings → Developer → Edit Config**

(or **Claude → Settings...** on macOS)

Example configuration:

```json
{
  "mcpServers": {
    "guest-server": {
      "command": "uv",
      "args": [
        "run",
        "python",
        "-m",
        "app.main",
        "--transport",
        "stdio"
      ],
      "cwd": "/absolute/path/to/share-bite/mcp/guest-server",
      "env": {
        "GUEST_API_BASE_URL": "http://localhost:3800",
        "TIMEOUT_SECONDS": "10",
        "GUEST_API_AUTH_TOKEN": ""
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
    "guest-server": {
      "command": "wsl.exe",
      "args": [
        "bash",
        "-c",
        "cd /absolute/path/to/share-bite/mcp/guest-server && .venv/bin/python -m app.main --transport stdio"
      ],
      "env": {
        "GUEST_API_BASE_URL": "http://localhost:3800",
        "TIMEOUT_SECONDS": "10",
        "GUEST_API_AUTH_TOKEN": "",
        "WSLENV": "GUEST_API_BASE_URL:GUEST_API_AUTH_TOKEN:TIMEOUT_SECONDS"
      }
    }
  }
}
```

Restart Claude Desktop after updating the configuration.
