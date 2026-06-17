import json

import pytest

from app.constants import URI_OPENAPI_SUMMARY
from app.resources.discovery import get_recommendation_signals, get_search_filters


@pytest.mark.asyncio
async def test_search_filters():
    result = await get_search_filters()
    data = json.loads(result)

    assert "posts" in data
    assert "customers" in data
    assert "collections" in data
    assert data["posts"]["parameters"]["limit"]["max"] == 100
    assert data["customers"]["followers"]["parameters"]["pageSize"]["max"] == 100
    assert data["customers"]["following"]["parameters"]["pageSize"]["max"] == 100
    assert data["customers"]["by_username"]["auth"] == "none"
    assert data["collections"]["list_mine"]["auth"] == "required"


@pytest.mark.asyncio
async def test_recommendation_signals():
    result = await get_recommendation_signals()
    data = json.loads(result)

    assert data["algorithm"] == "Weighted tag quotas with H3 spatial indexing"
    assert "h3_indexing" in data
    assert "tag_quotas" in data
    assert data["output_schema_reference"] == URI_OPENAPI_SUMMARY
