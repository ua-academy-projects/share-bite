from typing import Any
import httpx


def build_success_response(data: Any, status: int = 200) -> httpx.Response:
    return httpx.Response(status, json=data)


def build_error_response(status: int, message: str) -> httpx.Response:
    """Guest API error format: {"message": "..."}."""
    return httpx.Response(status, json={"message": message})


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


def build_posts_list_response(
    posts: list[dict[str, Any]] | None = None,
    total: int = 0,
    status: int = 200,
) -> httpx.Response:
    return build_success_response(
        {"posts": posts or [], "total": total},
        status=status,
    )


def build_post_detail_response(
    post: dict[str, Any] | None = None,
    status: int = 200,
) -> httpx.Response:
    return build_success_response(
        {"post": post or {}},
        status=status,
    )


def build_post_authors_response(
    authors: list[str] | None = None,
    count: int = 0,
    status: int = 200,
) -> httpx.Response:
    return build_success_response(
        {"authors": authors or [], "count": count},
        status=status,
    )


def build_customer_detail_response(
    customer: dict[str, Any] | None = None,
    status: int = 200,
) -> httpx.Response:
    return build_success_response(
        {"customer": customer or {}},
        status=status,
    )


def build_followers_list_response(
    customers: list[dict[str, Any]] | None = None,
    next_page_token: str = "",
    status: int = 200,
) -> httpx.Response:
    return build_success_response(
        {"customers": customers or [], "next_page_token": next_page_token},
        status=status,
    )


def build_collections_list_response(
    collections: list[dict[str, Any]] | None = None,
    next_page_token: str = "",
    status: int = 200,
) -> httpx.Response:
    return build_success_response(
        {"collections": collections or [], "nextPageToken": next_page_token},
        status=status,
    )


def build_collection_detail_response(
    collection: dict[str, Any] | None = None,
    status: int = 200,
) -> httpx.Response:
    return build_success_response(
        {"collection": collection or {}},
        status=status,
    )


def build_collection_venues_response(
    venues: list[dict[str, Any]] | None = None,
    status: int = 200,
) -> httpx.Response:
    return build_success_response(
        {"venues": venues or []},
        status=status,
    )
