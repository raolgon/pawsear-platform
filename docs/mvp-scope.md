# Pawsear MVP Scope

## Decision

The Pawsear MVP is an internal operations app for the current pet walking, pet sitting, and home-style boarding business. It should help the business owner and staff manage work that currently arrives through WhatsApp or Telegram without forcing clients to use a new app.

The MVP is confirmed around one core outcome:

> Turn informal client conversations into reliable operations: scheduled services, care tasks, completed work, charges, and payment tracking.

## Primary User

The first user is the business operator.

Secondary users are internal staff such as walkers or sitters. Clients are not app users in the MVP.

## MVP Service Types

The MVP supports these service types:

- Dog walk.
- Pet sitting at the client's home.
- Home-style boarding at the caregiver's home.
- Visit/check-in.

Other service types can be added later, but the data model should not hard-code assumptions that make them difficult.

## MVP Includes

### Household, Contact, And Pet Management

- Create and edit households.
- Create and edit contacts.
- Link multiple contacts to a household.
- Assign contact roles such as owner, partner, domestic worker, family member, payer, and emergency contact.
- Create and edit pets.
- Link pets to households instead of a single owner.
- Store pet care notes, including diet, food schedule, medicine, allergies, behavior, and vet notes.

### Scheduling And Operations

- Create bookings manually.
- View today's work as the main operational dashboard.
- Track service status: requested, confirmed, in progress, completed, cancelled.
- Assign bookings to staff.
- Add operational notes to bookings.
- Show active boarding stays, scheduled walks, visits, and care tasks.

### Care Tasks And Reminders

- Create care tasks for food, medicine, walks, cleaning, pickups, drop-offs, and photo updates.
- Mark care tasks as completed.
- Show overdue care tasks.
- Support repeated care routines for pets.

### Manual Charges And Payments

- Generate or create charges for completed services.
- Track payment status: unpaid, partially paid, paid.
- Register payments made by cash, transfer, or other manual methods.
- Allow any contact to make a payment.
- Allow one payment to cover multiple charges.
- Allow partial payments.
- Show balances by household.

### Assisted Message Capture

- Provide a manual paste/import screen for WhatsApp or Telegram messages.
- Detect possible date, time, pet names, service type, and contact from pasted text.
- Create a reviewable detected request before creating a real booking.
- Let the operator confirm, edit, ignore, or mark the request as needing more information.
- Store the original message text for context.

### Local-First Data And Export

- Store operational data locally.
- Keep the system usable without depending on a cloud service for the core workflow.
- Export data in practical formats such as JSON, CSV, SQLite, or ZIP.

## MVP Does Not Include

- Customer-facing app.
- Public booking page.
- Online payment processing.
- Marketplace for finding walkers or sitters.
- Fully automated WhatsApp booking.
- Fully automated Telegram booking.
- Automatic changes without human review.
- Multi-business SaaS administration.
- Native mobile apps.
- Advanced staff payroll.
- Customer self-service account management.

## Core Workflow

The MVP should support this workflow end to end:

1. A client or related contact sends a WhatsApp or Telegram message.
2. The operator manually creates a booking or pastes the message into Pawsear.
3. Pawsear helps identify the household, contact, pets, service type, date, and time.
4. The operator confirms or edits the detected request.
5. The booking appears in the daily operations dashboard.
6. Staff complete service and care tasks.
7. Pawsear records the completed work.
8. A charge is created or confirmed.
9. A manual payment is recorded when money is received.
10. The app shows what is paid, partially paid, or unpaid.

## Service Completion Definition

A service is considered operationally complete when:

- It has a service type.
- It is linked to the relevant household and pet or pets.
- It has a scheduled date and time or date range.
- It has a confirmed status before work starts.
- It has an assigned staff member when staff assignment is needed.
- Required care tasks are completed or explicitly skipped with a note.
- The service is marked completed.
- A charge exists or the service is marked as not billable.

A service is financially complete when:

- All related charges are paid, or
- The remaining balance is intentionally waived or adjusted.

## Data Model Boundaries

The MVP should preserve these modeling decisions:

- Household is the main grouping unit.
- Pet belongs to a household, not directly to one fixed owner.
- Contact can be linked to one or more households.
- Contact roles are contextual.
- Payment payer can be any contact.
- Payment allocation can cover multiple charges.
- Bookings and care tasks should keep source/context notes.

## Success Criteria

The MVP is successful if it can reliably answer:

- What work is scheduled today?
- Which pets need food, medicine, walks, pickup, drop-off, or other care?
- What has already been completed?
- What still needs attention?
- Who paid, how much, and by what method?
- Which services are unpaid or partially paid?
- What instructions matter for this pet right now?
- What conversation or note led to this booking?

## Open Questions

These questions should be answered during design, but they do not block confirming the MVP scope:

- Default city and timezone.
- Number of initial staff users.
- Pricing rules by service type.
- Whether boarding needs date-only stays or check-in/check-out times.
- Most common WhatsApp and Telegram message formats.
- Whether the first install runs on one laptop, a local network, or a small private server.
- Whether staff primarily use phone, tablet, or desktop.
