from fastmcp import FastMCP
from .constants import MCP_SERVER_NAME

mcp = FastMCP(MCP_SERVER_NAME, mask_error_details=True)
