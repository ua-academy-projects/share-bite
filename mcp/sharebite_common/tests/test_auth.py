from sharebite_common.auth import extract_auth_token, normalize_bearer_token


def test_normalize_bearer_token_adds_scheme():
    assert normalize_bearer_token("abc") == "Bearer abc"


def test_normalize_bearer_token_preserves_scheme():
    assert normalize_bearer_token("Bearer abc") == "Bearer abc"


def test_extract_auth_token_from_case_insensitive_header():
    assert extract_auth_token({"authorization": "Bearer abc"}) == "abc"
    assert extract_auth_token({"Authorization": "xyz"}) == "xyz"
