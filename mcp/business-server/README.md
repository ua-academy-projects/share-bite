# ShareBite Business MCP Server

Python-based MCP server for ShareBite business-owner operations.

The server wraps safe `business-api` workflows and exposes them as MCP tools and resources.

## Capabilities

Tools:

- `business_health_check`
- `get_business_api_status`

Resources:

- `sharebite://business/api-info`
- `sharebite://business/openapi-summary`

## Configuration

Required:

| Variable | Description | Example |
| --- | --- | --- |
| `BUSINESS_API_BASE_URL` | Base URL for `business-api` | `http://localhost:3900` |

Optional:

| Variable | Description | Default |
| --- | --- | --- |
| `BUSINESS_API_REQUEST_TIMEOUT_SECONDS` | Timeout for requests to `business-api` | `10` |
| `MCP_TRANSPORT` | MCP transport: `stdio` or `streamable-http` | `stdio` |
| `MCP_HOST` | Host for Streamable HTTP mode | `127.0.0.1` |
| `MCP_PORT` | Port for Streamable HTTP mode | `8000` |
| `MCP_PATH` | Streamable HTTP path | `/mcp` |

## Auth And Request Context

Tools forward `Authorization` to `business-api` when an `auth_token` argument is provided.

Tools propagate `X-Request-ID` when a `request_id` argument is provided. If no request ID is provided, the MCP server generates one.

Business IDs must never be guessed. Any tool that requires a business ID must receive it from authenticated context or explicit input.

## Local Setup

From this directory:

```powershell
cd mcp\business-server
py -m venv .venv
.\.venv\Scripts\python.exe -m pip install -e .
```

Set the required environment variable:

```powershell
$env:BUSINESS_API_BASE_URL="http://localhost:3900"
```

## Run With Stdio

```powershell
$env:MCP_TRANSPORT="stdio"
.\.venv\Scripts\python.exe -m app.main
```

Stdio mode is intended to be started by an MCP client.

## Run With Streamable HTTP

```powershell
$env:BUSINESS_API_BASE_URL="http://localhost:3900"
$env:MCP_TRANSPORT="streamable-http"
$env:MCP_HOST="127.0.0.1"
$env:MCP_PORT="8000"
$env:MCP_PATH="/mcp"
.\.venv\Scripts\python.exe -m app.main
```

The MCP endpoint is:

```text
http://127.0.0.1:8000/mcp
```

## Local Testing

Start `business-api` first. It should expose OpenAPI JSON at:

```text
http://localhost:3900/swagger/doc.json
```

Then verify that an MCP client can list tools and resources:

```powershell
$env:BUSINESS_API_BASE_URL="http://localhost:3900"

@'
import asyncio
import os

from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client


async def main():
    server_params = StdioServerParameters(
        command=os.path.abspath(".venv/Scripts/python.exe"),
        args=["-m", "app.main"],
        env={
            **os.environ,
            "BUSINESS_API_BASE_URL": "http://localhost:3900",
            "MCP_TRANSPORT": "stdio",
        },
    )

    async with stdio_client(server_params) as (read, write):
        async with ClientSession(read, write) as session:
            await session.initialize()

            tools = await session.list_tools()
            print("TOOLS:")
            for tool in tools.tools:
                print("-", tool.name)

            resources = await session.list_resources()
            print("RESOURCES:")
            for resource in resources.resources:
                print("-", resource.uri)


asyncio.run(main())
'@ | .\.venv\Scripts\python.exe -
```

Expected tools:

```text
business_health_check
get_business_api_status
```

Expected resources:

```text
sharebite://business/api-info
sharebite://business/openapi-summary
```

Verify that health and status calls reach `business-api`:

```powershell
@'
import asyncio
import os

from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client


async def main():
    server_params = StdioServerParameters(
        command=os.path.abspath(".venv/Scripts/python.exe"),
        args=["-m", "app.main"],
        env={
            **os.environ,
            "BUSINESS_API_BASE_URL": "http://localhost:3900",
            "MCP_TRANSPORT": "stdio",
        },
    )

    async with stdio_client(server_params) as (read, write):
        async with ClientSession(read, write) as session:
            await session.initialize()

            health = await session.call_tool("business_health_check", {})
            print("HEALTH:")
            print(health)

            status = await session.call_tool("get_business_api_status", {})
            print("STATUS:")
            print(status)


asyncio.run(main())
'@ | .\.venv\Scripts\python.exe -
```

A successful response includes:

```text
ShareBite Business API
```

## Error Handling

Errors returned by `business-api` are converted into clear MCP errors.

Timeout and connection failures are reported as:

```text
business-api request timed out
business-api request failed
```

HTTP error responses include the status code and the error message returned by `business-api` when available.
