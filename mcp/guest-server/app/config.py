from typing import ClassVar
from pydantic import AnyHttpUrl, Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    guest_api_base_url: AnyHttpUrl = Field(
        ..., description="Base URL of the Guest Go API"
    )
    timeout_seconds: int = Field(
        default=10, description="Global timeout for HTTP requests to Go API"
    )
    auth_token: str | None = Field(
        default=None, description="Admin JWT token for Guest API"
    )

    model_config: ClassVar[SettingsConfigDict] = SettingsConfigDict(
        env_file=".env", env_file_encoding="utf-8", extra="ignore"
    )


settings = Settings()
