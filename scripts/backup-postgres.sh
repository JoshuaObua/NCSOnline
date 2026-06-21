#!/usr/bin/env sh
set -eu
umask 077

: "${DATABASE_URL:?DATABASE_URL is required}"
BACKUP_DIR="${BACKUP_DIR:-/var/backups/ncs}"
RETENTION_DAYS="${BACKUP_RETENTION_DAYS:-7}"
mkdir -p "$BACKUP_DIR"
stamp="$(date -u +%Y%m%d-%H%M%S)"
tmp="$BACKUP_DIR/.ncs-$stamp.dump.tmp"
target="$BACKUP_DIR/ncs-$stamp.dump"

cleanup(){ rm -f "$tmp"; }
trap cleanup EXIT INT TERM
pg_dump --dbname="$DATABASE_URL" --format=custom --compress=9 --no-owner --no-privileges --file="$tmp"
pg_restore --list "$tmp" >/dev/null
mv "$tmp" "$target"
sha256sum "$target" > "$target.sha256"

if [ -n "${RCLONE_REMOTE:-}" ]; then
  rclone copy "$target" "$target.sha256" "$RCLONE_REMOTE"
fi

find "$BACKUP_DIR" -type f -name 'ncs-*.dump*' -mtime "+$RETENTION_DAYS" -delete
printf '%s\n' "$target"
