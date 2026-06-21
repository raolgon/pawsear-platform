# Agent Rules for Pawsear

These instructions apply to any agent or assistant working in this repository. The goal is to move the product forward with simple, verifiable changes that fit Pawsear's local-first direction.

## Project Context

- Pawsear is a local-first operations app for pet walking, pet sitting, and home-style pet boarding businesses.
- End customers should not need to install a new app. WhatsApp and Telegram remain the main communication channels.
- The first version is an internal tool for owners and staff, not a marketplace, payment processor, or public client portal.
- The MVP workflow is: message or manual request, reviewed booking, daily schedule, care tasks, completed work, charges, and manual payments.
- The model must support messy real-world relationships: households, contacts, pets, payers, services, charges, and payments are not always one-to-one.

## Technical Stack

- Monorepo using `pnpm`.
- Web app in `apps/web`.
- Frontend with SvelteKit, TypeScript, Tailwind, and Skeleton.dev.
- API in `apps/api`.
- Backend in Go 1.22.
- SQLite database with migrations.
- `sqlc` for typed Go database access.
- Planned exports: JSON, CSV, SQLite, and ZIP.

## Useful Commands

- API tests from the repo root:
  ```sh
  GOCACHE="$(pwd)/apps/api/.gocache" go test ./apps/api/...
  ```
- Local API:
  ```sh
  cd apps/api
  go run ./cmd/server
  ```
- Web check:
  ```sh
  cd apps/web
  pnpm check
  ```
- Web lint:
  ```sh
  cd apps/web
  pnpm lint
  ```
- Web build:
  ```sh
  cd apps/web
  pnpm build
  ```

## Collaboration Rules

- Do not perform any Git state-changing action until the developer explicitly asks for that
  specific action. This includes staging files, creating commits, pushing, pulling, merging,
  rebasing, creating or switching branches, creating tags, opening pull requests, or modifying
  Git history. Read-only commands such as `git status`, `git diff`, and `git log` are allowed.
- Do not create commits without explicit permission from the developer.
- Do not rewrite Git history, rebase, reset, force push, or delete branches without explicit permission.
- Before editing, read the nearby context and follow existing patterns.
- Keep changes focused on the requested task.
- Do not mix large refactors with functional changes unless the refactor is required to complete the task safely.
- If there are existing changes in the working tree that you did not make, do not revert them. Work around them or ask if they block the task.
- Document important decisions when they affect architecture, the data model, or product flows.

## Security and Secrets

- Never commit secrets, tokens, credentials, private keys, cookies, production dumps, or real `.env` files.
- Use example files such as `.env.example` to document environment variables without sensitive values.
- Do not print full secrets in logs, errors, screenshots, or fixtures.
- Validate inputs at trust boundaries: HTTP handlers, forms, imports, parsers, and automations.
- Treat customer data, contacts, pets, addresses, medical notes, and payments as sensitive information.
- Prefer local and exportable storage. Do not introduce cloud dependencies without discussing them first.

## Implementation Principles

- Prefer clear, direct code over premature abstraction.
- Avoid generic layers when there is only one real use case.
- Do not create internal frameworks, factories, registries, or global helpers without a concrete need.
- Every function should have a clear responsibility.
- No function should exceed 150 lines. If a function approaches that limit, split it by responsibility.
- If a component, page, handler, or service becomes difficult to read, split it into modules or create a new page/route when that better fits the product flow.
- Prefer explicit domain names over generic names such as `manager`, `processor`, `helper`, or `utils`.
- Keep domain logic out of the UI when it is shared or critical.
- Keep HTTP handlers small: parse the request, validate it, call a service/repository, and return a response.
- In Go, put domain rules in services, not handlers.
- In Svelte, keep components focused and extract subcomponents when a screen accumulates distinct states or responsibilities.

## Frontend

- Design mobile-first. Desktop should expand operational visibility, not become the only comfortable way to use the product.
- Prioritize dense, clear, operational interfaces over landing-page-style layouts.
- Use Skeleton.dev and existing patterns before introducing new styling approaches.
- Use TypeScript strictly. Avoid `any` unless there is a clear and localized reason.
- Use small Svelte components when a view accumulates too much logic or markup.
- Forms must have clear loading, error, success, and validation states.
- Do not hide network or persistence failures. Show actionable messages.
- Keep UI copy concrete and oriented toward daily operations.
- Do not assume dark mode for the MVP. The product direction is light mode first.

## Backend

- Keep the backend small and explicit.
- Use the Go standard library when it is enough.
- Keep SQLite migrations versioned in `apps/api/db/migrations`.
- Keep `sqlc` queries in `apps/api/db/queries` and generated code in `apps/api/internal/db/queries`.
- Do not edit generated code manually. Change queries or configuration and regenerate instead.
- Use `context.Context` for operations that depend on requests, the database, or IO.
- Return clear errors without leaking sensitive details to clients.
- Validate payloads before writing to SQLite.
- Preserve local-first compatibility: the app must work without required external services for the MVP flow.

## Data Model and Product

- Do not assume a pet has exactly one owner, caregiver, or payer.
- Do not assume the person who schedules is the person who pays.
- Do not assume a payment covers only one charge or that a charge belongs to only one pet.
- Keep households as the main grouping unit.
- Preserve relevant history for bookings, tasks, charges, and payments.
- Before changing central entities, review:
  - `docs/mvp-scope.md`
  - `docs/data-model.md`
  - `docs/task-list.md`
  - wireframes in `docs/wireframes/`

## Testing and Verification

- Run the relevant checks before finishing, based on the area touched.
- For backend changes, run `GOCACHE=apps/api/.gocache go test ./apps/api/...`.
- For frontend changes, run `pnpm check`, `pnpm lint`, or `pnpm build` from `apps/web` depending on the change.
- If a check cannot be run, explain why and state the remaining risk.
- Add or update tests when changing domain logic, validations, persistence, or flows with regression risk.

## Documentation

- Keep README files and docs in sync with important changes.
- Update wireframes or scope documents when a product decision changes.
- Write documentation that is brief and useful. Avoid repeating what the code already says clearly.
- Prefer concrete examples from the Pawsear domain.

## Dependencies

- Do not add dependencies without justifying the need.
- Prefer small, maintained dependencies that fit the current stack.
- Do not introduce required external services for the MVP without approval.
- Respect the repo package manager: `pnpm`.

## Delivery Criteria

- The change solves the requested task without unnecessarily expanding the scope.
- The code is readable for the next developer.
- Happy paths and main error paths are considered.
- Relevant checks pass or are clearly reported.
- No sensitive information, obvious dead code, debug logs, or temporary files are left behind.
