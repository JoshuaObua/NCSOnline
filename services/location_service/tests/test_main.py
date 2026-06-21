from app.main import platform_details


def test_mobile_platform_detection():
    platform, browser, device = platform_details(
        "Mozilla/5.0 (Linux; Android 14) AppleWebKit/537.36 Chrome/124 Mobile Safari/537.36"
    )
    assert platform == "Android"
    assert browser == "Chrome Mobile"
    assert device == "Mobile"
