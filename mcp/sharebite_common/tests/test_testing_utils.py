import pytest

from sharebite_common.testing import FakeContext, FakeMCP, assert_mcp_response_has_no_secrets
from sharebite_common.security import SecretLeakError


def test_fake_mcp_registers_tools_and_resources():
    mcp = FakeMCP()

    @mcp.tool()
    def ping():
        return "pong"

    @mcp.resource("sharebite://test")
    def resource():
        return {}

    assert mcp.tools["ping"]() == "pong"
    assert mcp.resources["sharebite://test"]() == {}


def test_fake_context_headers():
    ctx = FakeContext({"Authorization": "Bearer abc"})

    assert ctx.request_context.meta.model_dump()["headers"]["Authorization"] == "Bearer abc"


def test_assert_mcp_response_has_no_secrets():
    with pytest.raises(SecretLeakError):
        assert_mcp_response_has_no_secrets({"auth_token": "secret"})
