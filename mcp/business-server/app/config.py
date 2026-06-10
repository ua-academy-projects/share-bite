from dataclasses import dataclass
import os


@dataclass(frozen=True)
class Settings:
    """ Initial settings """
    business_api_base_url: str
    request_timeout_seconds: float
    transport: str
    host: str
    port: int
    path: str


def load_settings() -> Settings:
    """ Loading server configuration from .env config """
    return Settings(
        business_api_base_url=_required_env("BUSINESS_API_BASE_URL").rstrip("/"),
        request_timeout_seconds=float(os.getenv("BUSINESS_API_REQUEST_TIMEOUT_SECONDS", "10")),
        transport=os.getenv("MCP_TRANSPORT", "stdio"),
        host=os.getenv("MCP_HOST", "127.0.0.1"),
        port=int(os.getenv("MCP_PORT", "8000")),
        path=os.getenv("MCP_PATH", "/mcp"),
    )


def _required_env(name: str) -> str:
    """ Extracts enviromental variables """
    value = os.getenv(name)
    if not value:
        raise RuntimeError(f"{name} is required")
    return value