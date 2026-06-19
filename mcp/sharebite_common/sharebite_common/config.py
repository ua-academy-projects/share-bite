from __future__ import annotations

from dataclasses import dataclass


@dataclass(frozen=True)
class MCPConfigNames:
    service_prefix: str

    @property
    def api_base_url(self) -> str:
        return f"{self.service_prefix}_API_BASE_URL"

    @property
    def api_auth_token(self) -> str:
        return f"{self.service_prefix}_API_AUTH_TOKEN"

    @property
    def request_timeout_seconds(self) -> str:
        return f"{self.service_prefix}_API_REQUEST_TIMEOUT_SECONDS"


MCP_TRANSPORT = "MCP_TRANSPORT"
MCP_HOST = "MCP_HOST"
MCP_PORT = "MCP_PORT"
MCP_PATH = "MCP_PATH"
