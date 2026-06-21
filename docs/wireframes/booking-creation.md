# Booking Creation Wireframe

## Purpose

Booking creation turns a real-world request into scheduled work. It should be fast enough to use while replying to WhatsApp or Telegram, but structured enough to produce reliable operations.

The booking flow should support:

- Manual booking from dashboard.
- Booking from a household profile.
- Booking from a pet card.
- Booking from a reviewed WhatsApp/Telegram detected request.

The flow should not feel like one long form on mobile.

## Mobile-First Principle

Booking creation is a mobile sheet or focused wizard.

Design priorities:

- Fast service selection.
- Household and pet selection before details.
- Minimal typing.
- Smart defaults from context.
- A persistent summary before save.
- Human review when generated from a message.
- Ability to save as `requested` or `confirmed`.

## Entry Points

### Dashboard

From:

- `+` action.
- Empty state.
- Message review card.

Defaults:

- Date defaults to selected dashboard date.
- Status defaults to `requested` unless explicitly confirmed.

### Household Profile

From:

- Household `+ Booking`.
- Pet card `[Book]`.

Defaults:

- Household preselected.
- Pet preselected if launched from pet card.
- Address defaults to household address.

### Message Review

From:

- Detected request card.

Defaults:

- Source message linked.
- Suggested household/contact/pet/time/service prefilled.
- Confidence and missing fields visible.
- Save requires human confirmation.

## Mobile Flow Overview

```text
┌────────────────────────────┐
│ Create Booking         X    │
│ Step 1 of 5                 │
├────────────────────────────┤
│ What service?               │
│ [Walk]                      │
│ [Visit / Check-in]          │
│ [Pet sitting]               │
│ [Boarding]                  │
│ [Transport]                 │
│ [Other]                     │
├────────────────────────────┤
│ Summary                     │
│ No pets selected yet        │
├────────────────────────────┤
│ [Next]                      │
└────────────────────────────┘
```

Recommended steps:

1. Service.
2. Household and pets.
3. Date and time.
4. Location and instructions.
5. Confirm.

For message review, start on a prefilled review step instead of step 1.

## Step 1: Service

```text
┌────────────────────────────┐
│ Create Booking         X    │
│ Service                     │
├────────────────────────────┤
│ What service?               │
│ ┌────────────────────────┐ │
│ │ Walk                   │ │
│ │ Scheduled walk         │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ Visit / Check-in       │ │
│ │ Food, water, quick care │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ Pet sitting            │ │
│ │ Care at client home    │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ Boarding               │ │
│ │ At caregiver home      │ │
│ └────────────────────────┘ │
└────────────────────────────┘
```

Service types:

- Walk.
- Visit/check-in.
- Pet sitting at client home.
- Boarding.
- Transport.
- Other.

Behavior:

- Selecting `Boarding` asks for start and end date/time.
- Selecting `Walk` defaults to a shorter duration.
- Selecting `Visit` suggests care tasks.

## Step 2: Household And Pets

```text
┌────────────────────────────┐
│ Create Booking         X    │
│ Household & Pets            │
├────────────────────────────┤
│ Search household            │
│ [Casa García              ] │
├────────────────────────────┤
│ Pets                        │
│ [x] Luna                    │
│ [x] Max                     │
│ [ ] Nala                    │
├────────────────────────────┤
│ Contacts                    │
│ Requested by                │
│ [Mariana · owner        v]  │
├────────────────────────────┤
│ [Back]              [Next]  │
└────────────────────────────┘
```

Rules:

- Booking must have one household.
- Booking should have at least one pet for pet-specific services.
- Multi-pet selection must be easy.
- Requested-by contact is optional but useful.
- Contact role is shown to avoid assuming owner.

Empty state:

```text
No household found.
[Create Household]
```

## Step 3: Date And Time

### Walk / Visit / Pet Sitting

```text
┌────────────────────────────┐
│ Create Booking         X    │
│ Date & Time                 │
├────────────────────────────┤
│ Date                        │
│ Today · Wed Jun 3           │
│ [Change]                    │
├────────────────────────────┤
│ Start                       │
│ [09:00]                     │
│ Duration                    │
│ [45 min        v]           │
│ End                         │
│ 09:45                       │
├────────────────────────────┤
│ Status                      │
│ [Requested] [Confirmed]     │
├────────────────────────────┤
│ [Back]              [Next]  │
└────────────────────────────┘
```

### Boarding

```text
┌────────────────────────────┐
│ Create Booking         X    │
│ Stay Dates                  │
├────────────────────────────┤
│ Dropoff                     │
│ Wed Jun 3 · 10:00           │
│ [Change]                    │
├────────────────────────────┤
│ Pickup                      │
│ Fri Jun 6 · 10:00           │
│ [Change]                    │
├────────────────────────────┤
│ Nights                      │
│ 3                           │
├────────────────────────────┤
│ [Back]              [Next]  │
└────────────────────────────┘
```

Rules:

- Timezone comes from household or business default.
- Boarding should support date range and check-in/check-out times.
- Status defaults to `requested`.
- Confirmed status should be intentional.

## Step 4: Location And Instructions

```text
┌────────────────────────────┐
│ Create Booking         X    │
│ Location & Instructions     │
├────────────────────────────┤
│ Location                    │
│ (•) Household home          │
│ ( ) Caregiver home          │
│ ( ) Other                   │
├────────────────────────────┤
│ Address                     │
│ Calle demo 123              │
│ Roma Norte                  │
├────────────────────────────┤
│ Instructions                │
│ [Max pulls near bikes     ] │
│ [                         ] │
├────────────────────────────┤
│ Suggested care tasks        │
│ [x] Water refill            │
│ [ ] Photo update            │
│ [ ] Food                    │
│ [ ] Medicine                │
├────────────────────────────┤
│ [Back]              [Next]  │
└────────────────────────────┘
```

Rules:

- Address should be a snapshot on the booking.
- Household notes should be visible as context but not automatically copied into booking notes.
- Care task suggestions come from pet routines and service type.
- The operator can add/remove care tasks before saving.

## Step 5: Confirm

```text
┌────────────────────────────┐
│ Create Booking         X    │
│ Confirm                     │
├────────────────────────────┤
│ Walk                        │
│ Luna + Max                  │
│ Casa García                 │
│ Today 09:00-09:45           │
│ Household home              │
│ Rafa                        │
│ Status: Confirmed           │
├────────────────────────────┤
│ Care tasks                  │
│ Water refill                │
│ Photo update                │
├────────────────────────────┤
│ Charge                      │
│ ( ) Create charge now       │
│ (•) Create after completion │
├────────────────────────────┤
│ [Back]       [Save Booking] │
└────────────────────────────┘
```

Rules:

- Save button should be disabled until required fields are valid.
- Charge defaults to after completion.
- Creating charge now is optional for confirmed/prepaid use cases.
- Successful save returns to origin:
  - Dashboard selected day.
  - Household profile.
  - Message review queue.

## Message Review Flow

When creating from a detected message, use a review-first screen.

```text
┌────────────────────────────┐
│ Review Request         X    │
│ WhatsApp · Mariana          │
├────────────────────────────┤
│ Original message            │
│ "Puedes mañana a las 8      │
│ para Bruno?"                │
├────────────────────────────┤
│ Detected                    │
│ Service: Walk               │
│ Pet: Bruno                  │
│ Date: Tomorrow              │
│ Time: 08:00                 │
│ Household: Casa Rivera      │
│ Confidence: medium          │
├────────────────────────────┤
│ Missing / check             │
│ Staff not assigned          │
│ Status not confirmed        │
├────────────────────────────┤
│ [Edit] [Ask] [Ignore]       │
│ [Create Booking]            │
└────────────────────────────┘
```

Actions:

- `Create Booking`: continues to confirmation with detected fields.
- `Edit`: opens the normal wizard with values prefilled.
- `Ask`: marks request as needs more info.
- `Ignore`: marks request ignored.

Rules:

- Detected requests do not become confirmed bookings automatically.
- Original message remains linked through `booking_sources`.
- Missing required fields are clearly shown.

## Staff Assignment

Staff assignment can appear in Date/Time or Confirm step.

```text
Assigned staff
[Rafa                 v]
```

States:

- Unassigned.
- Assigned to active staff member.
- Staff unavailable warning, later feature.

MVP does not need advanced availability checking.

## Validation

Required fields:

- Service type.
- Household.
- At least one pet for pet services.
- Start date/time.
- Location type.
- Status.

Conditional requirements:

- End date/time for boarding.
- Address snapshot for household home or other location.
- Skip staff assignment allowed, but visible.

## Desktop Layout

Desktop can use a wider two-column modal or page.

```text
┌──────────────────────────────────────────────────────────────┐
│ Create Booking                                          X     │
├─────────────────────────────────────┬────────────────────────┤
│ Form                                │ Summary                │
│ Service                             │ Walk                   │
│ Household & Pets                    │ Luna + Max             │
│ Date & Time                         │ Today 09:00-09:45      │
│ Location                            │ Casa García            │
│ Instructions                        │ Confirmed · Rafa       │
│                                     │                        │
│ [Back] [Next]                       │ [Save Booking]         │
└─────────────────────────────────────┴────────────────────────┘
```

Desktop rules:

- Summary stays visible on the right.
- The same steps exist as mobile.
- Desktop should not introduce extra required fields.

## Empty And Edge States

No pets in household:

```text
This household has no pets yet.
[Add Pet]
```

Unknown pet from message:

```text
Detected pet not found.
[Match Existing Pet] [Create Pet] [Leave Unknown]
```

Unknown household from message:

```text
No matching household.
[Search] [Create Household]
```

Booking conflict:

```text
Rafa already has a booking at this time.
[Assign anyway] [Choose another staff]
```

MVP can show conflict warnings later. The first version can skip advanced conflict detection.

## First Implementation Target

The first static/mock-data version should support:

- Open create booking sheet.
- Select service.
- Pick preloaded household and pets.
- Pick date/time.
- Pick location.
- Add simple notes.
- Confirm summary.
- Show message-review variant with prefilled data.

This validates the booking flow before implementing full scheduling rules.
