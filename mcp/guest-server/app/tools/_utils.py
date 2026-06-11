import json

from app.http_client import APIResponse


def unwrap_api_result(result: APIResponse) -> str:
    """
    Convert a Guest API response into a JSON string for the MCP client.

    For 4xx errors the original API message is preserved so the AI assistant
    can surface actionable context to the user. 5xx and transport errors are
    sanitized to avoid leaking internal details.
    """
    if result.get("is_error"):
        status = result.get("status")
        msg = result.get("error_message", "Unknown error")

        if status == 401:
            return json.dumps({"error": "unauthorized", "message": msg})
        if status == 403:
            return json.dumps({"error": "forbidden", "message": msg})
        if status == 404:
            return json.dumps({"error": "not_found", "message": msg})
        if status is None or (isinstance(status, int) and status >= 500):
            return json.dumps(
                {
                    "error": "downstream_failure",
                    "message": "Service temporarily unavailable. Please try again later.",
                }
            )
        return json.dumps({"error": "api_error", "message": msg})

    return json.dumps(result.get("data", {}))
