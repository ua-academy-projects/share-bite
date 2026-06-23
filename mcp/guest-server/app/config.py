from typing import ClassVar
from pydantic import AliasChoices, AnyHttpUrl, Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    guest_api_base_url: AnyHttpUrl = Field(
        ..., description="Base URL of the Guest Go API"
    )

    guest_api_auth_token: str | None = Field(
        default=None, description="Optional JWT token for protected Guest API endpoints"
    )

    timeout_seconds: int = Field(
        default=10,
        validation_alias=AliasChoices("GUEST_API_REQUEST_TIMEOUT_SECONDS", "TIMEOUT_SECONDS"),
        description="Global timeout for HTTP requests to Go API",
    )

    model_config: ClassVar[SettingsConfigDict] = SettingsConfigDict(
        env_file=".env", env_file_encoding="utf-8", extra="ignore"
    )


settings = Settings()
