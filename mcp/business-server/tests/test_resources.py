import asyncio

from app.constants import URI_PROFILE_SCHEMA, URI_VENUE_HOURS_FORMAT, URI_VENUE_SCHEMA


def test_profile_schema_resource(registered_resources):
    fn = registered_resources[URI_PROFILE_SCHEMA]
    data = asyncio.run(fn())
    assert data["type"] == "object"
    assert "name" in data["properties"]


def test_venue_schema_resource(registered_resources):
    fn = registered_resources[URI_VENUE_SCHEMA]
    data = asyncio.run(fn())
    assert data["type"] == "object"
    assert "latitude" in data["properties"]
    assert "tagIds" in data["properties"]


def test_venue_hours_format_resource(registered_resources):
    fn = registered_resources[URI_VENUE_HOURS_FORMAT]
    data = asyncio.run(fn())
    assert "days" in data
    assert isinstance(data["days"], list)
    assert len(data["days"]) >= 1