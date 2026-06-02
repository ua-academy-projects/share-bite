import sys
from pathlib import Path

import pytest


SERVER_ROOT = Path(__file__).resolve().parents[1]
if str(SERVER_ROOT) not in sys.path:
    sys.path.insert(0, str(SERVER_ROOT))


class FakeMCP:
    def __init__(self) -> None:
        self.tools: dict[str, callable] = {}
        self.resources: dict[str, callable] = {}

    def tool(self, name: str | None = None, description: str = "", exclude_args: list[str] | None = None):
        def decorator(func):
            self.tools[name or func.__name__] = func
            return func
        return decorator

    def resource(
        self,
        uri: str,
        name: str | None = None,
        title: str = "",
        description: str = "",
        mime_type: str = "",
    ):
        def decorator(func):
            self.resources[uri] = func
            return func
        return decorator


class FakeAPIClient:
    def __init__(self) -> None:
        self.get_responses: list[dict] = []
        self.patch_responses: list[dict] = []
        self.get_calls: list[dict] = []
        self.patch_calls: list[dict] = []

    async def get(self, path, auth_token=None, request_id=None, params=None):
        self.get_calls.append(
            {"path": path, "auth_token": auth_token, "request_id": request_id, "params": params}
        )
        if not self.get_responses:
            raise AssertionError("Unexpected GET call: no fake response queued")
        return self.get_responses.pop(0)

    async def patch(self, path, json_data, auth_token=None, request_id=None):
        self.patch_calls.append(
            {"path": path, "json_data": json_data, "auth_token": auth_token, "request_id": request_id}
        )
        if not self.patch_responses:
            raise AssertionError("Unexpected PATCH call: no fake response queued")
        return self.patch_responses.pop(0)


class FakeSettings:
    business_api_base_url = "http://business-api.test"
    request_timeout_seconds = 10


@pytest.fixture
def fake_mcp():
    return FakeMCP()


@pytest.fixture
def api_client():
    return FakeAPIClient()


@pytest.fixture
def settings():
    return FakeSettings()


@pytest.fixture
def registered_tools(fake_mcp, settings, api_client):
    from app.tools import register_tools

    register_tools(fake_mcp, settings, api_client)
    return fake_mcp.tools


@pytest.fixture
def registered_resources(fake_mcp, settings):
    from app.resources import register_resources

    register_resources(fake_mcp, settings)
    return fake_mcp.resources
