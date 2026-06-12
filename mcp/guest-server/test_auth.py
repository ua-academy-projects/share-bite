from fastmcp import Context
from mcp.shared.context import RequestContext
import typing

class FakeRequest:
    headers = {"Authorization": "Bearer TEST_TOKEN"}

req_ctx = RequestContext(request_id="1", meta=None, session=None, lifespan_context=None, request=FakeRequest())
print(req_ctx.request.headers)
