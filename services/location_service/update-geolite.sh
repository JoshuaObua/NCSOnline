#!/usr/bin/env sh
set -eu
: "${MAXMIND_ACCOUNT_ID:?MAXMIND_ACCOUNT_ID is required}"
: "${MAXMIND_LICENSE_KEY:?MAXMIND_LICENSE_KEY is required}"
target="${GEOIP_DB_PATH:-/data/GeoLite2-City.mmdb}"
tmp="$(mktemp -d)"; trap 'rm -rf "$tmp"' EXIT INT TERM
url="https://download.maxmind.com/geoip/databases/GeoLite2-City/download?suffix=tar.gz"
curl --fail --silent --show-error --location --user "$MAXMIND_ACCOUNT_ID:$MAXMIND_LICENSE_KEY" "$url" -o "$tmp/db.tar.gz"
tar -xzf "$tmp/db.tar.gz" -C "$tmp"
db="$(find "$tmp" -name GeoLite2-City.mmdb -type f | head -1)"
test -n "$db"
install -m 0644 "$db" "$target.new"
mv "$target.new" "$target"
