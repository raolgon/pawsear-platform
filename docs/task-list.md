# Pawsear Task List

## Product Direction

Pawsear is a local-first operations app for pet walking and pet care businesses that currently run through WhatsApp or Telegram. The first version is not a customer portal. It is an internal tool that helps turn informal chats into confirmed appointments, care tasks, reminders, work records, charges, and payment tracking.

Task status describes usable product behavior, not only database or API capability. An unchecked item
may include a note when part of its foundation already exists.

## Phase 0: Product Definition

- [x] Define the first target user: internal use for the current business.
- [x] Define the first service types: dog walk, pet sitting at client home, home-style boarding, visit/check-in.
- [x] Define what "done" means for a service: scheduled, confirmed, assigned, completed, charged, paid.
- [x] Define what data must be local-first and exportable.
- [x] Define what is out of scope for the first version: customer app, online payments, marketplace, automated public booking.

## Phase 1: Core Data Model

- [x] Design `Household` as the main grouping unit instead of owner.
- [x] Design `Contact` with roles such as owner, partner, domestic worker, family member, payer, emergency contact.
- [x] Design many-to-many links between contacts and households.
- [x] Design `Pet` linked to a household.
- [x] Design pet care information: diet, food schedule, medicine, allergies, behavior notes, vet notes.
- [x] Design `Booking` for scheduled services.
- [x] Design `CareTask` for reminders such as food, medicine, walks, cleaning, photo updates.
- [x] Design `Charge` for work that should be paid.
- [x] Design `Payment` so one payment can cover multiple charges, pets, or households.
- [x] Design `Conversation` and `Message` records for imported WhatsApp/Telegram context.

## Phase 2: Internal Operations MVP

- [x] Build dashboard for today's work.
- [x] Show scheduled walks, visits, boarding stays, and care tasks for the selected day.
- [x] Add manual creation of households.
- [x] Add manual creation of contacts.
- [x] Add manual creation of pets.
- [x] Add manual creation of bookings.
- [ ] Add recurring care tasks for pets. The schema exists; API and web flows are pending.
- [x] Add validated status changes for bookings: requested, confirmed, in progress, completed, cancelled.
- [x] Add staff assignment during booking creation.
- [x] Add notes during booking and pet creation.

## Phase 3: Payment Tracking Without Processing

- [x] Create manual charges linked to completed services or directly to a household. Automatic generation remains future work.
- [x] Register payments by cash, transfer, external card, or another manual method.
- [x] Allow payments from any contact, not only a pet owner.
- [x] Allow one payment to cover multiple charges.
- [x] Allow partial payments and reject allocations above the outstanding balance.
- [x] Show unpaid, partially paid, and paid charges in the read-only payment view.
- [x] Show accurate outstanding balance by household after payment allocations.
- [ ] Show payment history by contact. Filtering exists in the API; a contact-facing view is pending.

## Phase 4: Message Capture And Assisted Scheduling

- [ ] Build a manual message paste/import screen.
- [ ] Parse pasted WhatsApp/Telegram text into possible service requests.
- [ ] Detect possible date, time, pet names, service type, and contact.
- [ ] Create a "detected request" review screen before making real bookings.
- [ ] Allow confirming, editing, ignoring, or asking for more information.
- [ ] Store original message text linked to the detected request.
- [ ] Add confidence states for ambiguous detections.

## Phase 5: Calendar And Reminders

- [ ] Build calendar view by day, week, and month. An experimental date navigator exists, but week and month are not full operational views yet.
- [ ] Add reminders for food, medicine, pickups, drop-offs, and walks.
- [ ] Add overdue task state.
- [x] Add completed and skipped task actions, with a required reason when skipped.
- [ ] Add staff-facing daily checklist.
- [ ] Explore calendar export or sync options.

## Phase 6: Automation Integrations

- [ ] Prototype Telegram bot ingestion.
- [ ] Prototype n8n workflow for message capture.
- [ ] Investigate WhatsApp Business API requirements and costs.
- [ ] Decide whether WhatsApp automation should start with manual import, n8n, or official API.
- [ ] Add automation queue for unreviewed detected requests.
- [ ] Add audit log for automated changes.

## Phase 7: Reports And Export

- [ ] Export households, contacts, pets, bookings, charges, and payments as JSON.
- [ ] Export operational reports as CSV.
- [ ] Export full local database as SQLite.
- [ ] Export backup bundle as ZIP.
- [ ] Build report for completed work by date range.
- [ ] Build report for unpaid work.
- [ ] Build report for staff workload.

## Phase 8: Future Expansion

- [ ] Keep architecture ready for multiple businesses.
- [x] Keep staff roles flexible for walkers and sitters.
- [x] Keep service area/location fields available for future matching.
- [x] Avoid marketplace assumptions in the MVP.
- [ ] Research marketplace requirements only after internal workflow works well.

## Key Open Questions

- [ ] What city/timezone should be the default?
- [ ] How many people currently perform walks or care work?
- [ ] What are the current service prices and pricing rules?
- [ ] Are prices per pet, per household, per visit, per day, or custom?
- [ ] Do boarding stays need check-in/check-out times or only dates?
- [ ] What is the most common WhatsApp message format today?
- [ ] Should the first version run on one laptop, a local network, or a small private server?
- [ ] Should staff access the app from phone, tablet, or desktop?

## Immediate Next Tasks

- [x] Confirm MVP scope.
- [x] Draft database schema.
- [x] Create first wireframe for daily operations dashboard.
- [x] Create first wireframe for household/contact/pet profile.
- [x] Create first wireframe for booking creation.
- [x] Create first wireframe for payment tracking.
- [x] Decide whether to update README with the refined product direction.
