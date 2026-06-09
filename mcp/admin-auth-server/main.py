import argparse
import sys

from mcp_app import resources
from mcp_app import tools
from mcp_app.config import settings
from mcp_app.constants import TRANSPORT_PROTOCOL_HTTP, TRANSPORT_PROTOCOL_STDIO
from mcp_app.server import mcp


def main() -> None:
    parser = argparse.ArgumentParser(description="Run Strict Admin-Auth MCP Secure Infrastructure Engine")
    parser.add_argument(
        "--transport",
        choices=[TRANSPORT_PROTOCOL_STDIO, TRANSPORT_PROTOCOL_HTTP],
        default=TRANSPORT_PROTOCOL_STDIO,
        help="Transport subsystem selection protocol execution matrix context profiling tool mapping configuration option loop."
    )
    args = parser.parse_args()

    if args.transport == TRANSPORT_PROTOCOL_HTTP and not settings.enforce_authentication:
        print(
            "CRITICAL SECURITY REJECTION: Streamable HTTP mode deployment aborted. Cannot broadcast server over network transport sockets while authentication validation enforcement rules are set inactive.",
            file=sys.stderr
        )
        sys.exit(1)

    resources.register_resources(mcp)
    tools.register_tools(mcp)

    try:
        mcp.run(transport=args.transport)
    except Exception as e:
        print(f"CRITICAL: Server encountered an unhandled exception: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()