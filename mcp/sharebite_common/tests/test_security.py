import pytest

from sharebite_common.security import (
    ConfirmationRequiredError,
    REDACTED,
    SecretLeakError,
    assert_no_secrets,
    redact_secrets,
    require_confirmation,
)


def test_redact_secrets_nested_payload():
    payload = {"user": {"access_token": "secret", "name": "Ada"}}

    redacted = redact_secrets(payload)

    assert redacted["user"]["access_token"] == REDACTED
    assert redacted["user"]["name"] == "Ada"


def test_assert_no_secrets_rejects_secret_fields():
    with pytest.raises(SecretLeakError):
        assert_no_secrets({"data": [{"refresh_token": "secret"}]})


def test_require_confirmation():
    with pytest.raises(ConfirmationRequiredError):
        require_confirmation(False, "delete-user")

    require_confirmation(True, "delete-user")
