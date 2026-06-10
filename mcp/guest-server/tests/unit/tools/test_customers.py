import json

import httpx
import pytest

from app.tools.customers import (
    get_customer_by_username,
    get_customer_followers,
    get_customer_following,
)
from tests.factories.guest_responses import (
    build_customer_detail_response,
    build_error_response,
    build_followers_list_response,
)


@pytest.mark.asyncio
class TestGetCustomerByUsername:
    async def test_success(self, mock_guest_api):
        mock_guest_api.get("/customers/alice").mock(
            return_value=build_customer_detail_response(
                customer={"id": "uuid-1", "userName": "alice"}
            )
        )
        result = await get_customer_by_username("alice")
        data = json.loads(result)
        assert data["customer"]["userName"] == "alice"

    async def test_not_found(self, mock_guest_api):
        mock_guest_api.get("/customers/fake").mock(
            return_value=build_error_response(
                404, "Customer with this username does not exist"
            )
        )
        result = await get_customer_by_username("fake")
        data = json.loads(result)
        assert data["error"] == "not_found"
        assert "does not exist" in data["message"]

    async def test_downstream_failure(self, mock_guest_api):
        mock_guest_api.get("/customers/x").mock(
            side_effect=httpx.ConnectError("mocked")
        )
        result = await get_customer_by_username("x")
        data = json.loads(result)
        assert data["error"] == "downstream_failure"


@pytest.mark.asyncio
class TestGetCustomerFollowers:
    async def test_success(self, mock_guest_api):
        mock_guest_api.get("/customers/uuid-1/followers").mock(
            return_value=build_followers_list_response(
                customers=[{"id": "f-1", "userName": "bob"}],
                next_page_token="token123",
            )
        )
        result = await get_customer_followers("uuid-1", page_size=10)
        data = json.loads(result)
        assert len(data["customers"]) == 1

    async def test_private_profile(self, mock_guest_api):
        mock_guest_api.get("/customers/uuid-1/followers").mock(
            return_value=build_error_response(403, "Followers list is private")
        )
        result = await get_customer_followers("uuid-1")
        data = json.loads(result)
        assert data["error"] == "forbidden"
        assert "private" in data["message"]


@pytest.mark.asyncio
class TestGetCustomerFollowing:
    async def test_success(self, mock_guest_api):
        mock_guest_api.get("/customers/uuid-1/following").mock(
            return_value=build_followers_list_response(
                customers=[{"id": "f-2", "userName": "charlie"}],
                next_page_token="",
            )
        )
        result = await get_customer_following("uuid-1")
        data = json.loads(result)
        assert data["customers"][0]["userName"] == "charlie"
