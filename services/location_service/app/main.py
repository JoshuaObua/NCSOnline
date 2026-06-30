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
    authenticated: bool = False
    user_id: str | None = None
    source: str = "geolite2"


app = FastAPI(title="NCS Geo-Sentinel", version="2.0.0")
INTERNAL_TOKEN = os.getenv("LOCATION_INTERNAL_TOKEN", "")
DB_PATH = Path(os.getenv("GEOIP_DB_PATH", "/data/GeoLite2-City.mmdb"))
WHITELIST_PATH = Path(os.getenv("GEO_WHITELIST_PATH", "/config/geo-whitelist.json"))
ALLOWED = {v.strip().upper() for v in os.getenv("GEO_ALLOWED_COUNTRIES", "UG").split(",") if v.strip()}
FAIL_CLOSED = os.getenv("GEO_FAIL_CLOSED", "false").lower() == "true"
PROVIDER_FALLBACK = os.getenv("LOCATION_PROVIDER_FALLBACK", "false").lower() == "true"
PROVIDER_URL = os.getenv("LOCATION_PROVIDER_URL", "https://ipwho.is/{ip}")
UPSTREAM_URL = os.getenv("GEO_UPSTREAM_URL", "http://backend:8080").rstrip("/")
FRONTEND_URL = os.getenv("GEO_FRONTEND_URL", "http://frontend:80").rstrip("/")
JWT_SECRET = os.getenv("JWT_SECRET", "")
TRUSTED_PROXIES = [ipaddress.ip_network(v.strip()) for v in os.getenv("GEO_TRUSTED_PROXIES", "10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,127.0.0.0/8").split(",") if v.strip()]
CACHE_TTL = int(os.getenv("LOCATION_CACHE_TTL_SECONDS", "86400"))
MAINTENANCE_CACHE_TTL = float(os.getenv("MAINTENANCE_CACHE_TTL_SECONDS", "2"))
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
    parsed = parse_user_agent(user_agent or "")
    device = "Mobile" if parsed.is_mobile else "Tablet" if parsed.is_tablet else "Desktop" if parsed.is_pc else "Bot" if parsed.is_bot else "Unknown"
    return parsed.os.family or "Unknown", parsed.browser.family or "Unknown", device


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
            _cache[ip] = (time.time() + CACHE_TTL, payload)
            return payload
        except Exception:
            pass
    if PROVIDER_FALLBACK:
        async with httpx.AsyncClient(timeout=3.0) as client:
            response = await client.get(PROVIDER_URL.format(ip=ip))
            response.raise_for_status()
            raw = response.json()
        security, connection, timezone = raw.get("security") or {}, raw.get("connection") or {}, raw.get("timezone") or {}
        payload = {"country": raw.get("country") or "Unknown", "country_code": raw.get("country_code") or "", "region": raw.get("region") or "", "city": raw.get("city") or "", "latitude": raw.get("latitude"), "longitude": raw.get("longitude"), "timezone": timezone.get("id") or "", "isp": connection.get("isp") or "", "is_proxy": bool(security.get("proxy")), "is_vpn": bool(security.get("vpn")), "is_tor": bool(security.get("tor")), "is_hosting": bool(security.get("hosting")), "source": "provider"}
        _cache[ip] = (time.time() + CACHE_TTL, payload)
        return payload
    return {"country": "Unknown", "country_code": "", "source": "unavailable"}


@app.get("/health")
def health() -> dict[str, Any]:
    return {"status": "ok", "geo_database_available": DB_PATH.exists(), "fail_closed": FAIL_CLOSED}


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
    return LocateResponse(ip=request.ip, platform=platform, browser=browser, device_type=device, authenticated=request.authenticated and bool(request.user_id), user_id=request.user_id if request.authenticated else None, **payload)


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
    roles = _verified_roles(authz)
    allowed_roles = {str(v).strip().lower() for v in rules.get("allowed_roles") or []}
    if roles and any(role in allowed_roles for role in roles):
        return True
    return False


def _b64url_decode(value: str) -> bytes:
    return base64.urlsafe_b64decode(value + "=" * (-len(value) % 4))


def _verified_roles(authz: str) -> set[str]:
    if not JWT_SECRET or not authz.lower().startswith("bearer "):
        return set()
    token = authz[7:].strip()
    parts = token.split(".")
    if len(parts) != 3:
        return set()
    signing_input = f"{parts[0]}.{parts[1]}".encode()
    expected = hmac.new(JWT_SECRET.encode(), signing_input, hashlib.sha256).digest()
    try:
        actual = _b64url_decode(parts[2])
        header = json.loads(_b64url_decode(parts[0]))
        payload = json.loads(_b64url_decode(parts[1]))
    except Exception:
        return set()
    if header.get("alg") != "HS256" or not hmac.compare_digest(expected, actual):
        return set()
    if payload.get("iss") != "ncsms-api":
        return set()
    if payload.get("exp") and float(payload["exp"]) < time.time():
        return set()
    return {str(v).strip().lower() for v in payload.get("roles") or []}


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
    html = f"""<!doctype html><html lang="en"><meta name="viewport" content="width=device-width"><title>{title}</title><body style="font-family:system-ui;background:#111827;color:white;display:grid;place-items:center;min-height:100vh;margin:0"><main style="max-width:42rem;text-align:center;padding:2rem"><h1>{title}</h1><p style="color:#d1d5db">{message}</p></main></body></html>"""
    return Response(content=html, status_code=503, media_type="text/html; charset=utf-8")


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
    async with httpx.AsyncClient(timeout=30.0, follow_redirects=False) as client:
        upstream = await client.request(request.method, target + original_path, params=request.query_params, content=body, headers=headers)
    response_headers = {k: v for k, v in upstream.headers.items() if k.lower() not in HOP_HEADERS}
    return Response(content=upstream.content, status_code=upstream.status_code, headers=response_headers)
