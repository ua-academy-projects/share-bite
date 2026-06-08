from typing import Any
import httpx


def build_success_response(data: Any, status: int = 200) -> httpx.Response:
    return httpx.Response(status, json=data)


def build_error_response(status: int, error: str) -> httpx.Response:
    return httpx.Response(status, json={"error": error})
