from types import SimpleNamespace

from app.main import platform_details, resolve_client_ip


def test_mobile_platform_detection():
    platform, browser, device = platform_details(
        "Mozilla/5.0 (Linux; Android 14) AppleWebKit/537.36 Chrome/124 Mobile Safari/537.36"
    )
    assert platform == "Android"
    assert browser == "Chrome Mobile"
    assert device == "Mobile"


def test_untrusted_peer_cannot_spoof_forwarded_header():
    request = SimpleNamespace(client=SimpleNamespace(host="8.8.8.8"), headers={"x-forwarded-for":"1.1.1.1"})
    assert resolve_client_ip(request) == "8.8.8.8"


def test_trusted_proxy_walks_chain_from_right_to_left():
    request = SimpleNamespace(client=SimpleNamespace(host="172.20.0.5"), headers={"x-forwarded-for":"8.8.8.8, 10.0.0.2"})
    assert resolve_client_ip(request) == "8.8.8.8"
