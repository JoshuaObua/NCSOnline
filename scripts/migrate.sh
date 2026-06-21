#!/usr/bin/env bash
set -euo pipefail

BASELINE="${1:-}"
if [[ "$BASELINE" == "--baseline" ]]; then BASELINE="${2:?baseline version required}"; else BASELINE="${MIGRATION_BASELINE:-0}"; fi
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

psql_cmd=(docker compose exec -T postgres sh -lc 'psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB"')
printf '%s\n' "CREATE TABLE IF NOT EXISTS schema_migrations(version INTEGER PRIMARY KEY,file_name TEXT NOT NULL UNIQUE,applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW());" | "${psql_cmd[@]}"

for file in backend/migrations/[0-9][0-9][0-9]_*.sql; do
  name="$(basename "$file")"; version="${name%%_*}"; version="$((10#$version))"
  applied="$(printf "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=%d);\n" "$version" | "${psql_cmd[@]}" -At)"
  [[ "$applied" == "t" ]] && continue
  if (( version <= BASELINE )); then
    printf "INSERT INTO schema_migrations(version,file_name) VALUES(%d,'%s') ON CONFLICT DO NOTHING;\n" "$version" "$name" | "${psql_cmd[@]}"
    echo "baselined $name"
    continue
  fi
  echo "applying $name"
  cat "$file" | "${psql_cmd[@]}"
  printf "INSERT INTO schema_migrations(version,file_name) VALUES(%d,'%s');\n" "$version" "$name" | "${psql_cmd[@]}"
done
