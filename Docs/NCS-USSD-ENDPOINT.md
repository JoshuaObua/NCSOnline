# NCS USSD Endpoint Guide

## Endpoint

```http
POST /ncs-ussd
Content-Type: application/x-www-form-urlencoded
X-USSD-Token: <configured callback secret>
```

The endpoint implements an aggregator-style USSD callback. Continuing responses start with `CON`; terminal responses start with `END`.

## Access model

| Menu | Service | Authentication |
|---|---|---|
| 1 | Check athlete registration | Public |
| 2 | Check licence or credential validity | Public |
| 3 | View the registered user’s latest applications | Registered phone number + NCS PIN |
| 4 | View protected renewal eligibility | Registered phone number + NCS PIN, then licence number |
| 5 | NCS help details | Public |

PINs must already have been established through authenticated `POST /api/v1/auth/pin/set`. The USSD endpoint never stores or returns a raw PIN. It compares the submitted 4–6 digit PIN to the account’s bcrypt hash.

Protected access is granted only when the entered phone is a valid Uganda `+256` number belonging to an active account, PIN setup is complete, the hash matches, and the provider-confirmed caller number matches the entered number. Authentication failures use one generic response.

## Provider fields

| Field | Description |
|---|---|
| `sessionId` | Provider session identifier |
| `serviceCode` | Dialled USSD code |
| `phoneNumber` | Provider-confirmed caller MSISDN |
| `text` | Full `*`-delimited input history; empty initially |

JSON with the same camel-case names is also accepted with `Content-Type: application/json`.

## Start a session

```bash
curl -X POST http://localhost:8080/ncs-ussd \
  -H "X-USSD-Token: $USSD_CALLBACK_SECRET" \
  -d "sessionId=test-001" \
  -d "serviceCode=*123#" \
  -d "phoneNumber=+256772123456" \
  -d "text="
```

```text
CON NCS Quick Services
1. Check athlete registration
2. Check licence validity
3. My applications
4. Renewal information
5. Help
```

## Public flows

Athlete registration:

```text
text=1
text=1*NCS-2026-AB12CD34
```

The result contains only masked identity, discipline, registration status, federation and verification date. It does not return DOB, contacts, identity documents or internal IDs.

Licence or credential validity:

```text
text=2
text=2*NCS-LIC-2026-001234
```

The lookup checks the `licences` registry, coach licence numbers and technical-official certification numbers. Past expiry dates are returned as `EXPIRED`.

## Protected flows

Latest applications:

```text
text=3
text=3*0772123456
text=3*0772123456*1234
```

Successful authentication returns at most three recent applications with their reference or draft identifier, status and payment status.

Renewal information:

```text
text=4
text=4*0772123456
text=4*0772123456*1234
text=4*0772123456*1234*NCS-LIC-2026-001234
```

Suspended, revoked, cancelled, superseded or explicitly non-renewable licences are refused. Forms, documents, signatures and payments remain on the authenticated website or assisted channel.

## Configuration

```env
USSD_CALLBACK_SECRET=<long-random-secret>
```

The gateway must send this value in `X-USSD-Token`. An empty value disables callback-token checking for local development only. Production must configure the secret and use HTTPS. Provider HMAC or mTLS should replace or supplement this header when the selected provider’s signing contract is available.

## Database requirements

Apply migrations through `020_ussd_quick_services.sql`. Migration 020 adds the general licence registry and a redacted USSD access-log table. Reconcile and import authoritative NCS-issued licences before enabling public validity checks.

## HTTP behaviour

- `200 text/plain` with `CON` or `END` for valid callbacks and ordinary not-found results.
- `400` for malformed encoding.
- `401` for a missing or incorrect configured callback token.
- `Cache-Control: no-store` on callback responses.
- Database failures produce a neutral temporary-unavailability message.

## Deployment checklist

1. Apply migrations 019 and 020.
2. Reconcile and populate `licences`.
3. Normalize registered user phone numbers and ensure users have set PINs.
4. Store `USSD_CALLBACK_SECRET` in the deployment secret manager.
5. Configure `https://<api-host>/ncs-ussd` at the provider.
6. Ensure the provider sends caller MSISDN and the callback token.
7. Exercise every flow in the provider sandbox and on supported networks.
8. Configure upstream USSD-specific rate limits and monitoring.

## Security requirements

- Never log callback bodies because protected `text` contains the PIN.
- Use TLS and trusted proxy-to-backend networking.
- Rotate the callback secret.
- Limit retries and session duration at the provider/proxy.
- Never return reviewer notes, documents, DOB, national IDs or full personal profiles.
- Phone/PIN authentication authorizes read-only summaries only. Profile changes, submissions, documents, signatures and payments require normal authenticated channels.
