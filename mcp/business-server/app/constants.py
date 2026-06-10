# HTTP Headers
HEADER_AUTH = "Authorization"
HEADER_REQUEST_ID = "X-Request-ID"
HEADER_CONTENT_TYPE = "Content-Type"

# MIME Types
CONTENT_TYPE_JSON = "application/json"

# Resource URIs
URI_PROFILE_SCHEMA = "sharebite://business/profile-schema"
URI_VENUE_SCHEMA = "sharebite://business/venue-schema"
URI_VENUE_HOURS_FORMAT = "sharebite://business/venue-hours-format"

# MCP Server
MCP_SERVER_NAME = "business-server"

# Tool names
TOOL_GET_BUSINESS_PROFILE = "get_business_profile"
TOOL_UPDATE_BUSINESS_PROFILE = "update_business_profile"
TOOL_LIST_BUSINESS_VENUES = "list_business_venues"
TOOL_GET_VENUE_DETAILS = "get_venue_details"
TOOL_UPDATE_VENUE_DETAILS = "update_venue_details"
TOOL_UPDATE_VENUE_HOURS = "update_venue_hours"
TOOL_RECOMMEND_VENUES_BY_CONTEXT = "recommend_venues_by_context"

# API paths
API_PATH_BUSINESS_PROFILE = "/business/{business_id}"
API_PATH_BUSINESS_VENUES = "/business/org-units/{business_id}/locations"
API_PATH_VENUE_DETAILS = "/business/org-units/{venue_id}"
API_PATH_UPDATE_VENUE_DETAILS = "/business/locations/{venue_id}"
API_PATH_UPDATE_VENUE_HOURS = "/business/locations/{venue_id}/hours"

# Roles
ROLE_BUSINESS = "business"