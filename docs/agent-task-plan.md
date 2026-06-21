# Pawsear Agent Task Plan

## Objective

Stabilize the current local MVP before expanding its scope. The immediate goal is to make the existing booking-to-payment workflow verifiable, accurately documented, and independent from hard-coded demo identifiers.

## Working Rules

- Work only inside WSL2. Do not invoke Windows executables or Windows paths.
- Read `AGENTS.md` and the nearby implementation before editing.
- Preserve existing uncommitted changes and do not revert unrelated work.
- Do not stage, commit, push, pull, merge, rebase, create or switch branches, create tags, or
  open pull requests until the developer explicitly requests that specific Git action.
- Keep each change focused and run the relevant checks before moving to the next task.
- Do not add dependencies unless the existing stack cannot solve the concrete need.
- Treat contacts, addresses, care notes, and payment information as sensitive data.
- Use the WSL2 NVM environment. In non-interactive shells, run commands through `bash -lic` or load `$HOME/.nvm/nvm.sh` explicitly.

## Phase 1: Establish A Green Baseline

- [x] Run the Go test suite from the repository root:
  `GOCACHE="$(pwd)/apps/api/.gocache" go test ./apps/api/...`
- [x] Run the API smoke test:
  `bash apps/api/scripts/smoke.sh`
- [x] Run `pnpm check` from `apps/web` in the WSL2 NVM environment.
- [x] Run `pnpm lint` from `apps/web` in the WSL2 NVM environment.
- [x] Run `pnpm build` from `apps/web` in the WSL2 NVM environment.
- [x] Record every failure before changing code; fix only failures related to the current MVP.

Acceptance criteria:

- Backend tests and smoke test pass.
- Frontend type checking, linting, formatting, and production build pass.
- No generated cache, build output, secret, or temporary file is added to version control.

## Phase 2: Align Documentation With Reality

- [x] Update `docs/task-list.md` so implemented work is checked and incomplete work remains open.
- [x] Update `apps/api/README.md` to acknowledge the current server-side frontend integration.
- [x] Distinguish functional screens from visual prototypes or inactive controls.
- [x] Document the current verification commands without machine-specific assumptions beyond the repository path already established by the project.

Acceptance criteria:

- A new agent can identify what works without reverse-engineering every route.
- Documentation does not claim implemented features are missing or prototype features are complete.

## Phase 3: Remove Hard-Coded Demo Coupling

- [x] Replace dashboard maps of demo household, pet, and staff IDs with API-provided display data.
- [x] Replace payment-page ID formatting with real household and payer names.
- [x] Review household profile fallback behavior and prevent an unknown household from silently showing Casa Garcia.
- [x] Keep demo seeding available, but treat seeded records like ordinary records in the UI.
- [x] Show explicit offline and empty states instead of silently substituting operational records.

Acceptance criteria:

- Records created through the UI display meaningful names without adding their IDs to frontend source code.
- An unknown household route returns a clear not-found state.
- API outages are visible and are not mistaken for real operational data.

## Phase 4: Complete One Vertical MVP Workflow

Target workflow:

```text
Create booking
  -> show it in the calendar/dashboard
  -> start and complete the service
  -> complete or skip required care tasks
  -> create or confirm a charge
  -> record and allocate a payment
  -> refresh household balance
```

- [x] Make booking status controls persist through the API.
- [x] Add functional care-task completion and skipped-with-reason actions.
- [x] Add a charge creation or confirmation flow.
- [x] Add a payment form with partial and multi-charge allocation support.
- [x] Refresh charge status and household balance after allocation.
- [x] Add regression tests for the principal happy path and validation failures.

Acceptance criteria:

- The complete workflow can be performed from the web UI against a fresh local SQLite database.
- Reloading any page preserves the resulting state.
- Partial payments produce `partially_paid`; complete allocation produces `paid`.

## Phase 5: Editing And Care Records

- [ ] Add edit flows for households, contacts, pets, and bookings using existing API capabilities.
- [ ] Expose structured pet diet and medication records through API and UI.
- [ ] Expose care-routine templates and generate concrete care tasks deliberately.
- [ ] Confirm mobile loading, validation, error, and success states for every form.

## Phase 6: Assisted Message Capture

- [ ] Build the manual WhatsApp/Telegram paste screen.
- [ ] Store source conversation and message context.
- [ ] Detect possible household, contact, pets, service type, date, and time.
- [ ] Require operator review before creating a booking.
- [ ] Support confirm, edit, ignore, and needs-more-information outcomes.

## Phase 7: Local Export And Recovery

- [ ] Export domain data as JSON.
- [ ] Export useful operational reports as CSV.
- [ ] Export or copy the SQLite database safely.
- [ ] Produce a ZIP backup with a small manifest.
- [ ] Document and test restoration into a clean local installation.

## Deferred Until The Core Workflow Is Stable

- Authentication and multi-user permissions.
- Telegram bot or WhatsApp Business API automation.
- n8n workflows.
- Calendar synchronization.
- Multi-business SaaS support.
- Marketplace, online payments, or customer-facing applications.

## Agent Handoff Template

At the end of each work session, report:

- Tasks completed from this document.
- Files changed.
- Checks run and their exact result.
- Known risks or checks that could not run.
- Recommended next unchecked task.
- Whether the working tree contains unrelated pre-existing changes.
