from pathlib import Path
from typing import ClassVar

from pydantic import AliasChoices, Field
from pydantic_settings import BaseSettings, SettingsConfigDict

BASE_DIR = Path(__file__).resolve().parent.parent


class Settings(BaseSettings):
    admin_auth_api_base_url: str | None = Field(
        default=None,
        alias="ADMIN_AUTH_API_BASE_URL",
        description="Base URL of the Admin/Auth Go API"
    )

    local_refresh_token: str | None = Field(
        default=None,
        alias="ADMIN_API_AUTH_REFRESH_TOKEN",
        description="Long-lived token to rotate access sessions"
    )

    timeout_seconds: int = Field(
        default=10,
        validation_alias=AliasChoices("ADMIN_AUTH_API_REQUEST_TIMEOUT_SECONDS", "REQUEST_TIMEOUT"),
        description="Global timeout for HTTP requests to Go API",
    )

    audit_log_destination: str = Field(
        default=str(BASE_DIR / "audit.log"),
        validation_alias=AliasChoices("MCP_AUDIT_LOG_PATH", "AUDIT_LOG_DESTINATION"),
        description="Path to the system audit log file",
    )

    enforce_authentication: bool = Field(
        default=True, alias="ENFORCE_AUTHENTICATION", description="Enforce strict authentication protocols"
    )

    model_config: ClassVar[SettingsConfigDict] = SettingsConfigDict(
        env_file=str(BASE_DIR / ".env"),
        env_file_encoding="utf-8",
        extra="ignore",
        populate_by_name=True
    )


settings = Settings()
