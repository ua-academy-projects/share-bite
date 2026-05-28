from __future__ import annotations

import json
import logging
import sys
from datetime import datetime, timezone
from typing import Any

from mcp.server.fastmcp import FastMCP

from app.config import Settings, load_settings
from app.resources import register_resources
from app.tools import register_tools


SUPPORTED_TRANSPORTS = {"stdio", "streamable-http"}


class JsonFormatter(logging.Formatter):
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
    handler = logging.StreamHandler(sys.stderr)
    handler.setFormatter(JsonFormatter())

    root_logger = logging.getLogger()
    root_logger.handlers.clear()
    root_logger.addHandler(handler)
    root_logger.setLevel(logging.INFO)


def create_server(settings: Settings) -> FastMCP:
    
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
    )

    register_tools(mcp, settings)
    register_resources(mcp, settings)

    return mcp


def main() -> None:
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

    mcp = create_server(settings)
    mcp.run(transport=settings.transport)


if __name__ == "__main__":
    main()