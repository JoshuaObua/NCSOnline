"""NCS Geo-Sentinel: local GeoIP enrichment and edge request gate."""

from __future__ import annotations

import ipaddress
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
TRUSTED_PROXIES = [ipaddress.ip_network(v.strip()) for v in os.getenv("GEO_TRUSTED_PROXIES", "10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,127.0.0.0/8").split(",") if v.strip()]
CACHE_TTL = int(os.getenv("LOCATION_CACHE_TTL_SECONDS", "86400"))

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
    async with httpx.AsyncClient(timeout=30.0, follow_redirects=False) as client:
        upstream = await client.request(request.method, UPSTREAM_URL + original_path, params=request.query_params, content=body, headers=headers)
    response_headers = {k: v for k, v in upstream.headers.items() if k.lower() not in HOP_HEADERS}
    return Response(content=upstream.content, status_code=upstream.status_code, headers=response_headers)
