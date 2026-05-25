import argparse
import sys

from app.constants import TRANSPORT_PROTOCOL_SSE, TRANSPORT_PROTOCOL_STDIO

from .server import mcp

from . import resources
from . import tools


def main() -> None:
    """
    Entry point for the application. Parses command line arguments to determine
    the transport mode (stdio or sse).
    """
    parser = argparse.ArgumentParser(description="Run the Guest MCP Server")
    _ = parser.add_argument(
        "--transport",
        choices=[TRANSPORT_PROTOCOL_STDIO, TRANSPORT_PROTOCOL_SSE],
        default="stdio",
        help="Transport protocol to use: 'stdio' for local client integration, 'sse' for HTTP deployment.",
    )
    args, _ = parser.parse_known_args()

    if args.transport == TRANSPORT_PROTOCOL_STDIO:
        # Run in standard input/output mode
        mcp.run(transport=TRANSPORT_PROTOCOL_STDIO)
    elif args.transport == TRANSPORT_PROTOCOL_SSE:
        # Run as a web server with Server-Sent Events
        # FastMCP automatically handles the Uvicorn/FastAPI startup under the hood
        mcp.run(transport=TRANSPORT_PROTOCOL_SSE)
    else:
        print("Unknown transport specified.", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
