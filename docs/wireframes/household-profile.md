# Household Profile Wireframe

## Purpose

The household profile is the operational home for a house, family, or care context. It should make it easy to understand:

- Which pets belong to this household.
- Who can communicate, authorize, hand off, or pay.
- What care instructions matter.
- What work is coming up.
- What is unpaid or needs attention.
- What messages or notes explain recent decisions.

This screen should not assume a single owner. A household can have multiple contacts with different roles.

## Mobile-First Principle

The mobile profile should be useful while standing at a door, preparing food, giving medicine, starting a walk, or answering a message.

Design priorities:

- Pet cards before administrative details.
- Contacts visible with roles and quick actions.
- Urgent tasks and upcoming bookings above history.
- Payment/balance visible but not louder than care needs.
- Notes and instructions easy to scan.
- Editing should happen through focused sheets, not large desktop-style forms.

## Primary Route

Suggested route:

`/app/households/:householdId`

Default mobile state:

- Opens on household overview.
- Shows pets, next/upcoming items, and contact shortcuts first.
- Uses tabs/sections to avoid an endless form.

## Mobile Layout

```text
┌────────────────────────────┐
│ < Back              Edit    │
│ Casa García                 │
│ Roma Norte · Active         │
├────────────────────────────┤
│ Alerts                      │
│ 1 overdue · 2 upcoming      │
│ Balance: $360 MXN unpaid    │
├────────────────────────────┤
│ Pets                        │
│ ┌────────────────────────┐ │
│ │ Luna                   │ │
│ │ Dog · medium           │ │
│ │ Next: walk today 09:00 │ │
│ │ Note: nervous w bikes  │ │
│ │ [Care] [Book]          │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ Max                    │ │
│ │ Dog · large            │ │
│ │ Next: walk today 09:00 │ │
│ │ Note: pulls on leash   │ │
│ │ [Care] [Book]          │ │
│ └────────────────────────┘ │
├────────────────────────────┤
│ Today / Next                │
│ ┌────────────────────────┐ │
│ │ 09:00 · Walk            │ │
│ │ Luna + Max              │ │
│ │ Rafa · Confirmed        │ │
│ │ [Start] [Details]       │ │
│ └────────────────────────┘ │
├────────────────────────────┤
│ Contacts                    │
│ Mariana · owner · WhatsApp  │
│ [Message] [Call]            │
│ Sofía · payer               │
│ [Message] [Call]            │
│ Rosa · domestic worker      │
│ [Message] [Call]            │
├────────────────────────────┤
│ Notes                       │
│ Gate code 4421. Leave keys  │
│ with front desk if needed.  │
├────────────────────────────┤
│ Tabs                        │
│ [Overview] [Care] [Money]   │
│ [History] [Messages]        │
└────────────────────────────┘
```

The first mobile screen should favor cards and summaries. Detailed lists live inside tabs or bottom sheets.

## Mobile Sections

### Header

The header should show:

- Back navigation.
- Household display name.
- Neighborhood or short location.
- Active/inactive state.
- Edit action.

Optional actions:

- Add pet.
- Add contact.
- Create booking.

These can live under a `+` action menu on small screens.

### Alerts

Alerts are operational summaries:

- Overdue tasks.
- Upcoming bookings today.
- Active boarding stays.
- Unpaid balance.
- Needs review messages.

Example:

```text
1 overdue · 2 upcoming
Balance: $360 MXN unpaid
```

Tapping an alert jumps to the relevant section.

### Pet Cards

Pet cards are the most important content in the profile.

```text
┌────────────────────────────┐
│ Luna                       │
│ Dog · medium · active      │
│ Next: walk today 09:00     │
│ Care: food 2x/day          │
│ Note: nervous with bikes   │
│ [Care] [Book] [Details]    │
└────────────────────────────┘
```

Required pet card fields:

- Pet name.
- Species and size.
- Active/inactive state.
- Next booking or care task.
- One important care/behavior note.
- Quick actions.

Multiple pets should be stacked vertically on mobile.

### Today / Next

Shows operational work for this household.

Includes:

- Next booking.
- Overdue care task.
- Active boarding stay.
- Message review related to this household.

Cards should use the same visual structure as the dashboard so the user does not relearn the UI.

### Contacts

Contacts should be role-first and action-oriented.

```text
Mariana
owner · primary · WhatsApp
[Message] [Call]

Sofía
payer
[Message] [Call]

Rosa
domestic worker · handoff contact
[Message] [Call]
```

Rules:

- Show roles attached to this household.
- One contact can appear with multiple roles.
- Do not imply every contact is an owner.
- Show payer clearly when known.
- Show emergency contact clearly when present.

### Household Notes

Short notes visible on overview:

- Access notes.
- Building/security notes.
- Handoff preferences.
- General care context.

Long notes should expand in a detail sheet.

### Tabs

Recommended tabs:

- `Overview`
- `Care`
- `Money`
- `History`
- `Messages`

On mobile, tabs can be horizontal scroll or a segmented control.

## Care Tab

The care tab shows pet-specific instructions and routines.

```text
┌────────────────────────────┐
│ Care                        │
├────────────────────────────┤
│ Luna                        │
│ Food                        │
│ 1 cup morning, 1 cup night  │
│ Medicine                    │
│ None                        │
│ Behavior                    │
│ Nervous around bikes        │
│ [Edit Care]                 │
├────────────────────────────┤
│ Max                         │
│ Food                        │
│ 1.5 cups morning            │
│ Medicine                    │
│ Joint supplement nightly    │
│ Behavior                    │
│ Pulls on leash              │
│ [Edit Care]                 │
└────────────────────────────┘
```

Care tab sections:

- Diet.
- Medicine.
- Allergies.
- Behavior.
- Vet notes.
- Active care routines.

## Money Tab

The money tab tracks manual charges and payments.

```text
┌────────────────────────────┐
│ Money                       │
│ Balance: $360 MXN unpaid    │
├────────────────────────────┤
│ Unpaid                      │
│ Walk · Luna + Max · $180    │
│ Visit · Luna · $180         │
│ [Record Payment]            │
├────────────────────────────┤
│ Recent Payments             │
│ Sofía · transfer · $500     │
│ Applied to 3 charges        │
└────────────────────────────┘
```

Rules:

- Balance is by household.
- Payment payer can be any contact.
- Payment can cover multiple charges.
- Show payment history by payer/contact.

## History Tab

History shows completed and cancelled work.

```text
Jun 2 · Walk · Luna + Max · completed · paid
Jun 1 · Food task · Max · completed
May 30 · Boarding · Coco · completed · partially paid
```

Filters:

- Service type.
- Pet.
- Status.
- Date range.
- Payment status.

## Messages Tab

Messages show imported or linked WhatsApp/Telegram context.

```text
WhatsApp · Mariana
"Puedes mañana a las 8?"
Detected: possible walk · Bruno · needs review
[Review]
```

Rules:

- Messages are context, not the source of truth.
- Detected requests must show review status.
- Converted bookings should link back to original message.

## Desktop Layout

Desktop expands the profile into columns but keeps the same content order.

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ < Back   Casa García · Roma Norte                         [Edit] [+ Booking] │
├──────────────────────────────┬──────────────────────────────┬────────────────┤
│ Pets                         │ Today / Next                 │ Contacts       │
│ [Luna Card]                  │ [Walk Card]                  │ Mariana owner  │
│ [Max Card]                   │ [Medicine Card]              │ Sofía payer    │
│                              │                              │ Rosa handoff   │
├──────────────────────────────┴──────────────────────────────┴────────────────┤
│ Tabs: Overview | Care | Money | History | Messages                           │
│                                                                              │
│ Selected tab content                                                          │
└──────────────────────────────────────────────────────────────────────────────┘
```

Desktop rules:

- Pets stay left or top-left.
- Today/Next remains prominent.
- Contacts are always visible enough for quick communication.
- Tabs can show richer tables, but mobile card logic remains the source of truth.

## Detail Sheets

Mobile editing should use focused sheets:

### Add/Edit Pet

Fields:

- Name.
- Species.
- Breed.
- Size.
- Sex.
- Birthdate.
- Behavior notes.
- Medical notes.
- Feeding notes.
- Vet notes.

### Add/Edit Contact

Fields:

- Name.
- Phone.
- WhatsApp ID/phone.
- Telegram ID.
- Email.
- Roles for this household.
- Primary contact toggle.
- Notes.

### Add/Edit Household

Fields:

- Display name.
- Address.
- Neighborhood.
- City.
- Timezone.
- Household notes.
- Active toggle.

## Empty States

No pets:

```text
No pets yet.
[Add Pet]
```

No contacts:

```text
No contacts yet.
[Add Contact]
```

No upcoming work:

```text
No upcoming work for this household.
[Create Booking]
```

No unpaid balance:

```text
No unpaid balance.
```

## First Implementation Target

The first static/mock-data version should support:

- Household header.
- Alert strip.
- Two pet cards.
- Today/Next cards.
- Contact list with roles.
- Notes preview.
- Tabs as static sections.
- Mobile layout first, desktop responsive expansion second.

This is enough to validate the household model before building booking creation and payment tracking.
