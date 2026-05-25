import logging
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
    error_message: str


APIResponse = APISuccessResponse | APIErrorResponse


class GuestAPIClient:
    def __init__(self) -> None:
        self.base_url: str = str(settings.guest_api_base_url).rstrip("/")
        self.timeout: int = settings.timeout_seconds

    def _build_headers(
        self,
        auth_token: str | None = None,
        request_id: str | None = None,
    ) -> dict[str, str]:
        headers: dict[str, str] = {HEADER_CONTENT_TYPE: CONTENT_TYPE_JSON}
        if auth_token:
            headers[HEADER_AUTH] = f"Bearer {auth_token}"
        if request_id:
            headers[HEADER_REQUEST_ID] = request_id
        return headers

    async def _request(
        self,
        method: str,
        path: str,
        auth_token: str | None = None,
        request_id: str | None = None,
        json_data: dict[str, Any] | None = None,
        params: dict[str, Any] | None = None,
    ) -> APIResponse:
        url = f"{self.base_url}/{path.lstrip('/')}"
        headers = self._build_headers(auth_token, request_id)

        async with httpx.AsyncClient(timeout=self.timeout) as client:
            try:
                response = await client.request(
                    method, url, headers=headers, json=json_data, params=params
                )

                _ = response.raise_for_status()
                data = response.json() if response.content else None
                return APISuccessResponse(
                    is_error=False, status=response.status_code, data=data
                )
            except httpx.HTTPStatusError as e:
                logger.error(
                    "HTTPStatusError %s for %s: %s",
                    e.response.status_code,
                    url,
                    e.response.text,
                )
                return APIErrorResponse(
                    is_error=True,
                    error_message=f"Guest API error ({e.response.status_code}): {e.response.text}",
                )
            except httpx.TimeoutException:
                logger.error("TimeoutException for %s", url)
                return APIErrorResponse(
                    is_error=True,
                    error_message=f"Request timed out after {self.timeout}s",
                )
            except httpx.RequestError as e:
                logger.error("RequestError for %s: %s", url, e)
                return APIErrorResponse(
                    is_error=True,
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
        json_data: dict[str, Any],
        auth_token: str | None = None,
        request_id: str | None = None,
    ) -> APIResponse:
        return await self._request(
            "POST",
            path=path,
            auth_token=auth_token,
            request_id=request_id,
            json_data=json_data,
        )


guest_client = GuestAPIClient()
