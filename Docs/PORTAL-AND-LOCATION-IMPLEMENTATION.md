# Portal, Organisation Provisioning, and Location Service

## Delivered journey

- Public pages no longer start applications. `/apply` only redirects to sign-in or `/my-portal/applications/new`.
- `/my-portal` is the individual workspace for federation, association, club, and renewal applications.
- Approval of Form 3 or Form 10 creates one organisation profile, an owner membership, and a secure invitation when the official email differs.
- Members select an organisation in `/my-portal` before opening its protected `/dashboard` context.
- Invitation tokens are random, stored as SHA-256 hashes, expire after 72 hours, and are bound to the signed-in email.

## Authenticated API

- `GET /api/v1/me/contexts`
- `GET /api/v1/organisations/{organisationID}`
- `POST /api/v1/organisation-invitations/accept` with `{ "token": "..." }`
- `POST /api/v1/admin/applications/{id}/approve` provisions Form 3/Form 10 registrations transactionally and idempotently.

## Location and audit

The Go API calls the internal Python service at `POST /v1/locate`. It validates the IP, caches provider results, parses the user agent, and returns country, region, city, coordinates, timezone, ISP, proxy/VPN/Tor/hosting signals, platform, browser, and device type.

Go remains authoritative for identity: it verifies the signed JWT and supplies the user ID. The outer audit middleware also validates signed bearer tokens on public routes, so a valid signed-in user is not stored as anonymous.

Configure `.env` from `.env.example`, run migration `021_portal_and_location.sql`, then start with `docker compose up --build`. IP coordinates are approximate and must not be presented as GPS-precise.
