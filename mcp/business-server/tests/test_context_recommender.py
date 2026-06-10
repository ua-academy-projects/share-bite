from app.context_recommender import parse_context, recommend_venues_by_context


VENUES = [
    {
        "id": 1,
        "name": "Candle Table",
        "description": "Quiet wine bar for dinner",
        "tags": ["romantic", "quiet", "wine"],
        "price_level": "$$$",
        "rating": 4.8,
        "latitude": 50.45,
        "longitude": 30.52,
    },
    {
        "id": 2,
        "name": "Campus Bowl",
        "description": "Affordable lunch and student deals",
        "tags": ["budget", "student-friendly", "lunch"],
        "price_level": "$",
        "rating": 4.2,
        "latitude": 50.46,
        "longitude": 30.53,
    },
    {
        "id": 3,
        "name": "Big Table Hub",
        "description": "Casual place for friends and teams",
        "tags": ["group-friendly", "casual", "beer"],
        "price_level": "$$",
        "rating": 4.4,
        "latitude": 50.5,
        "longitude": 30.6,
    },
]


def test_recommend_venues_detects_date_context():
    result = recommend_venues_by_context(
        {
            "query": "romantic date dinner",
            "budget_level": "premium",
            "party_size": 2,
            "lat": 50.45,
            "lon": 30.52,
        },
        VENUES,
    )

    assert result[0]["venue"]["id"] == 1
    assert result[0]["intent"] == "date"
    assert "matches a date-like context" in result[0]["reasons"]


def test_recommend_venues_detects_budget_context():
    result = recommend_venues_by_context(
        {
            "query": "cheap place for students",
            "budget": "low",
            "party_size": 3,
            "lat": 50.45,
            "lon": 30.52,
        },
        VENUES,
    )

    assert result[0]["venue"]["id"] == 2
    assert result[0]["intent"] == "budget"
    assert "fits the requested budget" in result[0]["reasons"]


def test_recommend_venues_detects_meetup_context():
    result = recommend_venues_by_context(
        {
            "query": "meetup with friends after work",
            "budget": "moderate",
            "party_size": 6,
        },
        VENUES,
    )

    assert result[0]["venue"]["id"] == 3
    assert result[0]["intent"] == "meetup"
    assert "works for the group size" in result[0]["reasons"]


def test_parse_context_uses_explicit_intent():
    context = parse_context({"intent": "date", "tags": ["Quiet"]})

    assert context.intent == "date"
    assert context.desired_tags == {"quiet"}
