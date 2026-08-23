"""NCS Geo-Sentinel: local GeoIP enrichment and edge request gate."""

from __future__ import annotations

import ipaddress
import base64
import hmac
import hashlib
import json
import os
import threading
import time
from collections import Counter
from datetime import datetime
from html import escape
from pathlib import Path
from typing import Any

import httpx
from fastapi import Depends, FastAPI, Header, HTTPException, Request, Response
from fastapi.responses import PlainTextResponse, RedirectResponse
from geoip2.database import Reader
from pydantic import BaseModel
from user_agents import parse as parse_user_agent


class LocateRequest(BaseModel):
    ip: str
    user_agent: str = ""
    authenticated: bool = False
    user_id: str | None = None


class LocateResponse(BaseModel):
    ip: str
    country: str = "Unknown"
    country_code: str = ""
    region: str = ""
    city: str = ""
    latitude: float | None = None
    longitude: float | None = None
    timezone: str = ""
    isp: str = ""
    is_proxy: bool = False
    is_vpn: bool = False
    is_tor: bool = False
    is_hosting: bool = False
    platform: str = "Unknown"
    browser: str = "Unknown"
    device_type: str = "Unknown"
    client_type: str = "Unknown"
    authenticated: bool = False
    user_id: str | None = None
    source: str = "geolite2"


app = FastAPI(title="NCS Geo-Sentinel", version="2.0.0")
INTERNAL_TOKEN = os.getenv("LOCATION_INTERNAL_TOKEN", "")
DB_PATH = Path(os.getenv("GEOIP_DB_PATH", "/data/GeoLite2-City.mmdb"))
WHITELIST_PATH = Path(os.getenv("GEO_WHITELIST_PATH", "/config/geo-whitelist.json"))
ALLOWED = {v.strip().upper() for v in os.getenv("GEO_ALLOWED_COUNTRIES", "UG").split(",") if v.strip()}
FAIL_CLOSED = os.getenv("GEO_FAIL_CLOSED", "false").lower() == "true"
PROVIDER_FALLBACK = os.getenv("LOCATION_PROVIDER_FALLBACK", "true").lower() == "true"
PROVIDER_URL = os.getenv("LOCATION_PROVIDER_URL", "")
UPSTREAM_URL = os.getenv("GEO_UPSTREAM_URL", os.getenv("BACKEND_INTERNAL_URL", "")).rstrip("/")
FRONTEND_URL = os.getenv("GEO_FRONTEND_URL", os.getenv("FRONTEND_INTERNAL_URL", "")).rstrip("/")
JWT_SECRET = os.getenv("JWT_SECRET", "")
TRUSTED_PROXIES = [ipaddress.ip_network(v.strip()) for v in os.getenv("GEO_TRUSTED_PROXIES", "").split(",") if v.strip()]
CACHE_TTL = int(os.getenv("LOCATION_CACHE_TTL_SECONDS", "86400"))
MAINTENANCE_CACHE_TTL = float(os.getenv("MAINTENANCE_CACHE_TTL_SECONDS", "2"))
PROVIDER_TIMEOUT = float(os.getenv("LOCATION_PROVIDER_TIMEOUT", "3.0"))
GATE_TIMEOUT = float(os.getenv("LOCATION_GATE_TIMEOUT", "30.0"))
_maintenance_cache: tuple[float, dict[str, Any]] | None = None

_reader: Reader | None = None
_reader_mtime = 0.0
_reader_lock = threading.Lock()
_cache: dict[str, tuple[float, dict[str, Any]]] = {}
_counters: Counter[tuple[str, str]] = Counter()
_whitelist_mtime = 0.0
_whitelist: list[ipaddress._BaseNetwork] = []


def require_internal_token(x_internal_token: str = Header(default="")) -> None:
    if INTERNAL_TOKEN and x_internal_token != INTERNAL_TOKEN:
        raise HTTPException(status_code=401, detail="invalid internal token")


def _trusted(value: str) -> bool:
    try:
        address = ipaddress.ip_address(value)
    except ValueError:
        return False
    return any(address in network for network in TRUSTED_PROXIES)


def resolve_client_ip(request: Request) -> str:
    peer = request.client.host if request.client else ""
    if not _trusted(peer):
        return peer
    cf_ip = request.headers.get("cf-connecting-ip", "").strip()
    if cf_ip:
        try:
            return str(ipaddress.ip_address(cf_ip))
        except ValueError:
            pass
    chain = [v.strip() for v in request.headers.get("x-forwarded-for", "").split(",") if v.strip()]
    chain.append(peer)
    for value in reversed(chain):
        try:
            normalized = str(ipaddress.ip_address(value))
        except ValueError:
            continue
        if not _trusted(normalized):
            return normalized
    real_ip = request.headers.get("x-real-ip", "").strip()
    try:
        return str(ipaddress.ip_address(real_ip))
    except ValueError:
        return peer


def _load_whitelist() -> list[ipaddress._BaseNetwork]:
    global _whitelist_mtime, _whitelist
    try:
        mtime = WHITELIST_PATH.stat().st_mtime
        if mtime != _whitelist_mtime:
            raw = json.loads(WHITELIST_PATH.read_text(encoding="utf-8"))
            now = time.time()
            networks = []
            for item in raw:
                if item.get("expires_at_epoch") and float(item["expires_at_epoch"]) < now:
                    continue
                networks.append(ipaddress.ip_network(item["cidr"], strict=False))
            _whitelist, _whitelist_mtime = networks, mtime
    except (OSError, ValueError, TypeError, json.JSONDecodeError):
        pass
    return _whitelist


def _open_reader() -> Reader | None:
    global _reader, _reader_mtime
    try:
        mtime = DB_PATH.stat().st_mtime
    except OSError:
        return None
    with _reader_lock:
        if _reader is None or mtime != _reader_mtime:
            if _reader is not None:
                _reader.close()
            _reader, _reader_mtime = Reader(str(DB_PATH)), mtime
    return _reader



def platform_details(user_agent: str) -> tuple[str, str, str]:
    value = user_agent or ""
    lower = value.lower()
    parsed = parse_user_agent(value)
    device = "Mobile" if parsed.is_mobile else "Tablet" if parsed.is_tablet else "Desktop" if parsed.is_pc else "Bot" if parsed.is_bot else "Unknown"

    # Priority browser detection (Opera/Edge/Samsung MUST precede Chrome)
    if "postman" in lower:
        browser = "Postman"
    elif "insomnia" in lower:
        browser = "Insomnia"
    elif "opr/" in lower or "opera" in lower or "opt/" in lower or "opios/" in lower:
        browser = "Opera Mobile" if device in {"Mobile", "Tablet"} or "mobile" in lower or "android" in lower else "Opera"
    elif "samsungbrowser/" in lower:
        browser = "Samsung Internet"
    elif "edg/" in lower or "edge/" in lower or "edga/" in lower or "edgios/" in lower:
        browser = "Edge Mobile" if device in {"Mobile", "Tablet"} or "mobile" in lower else "Edge"
    elif "brave" in lower:
        browser = "Brave"
    elif "vivaldi" in lower:
        browser = "Vivaldi"
    elif "ucbrowser" in lower or "ubrowser" in lower:
        browser = "UC Browser"
    elif "firefox/" in lower or "fxios/" in lower:
        browser = "Firefox Mobile" if device in {"Mobile", "Tablet"} or "mobile" in lower else "Firefox"
    elif "chrome/" in lower or "crios/" in lower:
        browser = "Chrome Mobile" if device in {"Mobile", "Tablet"} or "mobile" in lower else "Chrome"
    elif "safari/" in lower and "chrome/" not in lower and "crios/" not in lower:
        browser = "Mobile Safari" if device in {"Mobile", "Tablet"} or "mobile" in lower else "Safari"
    else:
        browser = parsed.browser.family or "Unknown"

    platform = parsed.os.family or "Unknown"
    if "android" in lower:
        platform = "Android"
    elif "iphone" in lower or "ipad" in lower or "ios" in lower:
        platform = "iOS"
    elif "windows nt 10" in lower or "windows nt 11" in lower or "windows 10" in lower or "windows 11" in lower:
        platform = "Windows"
    elif "mac os x" in lower or "macos" in lower:
        platform = "macOS"
    elif "linux" in lower:
        platform = "Linux"

    if platform == "Other":
        platform = "Unknown"
    if browser == "Other":
        browser = "Unknown"
    return platform, browser, device


def client_type(user_agent: str) -> str:
    value = (user_agent or "").lower()
    if "ussd" in value or "africastalking" in value or "africa's talking" in value:
        return "USSD Gateway"
    if "postman" in value:
        return "Postman"
    if "insomnia" in value:
        return "Insomnia"
    if "thunder-client" in value or "thunderclient" in value:
        return "Thunder Client"
    if "curl" in value:
        return "curl"
    if "httpie" in value:
        return "HTTPie"
    if "python-requests" in value:
        return "Python/requests"
    if "go-http-client" in value:
        return "Go/HTTP"
    if "okhttp" in value:
        return "Mobile App"
    if "mozilla" in value or "webkit" in value:
        return "Browser"
    return "Other" if value else "Unknown"


def _merge_location(primary: dict[str, Any], fallback: dict[str, Any]) -> dict[str, Any]:
    merged = dict(fallback)
    for key, value in primary.items():
        if value not in (None, "", "Unknown"):
            merged[key] = value
    if fallback and primary:
        merged["source"] = f"{primary.get('source', 'primary')}+{fallback.get('source', 'fallback')}"
    return merged


async def _provider_location(ip: str) -> dict[str, Any]:
    if not PROVIDER_FALLBACK:
        return {}
    urls = []
    if PROVIDER_URL:
        urls.append(PROVIDER_URL.format(ip=ip))
    urls.extend([
        f"https://ipwho.is/{ip}",
        f"https://freeipapi.com/api/json/{ip}",
    ])

    for url in urls:
        try:
            async with httpx.AsyncClient(timeout=PROVIDER_TIMEOUT) as client:
                response = await client.get(url)
                if response.status_code != 200:
                    continue
                raw = response.json()
            if raw.get("success") is False:
                continue
            
            # Format 1: ipwho.is / standard
            country = raw.get("country") or raw.get("countryName") or ""
            if not country or country == "Unknown":
                continue
            
            country_code = raw.get("country_code") or raw.get("countryCode") or ""
            region = raw.get("region") or raw.get("regionName") or ""
            city = raw.get("city") or raw.get("cityName") or ""
            latitude = raw.get("latitude")
            longitude = raw.get("longitude")
            
            timezone_raw = raw.get("timezone") or {}
            timezone_id = timezone_raw.get("id") if isinstance(timezone_raw, dict) else str(timezone_raw)
            
            connection = raw.get("connection") or {}
            isp = connection.get("isp") or raw.get("isp") or ""
            
            security = raw.get("security") or {}
            return {
                "country": country,
                "country_code": country_code,
                "region": region,
                "city": city,
                "latitude": latitude,
                "longitude": longitude,
                "timezone": timezone_id or "",
                "isp": isp,
                "is_proxy": bool(security.get("proxy")),
                "is_vpn": bool(security.get("vpn")),
                "is_tor": bool(security.get("tor")),
                "is_hosting": bool(security.get("hosting")),
                "source": "provider",
            }
        except Exception:
            continue
    return {}


async def locate_ip(ip: str) -> dict[str, Any]:
    address = ipaddress.ip_address(ip)
    if address.is_private or address.is_loopback or address.is_link_local:
        return {"country": "Uganda", "country_code": "UG", "region": "Internal network", "city": "Internal", "source": "internal"}
    cached = _cache.get(ip)
    if cached and cached[0] > time.time():
        return cached[1]
    reader = _open_reader()
    if reader:
        try:
            record = reader.city(ip)
            payload = {
                "country": record.country.name or "Unknown",
                "country_code": record.country.iso_code or "",
                "region": record.subdivisions.most_specific.name or "",
                "city": record.city.name or "",
                "latitude": record.location.latitude,
                "longitude": record.location.longitude,
                "timezone": record.location.time_zone or "",
                "source": "geolite2",
            }
            if not payload["city"] or not payload["region"] or not payload["timezone"] or payload["country"] == "Unknown":
                payload = _merge_location(payload, await _provider_location(ip))
            _cache[ip] = (time.time() + CACHE_TTL, payload)
            return payload
        except Exception:
            pass
    payload = await _provider_location(ip)
    if payload:
        _cache[ip] = (time.time() + CACHE_TTL, payload)
        return payload


@app.get("/metrics")
def metrics() -> Response:
    lines = ["# TYPE ncs_geo_requests_total counter"]
    for (country, action), count in sorted(_counters.items()):
        lines.append(f'ncs_geo_requests_total{{country="{country}",action="{action}"}} {count}')
    return PlainTextResponse("\n".join(lines) + "\n", media_type="text/plain; version=0.0.4")


@app.post("/v1/locate", response_model=LocateResponse, dependencies=[Depends(require_internal_token)])
async def locate(request: LocateRequest) -> LocateResponse:
    try:
        ipaddress.ip_address(request.ip)
    except ValueError as exc:
        raise HTTPException(status_code=422, detail="invalid IP address") from exc
    payload = await locate_ip(request.ip)
    platform, browser, device = platform_details(request.user_agent)
    return LocateResponse(ip=request.ip, platform=platform, browser=browser, device_type=device, client_type=client_type(request.user_agent), authenticated=request.authenticated and bool(request.user_id), user_id=request.user_id if request.authenticated else None, **payload)


HOP_HEADERS = {"connection", "keep-alive", "proxy-authenticate", "proxy-authorization", "te", "trailers", "transfer-encoding", "upgrade", "host", "content-length"}


def _is_admin_path(path: str) -> bool:
    return path.startswith(("/dashboard", "/admin", "/users", "/roles", "/applications", "/audit-logs", "/maintenance", "/settings", "/cms", "/nsmis"))


def _maintenance_scope(path: str) -> str:
    return "admin_dashboard" if _is_admin_path(path) or path.startswith("/api/v1/admin/") else "public_cms"


def _scope_active(scope: dict[str, Any]) -> bool:
    return bool(scope.get("is_active"))


def _maintenance_title(scope: dict[str, Any]) -> str:
    meta = scope.get("display_meta") or {}
    return meta.get("custom_title") or "Scheduled Maintenance"


def _maintenance_message(scope: dict[str, Any]) -> str:
    meta = scope.get("display_meta") or {}
    return meta.get("custom_message") or scope.get("reason") or "The platform is temporarily unavailable during scheduled maintenance."


async def _maintenance_status() -> dict[str, Any]:
    global _maintenance_cache
    now = time.time()
    if _maintenance_cache and _maintenance_cache[0] > now:
        return _maintenance_cache[1]
    try:
        async with httpx.AsyncClient(timeout=2.0) as client:
            response = await client.get(UPSTREAM_URL + "/api/v1/system/maintenance-status")
            response.raise_for_status()
            payload = response.json().get("data") or {}
    except Exception:
        payload = {}
    _maintenance_cache = (now + MAINTENANCE_CACHE_TTL, payload)
    return payload


def _has_bypass(request: Request, scope: dict[str, Any]) -> bool:
    rules = scope.get("bypass_rules") or {}
    token = (rules.get("secret_query_param") or "").strip()
    if token and (request.query_params.get("maintenance_bypass") == token or token in request.query_params):
        return True
    authz = request.headers.get("authorization", "")
    if not authz and request.cookies.get("ncsms_access"):
        authz = "Bearer " + request.cookies["ncsms_access"]
    claims = _verified_claims(authz)
    roles = {str(v).strip().lower() for v in claims.get("roles") or []}
    allowed_roles = {str(v).strip().lower() for v in rules.get("allowed_roles") or []}
    if roles and any(role in allowed_roles for role in roles):
        return True
    user_id = str(claims.get("user_id") or claims.get("sub") or "").strip()
    allowed_users = {str(v).strip() for v in rules.get("allowed_user_ids") or []}
    if user_id and user_id in allowed_users:
        return True
    return False


def _b64url_decode(value: str) -> bytes:
    return base64.urlsafe_b64decode(value + "=" * (-len(value) % 4))


def _verified_roles(authz: str) -> set[str]:
    payload = _verified_claims(authz)
    return {str(v).strip().lower() for v in payload.get("roles") or []}


def _verified_claims(authz: str) -> dict[str, Any]:
    if not JWT_SECRET or not authz.lower().startswith("bearer "):
        return {}
    token = authz[7:].strip()
    parts = token.split(".")
    if len(parts) != 3:
        return {}
    signing_input = f"{parts[0]}.{parts[1]}".encode()
    expected = hmac.new(JWT_SECRET.encode(), signing_input, hashlib.sha256).digest()
    try:
        actual = _b64url_decode(parts[2])
        header = json.loads(_b64url_decode(parts[0]))
        payload = json.loads(_b64url_decode(parts[1]))
    except Exception:
        return {}
    if header.get("alg") != "HS256" or not hmac.compare_digest(expected, actual):
        return {}
    if payload.get("iss") != "ncsms-api":
        return {}
    if payload.get("exp") and float(payload["exp"]) < time.time():
        return {}
    return payload


def _maintenance_response(request: Request, scope: dict[str, Any]) -> Response:
    accepts_json = "application/json" in request.headers.get("accept", "") or request.url.path.startswith("/api/")
    if accepts_json:
        return Response(
            content=json.dumps({"success": False, "error": {"code": "MAINTENANCE_MODE", "message": _maintenance_message(scope)}}),
            status_code=503,
            media_type="application/json",
        )
    title = _maintenance_title(scope)
    message = _maintenance_message(scope)
    html = _public_maintenance_html(title, message, scope.get("expected_end"))
    return Response(content=html, status_code=503, media_type="text/html; charset=utf-8")


def _public_maintenance_html(title: str, message: str, expected_end: Any) -> str:
    expected = "Shortly"
    if expected_end:
        try:
            parsed = datetime.fromisoformat(str(expected_end).replace("Z", "+00:00"))
            expected = parsed.astimezone().strftime("%a, %d %b at %H:%M %Z")
        except (TypeError, ValueError):
            expected = str(expected_end)
    document = """<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>__TITLE__ | NCS Uganda</title>
<style>
*{box-sizing:border-box}body{min-height:100vh;margin:0;display:flex;flex-direction:column;justify-content:center;overflow-x:hidden;background:linear-gradient(135deg,#102b4d 0%,#1a365d 52%,#0d223d 100%);color:#fff;font-family:Arial,sans-serif;padding:clamp(20px,4vw,56px)}body:before{position:fixed;inset:0;content:"";pointer-events:none;background-image:linear-gradient(rgba(255,255,255,.045) 1px,transparent 1px),linear-gradient(90deg,rgba(255,255,255,.045) 1px,transparent 1px);background-size:52px 52px;mask-image:linear-gradient(to bottom right,#000,transparent 78%)}.shell{position:relative;width:min(100%,1160px);margin:auto;display:grid;grid-template-columns:minmax(0,1.08fr) minmax(320px,.72fr);align-items:center;gap:clamp(32px,7vw,104px)}.brand{display:flex;align-items:center;gap:14px;margin-bottom:clamp(36px,7vh,80px);color:rgba(255,255,255,.86);font-size:13px;font-weight:700;text-transform:uppercase}.brand img{width:auto;max-width:168px;height:60px;object-fit:contain;border-radius:8px;background:#fff;padding:6px 11px}.kicker{display:flex;align-items:center;gap:9px;margin:0 0 18px;color:#f5a623;font-size:12px;font-weight:800;text-transform:uppercase}.kicker:before{width:28px;height:2px;content:"";background:#f5a623}h1{max-width:720px;margin:0;font-size:clamp(42px,6vw,86px);line-height:1.02;overflow-wrap:anywhere}p.message{max-width:650px;margin:24px 0 0;color:rgba(255,255,255,.76);font-size:clamp(16px,1.8vw,19px);line-height:1.75}.actions{display:flex;flex-wrap:wrap;gap:12px;margin-top:32px}.actions a{min-height:44px;display:inline-flex;align-items:center;justify-content:center;border:1px solid #f5a623;border-radius:6px;background:#f5a623;color:#102b4d;padding:12px 18px;text-decoration:none;font-size:14px;font-weight:800}.actions a.alt{border-color:rgba(255,255,255,.34);background:transparent;color:#fff}.panel{border:1px solid rgba(255,255,255,.22);border-radius:8px;background:#fff;color:#1a365d;padding:clamp(24px,4vw,38px);box-shadow:0 28px 80px rgba(4,18,36,.34)}.icon{width:56px;height:56px;display:grid;place-items:center;border-radius:8px;background:#fff5e4;color:#e2920f;font-size:29px;font-weight:700}.label{margin:24px 0 8px;color:#e2920f;font-size:12px;font-weight:800;text-transform:uppercase}.panel h2{margin:0;font-size:clamp(22px,2.6vw,28px);line-height:1.25}.meta{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:1px;margin-top:24px;overflow:hidden;border:1px solid #e5e9ef;border-radius:6px;background:#e5e9ef}.meta div{min-width:0;background:#f8fafc;padding:16px}.meta span,.meta strong{display:block}.meta span{margin-bottom:6px;color:#6b7280;font-size:11px;font-weight:700;text-transform:uppercase}.meta strong{font-size:14px;line-height:1.45;overflow-wrap:anywhere}.steps{display:grid;gap:14px;margin:24px 0 0;padding:22px 0 0;border-top:1px solid #e5e9ef;list-style:none}.steps li{display:grid;grid-template-columns:12px minmax(0,1fr);align-items:center;gap:11px;color:#8a94a3;font-size:14px;font-weight:700}.steps i{width:10px;height:10px;border:2px solid #cbd5e1;border-radius:50%}.steps .done{color:#237a4b}.steps .done i{border-color:#2fa96b;background:#2fa96b}.steps .active{color:#1a365d}.steps .active i{border-color:#f5a623;background:#f5a623;box-shadow:0 0 0 4px rgba(245,166,35,.18)}footer{position:relative;width:min(100%,1160px);margin:clamp(32px,7vh,72px) auto 0;color:rgba(255,255,255,.48);font-size:12px}@media(max-width:820px){body{justify-content:flex-start}.shell{grid-template-columns:1fr;gap:36px}.brand{margin-bottom:40px}.panel{max-width:620px}}@media(max-width:480px){body{padding:18px}.brand{align-items:flex-start;flex-direction:column}.brand img{height:52px;max-width:150px}.actions{display:grid}.actions a{width:100%}.meta{grid-template-columns:1fr}}
</style>
</head>
<body>
<main class="shell">
<section>
<div class="brand"><img src="/main-logo.png" alt="National Council of Sports logo"><span>NCS Uganda</span></div>
<p class="kicker">Public Portal Maintenance</p>
<h1>__TITLE__</h1>
<p class="message">__MESSAGE__</p>
<div class="actions"><a href="mailto:info@ncs.go.ug">Email NCS</a><a class="alt" href="tel:+256414254477">Call NCS</a></div>
</section>
<aside class="panel">
<div class="icon" aria-hidden="true">&#9881;</div>
<p class="label">Current Status</p>
<h2>Scheduled upgrades are in progress</h2>
<div class="meta"><div><span>Expected return</span><strong>__EXPECTED__</strong></div><div><span>System status</span><strong>In progress</strong></div></div>
<ol class="steps"><li class="done"><i></i>Updates started</li><li class="active"><i></i>Quality checks</li><li><i></i>Website restored</li></ol>
</aside>
</main>
<footer>&copy; 1964 - __YEAR__ National Council of Sports. All Rights Reserved.</footer>
</body>
</html>"""
    return (
        document
        .replace("__TITLE__", escape(title or "We'll be right back"))
        .replace("__MESSAGE__", escape(message))
        .replace("__EXPECTED__", escape(expected))
        .replace("__YEAR__", str(datetime.now().year))
    )


@app.api_route("/gate/{path:path}", methods=["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"])
async def gate(path: str, request: Request) -> Response:
    ip = resolve_client_ip(request)
    try:
        address = ipaddress.ip_address(ip)
    except ValueError:
        return Response(status_code=403)
    whitelisted = any(address in network for network in _load_whitelist())
    location = await locate_ip(ip)
    code = str(location.get("country_code") or "").upper()
    allowed = whitelisted or code in ALLOWED or (not code and not FAIL_CLOSED)
    _counters[(code or "UNKNOWN", "allow" if allowed else "block")] += 1
    original_path = "/" + path
    # Static SPA assets must always pass — otherwise the maintenance page itself
    # can never load the SPA, and authorised users with the bypass cookie/role
    # still see a white screen because chunks 503.
    _is_static_asset = (
        original_path.startswith("/assets/")
        or original_path.startswith("/static/")
        or original_path.startswith("/fonts/")
        or original_path.startswith("/img/")
        or original_path.startswith("/images/")
        or original_path == "/favicon.ico"
        or original_path == "/favicon.png"
        or original_path == "/main-logo.png"
        or original_path == "/main-logo-white.png"
        or original_path == "/robots.txt"
        or original_path == "/manifest.json"
        or original_path == "/manifest.webmanifest"
        or original_path == "/sw.js"
    )
    if (
        not _is_static_asset
        and request.method in {"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE"}
        and original_path not in {"/api/v1/auth/login", "/api/v1/auth/refresh", "/health", "/healthz", "/readyz", "/login", "/register"}
    ):
        status = await _maintenance_status()
        scope = status.get(_maintenance_scope(original_path)) or {}
        if _scope_active(scope) and not _has_bypass(request, scope):
            return _maintenance_response(request, scope)
    if not allowed:
        if original_path == "/ncs-ussd":
            return PlainTextResponse("END Secure System Violation", status_code=403)
        accepts_json = "application/json" in request.headers.get("accept", "") or original_path.startswith("/api/")
        if accepts_json:
            return Response(status_code=403)
        return RedirectResponse("https://www.google.com/", status_code=301)
    body = await request.body()
    headers = {k: v for k, v in request.headers.items() if k.lower() not in HOP_HEADERS}
    headers["x-sentinel-client-ip"] = ip
    target = UPSTREAM_URL if original_path.startswith("/api/") or original_path == "/ncs-ussd" else FRONTEND_URL
    async with httpx.AsyncClient(timeout=GATE_TIMEOUT, follow_redirects=False) as client:
        upstream = await client.request(request.method, target + original_path, params=request.query_params, content=body, headers=headers)
    response_headers = {k: v for k, v in upstream.headers.items() if k.lower() not in HOP_HEADERS}
    return Response(content=upstream.content, status_code=upstream.status_code, headers=response_headers)
