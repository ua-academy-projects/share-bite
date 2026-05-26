from .constants import (
    CONTENT_TYPE_JSON,
    URI_PROFILE_SCHEMA,
    URI_VENUE_HOURS_FORMAT,
    URI_VENUE_SCHEMA,
)


def register_resources(mcp) -> None:
    @mcp.resource(
        uri=URI_PROFILE_SCHEMA,
        name="business_profile_schema",
        title="Business Profile Schema",
        description="Allowed fields for update_business_profile",
        mime_type=CONTENT_TYPE_JSON,
    )
    async def profile_schema() -> dict:
        return {
            "type": "object",
            "properties": {
                "name": {"type": "string", "minLength": 3, "maxLength": 40},
                "avatar": {"type": ["string", "null"]},
                "banner": {"type": ["string", "null"]},
                "description": {"type": ["string", "null"]},
            },
            "additionalProperties": False,
        }

    @mcp.resource(
        uri=URI_VENUE_SCHEMA,
        name="business_venue_schema",
        title="Venue Schema",
        description="Allowed fields for update_venue_details",
        mime_type=CONTENT_TYPE_JSON,
    )
    async def venue_schema() -> dict:
        return {
            "type": "object",
            "properties": {
                "name": {"type": ["string", "null"]},
                "avatar": {"type": ["string", "null"]},
                "banner": {"type": ["string", "null"]},
                "description": {"type": ["string", "null"]},
                "latitude": {"type": ["number", "null"], "minimum": -90, "maximum": 90},
                "longitude": {"type": ["number", "null"], "minimum": -180, "maximum": 180},
                "tagIds": {"type": ["array", "null"], "items": {"type": "integer"}, "maxItems": 5},
            },
            "additionalProperties": False,
        }

    @mcp.resource(
        uri=URI_VENUE_HOURS_FORMAT,
        name="business_venue_hours_format",
        title="Venue Hours Format",
        description="Format for update_venue_hours",
        mime_type=CONTENT_TYPE_JSON,
    )
    async def venue_hours_format() -> dict:
        return {
            "days": [
                {"weekday": 1, "openTime": "09:00", "closeTime": "18:00"},
                {"weekday": 2, "openTime": "09:00", "closeTime": "18:00"},
                {"weekday": 3, "openTime": "09:00", "closeTime": "18:00"},
                {"weekday": 4, "openTime": "09:00", "closeTime": "18:00"},
                {"weekday": 5, "openTime": "09:00", "closeTime": "18:00"},
                {"weekday": 6, "openTime": None, "closeTime": None},
                {"weekday": 7, "openTime": None, "closeTime": None},
            ]
        }