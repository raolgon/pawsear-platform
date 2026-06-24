# Pawsear

Local-first operations app for pet walking, pet sitting, and home-style pet boarding businesses.

Pawsear starts as an internal tool for businesses that currently run through WhatsApp or Telegram. Clients should not need to install a new app or change their habits. The app's job is to help turn informal conversations into reliable scheduling, care reminders, work records, charges, and payment tracking.

## Product Direction

The first version is not a customer portal, payment processor, or marketplace. It is an operations layer for the business owner and staff.

Core goals:

- Keep WhatsApp and Telegram as the main client communication channels.
- Design mobile-first because the app will likely be used most while moving between services.
- Help schedule walks, visits, pet sitting, and home-style boarding.
- Track care instructions such as food, diet, medicine, allergies, behavior notes, and routines.
- Remind staff about time-sensitive care tasks.
- Record completed services.
- Generate manual charges from completed work.
- Track payments received by cash, transfer, or other methods.
- Support messy real-world relationships between households, contacts, pets, and payers.

## Current Planning Artifacts

The product direction is being captured in versioned docs before large implementation work:

- [MVP scope](docs/mvp-scope.md)
- [Data model](docs/data-model.md)
- [Design system](docs/design-system.md)
- [Daily calendar dashboard wireframe](docs/wireframes/daily-calendar-dashboard.md)
- [Household profile wireframe](docs/wireframes/household-profile.md)
- [Booking creation wireframe](docs/wireframes/booking-creation.md)
- [Payment tracking wireframe](docs/wireframes/payment-tracking.md)
- [Task list](docs/task-list.md)
- [Agent task plan](docs/agent-task-plan.md)

## MVP Scope

The initial MVP should focus on:

- Household, contact, and pet records.
- Manual booking creation.
- Daily operations dashboard.
- Care tasks and reminders.
- Booking statuses: requested, confirmed, in progress, completed, cancelled.
- Manual charge and payment tracking.
- Message paste/import from WhatsApp or Telegram as an assisted scheduling workflow.
- Local-first data storage and export.

Out of scope for the first version:

- Customer-facing app.
- Online payment processing.
- Marketplace for finding walkers or sitters.
- Fully automated booking without human review.
- Public self-serve booking.

## Domain Model

Pawsear should not assume that every pet has exactly one owner who schedules and pays for everything.

The main grouping unit is a household:

- A household can have many pets.
- A household can have many contacts.
- A contact can be an owner, partner, family member, domestic worker, emergency contact, payer, or another relevant person.
- A payment can come from any contact.
- A payment can cover one or many charges.
- A charge can be related to one or many pets or services.

This keeps the app flexible for real cases where one person schedules care, another person hands off the pet, and someone else pays.

## Stack

- SvelteKit + TypeScript
- Tailwind + Skeleton.dev
- Go backend in `apps/api`
- SQLite with migrations
- `sqlc` for typed Go database access
- n8n automations
- WhatsApp / Telegram integrations
- Exportable data: JSON, CSV, SQLite, ZIP

## Current Implementation Status

Functional now:

- Go API, SQLite migration, typed `sqlc` queries, generic demo seed, and an HTTP smoke test.
- API endpoints for households, contacts, pets, staff, bookings, care tasks, charges, payments, and the daily dashboard.
- Server-rendered dashboard, calendar date navigation, household list/profile, and payment history reads.
- Web forms that create households, contacts, pets, and bookings.
- API payment allocation with partial and multi-charge payments and automatic charge-status refresh.
- Dashboard, calendar, and payment views resolve household, pet, staff, and payer names from API data.
- API outages and empty datasets have explicit states; the web app does not substitute demo operational records.
- Household operations for confirming, starting, completing, and cancelling bookings.
- Care-task completion or skipping with a required reason.
- Manual charge creation and payment entry with explicit per-charge allocations.
- Outstanding balances that subtract previous payment allocations.
- Telegram ingestion through n8n with authenticated, idempotent delivery into the review queue.
- Conservative local message suggestions for service, date, time, duration, and known pets.
- Pending/history request views with reviewed conversion into the agenda.

Partial or presentation-only:

- The calendar has experimental day, week, and month navigation; only the selected day's operational data is loaded.
- Editing existing household, contact, pet, and booking details is still API-only.
- Message interpretation is intentionally limited to clear Spanish phrases; ambiguous language remains for manual review.
- A pinned local n8n Compose setup, generic webhook workflow, and Telegram ingestion workflow are available in `automation/n8n`; channel credentials are not stored in the repository.

Not implemented yet:

- Official WhatsApp ingestion.
- Local export and restoration flows.
- Authentication and multi-user permissions.

## Local Backend

The backend is intentionally small for now: it starts a Go HTTP server, opens SQLite, applies local migrations, and exposes health/meta endpoints.

```sh
cd apps/api
go run ./cmd/server
```

Environment variables:

- `PAWSEAR_HTTP_ADDR`: defaults to `:8080`.
- `PAWSEAR_DB_PATH`: defaults to `../../data/pawsear-local.db` when running from `apps/api`.
- `PAWSEAR_SEED_DEMO`: defaults to `false`; set to `true` only for an explicit disposable demo.
- `PAWSEAR_AUTOMATION_TOKEN`: optional bearer token shared with trusted automation clients.

The default local database starts with an empty migrated schema. Demo data is opt-in and should use
a separate path such as `../../data/pawsear-demo.db`.

Run the web app against the local API:

```sh
cd apps/web
PAWSEAR_API_BASE_URL=http://127.0.0.1:8080 pnpm dev
```

Verify the Go backend from the repo root:

```sh
GOCACHE="$(pwd)/apps/api/.gocache" go test ./apps/api/...
```

Verify the web app from `apps/web`:

```sh
pnpm check
pnpm lint
pnpm build
```

## Unified Local Stack

The preferred WSL2 setup runs Pawsear and n8n on one private Docker network while
preserving SQLite in `data/` and n8n state in its existing Docker volume:

```sh
DOCKER_BUILDKIT=0 docker compose --env-file automation/n8n/.env -f compose.local.yaml up -d
```

Open Pawsear at `http://localhost:5173` and n8n at `http://localhost:5678`. Inside
n8n, Pawsear is available at `http://api:8080`; HTTP Request nodes must not use
`localhost`, because that would refer to the n8n container itself.

Check or stop the stack without deleting local data:

```sh
docker compose --env-file automation/n8n/.env -f compose.local.yaml ps
docker compose --env-file automation/n8n/.env -f compose.local.yaml stop
```

Do not add `-v` when stopping or removing containers. The n8n volume contains the
local account, credentials, workflows, and encryption key.

## Product Workflow

The intended MVP workflow is:

```text
WhatsApp/Telegram message
  -> detected request or manual booking
  -> reviewed booking
  -> daily operations queue
  -> care tasks and completed work
  -> charges
  -> manual payments and allocations
```

The app should keep clients in their current communication habits while giving the business a structured operational record.

## Philosophy

Local by default. Open by design. Cloud optional later.

Mobile-first by default. Desktop should expand operational visibility, not become the only comfortable way to use the product.

Light mode first because the product will likely be used mostly during the day, often while moving between services or checking tasks in bright environments. Dark mode can come later as a dedicated pass.

Pawsear should make the current business easier to run before trying to become a platform for other businesses. If the internal workflow becomes strong, the product can later grow into a broader tool for other walkers, sitters, and home-style pet care providers.
