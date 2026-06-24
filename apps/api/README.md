# Pawsear API

Go backend for Pawsear's local-first operations app.

## Local Development

```sh
go run ./cmd/server
```

Environment variables:

- `PAWSEAR_HTTP_ADDR`: HTTP listen address. Defaults to `:8080`.
- `PAWSEAR_DB_PATH`: SQLite database path. Defaults to `../../data/pawsear-local.db` when running from `apps/api`.
- `PAWSEAR_SEED_DEMO`: Defaults to `false`. Set it to `true` only when a disposable demo database should be populated.
- `PAWSEAR_AUTOMATION_TOKEN`: Optional bearer token protecting automated message imports. Set it before exposing an n8n webhook outside the local machine.

The default database is intended for real local data and starts empty except for its schema. Existing
`data/pawsear.db` files from the demo setup are not opened or modified by default.

Explicit demo mode:

```sh
PAWSEAR_DB_PATH=../../data/pawsear-demo.db PAWSEAR_SEED_DEMO=true go run ./cmd/server
```

Initial endpoints:

- `GET /health`
- `GET /api/meta`
- `GET /api/households`
- `POST /api/households`
- `GET /api/households/{id}`
- `PATCH /api/households/{id}`
- `DELETE /api/households/{id}` (requires the exact household name and deletes its dependent operational data)
- `GET /api/households/{id}/contacts`
- `POST /api/households/{id}/contacts`
- `POST /api/detected-requests/{id}/replies`
- `POST /api/detected-requests/{id}/household-link`
- `GET /api/outbound-messages`
- `GET /api/automation/outbound-messages` (automation token required)
- `PATCH /api/automation/outbound-messages/{id}` (automation token required)
- `GET /api/contacts`
- `GET /api/contact-household-links`
- `POST /api/contacts`
- `GET /api/contacts/{id}`
- `PATCH /api/contacts/{id}`
- `GET /api/pets`
- `POST /api/pets`
- `GET /api/pets/{id}`
- `PATCH /api/pets/{id}`
- `GET /api/staff`
- `POST /api/staff`
- `GET /api/staff/{id}`
- `PATCH /api/staff/{id}`
- `GET /api/bookings`
- `POST /api/bookings`
- `GET /api/bookings/{id}`
- `PATCH /api/bookings/{id}`
- `GET /api/care-tasks`
- `POST /api/care-tasks`
- `GET /api/care-tasks/{id}`
- `PATCH /api/care-tasks/{id}`
- `GET /api/charges`
- `POST /api/charges`
- `GET /api/charges/{id}`
- `PATCH /api/charges/{id}`
- `GET /api/payments`
- `POST /api/payments`
- `GET /api/payments/{id}`
- `POST /api/payments/{id}/receipt`
- `GET /api/payments/{id}/receipt`
- `GET /api/payments/{id}/receipt/png`
- `GET /api/payments/{id}/receipt/pdf`
- `GET /api/dashboard/today`
- `POST /api/message-imports`
- `GET /api/detected-requests`
- `GET /api/detected-requests/{id}`
- `PATCH /api/detected-requests/{id}`
- `POST /api/detected-requests/{id}/bookings`

MVP backend coverage:

- Manual household, contact, pet, staff, booking, care task, charge, and payment records.
- Household-contact links with contextual roles.
- Bookings can include multiple pets.
- Payments can allocate money to one or more charges.
- Payments can issue one immutable internal receipt and download matching PDF and PNG artifacts.
- Charge status refreshes to `unpaid`, `partially_paid`, or `paid` after payment allocation.
- Charge responses include allocated and outstanding amounts, and over-allocation is rejected.
- Booking and care-task terminal state transitions are validated.
- Daily dashboard returns bookings, care tasks, and open charges for the selected date.
- Message imports preserve source text, match known channel contacts, avoid duplicate external events, and create a human-review queue.
- Reviewed requests convert transactionally into bookings with `booking_sources` traceability.

OpenAPI and integration guide:

- `openapi.yaml`
- `../../docs/api/message-automation.md`

Not included yet:

- Automatic natural-language message parsing.
- Export endpoints.
- Authentication or multi-user permission checks.

Current frontend integration:

- SvelteKit server loads read the dashboard, calendar, households, household profiles, charges,
  and payments from the API.
- Web forms create households, contacts, pets, and bookings through the API.
- Web actions cover booking status changes, care-task completion/skipping, charge creation, and
  partial or multi-charge payment allocation. Editing existing records is not complete yet.
- The request-review page supports manual message capture and review decisions.

## Database

The API uses SQLite and runs pending migrations on startup.

```sh
go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate
```

From this repository, the test command is:

```sh
GOCACHE="$(pwd)/apps/api/.gocache" go test ./apps/api/...
```

Real HTTP smoke test:

```sh
bash apps/api/scripts/smoke.sh
```

The smoke test starts the API on `127.0.0.1:18080` with a temporary SQLite database, creates a full MVP workflow, checks the dashboard, and verifies a few validation failures.

## Planned Backend Shape

- SQLite migrations live in `db/migrations`.
- `sqlc` queries live in `db/queries`.
- Generated `sqlc` code should live in `internal/db/queries`.
- Domain rules should live in Go services instead of HTTP handlers.
