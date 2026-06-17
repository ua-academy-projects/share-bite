from __future__ import annotations

from collections.abc import Callable
from dataclasses import dataclass
from typing import Any

from .security import assert_no_secrets


class FakeMCP:
    def __init__(self) -> None:
        self.tools: dict[str, Callable[..., Any]] = {}
        self.resources: dict[str, Callable[..., Any]] = {}

    def tool(self, name: str | None = None, **_kwargs: Any):
        def decorator(func: Callable[..., Any]):
            self.tools[name or func.__name__] = func
            return func
        return decorator

    def resource(self, uri: str, **_kwargs: Any):
        def decorator(func: Callable[..., Any]):
            self.resources[uri] = func
            return func
        return decorator


@dataclass
class FakeMeta:
    headers: dict[str, str]

    def model_dump(self) -> dict[str, Any]:
        return {"headers": self.headers}


@dataclass
class FakeRequestContext:
    meta: FakeMeta


class FakeContext:
    def __init__(self, headers: dict[str, str] | None = None) -> None:
        self.request_context = FakeRequestContext(FakeMeta(headers or {}))


def assert_mcp_response_has_no_secrets(response: Any) -> None:
    assert_no_secrets(response)
