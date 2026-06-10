from __future__ import annotations

from typing import Any
from uuid import uuid4

import httpx


class BusinessApiError(Exception):
    def __init__(self, message: str, status_code: int | None = None) -> None:
        super().__init__(message)
        self.status_code = status_code


class BusinessApiClient:
    """Asynchronous HTTP client for interacting with the Business API."""
    def __init__(self, base_url: str, timeout_seconds: float) -> None:
        self._base_url = base_url.rstrip("/")
        self._timeout = httpx.Timeout(timeout_seconds)
        self._client = httpx.AsyncClient(base_url=self._base_url, timeout=self._timeout)

    async def close(self):
        await self._client.aclose()

    async def get(
        self,
        path: str,
        *,
        auth_token: str | None = None,
        request_id: str | None = None,
        params: dict[str, Any] | None = None,
    ) -> dict[str, Any]:
        """Performs a GET request to the specified API path."""
        headers = self._build_headers(auth_token=auth_token, request_id=request_id)

        try:
            async with httpx.AsyncClient(base_url=self._base_url, timeout=self._timeout) as client:
                response = await client.get(path, headers=headers, params=params)
        except httpx.TimeoutException as exc:
            raise BusinessApiError("business-api request timed out") from exc
        except httpx.HTTPError as exc:
            raise BusinessApiError(f"business-api request failed: {exc}") from exc

        return self._parse_response(response)

    async def patch(
        self,
        path: str,
        *,
        json_data: dict[str, Any],
        auth_token: str | None = None,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        headers = self._build_headers(auth_token=auth_token, request_id=request_id)

        try:
            async with httpx.AsyncClient(base_url=self._base_url, timeout=self._timeout) as client:
                response = await client.patch(path, headers=headers, json=json_data)
        except httpx.TimeoutException as exc:
            raise BusinessApiError("business-api request timed out") from exc
        except httpx.HTTPError as exc:
            raise BusinessApiError(f"business-api request failed: {exc}") from exc

        return self._parse_response(response)

    def _build_headers(
        self,
        *,
        auth_token: str | None,
        request_id: str | None,
    ) -> dict[str, str]:
        """Header builder"""
        headers = {
            "Accept": "application/json",
            "X-Request-ID": request_id or str(uuid4()),
        }

        if auth_token:
            headers["Authorization"] = _normalize_bearer_token(auth_token)

        return headers

    def _parse_response(self, response: httpx.Response) -> dict[str, Any]:
        """Function to parse server response"""
        if response.status_code >= 400:
            raise BusinessApiError(
                _extract_error_message(response),
                status_code=response.status_code,
            )

        if not response.content:
            return {}

        try:
            parsed = response.json()
        except ValueError as exc:
            raise BusinessApiError("business-api returned non-JSON response") from exc

        if not isinstance(parsed, dict):
            return {"data": parsed}

        return parsed


def _normalize_bearer_token(auth_token: str) -> str:
    if auth_token.lower().startswith("bearer "):
        return auth_token
    return f"Bearer {auth_token}"


def _extract_error_message(response: httpx.Response) -> str:
    """Extract error message to insert it into response"""
    try:
        body = response.json()
    except ValueError:
        body_text = response.text.strip()
        if body_text:
            return f"business-api returned {response.status_code}: {body_text}"
        return f"business-api returned {response.status_code}"

    if isinstance(body, dict):
        error = body.get("error") or body.get("message")
        if error:
            return f"business-api returned {response.status_code}: {error}"

    return f"business-api returned {response.status_code}"