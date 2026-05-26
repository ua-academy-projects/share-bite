from .resources import register_resources
from .tools import register_tools


def register_business_profile_venue_workflows(mcp, api_client) -> None:
    register_resources(mcp)
    register_tools(mcp, api_client)