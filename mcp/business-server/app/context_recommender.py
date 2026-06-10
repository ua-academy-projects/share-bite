from __future__ import annotations

from dataclasses import dataclass
from math import atan2, cos, radians, sin, sqrt
from typing import Any


INTENT_KEYWORDS: dict[str, set[str]] = {
    "date": {"date", "romantic", "anniversary", "couple", "cozy", "quiet", "wine", "dinner"},
    "meetup": {"meetup", "friends", "group", "team", "coworking", "casual", "lunch", "brunch"},
    "budget": {"budget", "cheap", "affordable", "student", "deal", "discount", "value"},
}

INTENT_TAGS: dict[str, set[str]] = {
    "date": {"romantic", "cozy", "quiet", "wine", "dinner", "dessert"},
    "meetup": {"group-friendly", "casual", "brunch", "lunch", "beer", "outdoor"},
    "budget": {"budget", "affordable", "student-friendly", "fast-food", "street-food"},
}

INTENT_REASON: dict[str, str] = {
    "date": "matches a date-like context",
    "meetup": "fits a meetup or group plan",
    "budget": "fits a budget-conscious plan",
}


@dataclass(frozen=True)
class RecommendationContext:
    query: str
    intent: str
    budget_level: str | None
    party_size: int | None
    latitude: float | None
    longitude: float | None
    desired_tags: set[str]


def recommend_venues_by_context(
    context: dict[str, Any],
    venues: list[dict[str, Any]],
    limit: int = 10,
) -> list[dict[str, Any]]:
    """
    Rank venue candidates using a lightweight, deterministic AI-style model.

    It extracts intent from natural-language context and combines semantic tag
    match, budget fit, group fit, distance, rating, and popularity. The model is
    explainable and does not require network calls or external API keys.
    """
    parsed_context = parse_context(context)
    ranked = [_score_venue(parsed_context, venue) for venue in venues if isinstance(venue, dict)]
    ranked.sort(key=lambda item: item["score"], reverse=True)
    return ranked[: max(1, min(limit, 50))]


def parse_context(context: dict[str, Any]) -> RecommendationContext:
    text_parts = [
        str(context.get("query") or ""),
        str(context.get("occasion") or ""),
        str(context.get("notes") or ""),
    ]
    raw_tags = context.get("tags") or context.get("desired_tags") or []
    if isinstance(raw_tags, str):
        raw_tags = [raw_tags]

    normalized_text = " ".join(text_parts).lower()
    desired_tags = {_normalize_token(tag) for tag in raw_tags if str(tag).strip()}
    desired_tags.discard("")

    intent = str(context.get("intent") or "").lower().strip()
    if intent not in INTENT_KEYWORDS:
        intent = _infer_intent(normalized_text, desired_tags)

    budget_level = context.get("budget_level") or context.get("budget")
    if budget_level is not None:
        budget_level = str(budget_level).lower().strip()

    return RecommendationContext(
        query=normalized_text,
        intent=intent,
        budget_level=budget_level,
        party_size=_optional_int(context.get("party_size")),
        latitude=_optional_float(context.get("latitude") or context.get("lat")),
        longitude=_optional_float(context.get("longitude") or context.get("lon")),
        desired_tags=desired_tags,
    )


def _infer_intent(text: str, desired_tags: set[str]) -> str:
    scores: dict[str, int] = {}
    for intent, keywords in INTENT_KEYWORDS.items():
        tag_hits = len(desired_tags & INTENT_TAGS[intent])
        text_hits = sum(1 for keyword in keywords if keyword in text)
        scores[intent] = tag_hits * 2 + text_hits

    best_intent, best_score = max(scores.items(), key=lambda item: item[1])
    return best_intent if best_score > 0 else "general"


def _score_venue(context: RecommendationContext, venue: dict[str, Any]) -> dict[str, Any]:
    tags = {_normalize_token(tag) for tag in venue.get("tags", []) if str(tag).strip()}
    name = str(venue.get("name") or "")
    description = str(venue.get("description") or "")
    text = f"{name} {description}".lower()

    score = 0.0
    reasons: list[str] = []

    intent_hits = tags & INTENT_TAGS.get(context.intent, set())
    if intent_hits:
        score += 36 + len(intent_hits) * 4
        reasons.append(INTENT_REASON.get(context.intent, "matches the requested context"))

    desired_hits = tags & context.desired_tags
    if desired_hits:
        score += 28 + len(desired_hits) * 3
        reasons.append("matches requested tags")

    if context.intent != "general" and context.intent in text:
        score += 8

    budget_score = _budget_score(context.budget_level, venue.get("price_level"))
    if budget_score:
        score += budget_score
        reasons.append("fits the requested budget")

    party_score = _party_score(context.party_size, tags)
    if party_score:
        score += party_score
        reasons.append("works for the group size")

    distance_km = _distance_from_context(context, venue)
    if distance_km is not None:
        score += max(0.0, 18.0 - min(distance_km, 18.0))
        if distance_km <= 3:
            reasons.append("near the requested area")

    rating = _optional_float(venue.get("rating"))
    if rating is not None:
        score += max(0.0, min(rating, 5.0)) * 3

    popularity = _optional_float(venue.get("popularity_score"))
    if popularity is not None:
        score += max(0.0, min(popularity, 100.0)) / 10

    if not reasons:
        reasons.append("best available general match")

    return {
        "venue": venue,
        "score": round(score, 2),
        "intent": context.intent,
        "distance_km": round(distance_km, 2) if distance_km is not None else None,
        "reasons": reasons[:4],
    }


def _budget_score(requested: str | None, venue_price: Any) -> float:
    if requested is None or venue_price is None:
        return 0.0

    requested_rank = _price_rank(requested)
    venue_rank = _price_rank(venue_price)
    if requested_rank is None or venue_rank is None:
        return 0.0

    diff = abs(requested_rank - venue_rank)
    if diff == 0:
        return 18.0
    if diff == 1:
        return 8.0
    return -8.0


def _party_score(party_size: int | None, tags: set[str]) -> float:
    if party_size is None:
        return 0.0
    if party_size >= 4 and {"group-friendly", "casual", "outdoor"} & tags:
        return 14.0
    if party_size <= 2 and {"cozy", "quiet", "romantic"} & tags:
        return 12.0
    return 0.0


def _distance_from_context(context: RecommendationContext, venue: dict[str, Any]) -> float | None:
    venue_lat = _optional_float(venue.get("latitude") or venue.get("lat"))
    venue_lon = _optional_float(venue.get("longitude") or venue.get("lon"))
    if None in (context.latitude, context.longitude, venue_lat, venue_lon):
        return None
    return _haversine_km(context.latitude, context.longitude, venue_lat, venue_lon)


def _haversine_km(lat1: float, lon1: float, lat2: float, lon2: float) -> float:
    radius_km = 6371.0
    dlat = radians(lat2 - lat1)
    dlon = radians(lon2 - lon1)
    a = sin(dlat / 2) ** 2 + cos(radians(lat1)) * cos(radians(lat2)) * sin(dlon / 2) ** 2
    return radius_km * 2 * atan2(sqrt(a), sqrt(1 - a))


def _price_rank(value: Any) -> int | None:
    if isinstance(value, int | float):
        return int(value)

    normalized = _normalize_token(value)
    mapping = {
        "$": 1,
        "cheap": 1,
        "budget": 1,
        "low": 1,
        "$$": 2,
        "moderate": 2,
        "medium": 2,
        "$$$": 3,
        "expensive": 3,
        "high": 3,
        "$$$$": 4,
        "premium": 4,
        "luxury": 4,
    }
    return mapping.get(normalized)


def _optional_float(value: Any) -> float | None:
    if value is None or value == "":
        return None
    try:
        return float(value)
    except (TypeError, ValueError):
        return None


def _optional_int(value: Any) -> int | None:
    if value is None or value == "":
        return None
    try:
        return int(value)
    except (TypeError, ValueError):
        return None


def _normalize_token(value: Any) -> str:
    return str(value).lower().strip().replace("_", "-").replace(" ", "-")
