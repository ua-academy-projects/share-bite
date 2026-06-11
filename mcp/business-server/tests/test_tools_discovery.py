import asyncio
from unittest.mock import patch

from app.constants import (
    TOOL_GET_FEED_ITEMS,
    TOOL_GET_RECOMMENDED_VENUES,
    TOOL_SEARCH_BOXES,
    TOOL_SEARCH_VENUES,
)

AUTH_TOKEN = "token-123"
REQUEST_ID = "req-1"


class FakeContext:
    class _RequestContext:
        def __init__(self, headers):
            self.meta = type(
                "Meta",
                (),
                {
                    "model_dump": lambda self: {"headers": headers or {}},
                    "dict": lambda self: {"headers": headers or {}},
                },
            )()

    def __init__(self, headers=None):
        self.request_context = self._RequestContext(headers or {})


def _ctx(headers=None):
    return FakeContext(headers=headers)


def test_search_venues_success(registered_tools, api_client):
    api_client.get_responses = [
        {"data": {"items": [{"id": 1, "name": "Venue A"}], "total": 1}}
    ]

    res = asyncio.run(
        registered_tools[TOOL_SEARCH_VENUES](
            ctx=_ctx({"Authorization": "Bearer " + AUTH_TOKEN}),
            q="coffee",
            tags="vegan",
            skip=0,
            limit=10,
            request_id=REQUEST_ID,
        )
    )
    assert res["ok"] is True
    assert res["result"]["items"][0]["name"] == "Venue A"
    assert api_client.get_calls[0]["params"] == {
        "q": "coffee",
        "tags": "vegan",
        "skip": 0,
        "limit": 10,
    }
    assert api_client.get_calls[0]["auth_token"] == AUTH_TOKEN
    assert api_client.get_calls[0]["request_id"] == REQUEST_ID


def test_search_venues_anonymous(registered_tools, api_client):
    api_client.get_responses = [{"data": {"items": [], "total": 0}}]

    res = asyncio.run(
        registered_tools[TOOL_SEARCH_VENUES](
            ctx=_ctx(),
            q="pizza",
            skip=0,
            limit=5,
        )
    )
    assert res["ok"] is True
    assert api_client.get_calls[0]["auth_token"] is None
    assert api_client.get_calls[0]["params"]["q"] == "pizza"


def test_search_venues_empty_filters(registered_tools, api_client):
    res = asyncio.run(
        registered_tools[TOOL_SEARCH_VENUES](
            ctx=_ctx(),
            q="",
            tags=None,
            skip=0,
            limit=10,
        )
    )
    assert res["ok"] is False
    assert "at least one search filter" in res["error"]
    assert len(api_client.get_calls) == 0


def test_search_venues_downstream_failure(registered_tools, api_client):
    api_client.get_responses = [{"is_error": True, "error_message": "internal error"}]

    res = asyncio.run(
        registered_tools[TOOL_SEARCH_VENUES](
            ctx=_ctx(),
            q="test",
            skip=0,
            limit=10,
        )
    )
    assert res["ok"] is False
    assert "internal error" in res["error"]


def test_get_recommended_venues_success(registered_tools, api_client):
    api_client.get_responses = [
        {"data": {"items": [{"id": 2, "name": "Venue B", "distance": 1.5}], "total": 1}}
    ]

    res = asyncio.run(
        registered_tools[TOOL_GET_RECOMMENDED_VENUES](
            ctx=_ctx(),
            lat=50.45,
            lon=30.52,
            skip=0,
            limit=10,
            request_id=REQUEST_ID,
        )
    )
    assert res["ok"] is True
    assert res["result"]["items"][0]["distance"] == 1.5
    assert api_client.get_calls[0]["params"] == {
        "lat": 50.45,
        "lon": 30.52,
        "skip": 0,
        "limit": 10,
    }


def test_get_recommended_venues_invalid_coords(registered_tools, api_client):
    res = asyncio.run(
        registered_tools[TOOL_GET_RECOMMENDED_VENUES](
            ctx=_ctx(),
            lat=999,
            lon=30.52,
            skip=0,
            limit=10,
        )
    )
    assert res["ok"] is False
    assert any(e["field"] == "lat" for e in res["validation_errors"])
    assert len(api_client.get_calls) == 0


def test_get_feed_items_success(registered_tools, api_client):
    api_client.get_responses = [
        {"data": {"items": [{"id": 101, "content": "Post 1"}], "total": 1}}
    ]

    res = asyncio.run(
        registered_tools[TOOL_GET_FEED_ITEMS](
            ctx=_ctx({"Authorization": "Bearer " + AUTH_TOKEN}),
            lat=50.45,
            lon=30.52,
            skip=0,
            limit=24,
            request_id=REQUEST_ID,
        )
    )
    assert res["ok"] is True
    assert res["result"]["items"][0]["content"] == "Post 1"
    assert api_client.get_calls[0]["auth_token"] == AUTH_TOKEN
    assert api_client.get_calls[0]["params"] == {
        "lat": 50.45,
        "lon": 30.52,
        "skip": 0,
        "limit": 24,
    }


def test_get_feed_items_unauthorized(registered_tools, api_client):
    with patch("app.tools.resolve_auth_token", return_value=None):
        res = asyncio.run(
            registered_tools[TOOL_GET_FEED_ITEMS](
                ctx=_ctx(),
                lat=50.45,
                lon=30.52,
                skip=0,
                limit=10,
            )
        )
    assert res["ok"] is False
    assert "Unauthorized" in res["error"]
    assert len(api_client.get_calls) == 0


def test_search_boxes_success(registered_tools, api_client):
    api_client.get_responses = [
        {"data": {"items": [{"id": 10, "full_price": 100}], "total": 1}}
    ]

    res = asyncio.run(
        registered_tools[TOOL_SEARCH_BOXES](
            ctx=_ctx({"Authorization": "Bearer " + AUTH_TOKEN}),
            lat=50.45,
            lon=30.52,
            skip=0,
            limit=10,
            org_id=5,
            category_id=2,
            request_id=REQUEST_ID,
        )
    )
    assert res["ok"] is True
    assert res["result"]["items"][0]["full_price"] == 100
    assert api_client.get_calls[0]["params"] == {
        "lat": 50.45,
        "lon": 30.52,
        "skip": 0,
        "limit": 10,
        "org_id": 5,
        "category_id": 2,
    }


def test_search_boxes_anonymous(registered_tools, api_client):
    api_client.get_responses = [{"data": {"items": [], "total": 0}}]

    res = asyncio.run(
        registered_tools[TOOL_SEARCH_BOXES](
            ctx=_ctx(),
            lat=50.45,
            lon=30.52,
            skip=0,
            limit=10,
        )
    )
    assert res["ok"] is True
    assert api_client.get_calls[0]["auth_token"] is None


def test_search_boxes_pagination_validation(registered_tools, api_client):
    res = asyncio.run(
        registered_tools[TOOL_SEARCH_BOXES](
            ctx=_ctx(),
            lat=50.45,
            lon=30.52,
            skip=-1,
            limit=200,
        )
    )
    assert res["ok"] is False
    assert any(e["field"] == "skip" for e in res["validation_errors"])
    assert any(e["field"] == "limit" for e in res["validation_errors"])
    assert len(api_client.get_calls) == 0
