import argparse
import asyncio

from app.constants import (
    TRANSPORT_PROTOCOL_HTTP,
    TRANSPORT_PROTOCOL_STDIO,
)

from . import resources
from . import tools

from .http_client import guest_client
from .server import mcp


def main() -> None:
    """
    Entry point for the application.
    Parses command line arguments to determine transport mode.
    """
    parser = argparse.ArgumentParser(description="Run the Guest MCP Server")

    _ = parser.add_argument(
        "--transport",
        choices=[
            TRANSPORT_PROTOCOL_STDIO,
            TRANSPORT_PROTOCOL_HTTP,
        ],
        default=TRANSPORT_PROTOCOL_STDIO,
        help=(
            "Transport protocol to use: "
            "'stdio' for local client integration, "
            "'http' for HTTP deployment."
        ),
    )

    args = parser.parse_args()

    try:
        mcp.run(transport=args.transport)
    finally:
        asyncio.run(guest_client.close())


if __name__ == "__main__":
    main()
