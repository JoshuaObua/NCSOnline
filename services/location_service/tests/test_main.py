from types import SimpleNamespace

import pytest

from app import main
from app.main import _merge_location, client_type, platform_details, resolve_client_ip


def test_mobile_platform_detection():
    platform, browser, device = platform_details(
        "Mozilla/5.0 (Linux; Android 14) AppleWebKit/537.36 Chrome/124 Mobile Safari/537.36"
    )
    assert platform == "Android"
    assert browser == "Chrome Mobile"
    assert device == "Mobile"


def test_api_and_ussd_client_detection():
    assert client_type("PostmanRuntime/7.43.0") == "Postman"
    assert client_type("Africa's Talking USSD Gateway") == "USSD Gateway"
    assert client_type("curl/8.11.0") == "curl"


def test_location_merge_fills_missing_city_without_losing_primary_country():
    merged = _merge_location(
        {"country": "Uganda", "country_code": "UG", "city": "", "source": "geolite2"},
        {"country": "Uganda", "city": "Kampala", "region": "Central Region", "source": "provider"},
    )
    assert merged["country"] == "Uganda"
    assert merged["city"] == "Kampala"
    assert merged["region"] == "Central Region"
    assert merged["source"] == "geolite2+provider"


def test_untrusted_peer_cannot_spoof_forwarded_header():
    request = SimpleNamespace(client=SimpleNamespace(host="8.8.8.8"), headers={"x-forwarded-for":"1.1.1.1"})
    assert resolve_client_ip(request) == "8.8.8.8"


def test_trusted_proxy_walks_chain_from_right_to_left(monkeypatch):
    monkeypatch.setattr(
        main,
        "TRUSTED_PROXIES",
        [main.ipaddress.ip_network("10.0.0.0/8"), main.ipaddress.ip_network("172.16.0.0/12")],
    )
    request = SimpleNamespace(client=SimpleNamespace(host="172.20.0.5"), headers={"x-forwarded-for":"8.8.8.8, 10.0.0.2"})
    assert resolve_client_ip(request) == "8.8.8.8"


@pytest.mark.asyncio
async def test_provider_fallback_returns_city_when_database_is_unavailable(monkeypatch):
    class Response:
        def raise_for_status(self):
            return None

        def json(self):
            return {
                "success": True,
                "country": "Uganda",
                "country_code": "UG",
                "region": "Central Region",
                "city": "Kampala",
                "latitude": 0.3476,
                "longitude": 32.5825,
                "timezone": {"id": "Africa/Kampala"},
            }

    class Client:
        async def __aenter__(self):
            return self

        async def __aexit__(self, *_):
            return None

        async def get(self, *_args, **_kwargs):
            return Response()

    monkeypatch.setattr(main, "_open_reader", lambda: None)
    monkeypatch.setattr(main, "PROVIDER_FALLBACK", True)
    monkeypatch.setattr(main, "PROVIDER_URL", "https://location.test/{ip}")
    monkeypatch.setattr(main.httpx, "AsyncClient", lambda **_kwargs: Client())
    main._cache.clear()

    result = await main.locate_ip("8.8.8.8")

    assert result["city"] == "Kampala"
    assert result["region"] == "Central Region"
    assert result["source"] == "provider"
