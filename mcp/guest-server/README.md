# Guest API MCP Server

A Python-based Model Context Protocol (MCP) server built with **FastMCP**. It wraps the Go-based Guest API observability endpoints and exposes them as native MCP tools and resources.

## Project Structure

```text
guest-server/
├── .env.example         # Example environment variables template
├── pyproject.toml       # Project metadata and dependencies
├── README.md            # Service documentation
└── app/
    ├── config.py         # Pydantic configuration
    ├── constants.py      # Global constants
    ├── http_client.py    # Go API HTTP client
    ├── server.py         # FastMCP instance
    ├── resources.py      # Read-only contexts
    ├── tools.py          # Executable actions
    └── main.py           # Application entrypoint
```

**Scaling:** To expand the server, add new domain-specific modules for tools or resources (e.g., `app/tools_users.py` or `app/resources_posts.py`). Decorate your functions with `@mcp.tool()` or `@mcp.resource()`, and register them by adding a blank import (e.g., `from . import tools_users`) in `main.py`.

## Prerequisites

- [uv](https://docs.astral.sh/uv/)
- Python 3.12+
- Guest Go API running

## Configuration

Copy the provided `.env.example` file to create your local configuration:
```bash
cp .env.example .env
```

Then, define the required variables in your new `.env` file:
```env
GUEST_API_BASE_URL=http://localhost:3800
TIMEOUT_SECONDS=10
AUTH_TOKEN=your_admin_jwt_token_here
```

## Installation & Run

Sync dependencies and create the virtual environment:
```bash
uv sync
```

**Run in Stdio Mode (for LLM Desktop Clients):**
```bash
uv run python -m app.main --transport stdio
```

**Run in SSE Mode (HTTP microservice):**
```bash
uv run python -m app.main --transport sse
```

## Testing (MCP Inspector)

Use the MCP Inspector to locally debug and test available tools and resources without consuming LLM API credits.

```bash
npx @modelcontextprotocol/inspector uv run python -m app.main --transport stdio
```
*Once started, open the provided localhost URL in your browser, add `GUEST_API_BASE_URL` in the Environment Variables tab, and click Connect.*

## Claude Desktop Integration

To connect this MCP server to your local Claude Desktop app, you need to update its configuration file. Open Claude Desktop and navigate to **File > Settings > Developer > Edit Config** (or **Claude > Settings...** on macOS).

Add the following standard configuration, replacing the absolute path with your actual project location:

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
        "AUTH_TOKEN": "your_admin_jwt_token_here",
        "TIMEOUT_SECONDS": "10"
      }
    }
  }
}
```

### Windows WSL Option
If you are running Claude Desktop on Windows but your codebase resides inside WSL (Ubuntu) *(because why choose one operating system?)*, you must bridge the connection using `wsl.exe`. Use this configuration instead:

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
        "AUTH_TOKEN": "your_admin_jwt_token_here",
        "TIMEOUT_SECONDS": "10",
        "WSLENV": "GUEST_API_BASE_URL/u:AUTH_TOKEN/u:TIMEOUT_SECONDS/u"
      }
    }
  }
}
```

*Note: Restart Claude Desktop completely for the new configuration to take effect.*
