import asyncio
import inspect

from app.constants import URI_PROFILE_SCHEMA, URI_VENUE_HOURS_FORMAT, URI_VENUE_SCHEMA


def test_profile_schema_resource(registered_resources):
    fn = registered_resources[URI_PROFILE_SCHEMA]
    data = _run_resource(fn)
    assert data["type"] == "object"
    assert "name" in data["properties"]


def test_venue_schema_resource(registered_resources):
    fn = registered_resources[URI_VENUE_SCHEMA]
    data = _run_resource(fn)
    assert data["type"] == "object"
    assert "latitude" in data["properties"]
    assert "tagIds" in data["properties"]


def test_venue_hours_format_resource(registered_resources):
    fn = registered_resources[URI_VENUE_HOURS_FORMAT]
    data = _run_resource(fn)
    assert "days" in data
    assert isinstance(data["days"], list)
    assert len(data["days"]) >= 1


def _run_resource(fn):
    result = fn()
    if inspect.isawaitable(result):
        return asyncio.run(result)
    return result
