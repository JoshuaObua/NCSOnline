#!/usr/bin/env sh
set -eu
forbidden='/var/www/ncs/uploads|docker compose down -v|docker-compose down -v|rm -rf.*uploads|volume rm'
if grep -R -n -E "$forbidden" scripts .github --exclude='check-deploy-safety.sh'; then
  echo 'Unsafe deployment command found.' >&2
  exit 1
fi
echo 'Deployment scripts preserve uploads and persistent volumes.'
