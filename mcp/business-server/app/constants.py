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

# Discovery tool names
TOOL_SEARCH_VENUES = "search_venues"
TOOL_GET_RECOMMENDED_VENUES = "get_recommended_venues"
TOOL_GET_FEED_ITEMS = "get_feed_items"
TOOL_SEARCH_BOXES = "search_boxes"

# API paths
API_PATH_BUSINESS_PROFILE = "/business/{business_id}"
API_PATH_BUSINESS_VENUES = "/business/org-units/{business_id}/locations"
API_PATH_VENUE_DETAILS = "/business/org-units/{venue_id}"
API_PATH_UPDATE_VENUE_DETAILS = "/business/locations/{venue_id}"
API_PATH_UPDATE_VENUE_HOURS = "/business/locations/{venue_id}/hours"

# Discovery API paths
API_PATH_SEARCH_VENUES = "/business/venues/search"
API_PATH_NEARBY_VENUES = "/business/locations/nearby"
API_PATH_RECOMMEND_POSTS = "/business/posts/recommend"
API_PATH_NEARBY_BOXES = "/business/nearby-boxes"

# Roles
ROLE_BUSINESS = "business"

API_PATH_DAILY_SUMMARY = "/business/analytics/daily-summary"
API_PATH_RESERVATION_SUMMARY = "/business/analytics/reservation-summary"
API_PATH_FOOD_BOX_PERFORMANCE = "/business/analytics/food-box-performance"
API_PATH_ENGAGEMENT_SUMMARY = "/business/analytics/engagement-summary"
API_PATH_VENUE_ACTIVITY = "/business/analytics/venues/{venue_id}/activity"

# Analytics Resources
URI_ANALYTICS_METRICS = "sharebite://business/analytics-metrics"
URI_REPORTING_PERIODS = "sharebite://business/reporting-periods"
