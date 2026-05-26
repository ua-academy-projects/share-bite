import asyncio

from app.constants import (
    TOOL_GET_BUSINESS_PROFILE,
    TOOL_GET_VENUE_DETAILS,
    TOOL_LIST_BUSINESS_VENUES,
    TOOL_UPDATE_BUSINESS_PROFILE,
    TOOL_UPDATE_VENUE_DETAILS,
    TOOL_UPDATE_VENUE_HOURS,
)


def test_get_business_profile_success(registered_tools, api_client, auth_context):
    api_client.get_responses = [{"data": {"id": 10, "name": "Brand 10"}}]

    res = asyncio.run(registered_tools[TOOL_GET_BUSINESS_PROFILE](auth_context=auth_context))
    assert res["ok"] is True
    assert res["business_id"] == 10
    assert res["result"]["name"] == "Brand 10"


def test_get_business_profile_forbidden_on_business_id_mismatch(registered_tools, auth_context):
    res = asyncio.run(
        registered_tools[TOOL_GET_BUSINESS_PROFILE](business_id=11, auth_context=auth_context)
    )
    assert res["ok"] is False
    assert "mismatch" in res["error"]


def test_list_business_venues_success(registered_tools, api_client, auth_context):
    api_client.get_responses = [{"data": {"items": [{"id": 7}, {"id": 8}], "total": 2}}]

    res = asyncio.run(
        registered_tools[TOOL_LIST_BUSINESS_VENUES](skip=0, limit=10, auth_context=auth_context)
    )
    assert res["ok"] is True
    assert len(res["result"]["items"]) == 2


def test_get_venue_details_forbidden_for_foreign_venue(registered_tools, api_client, auth_context):
    api_client.get_responses = [{"data": {"id": 7, "name": "Venue", "brand": {"id": 999}}}]

    res = asyncio.run(
        registered_tools[TOOL_GET_VENUE_DETAILS](venue_id=7, auth_context=auth_context)
    )
    assert res["ok"] is False
    assert "unauthorized" in res["error"]


def test_update_business_profile_success_changed_fields(registered_tools, api_client, auth_context):
    api_client.get_responses = [{"data": {"id": 10, "name": "Old Name", "description": "old"}}]
    api_client.patch_responses = [{"data": {"id": 10, "name": "New Name", "description": "old"}}]

    res = asyncio.run(
        registered_tools[TOOL_UPDATE_BUSINESS_PROFILE](
            payload={"name": "New Name"},
            auth_context=auth_context,
        )
    )
    assert res["ok"] is True
    assert "name" in res["changed_fields"]


def test_update_business_profile_validation_failure(registered_tools, api_client, auth_context):
    res = asyncio.run(
        registered_tools[TOOL_UPDATE_BUSINESS_PROFILE](payload={}, auth_context=auth_context)
    )
    assert res["ok"] is False
    assert res["validation_errors"]
    assert len(api_client.get_calls) == 0
    assert len(api_client.patch_calls) == 0


def test_update_venue_details_success_changed_fields(registered_tools, api_client, auth_context):
    api_client.get_responses = [
        {"data": {"id": 7, "name": "Old Venue", "brand": {"id": 10}, "description": "old"}}
    ]
    api_client.patch_responses = [
        {"data": {"id": 7, "name": "New Venue", "brand": {"id": 10}, "description": "old"}}
    ]

    res = asyncio.run(
        registered_tools[TOOL_UPDATE_VENUE_DETAILS](
            venue_id=7,
            payload={"name": "New Venue"},
            auth_context=auth_context,
        )
    )
    assert res["ok"] is True
    assert "name" in res["changed_fields"]


def test_update_venue_details_forbidden_for_foreign_venue(registered_tools, api_client, auth_context):
    api_client.get_responses = [{"data": {"id": 7, "name": "Venue", "brand": {"id": 999}}}]

    res = asyncio.run(
        registered_tools[TOOL_UPDATE_VENUE_DETAILS](
            venue_id=7,
            payload={"name": "X"},
            auth_context=auth_context,
        )
    )
    assert res["ok"] is False
    assert "unauthorized" in res["error"]


def test_update_venue_hours_success(registered_tools, api_client, auth_context):
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
            venue_id=7,
            payload=payload,
            auth_context=auth_context,
        )
    )
    assert res["ok"] is True
    assert res["changed_fields"] == ["days"]


def test_update_venue_hours_validation_failure(registered_tools, api_client, auth_context):
    bad_payload = {"days": [{"weekday": 1, "openTime": "09:00", "closeTime": None}]}

    res = asyncio.run(
        registered_tools[TOOL_UPDATE_VENUE_HOURS](
            venue_id=7,
            payload=bad_payload,
            auth_context=auth_context,
        )
    )
    assert res["ok"] is False
    assert res["validation_errors"]
    assert len(api_client.get_calls) == 0
    assert len(api_client.patch_calls) == 0


def test_update_venue_hours_forbidden_for_foreign_venue(registered_tools, api_client, auth_context):
    payload = {"days": [{"weekday": 1, "openTime": "09:00", "closeTime": "18:00"}]}
    api_client.get_responses = [{"data": {"id": 7, "brand": {"id": 999}}}]

    res = asyncio.run(
        registered_tools[TOOL_UPDATE_VENUE_HOURS](
            venue_id=7,
            payload=payload,
            auth_context=auth_context,
        )
    )
    assert res["ok"] is False
    assert "unauthorized" in res["error"]