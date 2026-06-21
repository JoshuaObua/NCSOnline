"""Internal IP location and client-platform enrichment service."""

from __future__ import annotations

import ipaddress
import os
import time
from dataclasses import dataclass
from typing import Any

import httpx
from fastapi import Depends, FastAPI, Header, HTTPException
from pydantic import BaseModel, Field
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
    source: str = "ipwho.is"


@dataclass
class CacheValue:
    expires_at: float
    payload: dict[str, Any]


app = FastAPI(title="NCS Location Service", version="1.0.0")
_cache: dict[str, CacheValue] = {}
_ttl = int(os.getenv("LOCATION_CACHE_TTL_SECONDS", "86400"))
_provider_url = os.getenv("LOCATION_PROVIDER_URL", "https://ipwho.is/{ip}")
_internal_token = os.getenv("LOCATION_INTERNAL_TOKEN", "")


def require_internal_token(x_internal_token: str = Header(default="")) -> None:
    if _internal_token and x_internal_token != _internal_token:
        raise HTTPException(status_code=401, detail="invalid internal token")


def platform_details(user_agent: str) -> tuple[str, str, str]:
    parsed = parse_user_agent(user_agent or "")
    platform = parsed.os.family or "Unknown"
    browser = parsed.browser.family or "Unknown"
    if parsed.is_mobile:
        device = "Mobile"
    elif parsed.is_tablet:
        device = "Tablet"
    elif parsed.is_pc:
        device = "Desktop"
    elif parsed.is_bot:
        device = "Bot"
    else:
        device = "Unknown"
    return platform, browser, device


async def provider_lookup(ip: str) -> dict[str, Any]:
    cached = _cache.get(ip)
    if cached and cached.expires_at > time.time():
        return cached.payload

    address = ipaddress.ip_address(ip)
    if address.is_private or address.is_loopback or address.is_link_local:
        payload = {
            "country": "Uganda",
            "country_code": "UG",
            "region": "Internal network",
            "city": "Internal",
            "source": "internal",
        }
    else:
        async with httpx.AsyncClient(timeout=3.0) as client:
            response = await client.get(_provider_url.format(ip=ip))
            response.raise_for_status()
            raw = response.json()
        if raw.get("success") is False:
            raise HTTPException(status_code=502, detail="location provider rejected the address")
        security = raw.get("security") or {}
        connection = raw.get("connection") or {}
        timezone = raw.get("timezone") or {}
        payload = {
            "country": raw.get("country") or "Unknown",
            "country_code": raw.get("country_code") or "",
            "region": raw.get("region") or "",
            "city": raw.get("city") or "",
            "latitude": raw.get("latitude"),
            "longitude": raw.get("longitude"),
            "timezone": timezone.get("id") or "",
            "isp": connection.get("isp") or "",
            "is_proxy": bool(security.get("proxy")),
            "is_vpn": bool(security.get("vpn")),
            "is_tor": bool(security.get("tor")),
            "is_hosting": bool(security.get("hosting")),
            "source": "ipwho.is",
        }
    _cache[ip] = CacheValue(time.time() + _ttl, payload)
    return payload


@app.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok"}


@app.post("/v1/locate", response_model=LocateResponse, dependencies=[Depends(require_internal_token)])
async def locate(request: LocateRequest) -> LocateResponse:
    try:
        ipaddress.ip_address(request.ip)
    except ValueError as exc:
        raise HTTPException(status_code=422, detail="invalid IP address") from exc
    payload = await provider_lookup(request.ip)
    platform, browser, device = platform_details(request.user_agent)
    # user_id is echoed only for correlation. Authentication remains owned by the Go API.
    return LocateResponse(
        ip=request.ip,
        platform=platform,
        browser=browser,
        device_type=device,
        authenticated=request.authenticated and bool(request.user_id),
        user_id=request.user_id if request.authenticated else None,
        **payload,
    )
