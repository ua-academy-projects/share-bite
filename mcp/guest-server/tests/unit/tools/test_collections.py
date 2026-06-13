import json
from unittest.mock import patch

import pytest

from app.tools.collections import (
    get_collection,
    get_collection_venues,
    list_my_collections,
)
from tests.factories.guest_responses import (
    build_collection_detail_response,
    build_collection_venues_response,
    build_collections_list_response,
    build_error_response,
)


@pytest.mark.asyncio
class TestListMyCollections:
    async def test_success(self, mock_guest_api, fake_auth_token):
        route = mock_guest_api.get("/collections/me").mock(
            return_value=build_collections_list_response(
                collections=[{"id": "col-1", "name": "Favorites"}],
                next_page_token="tok",
            )
        )
        result = await list_my_collections(None, page_size=10)
        data = json.loads(result)
        assert data["collections"][0]["name"] == "Favorites"
        assert route.called
        sent_auth = route.calls.last.request.headers.get("Authorization", "")
        assert "fake-jwt" in sent_auth

    async def test_unauthorized_no_token(self):
        with patch("app.auth.resolve_auth_token", return_value=None):
            result = await list_my_collections(None)
        data = json.loads(result)
        assert data["error"] == "unauthorized"
        assert "Authentication required" in data["message"]

    async def test_downstream_failure(self, mock_guest_api, fake_auth_token):
        mock_guest_api.get("/collections/me").mock(
            return_value=build_error_response(500, "Internal server error")
        )
        result = await list_my_collections(None)
        data = json.loads(result)
        assert data["error"] == "downstream_failure"


@pytest.mark.asyncio
class TestGetCollection:
    async def test_success_public(self, mock_guest_api):
        mock_guest_api.get("/collections/col-1").mock(
            return_value=build_collection_detail_response(
                collection={"id": "col-1", "name": "Public", "isPublic": True}
            )
        )
        result = await get_collection(None, "col-1")
        data = json.loads(result)
        assert data["collection"]["isPublic"] is True

    async def test_not_found(self, mock_guest_api):
        mock_guest_api.get("/collections/col-1").mock(
            return_value=build_error_response(404, "Collection does not exist")
        )
        result = await get_collection(None, "col-1")
        data = json.loads(result)
        assert data["error"] == "not_found"


@pytest.mark.asyncio
class TestGetCollectionVenues:
    async def test_success(self, mock_guest_api):
        mock_guest_api.get("/collections/col-1/venues").mock(
            return_value=build_collection_venues_response(
                venues=[{"id": 1, "name": "Venue A", "sortOrder": 1.0}]
            )
        )
        result = await get_collection_venues(None, "col-1")
        data = json.loads(result)
        assert data["venues"][0]["name"] == "Venue A"

    async def test_empty(self, mock_guest_api):
        mock_guest_api.get("/collections/col-1/venues").mock(
            return_value=build_collection_venues_response()
        )
        result = await get_collection_venues(None, "col-1")
        data = json.loads(result)
        assert data["venues"] == []
