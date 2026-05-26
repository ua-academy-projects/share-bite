import pytest

from app.utils import (
    AccessError,
    ForbiddenError,
    ensure_venue_owned_by_business,
    resolve_business_access,
    validate_profile_update,
    validate_venue_hours,
    validate_venue_update,
)


def test_resolve_business_access_success():
    out = resolve_business_access(
        auth_context={"role": "business", "auth_token": "t", "business_id": 7, "request_id": "r1"},
        explicit_business_id=None,
    )
    assert out["business_id"] == 7
    assert out["auth_token"] == "t"
    assert out["request_id"] == "r1"


def test_resolve_business_access_missing_auth_context():
    with pytest.raises(AccessError):
        resolve_business_access(None, None)


def test_resolve_business_access_forbidden_on_business_mismatch():
    with pytest.raises(ForbiddenError):
        resolve_business_access(
            auth_context={"role": "business", "auth_token": "t", "business_id": 7},
            explicit_business_id=8,
        )


def test_validate_profile_update_unknown_field():
    errs = validate_profile_update({"badField": "x"})
    assert any(e["field"] == "badField" for e in errs)


def test_validate_venue_update_tag_ids_limit():
    errs = validate_venue_update({"tagIds": [1, 2, 3, 4, 5, 6]})
    assert any(e["field"] == "tagIds" for e in errs)


def test_validate_venue_hours_partial_pair_error():
    errs = validate_venue_hours({"days": [{"weekday": 1, "openTime": "09:00", "closeTime": None}]})
    assert any("both openTime and closeTime" in e["message"] for e in errs)


def test_validate_venue_hours_invalid_order_error():
    errs = validate_venue_hours({"days": [{"weekday": 1, "openTime": "18:00", "closeTime": "09:00"}]})
    assert any("openTime must be before closeTime" in e["message"] for e in errs)


def test_validate_venue_hours_closed_day_allowed():
    errs = validate_venue_hours({"days": [{"weekday": 7, "openTime": None, "closeTime": None}]})
    assert errs == []


def test_ensure_venue_owned_by_business_forbidden():
    with pytest.raises(ForbiddenError):
        ensure_venue_owned_by_business({"brand": {"id": 999}}, business_id=10)