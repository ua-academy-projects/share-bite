from typing import Any
import httpx


def build_success_response(data: Any, status: int = 200) -> httpx.Response:
    return httpx.Response(status, json=data)


def build_error_response(status: int, error: str) -> httpx.Response:
    return httpx.Response(status, json={"error": error})


def build_health_response() -> httpx.Response:
    return build_success_response({"status": "OK"})


def build_status_response(
    app: str = "share-bite",
    db: str = "connected",
    redis: str = "connected",
    status: int = 200,
) -> httpx.Response:
    return build_success_response(
        {"app": app, "database": db, "redis": redis},
        status=status,
    )


def build_info_response(
    version: str = "1.0.0",
    commit: str = "abc123",
    build_time: str = "2026-01-01T00:00:00Z",
    environment: str = "dev",
) -> httpx.Response:
    return build_success_response(
        {
            "version": version,
            "commit": commit,
            "buildTime": build_time,
            "environment": environment,
        }
    )


def build_openapi_response() -> httpx.Response:
    return build_success_response(
        {
            "swagger": "2.0",
            "info": {"title": "Share Bite - Guest Service API", "version": "1.0"},
            "paths": {},
        }
    )
