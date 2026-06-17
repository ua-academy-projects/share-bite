import asyncio
import inspect

from app.constants import (
    URI_PROFILE_SCHEMA,
    URI_VENUE_HOURS_FORMAT,
    URI_VENUE_SCHEMA,
    URI_ANALYTICS_METRICS,
    URI_REPORTING_PERIODS,
)


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
    assert data["type"] == "object"
    assert "days" in data["required"]
    assert "properties" in data
    assert "days" in data["properties"]
    assert data["properties"]["days"]["type"] == "array"
    assert "example" in data
    assert isinstance(data["example"]["days"], list)


def _run_resource(fn):
    result = fn()
    if inspect.iscoroutine(result):
        return asyncio.run(result)
    return result


def test_business_analytics_metrics_resource(registered_resources):
    fn = registered_resources[URI_ANALYTICS_METRICS]
    data = _run_resource(fn)

    assert data["title"] == "Share-Bite Analytics Metrics"
    assert "metrics" in data
    assert len(data["metrics"]) == 5
    assert data["metrics"][0]["name"] == "Sell-Through Rate"


def test_business_reporting_periods_resource(registered_resources):
    fn = registered_resources[URI_REPORTING_PERIODS]
    data = _run_resource(fn)

    assert data["title"] == "Reporting Periods & Constraints"
    assert data["constraints"]["max_days"] == 90
    assert data["constraints"]["format"] == "YYYY-MM-DD"

