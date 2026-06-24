# Message automation API

This API is the boundary between n8n/channel integrations and Pawsear's local SQLite data. n8n transports and normalizes events; Pawsear validates them, deduplicates them, stores source context, and enforces human review.

Machine-readable contract: [`apps/api/openapi.yaml`](../../apps/api/openapi.yaml).

## Local base URL

```text
http://127.0.0.1:8080
```

From the n8n container, use:

```text
http://api:8080
```

## Authentication

Set `PAWSEAR_AUTOMATION_TOKEN` on the API before accepting events from a public webhook. Send it only to `POST /api/message-imports`:

The unified Compose stack injects this token into API, web, and n8n. Generate or rotate
the ignored value with `automation/n8n/scripts/configure_automation_token.sh`; workflows
reference the environment value and never contain the secret itself.

```http
Authorization: Bearer <local automation token>
```

The local review endpoints are currently intended for the Pawsear web server and do not implement multi-user authorization. Do not expose the Go API directly to the internet.

## Operator-approved Telegram replies

The web app queues one of three reviewed templates with
`POST /api/detected-requests/{id}/replies`: request more details, confirm a converted
booking, or decline and notify. Pawsear renders and persists the exact message before
delivery.

n8n polls the protected `GET /api/automation/outbound-messages?status=pending` endpoint,
sends each Telegram message, then calls
`PATCH /api/automation/outbound-messages/{id}` with `{"status":"sent"}`. If delivery
does not complete, the item remains pending for a later retry. The public
`GET /api/outbound-messages` endpoint is local-only and supplies delivery status to the UI.

## Import contract and idempotency

Telegram and WhatsApp imports require:

- `channel`
- `externalConversationId`
- `externalMessageId`
- `body`

The tuple `(channel, externalConversationId, externalMessageId)` identifies one external event. The first request returns `201 Created`; retries return `200 OK` with `duplicate: true` and the original Pawsear IDs.

```sh
curl -X POST http://127.0.0.1:8080/api/message-imports \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer YOUR_LOCAL_TOKEN' \
  -d '{
    "channel": "telegram",
    "externalConversationId": "9988123",
    "externalMessageId": "184",
    "senderExternalId": "5511223",
    "direction": "inbound",
    "body": "¿Puedes pasear a Luna mañana a las 8?",
    "sentAt": "2026-06-21T14:00:00Z"
  }'
```

If `senderExternalId` matches a contact's Telegram or WhatsApp ID, Pawsear suggests that contact and its primary household. A trusted parser may include a `suggestion`, but all suggested values remain reviewable.

Sender and household resolution are separate. Pawsear resolves an external identity to
a contact, then chooses a household using this order: a household remembered for the
conversation, the contact's only active household, or a unique pet-name match across
the contact's households. If several households remain possible, the request keeps the
known contact but leaves `householdId` empty. The operator selects one through
`POST /api/detected-requests/{id}/household-link`; Pawsear validates membership and
remembers that household for later messages in the same chat.

This supports multiple Telegram users in one household, different users in different
households, and group chats where `externalConversationId` identifies the chat while
`senderExternalId` identifies each participant.

When no parser suggestion is supplied, Pawsear applies a small local Spanish-language
parser. It recognizes common service words, `hoy`, `mañana`, `pasado mañana`, ISO or
`DD/MM` dates, times such as `a las 10:30`, and durations in minutes or hours. It does
not create bookings automatically. The booking review form preselects these values and
matches pet names only against active pets in the recognized household.

Ambiguous messages deliberately remain incomplete. For example, “un paseo cuando
puedas” suggests a walk but does not invent a date or time.

## Review lifecycle

```text
needs_review ──► confirmed ──► converted_to_booking
      │               │
      ├──► needs_more_info
      └──► ignored ──► needs_review
```

`converted_to_booking` is terminal. Conversion creates the booking, pet links, source link, and request update in one SQLite transaction.

List or filter the queue:

```sh
curl 'http://127.0.0.1:8080/api/detected-requests?status=needs_review'
```

The Pawsear UI groups this response into pending and history views. Converted requests
include the resulting booking start and household IDs so operators can jump directly
to the correct agenda date. Ignored requests remain visible in history and can be
reopened to `needs_review`.

Record a review decision:

```sh
curl -X PATCH http://127.0.0.1:8080/api/detected-requests/REQUEST_ID \
  -H 'Content-Type: application/json' \
  -d '{"status":"needs_more_info","reviewNotes":"Confirm the pickup date"}'
```

Remember an unknown sender after choosing the correct existing contact:

```sh
curl -X POST http://127.0.0.1:8080/api/detected-requests/REQUEST_ID/contact-link \
  -H 'Content-Type: application/json' \
  -d '{"contactId":"CONTACT_ID"}'
```

This stores the Telegram or WhatsApp sender ID on the contact and refreshes the
request, source message, and conversation context. Later messages from that sender
automatically suggest the contact and its primary household. Conflicting links are
rejected instead of silently moving an identity between contacts.

If the sender is not yet a contact, the same endpoint can create and link it atomically:

```json
{
  "displayName": "Rafael",
  "householdId": "HOUSEHOLD_ID",
  "role": "owner"
}
```

Convert after review:

```sh
curl -X POST http://127.0.0.1:8080/api/detected-requests/REQUEST_ID/bookings \
  -H 'Content-Type: application/json' \
  -d '{
    "status": "confirmed",
    "startAt": "2026-06-22T14:00:00Z",
    "petIds": ["PET_ID"],
    "reviewNotes": "Confirmed with the primary contact"
  }'
```

The conversion body can override detected household, contact, service, start, and end values. A valid household and start time must exist after detected values and overrides are combined.

## HTTP behavior

| Status | Meaning |
| --- | --- |
| `200` | Successful read/update or idempotent duplicate import |
| `201` | New import or booking conversion created |
| `400` | Invalid JSON, enum, timestamp, or required domain value |
| `401` | Missing/invalid automation token when configured |
| `404` | Record not found |
| `409` | Illegal lifecycle transition or repeated conversion |
| `500` | Local persistence failure |

Errors use:

```json
{
  "error": "invalid_request",
  "message": "body is required"
}
```
