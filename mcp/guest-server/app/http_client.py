import logging
import uuid
import httpx
from typing import Any, Literal, TypedDict

from .config import settings
from .constants import (
    CONTENT_TYPE_JSON,
    HEADER_AUTH,
    HEADER_CONTENT_TYPE,
    HEADER_REQUEST_ID,
)

logger = logging.getLogger(__name__)


class APISuccessResponse(TypedDict):
    is_error: Literal[False]
    status: int
    data: Any


class APIErrorResponse(TypedDict):
    is_error: Literal[True]
    status: int | None
    error_message: str


APIResponse = APISuccessResponse | APIErrorResponse


class GuestAPIClient:
    """
    A robust asynchronous HTTP client tailored for the Guest API.
    Handles authorization, timeouts, safe error parsing, and both JSON/Multipart payloads.
    """

    def __init__(self) -> None:
        self.base_url: str = str(settings.guest_api_base_url).rstrip("/")
        self.timeout: int = settings.timeout_seconds

        self._client: httpx.AsyncClient = httpx.AsyncClient(
            timeout=self.timeout, base_url=self.base_url
        )

    def _build_headers(
        self,
        auth_token: str | None = None,
        request_id: str | None = None,
        is_multipart: bool = False,
    ) -> dict[str, str]:
        headers: dict[str, str] = {}

        # For JSON requests, we must specify the Content-Type.
        # For multipart/form-data (e.g., file uploads), we intentionally omit the Content-Type
        # header here so that httpx can automatically generate it with the correct boundary string.
        if not is_multipart:
            headers[HEADER_CONTENT_TYPE] = CONTENT_TYPE_JSON

        if auth_token:
            headers[HEADER_AUTH] = f"Bearer {auth_token}"
        if request_id:
            headers[HEADER_REQUEST_ID] = request_id

        return headers

    async def close(self) -> None:
        await self._client.aclose()

    async def _request(
        self,
        method: str,
        path: str,
        auth_token: str | None = None,
        request_id: str | None = None,
        json_data: dict[str, Any] | None = None,
        data: (
            dict[str, Any] | None
        ) = None,  # Used for form fields in multipart requests
        files: (
            dict[str, Any] | None
        ) = None,  # Used for file uploads (e.g., avatars, post images)
        params: dict[str, Any] | None = None,
    ) -> APIResponse:
        is_multipart = bool(data or files)
        headers = self._build_headers(
            auth_token,
            request_id or str(uuid.uuid4()),
            is_multipart,
        )

        try:
            response = await self._client.request(
                method,
                path,
                headers=headers,
                json=json_data,
                data=data,
                files=files,
                params=params,
            )

            _ = response.raise_for_status()
            if not response.content:
                response_data = None
            else:
                try:
                    response_data = response.json()
                except ValueError:
                    response_data = response.text

            return APISuccessResponse(
                is_error=False, status=response.status_code, data=response_data
            )
        except httpx.HTTPStatusError as e:
            try:
                error_payload = e.response.json()
                parsed_error = (
                    error_payload.get("message")
                    or error_payload.get("error")
                    or "Unknown API Error"
                )
            except ValueError:
                parsed_error = "Non-JSON response received (possible server fault)"

            logger.error(
                "HTTPStatusError %s for %s: %s",
                e.response.status_code,
                path,
                parsed_error,
            )

            return APIErrorResponse(
                is_error=True,
                status=e.response.status_code,
                error_message=parsed_error,
            )
        except httpx.TimeoutException:
            logger.error("TimeoutException for %s", path)
            return APIErrorResponse(
                is_error=True,
                status=408,
                error_message=f"Request timed out after {self.timeout}s",
            )
        except httpx.RequestError as e:
            logger.error("RequestError for %s: %s", path, e)
            return APIErrorResponse(
                is_error=True,
                status=None,
                error_message=f"Failed to connect to Guest API: {e}",
            )

    async def get(
        self,
        path: str,
        auth_token: str | None = None,
        request_id: str | None = None,
        params: dict[str, Any] | None = None,
    ) -> APIResponse:
        return await self._request(
            "GET",
            path=path,
            auth_token=auth_token,
            request_id=request_id,
            params=params,
        )

    async def post(
        self,
        path: str,
        json_data: dict[str, Any] | None = None,
        data: dict[str, Any] | None = None,
        files: dict[str, Any] | None = None,
        auth_token: str | None = None,
        request_id: str | None = None,
    ) -> APIResponse:
        return await self._request(
            "POST",
            path=path,
            auth_token=auth_token,
            request_id=request_id,
            json_data=json_data,
            data=data,
            files=files,
        )

    async def patch(
        self,
        path: str,
        json_data: dict[str, Any] | None = None,
        data: dict[str, Any] | None = None,
        files: dict[str, Any] | None = None,
        auth_token: str | None = None,
        request_id: str | None = None,
    ) -> APIResponse:
        return await self._request(
            "PATCH",
            path=path,
            auth_token=auth_token,
            request_id=request_id,
            json_data=json_data,
            data=data,
            files=files,
        )

    async def delete(
        self,
        path: str,
        auth_token: str | None = None,
        request_id: str | None = None,
        params: dict[str, Any] | None = None,
    ) -> APIResponse:
        return await self._request(
            "DELETE",
            path=path,
            auth_token=auth_token,
            request_id=request_id,
            params=params,
        )


guest_client = GuestAPIClient()
