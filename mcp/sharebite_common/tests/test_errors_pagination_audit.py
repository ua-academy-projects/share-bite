from sharebite_common.audit import make_audit_event
from sharebite_common.errors import map_http_error
from sharebite_common.pagination import clamp_pagination
from sharebite_common.request_id import get_or_create_request_id


def test_map_http_error_uses_safe_messages():
    assert map_http_error(401, "raw token failure").to_dict()["error"] == "unauthorized"
    assert map_http_error(503, "database password leaked").to_dict()["error"] == "upstream service error"
    assert "details" not in map_http_error(503, "database password leaked").to_dict()


def test_clamp_pagination():
    page = clamp_pagination(limit=500, offset=-10, max_limit=100)

    assert page.limit == 100
    assert page.offset == 0


def test_get_or_create_request_id():
    assert get_or_create_request_id("req-1") == "req-1"
    assert get_or_create_request_id()


def test_make_audit_event_redacts_details():
    event = make_audit_event(
        action="admin_write",
        status="SUCCESS",
        actor_id="admin-1",
        details={"token": "secret", "field": "value"},
    ).to_dict()

    assert event["details"]["token"] == "[REDACTED]"
    assert event["details"]["field"] == "value"
