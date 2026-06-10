import json

import httpx
import pytest

from app.tools.posts import get_post, get_post_authors, search_posts
from tests.factories.guest_responses import (
    build_error_response,
    build_post_authors_response,
    build_post_detail_response,
    build_posts_list_response,
)


@pytest.mark.asyncio
class TestSearchPosts:
    async def test_success(self, mock_guest_api):
        mock_guest_api.get("/posts/").mock(
            return_value=build_posts_list_response(
                posts=[{"id": "1", "text": "hello"}],
                total=1,
            )
        )
        result = await search_posts(limit=10, offset=0)
        data = json.loads(result)
        assert data["posts"][0]["id"] == "1"
        assert data["total"] == 1

    async def test_empty_result(self, mock_guest_api):
        mock_guest_api.get("/posts/").mock(
            return_value=build_posts_list_response(total=0)
        )
        result = await search_posts()
        data = json.loads(result)
        assert data["posts"] == []
        assert data["total"] == 0

    async def test_unauthorized(self, mock_guest_api):
        mock_guest_api.get("/posts/").mock(
            return_value=build_error_response(401, "invalid or expired token")
        )
        result = await search_posts()
        data = json.loads(result)
        assert data["error"] == "unauthorized"
        assert "invalid or expired token" in data["message"]

    async def test_downstream_failure(self, mock_guest_api):
        mock_guest_api.get("/posts/").mock(side_effect=httpx.ConnectTimeout("mocked"))
        result = await search_posts()
        data = json.loads(result)
        assert data["error"] == "downstream_failure"


@pytest.mark.asyncio
class TestGetPost:
    async def test_success(self, mock_guest_api):
        mock_guest_api.get("/posts/42").mock(
            return_value=build_post_detail_response(
                post={"id": "42", "text": "great venue"}
            )
        )
        result = await get_post(id=42)
        data = json.loads(result)
        assert data["post"]["id"] == "42"

    async def test_not_found(self, mock_guest_api):
        mock_guest_api.get("/posts/99").mock(
            return_value=build_error_response(404, "Post not found")
        )
        result = await get_post(id=99)
        data = json.loads(result)
        assert data["error"] == "not_found"
        assert data["message"] == "Post not found"

    async def test_downstream_timeout(self, mock_guest_api):
        mock_guest_api.get("/posts/1").mock(side_effect=httpx.ConnectTimeout("mocked"))
        result = await get_post(id=1)
        data = json.loads(result)
        assert data["error"] == "downstream_failure"


@pytest.mark.asyncio
class TestGetPostAuthors:
    async def test_success(self, mock_guest_api):
        mock_guest_api.get("/posts/1/authors").mock(
            return_value=build_post_authors_response(authors=["uuid-1"], count=1)
        )
        result = await get_post_authors(id=1)
        data = json.loads(result)
        assert data["count"] == 1

    async def test_not_found(self, mock_guest_api):
        mock_guest_api.get("/posts/1/authors").mock(
            return_value=build_error_response(404, "Post not found")
        )
        result = await get_post_authors(id=1)
        data = json.loads(result)
        assert data["error"] == "not_found"
