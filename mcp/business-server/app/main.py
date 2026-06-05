from __future__ import annotations

import json
import asyncio
import logging
import sys
from contextlib import asynccontextmanager
from datetime import datetime, timezone
from typing import Any

from mcp.server.fastmcp import FastMCP
from app.client import BusinessApiClient

from app.config import Settings, load_settings
from app.resources import register_resources
from app.tools import register_tools


SUPPORTED_TRANSPORTS = {"stdio", "streamable-http"}


class JsonFormatter(logging.Formatter):
    """ Class to formate data to json """
    def format(self, record: logging.LogRecord) -> str:
        payload: dict[str, Any] = {
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "level": record.levelname,
            "logger": record.name,
            "message": record.getMessage(),
        }

        if record.exc_info:
            payload["exception"] = self.formatException(record.exc_info)

        return json.dumps(payload, ensure_ascii=False)


def configure_logging() -> None:
    """ Configuring logger setup """
    handler = logging.StreamHandler(sys.stderr)
    handler.setFormatter(JsonFormatter())

    root_logger = logging.getLogger()
    root_logger.handlers.clear()
    root_logger.addHandler(handler)
    root_logger.setLevel(logging.INFO)


def create_server(settings: Settings, client: BusinessApiClient) -> FastMCP:
    """ Initializing MCP server """
    @asynccontextmanager
    async def server_lifespan(ctx: FastMCP):
        yield
        logging.info("Closing HTTP client connection pool...")
        await client.close()
    
    mcp = FastMCP(
        "sharebite-business-mcp-server",
        instructions=(
            "MCP server for ShareBite business-owner operations. "
            "Business IDs must never be guessed; tools must use authenticated context "
            "or explicit input when a business ID is required."
        ),
        host=settings.host,
        port=settings.port,
        streamable_http_path=settings.path,
        stateless_http=True,
        json_response=True,
        lifespan=server_lifespan,
    )

    register_tools(mcp, settings,client)
    register_resources(mcp, settings, client)

    return mcp


def main() -> None:
    """Entry point for the ShareBite Business MCP server."""
    configure_logging()

    settings = load_settings()
    if settings.transport not in SUPPORTED_TRANSPORTS:
        raise RuntimeError(
            f"Unsupported MCP_TRANSPORT={settings.transport}. "
            "Use 'stdio' or 'streamable-http'."
        )

    logging.info(
        "starting MCP business server",
        extra={
            "transport": settings.transport,
            "business_api_base_url": settings.business_api_base_url,
        },
    )

    client = BusinessApiClient(
        base_url=settings.business_api_base_url,
        timeout_seconds=settings.request_timeout_seconds,
    )

    mcp = create_server(settings, client)
    mcp.run(transport=settings.transport)


if __name__ == "__main__":
    main()
