import httpx
import logging
from typing import Any, Dict
from mcp_app.config import settings, BASE_DIR

logger = logging.getLogger("mcp_server")

def write_new_refresh_to_env(new_token: str):
    env_path = BASE_DIR / ".env"
    if not env_path.exists():
        return
    lines = env_path.read_text(encoding="utf-8").splitlines()
    updated = False
    with open(env_path, "w", encoding="utf-8") as f:
        for line in lines:
            if line.startswith("ADMIN_API_AUTH_REFRESH_TOKEN="):
                f.write(f"ADMIN_API_AUTH_REFRESH_TOKEN={new_token}\n")
                updated = True
            else:
                f.write(f"{line}\n")
        if not updated:
            f.write(f"ADMIN_API_AUTH_REFRESH_TOKEN={new_token}\n")

class AdminHttpClient:
    def __init__(self):
        if not settings.admin_auth_api_base_url:
            raise RuntimeError("CRITICAL: ADMIN_AUTH_API_BASE_URL is not set.")
        self.base_url = str(settings.admin_auth_api_base_url).rstrip("/")
        self.timeout = settings.timeout_seconds
        self._access_token: str | None = None

    async def _refresh_access_token(self) -> str:
        if not settings.local_refresh_token:
            raise ValueError("No refresh token stored in settings/env.")

        async with httpx.AsyncClient(timeout=self.timeout) as client:
            response = await client.post(
                f"{self.base_url}/auth/refresh",
                json={"refresh_token": settings.local_refresh_token}
            )
            if response.status_code != 200:
                raise ValueError("Refresh token expired or invalid on Go backend.")

            data = response.json()
            self._access_token = data.get("access_token")
            new_refresh = data.get("refresh_token")

            if new_refresh:
                settings.local_refresh_token = new_refresh
                write_new_refresh_to_env(new_refresh)

            return self._access_token

    async def get_token(self, explicit_token: str | None = None) -> str:
        if explicit_token:
            return explicit_token
        if self._access_token:
            return self._access_token
        return await self._refresh_access_token()

    async def request(self, method: str, path: str, auth_token: str | None = None, **kwargs) -> Dict[str, Any]:
        try:
            token = await self.get_token(auth_token)
        except Exception as e:
            return {"is_error": True, "error_message": f"Authentication loop failure: {str(e)}"}

        headers = kwargs.pop("headers", {})
        headers["Authorization"] = f"Bearer {token}"
        headers["Content-Type"] = "application/json"

        async with httpx.AsyncClient(timeout=self.timeout) as client:
            url = f"{self.base_url}{path}"
            try:
                response = await client.request(method, url, headers=headers, **kwargs)
                if response.status_code == 401:
                    logger.warning("Access token expired dynamically. Retrying with fresh refresh token...")
                    try:
                        new_token = await self._refresh_access_token()
                        headers["Authorization"] = f"Bearer {new_token}"
                        response = await client.request(method, url, headers=headers, **kwargs)
                    except Exception as e:
                        return {"is_error": True, "error_message": f"Session completely expired: {str(e)}"}

                if response.status_code >= 400:
                    return {"is_error": True, "error_message": f"Go API Error {response.status_code}"}

                return {"is_error": False, "data": response.json()}
            except Exception as e:
                return {"is_error": True, "error_message": f"Network error: {str(e)}"}

    async def get(self, path: str, auth_token: str | None = None, **kwargs) -> Dict[str, Any]:
        return await self.request("GET", path, auth_token, **kwargs)

    async def post(self, path: str, auth_token: str | None = None, **kwargs) -> Dict[str, Any]:
        return await self.request("POST", path, auth_token, **kwargs)

admin_client = AdminHttpClient()