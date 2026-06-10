import json

from ..constants import (
    CONTENT_TYPE_JSON,
    URI_OPENAPI_SUMMARY,
    URI_SEARCH_FILTERS,
    URI_RECOMMENDATION_SIGNALS,
)
from ..server import mcp


@mcp.resource(
    uri=URI_SEARCH_FILTERS,
    name="guest_search_filters",
    title="Guest Search Filters",
    description="Reference for available search filters across guest discovery tools.",
    mime_type=CONTENT_TYPE_JSON,
)
async def get_search_filters() -> str:
    filters = {
        "posts": {
            "endpoint": "GET /posts/",
            "auth": "optional",
            "parameters": {
                "limit": {"type": "integer", "default": 20, "min": 1, "max": 100},
                "offset": {"type": "integer", "default": 0, "min": 0, "max": 1000},
                "author_id": {
                    "type": "string",
                    "format": "uuid",
                    "optional": True,
                    "description": "Filter by post author customer ID",
                },
            },
        },
        "customers": {
            "by_username": {
                "endpoint": "GET /customers/{username}",
                "auth": "none",
                "description": "Exact username match",
            },
            "followers": {
                "endpoint": "GET /customers/{id}/followers",
                "auth": "optional",
                "parameters": {
                    "pageSize": {
                        "type": "integer",
                        "default": 20,
                        "min": 1,
                        "max": 100,
                    },
                    "pageToken": {"type": "string", "optional": True},
                },
            },
            "following": {
                "endpoint": "GET /customers/{id}/following",
                "auth": "optional",
                "parameters": {
                    "pageSize": {
                        "type": "integer",
                        "default": 20,
                        "min": 1,
                        "max": 100,
                    },
                    "pageToken": {"type": "string", "optional": True},
                },
            },
        },
        "collections": {
            "list_mine": {
                "endpoint": "GET /collections/me",
                "auth": "required",
                "parameters": {
                    "pageSize": {
                        "type": "integer",
                        "default": 20,
                        "min": 1,
                        "max": 100,
                    },
                    "pageToken": {"type": "string", "optional": True},
                },
            },
            "get_by_id": {
                "endpoint": "GET /collections/{collectionId}",
                "auth": "optional",
            },
            "venues": {
                "endpoint": "GET /collections/{collectionId}/venues",
                "auth": "optional",
            },
        },
    }
    return json.dumps(filters)


@mcp.resource(
    uri=URI_RECOMMENDATION_SIGNALS,
    name="guest_recommendation_signals",
    title="Guest Recommendation Signals",
    description="Explains how behavioral recommendations are generated in the ShareBite ecosystem.",
    mime_type=CONTENT_TYPE_JSON,
)
async def get_recommendation_signals() -> str:
    signals = {
        "algorithm": "Weighted tag quotas with H3 spatial indexing",
        "h3_indexing": {
            "description": (
                "Posts and venues are indexed using H3 hexagonal grid for efficient "
                "geospatial queries. Neighbor radius is used to find nearby content."
            ),
        },
        "tag_quotas": {
            "description": (
                "Based on user's top 5 liked tags, posts are fetched using weighted "
                "quotas: 5-3-2-1-1 distribution. If user has no like history, "
                "random nearby posts are returned and shuffled."
            ),
            "fallback": "Random nearby posts when no behavioral history exists.",
        },
        "inputs_required": {
            "lat": "User latitude",
            "lon": "User longitude",
            "auth": "Bearer token required for personalized recommendations",
        },
        "output_schema_reference": URI_OPENAPI_SUMMARY,
    }
    return json.dumps(signals)
