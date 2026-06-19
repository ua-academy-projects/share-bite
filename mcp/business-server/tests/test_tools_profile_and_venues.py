import asyncio

from app.constants import (
    TOOL_GET_BUSINESS_PROFILE,
    TOOL_GET_VENUE_DETAILS,
    TOOL_LIST_BUSINESS_VENUES,
    TOOL_RECOMMEND_VENUES_BY_CONTEXT,
    TOOL_UPDATE_BUSINESS_PROFILE,
    TOOL_UPDATE_VENUE_DETAILS,
    TOOL_UPDATE_VENUE_HOURS,
    TOOL_PREVIEW_VENUE_HOURS_UPDATE,
)


AUTH_TOKEN = "token-123"
REQUEST_ID = "req-1"
BUSINESS_ID = 10


def test_get_business_profile_success(registered_tools, api_client):
    api_client.get_responses = [{"data": {"id": 10, "name": "Brand 10"}}]

    res = asyncio.run(
        registered_tools[TOOL_GET_BUSINESS_PROFILE](
            business_id=BUSINESS_ID,
            auth_token=AUTH_TOKEN,
            request_id=REQUEST_ID,
        )
    )
    assert res["ok"] is True
    assert res["business_id"] == 10
    assert res["result"]["name"] == "Brand 10"
    assert api_client.get_calls[0]["auth_token"] == AUTH_TOKEN
    assert api_client.get_calls[0]["request_id"] == REQUEST_ID


def test_list_business_venues_success(registered_tools, api_client):
    api_client.get_responses = [{"data": {"items": [{"id": 7}, {"id": 8}], "total": 2}}]

    res = asyncio.run(
        registered_tools[TOOL_LIST_BUSINESS_VENUES](
            business_id=BUSINESS_ID,
            skip=0,
            limit=10,
            auth_token=AUTH_TOKEN,
            request_id=REQUEST_ID,
        )
    )
    assert res["ok"] is True
    assert len(res["result"]["items"]) == 2
    assert api_client.get_calls[0]["params"] == {"skip": 0, "limit": 10}


def test_recommend_venues_by_context_tool(registered_tools):
    res = asyncio.run(
        registered_tools[TOOL_RECOMMEND_VENUES_BY_CONTEXT](
            context={"query": "budget lunch", "budget": "low"},
            venues=[
                {"id": 1, "name": "Fine Dining", "tags": ["romantic"], "price_level": "$$$"},
                {"id": 2, "name": "Student Lunch", "tags": ["budget", "lunch"], "price_level": "$"},
            ],
            limit=1,
        )
    )

    assert res["ok"] is True
    assert res["result"]["items"][0]["venue"]["id"] == 2


def test_get_venue_details_forbidden_for_foreign_venue(registered_tools, api_client):
    api_client.get_responses = [{"data": {"id": 7, "name": "Venue", "brand": {"id": 999}}}]

    res = asyncio.run(
        registered_tools[TOOL_GET_VENUE_DETAILS](
            business_id=BUSINESS_ID,
            venue_id=7,
            auth_token=AUTH_TOKEN,
            request_id=REQUEST_ID,
        )
    )
    assert res["ok"] is False
    assert "unauthorized" in res["error"]


def test_update_business_profile_success_changed_fields(registered_tools, api_client):
    api_client.get_responses = [{"data": {"id": 10, "name": "Old Name", "description": "old"}}]
    api_client.patch_responses = [{"data": {"id": 10, "name": "New Name", "description": "old"}}]

    res = asyncio.run(
        registered_tools[TOOL_UPDATE_BUSINESS_PROFILE](
            business_id=BUSINESS_ID,
            payload={"name": "New Name"},
            auth_token=AUTH_TOKEN,
            request_id=REQUEST_ID,
        )
    )
    assert res["ok"] is True
    assert "name" in res["changed_fields"]


def test_update_business_profile_validation_failure(registered_tools, api_client):
    res = asyncio.run(
        registered_tools[TOOL_UPDATE_BUSINESS_PROFILE](
            business_id=BUSINESS_ID,
            payload={},
            auth_token=AUTH_TOKEN,
            request_id=REQUEST_ID,
        )
    )
    assert res["ok"] is False
    assert res["validation_errors"]
    assert len(api_client.get_calls) == 0
    assert len(api_client.patch_calls) == 0


def test_update_venue_details_success_changed_fields(registered_tools, api_client):
    api_client.get_responses = [
        {"data": {"id": 7, "name": "Old Venue", "brand": {"id": 10}, "description": "old"}}
    ]
    api_client.patch_responses = [
        {"data": {"id": 7, "name": "New Venue", "brand": {"id": 10}, "description": "old"}}
    ]

    res = asyncio.run(
        registered_tools[TOOL_UPDATE_VENUE_DETAILS](
            business_id=BUSINESS_ID,
            venue_id=7,
            payload={"name": "New Venue"},
            auth_token=AUTH_TOKEN,
            request_id=REQUEST_ID,
        )
    )
    assert res["ok"] is True
    assert "name" in res["changed_fields"]


def test_update_venue_details_forbidden_for_foreign_venue(registered_tools, api_client):
    api_client.get_responses = [{"data": {"id": 7, "name": "Venue", "brand": {"id": 999}}}]

    res = asyncio.run(
        registered_tools[TOOL_UPDATE_VENUE_DETAILS](
            business_id=BUSINESS_ID,
            venue_id=7,
            payload={"name": "X"},
            auth_token=AUTH_TOKEN,
            request_id=REQUEST_ID,
        )
    )
    assert res["ok"] is False
    assert "unauthorized" in res["error"]


def test_update_venue_hours_success(registered_tools, api_client):
    payload = {
        "days": [
            {"weekday": 1, "openTime": "09:00", "closeTime": "18:00"},
            {"weekday": 7, "openTime": None, "closeTime": None},
        ]
    }
    api_client.get_responses = [{"data": {"id": 7, "brand": {"id": 10}}}]
    api_client.patch_responses = [{"data": {"venueId": 7, "days": payload["days"]}}]

    res = asyncio.run(
        registered_tools[TOOL_UPDATE_VENUE_HOURS](
            business_id=BUSINESS_ID,
            venue_id=7,
            payload=payload,
            auth_token=AUTH_TOKEN,
            request_id=REQUEST_ID,
        )
    )
    assert res["ok"] is True
    assert res["changed_fields"] == ["days"]


def test_update_venue_hours_validation_failure(registered_tools, api_client):
    bad_payload = {"days": [{"weekday": 1, "openTime": "09:00", "closeTime": None}]}

    res = asyncio.run(
        registered_tools[TOOL_UPDATE_VENUE_HOURS](
            business_id=BUSINESS_ID,
            venue_id=7,
            payload=bad_payload,
            auth_token=AUTH_TOKEN,
            request_id=REQUEST_ID,
        )
    )
    assert res["ok"] is False
    assert res["validation_errors"]
    assert len(api_client.get_calls) == 0
    assert len(api_client.patch_calls) == 0


def test_update_venue_hours_forbidden_for_foreign_venue(registered_tools, api_client):
    payload = {"days": [{"weekday": 1, "openTime": "09:00", "closeTime": "18:00"}]}
    api_client.get_responses = [{"data": {"id": 7, "brand": {"id": 999}}}]

    res = asyncio.run(
        registered_tools[TOOL_UPDATE_VENUE_HOURS](
            business_id=BUSINESS_ID,
            venue_id=7,
            payload=payload,
            auth_token=AUTH_TOKEN,
            request_id=REQUEST_ID,
        )
    )
    assert res["ok"] is False
    assert "unauthorized" in res["error"]


def test_preview_venue_hours_update_success(registered_tools, api_client):
    api_client.get_responses = [
        {
            "data": {
                "id": 7,
                "brand": {"id": 10},
                "days": [
                    {"weekday": 1, "openTime": "09:00", "closeTime": "18:00"},
                    {"weekday": 2, "openTime": "09:00", "closeTime": "18:00"},
                ],
            }
        }
    ]

    payload = {
        "days": [
            {"weekday": 1, "openTime": "10:00", "closeTime": "18:00"},
            {"weekday": 2, "openTime": "09:00", "closeTime": "18:00"},
            {"weekday": 7, "openTime": None, "closeTime": None},
        ]
    }

    res = asyncio.run(
        registered_tools[TOOL_PREVIEW_VENUE_HOURS_UPDATE](
            business_id=BUSINESS_ID,
            venue_id=7,
            payload=payload,
            auth_token=AUTH_TOKEN,
            request_id=REQUEST_ID,
        )
    )

    assert res["ok"] is True
    assert res["changed_fields"] == ["days"]
    assert "summary" in res["result"]
    assert len(api_client.patch_calls) == 0


def test_preview_venue_hours_update_noop(registered_tools, api_client):
    current_days = [
        {"weekday": 1, "openTime": "09:00", "closeTime": "18:00"},
        {"weekday": 7, "openTime": None, "closeTime": None},
    ]
    api_client.get_responses = [
        {
            "data": {
                "id": 7,
                "brand": {"id": 10},
                "days": current_days,
            }
        }
    ]

    res = asyncio.run(
        registered_tools[TOOL_PREVIEW_VENUE_HOURS_UPDATE](
            business_id=BUSINESS_ID,
            venue_id=7,
            payload={"days": current_days},
            auth_token=AUTH_TOKEN,
            request_id=REQUEST_ID,
        )
    )

    assert res["ok"] is True
    assert res["changed_fields"] == []
    assert res["result"]["summary"] == "No venue-hours changes detected."
    assert len(api_client.patch_calls) == 0


def test_preview_venue_hours_update_validation_failure(registered_tools, api_client):
    bad_payload = {"days": [{"weekday": 1, "openTime": "09:00", "closeTime": None}]}

    res = asyncio.run(
        registered_tools[TOOL_PREVIEW_VENUE_HOURS_UPDATE](
            business_id=BUSINESS_ID,
            venue_id=7,
            payload=bad_payload,
            auth_token=AUTH_TOKEN,
            request_id=REQUEST_ID,
        )
    )

    assert res["ok"] is False
    assert res["validation_errors"]
    assert len(api_client.get_calls) == 0
    assert len(api_client.patch_calls) == 0


def test_preview_venue_hours_update_forbidden_for_foreign_venue(registered_tools, api_client):
    api_client.get_responses = [
        {
            "data": {
                "id": 7,
                "brand": {"id": 999},
                "days": [],
            }
        }
    ]

    res = asyncio.run(
        registered_tools[TOOL_PREVIEW_VENUE_HOURS_UPDATE](
            business_id=BUSINESS_ID,
            venue_id=7,
            payload={"days": [{"weekday": 1, "openTime": "09:00", "closeTime": "18:00"}]},
            auth_token=AUTH_TOKEN,
            request_id=REQUEST_ID,
        )
    )

    assert res["ok"] is False
    assert "unauthorized" in res["error"]
    assert len(api_client.patch_calls) == 0
    

def test_preview_venue_hours_update_fails_when_current_hours_unavailable(registered_tools, api_client):
    api_client.get_responses = [
        {
            "data": {
                "id": 7,
                "brand": {"id": 10},
            }
        }
    ]

    res = asyncio.run(
        registered_tools[TOOL_PREVIEW_VENUE_HOURS_UPDATE](
            business_id=BUSINESS_ID,
            venue_id=7,
            payload={"days": [{"weekday": 1, "openTime": "09:00", "closeTime": "18:00"}]},
            auth_token=AUTH_TOKEN,
            request_id=REQUEST_ID,
        )
    )

    assert res["ok"] is False
    assert res["error"] == "current venue hours are unavailable for preview"
    assert len(api_client.patch_calls) == 0