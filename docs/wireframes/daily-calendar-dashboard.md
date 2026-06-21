# Daily Calendar Dashboard Wireframe

## Purpose

This screen is the first operational dashboard for Pawsear. Its job is not to show everything at once. Its job is to help the operator answer one question quickly:

What needs attention for the selected day?

The dashboard should feel calm, mobile-first, and focused. The previous direction had too many visible elements on mobile, so this version reduces the default screen to:

- One selected day.
- One compact daily summary.
- The next important item.
- A short task list grouped by priority.
- Calendar mode controls hidden behind simple day/week/month switching.

Desktop may show more context, but mobile remains the source of truth.

## Design Direction

Core idea:

The default dashboard is a daily command center, not a full calendar grid.

Principles:

- Show one day at a time by default.
- Keep mobile screens short and scannable.
- Avoid showing neighboring days in the default mobile state.
- Make date navigation simple with previous/next controls.
- Offer Day, Week, and Month views as explicit modes.
- Keep cards compact until the user opens details.
- Put only the most urgent context above the fold.
- Use full detail sheets instead of dense cards when more information is needed.

## Default Route

Suggested route:

```text
/app/calendar
```

Default mode:

```text
Day
```

Default date:

```text
Today
```

## Mobile Default: Day Focus

The mobile day view should not show the week strip by default. Date movement happens through compact controls:

- Previous day.
- Selected date button.
- Next day.
- View mode menu: Day, Week, Month.

The selected date button opens a date picker or bottom sheet.

```text
┌────────────────────────────┐
│ Pawsear              [+]    │
│ Today                        │
├────────────────────────────┤
│ ‹  Wed Jun 3, 2026      ›   │
│ [Day ▾]                     │
├────────────────────────────┤
│ 1 Overdue · 2 Review · 5 Due│
├────────────────────────────┤
│ Next                        │
│ 08:00 Medicine              │
│ Mía · Casa Torres           │
│ [Complete]                  │
├────────────────────────────┤
│ Needs attention        3    │
│ ┌────────────────────────┐ │
│ │ Mía                    │ │
│ │ Medicine · 08:00       │ │
│ │ Overdue · Casa Torres  │ │
│ │ [Complete] [Details]   │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ Bruno                  │ │
│ │ WhatsApp request       │ │
│ │ Needs review           │ │
│ │ [Review]               │ │
│ └────────────────────────┘ │
├────────────────────────────┤
│ Later today            4    │
│ ┌────────────────────────┐ │
│ │ Luna + Max             │ │
│ │ Walk · 09:00           │ │
│ │ Confirmed · Rafa       │ │
│ │ [Start] [Details]      │ │
│ └────────────────────────┘ │
│ [Show 3 more]              │
├────────────────────────────┤
│ Staying with us        1    │
│ ┌────────────────────────┐ │
│ │ Coco                   │ │
│ │ Boarding · Jun 3-6     │ │
│ │ Next food 19:00        │ │
│ │ [Open stay]            │ │
│ └────────────────────────┘ │
└────────────────────────────┘
```

### Why This Is Simpler

This design removes these from the default mobile screen:

- Week strip.
- Multiple stat cards.
- Desktop-style side panels.
- Repeated filters.
- Full timeline.
- Payment summary blocks.

Those things still exist, but behind actions, detail sheets, or secondary views.

## Mobile Header

Header should stay compact:

```text
Pawsear                         +
Today
```

Required controls:

- `+` create button.
- Optional search icon only after the dashboard gets real data.
- No visible desktop-style filter button in the first version unless needed.

The `+` button opens a bottom sheet:

```text
Create
- Booking
- Care task
- Payment note
- Manual message review
```

## Date Navigation

Date navigation should be its own small bar:

```text
‹  Wed Jun 3, 2026  ›
[Day ▾]
```

Behavior:

- Tap `‹` moves to previous day.
- Tap `›` moves to next day.
- Tap selected date opens date picker.
- Tap `Day` opens view mode menu.

View mode menu:

```text
View
- Day
- Week
- Month
```

This keeps calendar power available without making the default mobile layout noisy.

## Daily Summary

The summary should be one compact text row, not four separate boxes.

```text
1 Overdue · 2 Review · 5 Due
```

Rules:

- Show only non-zero counts when possible.
- Make overdue visually stronger.
- Tapping a count filters the visible cards.
- If everything is calm, show:

```text
All clear · 5 scheduled
```

## Next Item

The next item is the only prominent card above the task list.

```text
Next
08:00 Medicine
Mía · Casa Torres
[Complete]
```

Rules:

- If there is overdue work, overdue wins over upcoming.
- If there is an active walk/visit, active work wins over upcoming.
- If there are no pending items, show the next confirmed booking.
- If today is empty, hide this block.

## Day Sections

The day view should use three primary sections:

1. `Needs attention`
2. `Later today`
3. `Staying with us`

Optional collapsed sections:

- `Completed`
- `Payments to record`

### Needs Attention

Includes:

- Overdue care tasks.
- Imported messages needing review.
- Bookings missing key info.
- Active work that needs check-in.

Only show the first two cards by default on mobile. If there are more:

```text
[Show 4 more]
```

### Later Today

Includes confirmed bookings and routine care tasks for the selected day.

Only show the next one or two cards by default. This prevents the dashboard from becoming a long scroll before the operator understands the day.

### Staying With Us

Includes active boarding/hotel stays.

This section should stay small:

- Pet name.
- Date range.
- Next care task.
- Pickup/dropoff context if urgent.

## Compact Card Design

Cards should be shorter than the current mockup. The card shows only the minimum needed to act.

```text
┌────────────────────────────┐
│ Mía                 Overdue│
│ Medicine · 08:00           │
│ Casa Torres                │
│ [Complete] [Details]       │
└────────────────────────────┘
```

Required fields:

- Pet or pet group.
- Type and time.
- Status.
- Household.
- One primary action.
- Details action if needed.

Avoid showing long notes in the card body. Notes belong in the detail sheet.

## Detail Sheet

Tapping a card opens a detail sheet on mobile.

```text
┌────────────────────────────┐
│ Mía                        │
│ Medicine · Overdue         │
├────────────────────────────┤
│ 08:00                      │
│ Casa Torres                │
│ Assigned: Ana              │
├────────────────────────────┤
│ Instructions               │
│ 1/2 tablet after food       │
├────────────────────────────┤
│ Source / history            │
│ Created from care routine   │
├────────────────────────────┤
│ [Complete] [Skip]           │
└────────────────────────────┘
```

The dashboard stays calm because details move into the sheet.

## Week View

Week view is for planning, not the default day workflow.

Mobile week view should be a compact list of days, not a seven-column calendar grid.

```text
┌────────────────────────────┐
│ ‹  Jun 1-7, 2026       ›   │
│ [Week ▾]                   │
├────────────────────────────┤
│ Mon 1                      │
│ 3 jobs · 1 review           │
├────────────────────────────┤
│ Tue 2                      │
│ 2 jobs                      │
├────────────────────────────┤
│ Wed 3                      │
│ 5 jobs · 1 overdue          │
├────────────────────────────┤
│ Thu 4                      │
│ Coco boarding · 4 jobs      │
└────────────────────────────┘
```

Tap a day to open Day view for that date.

Desktop may show a real weekly grid, but mobile should stay list-based.

## Month View

Month view is for choosing a date and seeing workload density.

Mobile month view can show a small calendar, but it should not be the default screen.

```text
┌────────────────────────────┐
│ ‹  June 2026           ›   │
│ [Month ▾]                  │
├────────────────────────────┤
│ S  M  T  W  T  F  S        │
│    1  2  3  4  5  6        │
│       • 3  • 5  • 2        │
│ 7  8  9 10 11 12 13        │
│    • 1    • 4              │
├────────────────────────────┤
│ Selected: Wed Jun 3         │
│ 5 jobs · 1 overdue          │
│ [Open day]                  │
└────────────────────────────┘
```

Month view indicators:

- Dot for scheduled work.
- Small number for count when useful.
- Warning marker for overdue/unresolved items.
- Boarding span can be simplified in MVP.

## Desktop Layout

Desktop should expand context without changing the core model.

```text
┌──────────────────────────────────────────────────────────────┐
│ Pawsear                                       [+ Create]      │
├──────────────────────────────────────────────────────────────┤
│ ‹ Wed Jun 3, 2026 ›             [Day] [Week] [Month]          │
├───────────────────────────────┬──────────────────────────────┤
│ Main day queue                │ Context                      │
│                               │                              │
│ 1 Overdue · 2 Review · 5 Due  │ Mini month                   │
│                               │ Filters                      │
│ Next                          │ Staff                        │
│ [Next card]                   │ Type                         │
│                               │ Status                       │
│ Needs attention               │                              │
│ [2 compact cards]             │ Urgent                       │
│                               │ Mía medicine overdue         │
│ Later today                   │ Bruno message review         │
│ [timeline/list cards]         │                              │
│                               │                              │
│ Staying with us               │                              │
│ [boarding card]               │                              │
└───────────────────────────────┴──────────────────────────────┘
```

Desktop may include:

- Mini month calendar.
- Filters.
- Urgent list.
- More visible cards.

Desktop should not force mobile to inherit all that visual weight.

## Filters

Filters should be hidden by default on mobile.

Mobile filter access:

```text
[Filter]
```

Filter sheet:

- Type.
- Status.
- Staff.
- Household.
- Pet search.

Default filters:

- Selected date: today.
- Status: hide cancelled, show everything else.
- Staff: all.

## Card Types

The same compact card model should work for:

- Walk.
- Visit.
- Boarding.
- Food.
- Medicine.
- Pickup/dropoff.
- Message review.
- Payment reminder.

Examples:

```text
Luna + Max            Confirmed
Walk · 09:00
Casa García
[Start] [Details]
```

```text
Bruno                 Review
WhatsApp request
Casa Rivera?
[Review]
```

```text
Coco                  Active
Boarding · Jun 3-6
Next food 19:00
[Open stay]
```

## Status System

Use text, color, and icon when possible.

Statuses:

- Needs review.
- Requested.
- Confirmed.
- In progress.
- Pending.
- Overdue.
- Completed.
- Skipped.
- Cancelled.
- Unpaid.
- Partially paid.

Priority display:

- Overdue.
- Active/in progress.
- Needs review.
- Next upcoming.
- Later pending.
- Completed.

## Empty States

No work for selected day:

```text
No scheduled work today.
[Create booking] [Import message]
```

Everything completed:

```text
All done for today.
[Review completed] [Create task]
```

No matching filters:

```text
No items match these filters.
[Clear filters]
```

## First Implementation Target

Replace the current mockup with a simpler static version that validates this direction:

- Header.
- Date navigation row with previous/next day.
- Day/Week/Month mode switch.
- One compact daily summary row.
- One next item block.
- Three sections:
  - Needs attention.
  - Later today.
  - Staying with us.
- Compact cards with fewer visible fields.
- Detail behavior can remain non-functional in the static mockup.

Do not connect real API data until this layout feels right.
