import asyncio
import pytest

from app.constants import API_PATH_FOOD_BOX_PERFORMANCE

TOOL_GET_FOOD_BOX_PERFORMANCE = "get_food_box_performance" 
TOOL_DAILY_SUMMARY = "get_business_daily_summary"
TOOL_RESERVATION_SUMMARY = "get_reservation_summary"
TOOL_ENGAGEMENT_SUMMARY = "get_engagement_summary"
TOOL_VENUE_ACTIVITY = "get_venue_activity_summary"

AUTH_TOKEN = "token-123"
REQUEST_ID = "req-1"

class MockMeta:
    def __init__(self, headers: dict[str, str]):
        self.headers = headers

    def model_dump(self):
        return {"headers": self.headers}

class MockRequestContext:
    def __init__(self, meta: MockMeta):
        self.meta = meta

class MockContext:
    def __init__(self, headers: dict[str, str]):
        self.request_context = MockRequestContext(MockMeta(headers))


def test_get_food_box_performance_success(registered_tools, api_client):
    api_client.get_responses = [
        {"data": {"waste_rate": 0.1, "sell_through_rate": 0.9}}
    ]

    ctx = MockContext({"Authorization": f"Bearer {AUTH_TOKEN}"})

    res = asyncio.run(
        registered_tools[TOOL_GET_FOOD_BOX_PERFORMANCE](
            ctx=ctx,
            start_date="2026-05-01",
            end_date="2026-05-10",
            request_id=REQUEST_ID,
        )
    )
    assert res["ok"] is True
    assert res["result"]["waste_rate"] == 0.1
    assert api_client.get_calls[0]["params"]["start_date"] == "2026-05-01"

def test_get_food_box_performance_invalid_date_range(registered_tools, api_client):
    ctx = MockContext({"Authorization": f"Bearer {AUTH_TOKEN}"})

    res = asyncio.run(
        registered_tools[TOOL_GET_FOOD_BOX_PERFORMANCE](
            ctx=ctx,
            start_date="2026-01-01",
            end_date="2026-06-10",
            request_id=REQUEST_ID,
        )
    )
    assert res["ok"] is False
    assert "exceeds maximum allowed period" in res["error"]
    assert len(api_client.get_calls) == 0

def test_get_food_box_performance_unauthorized(registered_tools, api_client):
    ctx = MockContext({}) 

    res = asyncio.run(
        registered_tools[TOOL_GET_FOOD_BOX_PERFORMANCE](
            ctx=ctx,
            start_date="2026-05-01",
            end_date="2026-05-10",
            request_id=REQUEST_ID,
        )
    )
    assert res["ok"] is False
    assert "Unauthorized" in res["error"]


def test_get_business_daily_summary_success(registered_tools, api_client):
    api_client.get_responses = [{"data": {"boxes_created": 10, "posts_created": 5, "total_venues": 2}}]
    ctx = MockContext({"Authorization": f"Bearer {AUTH_TOKEN}"})

    res = asyncio.run(
        registered_tools[TOOL_DAILY_SUMMARY](
            ctx=ctx,
            start_date="2026-05-01",
            end_date="2026-05-10",
            request_id=REQUEST_ID,
        )
    )
    
    assert res["ok"] is True
    assert res["result"]["boxes_created"] == 10
    assert api_client.get_calls[0]["auth_token"] == AUTH_TOKEN
    assert api_client.get_calls[0]["params"]["start_date"] == "2026-05-01"


def test_get_business_daily_summary_empty_period(registered_tools, api_client):
    api_client.get_responses = [{"data": {}}] 
    ctx = MockContext({"Authorization": f"Bearer {AUTH_TOKEN}"})

    res = asyncio.run(
        registered_tools[TOOL_DAILY_SUMMARY](
            ctx=ctx,
            start_date="2026-05-01",
            end_date="2026-05-10",
            request_id=REQUEST_ID,
        )
    )
    
    assert res["ok"] is True
    assert res["result"] == {}


def test_get_business_daily_summary_invalid_date_range(registered_tools, api_client):
    ctx = MockContext({"Authorization": f"Bearer {AUTH_TOKEN}"})

    res = asyncio.run(
        registered_tools[TOOL_DAILY_SUMMARY](
            ctx=ctx,
            start_date="2026-01-01",
            end_date="2026-06-10",
            request_id=REQUEST_ID,
        )
    )
    
    assert res["ok"] is False
    assert "Validation failed" in res["error"]
    assert len(api_client.get_calls) == 0


def test_get_business_daily_summary_unauthorized(registered_tools, api_client):
    ctx = MockContext({})

    res = asyncio.run(
        registered_tools[TOOL_DAILY_SUMMARY](
            ctx=ctx,
            start_date="2026-05-01",
            end_date="2026-05-10",
            request_id=REQUEST_ID,
        )
    )
    
    assert res["ok"] is False
    assert "Unauthorized" in res["error"]
    assert len(api_client.get_calls) == 0


def test_get_reservation_summary_success(registered_tools, api_client):
    api_client.get_responses = [{"data": {"items_sold": 50, "items_reserved": 20}}]
    ctx = MockContext({"Authorization": f"Bearer {AUTH_TOKEN}"})

    res = asyncio.run(
        registered_tools[TOOL_RESERVATION_SUMMARY](
            ctx=ctx,
            start_date="2026-05-01",
            end_date="2026-05-10",
            venue_id=5,
            request_id=REQUEST_ID,
        )
    )
    
    assert res["ok"] is True
    assert res["result"]["items_sold"] == 50
    assert res["venue_id"] == 5
    assert api_client.get_calls[0]["params"]["venue_id"] == 5

def test_get_reservation_summary_invalid_date(registered_tools, api_client):
    ctx = MockContext({"Authorization": f"Bearer {AUTH_TOKEN}"})
    res = asyncio.run(
        registered_tools[TOOL_RESERVATION_SUMMARY](
            ctx=ctx, 
            start_date="2026-05-10",
            end_date="2026-05-01", 
            request_id=REQUEST_ID
        )
    )
    assert res["ok"] is False
    assert "Validation failed" in res["error"]
    assert len(api_client.get_calls) == 0


def test_get_engagement_summary_success(registered_tools, api_client):
    api_client.get_responses = [{"data": {"total_likes": 150, "total_comments": 30}}]
    ctx = MockContext({"Authorization": f"Bearer {AUTH_TOKEN}"})

    res = asyncio.run(
        registered_tools[TOOL_ENGAGEMENT_SUMMARY](
            ctx=ctx,
            start_date="2026-05-01",
            end_date="2026-05-10",
            request_id=REQUEST_ID,
        )
    )
    
    assert res["ok"] is True
    assert res["result"]["total_likes"] == 150

def test_get_engagement_summary_unauthorized(registered_tools, api_client):
    ctx = MockContext({})
    res = asyncio.run(
        registered_tools[TOOL_ENGAGEMENT_SUMMARY](
            ctx=ctx, start_date="2026-05-01", end_date="2026-05-10", request_id=REQUEST_ID
        )
    )
    assert res["ok"] is False
    assert "Unauthorized" in res["error"]


def test_get_venue_activity_summary_success(registered_tools, api_client):
    api_client.get_responses = [{"data": {"venue_name": "Kyiv Central", "boxes_created": 12}}]
    ctx = MockContext({"Authorization": f"Bearer {AUTH_TOKEN}"})

    res = asyncio.run(
        registered_tools[TOOL_VENUE_ACTIVITY](
            ctx=ctx,
            start_date="2026-05-01",
            end_date="2026-05-10",
            venue_id=42,
            request_id=REQUEST_ID,
        )
    )
    
    assert res["ok"] is True
    assert res["result"]["venue_name"] == "Kyiv Central"
    assert res["venue_id"] == 42
    assert api_client.get_calls[0]["path"] == "/business/analytics/venues/42/activity"

def test_get_venue_activity_summary_empty_period(registered_tools, api_client):
    api_client.get_responses = [{"data": {}}]
    ctx = MockContext({"Authorization": f"Bearer {AUTH_TOKEN}"})
    res = asyncio.run(
        registered_tools[TOOL_VENUE_ACTIVITY](
            ctx=ctx, start_date="2026-05-01", end_date="2026-05-10", venue_id=42, request_id=REQUEST_ID
        )
    )
    assert res["ok"] is True
    assert res["result"] == {}