# Payment Tracking Wireframe

## Purpose

Payment tracking helps the operator record money received and understand what is still unpaid. The MVP does not process payments online. It only tracks manual payments such as cash, bank transfer, or another external method.

Core rule:

```text
Payment -> Payment Allocations -> Charges
```

Not:

```text
Payment -> Owner
Payment -> Pet
```

A payment can come from any contact and can cover one or many charges.

## Mobile-First Principle

Payment tracking should work from a phone when the operator receives cash, sees a transfer, or checks what is pending after a service.

Design priorities:

- Fast record-payment flow.
- Clear unpaid/partial/paid states.
- Payer can be any contact.
- Payment can cover multiple charges.
- Household balance is visible.
- Allocation is explicit but not tedious.
- No online payment assumptions.

## Entry Points

### Dashboard

From:

- Payment attention card.
- Completed unpaid service.
- `+` action menu.

Defaults:

- Selected date defaults to today.
- Related charge preselected when launched from payment attention card.

### Household Profile

From:

- Money tab.
- Unpaid charge.
- Balance alert.

Defaults:

- Household preselected.
- Open charges for household shown first.

### Booking Completion

From:

- Complete booking flow.

Defaults:

- Charge generated or selected.
- Payment optional.
- Default is create charge after completion, not record payment.

## Money Dashboard

Suggested route:

`/app/money`

Mobile layout:

```text
┌────────────────────────────┐
│ Money                 [+]   │
│ Today · Wed Jun 3           │
├────────────────────────────┤
│ Summary                     │
│ Unpaid: $1,240 MXN          │
│ Partial: $360 MXN           │
│ Received today: $800 MXN    │
├────────────────────────────┤
│ Tabs                        │
│ [Unpaid] [Payments] [All]   │
├────────────────────────────┤
│ Unpaid                      │
│ ┌────────────────────────┐ │
│ │ Casa García             │ │
│ │ Luna + Max · Walk       │ │
│ │ Jun 3 · $180 MXN        │ │
│ │ [Record Payment]        │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ Casa Rivera             │ │
│ │ Bruno · Visit           │ │
│ │ Jun 2 · $220 / $360     │ │
│ │ Partially paid          │ │
│ │ [Record Payment]        │ │
│ └────────────────────────┘ │
├────────────────────────────┤
│ Recent Payments             │
│ Sofía · Transfer · $500     │
│ Applied to 3 charges        │
└────────────────────────────┘
```

Default mobile tab:

- `Unpaid`

The money dashboard is secondary to daily operations, but it should be easy to open after services are completed.

## Household Money Tab

Within household profile:

```text
┌────────────────────────────┐
│ Casa García · Money         │
│ Balance: $360 MXN unpaid    │
├────────────────────────────┤
│ Open Charges                │
│ ┌────────────────────────┐ │
│ │ Walk · Luna + Max       │ │
│ │ Jun 3 · $180 unpaid     │ │
│ │ [Record Payment]        │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ Visit · Luna            │ │
│ │ Jun 1 · $180 unpaid     │ │
│ │ [Record Payment]        │ │
│ └────────────────────────┘ │
├────────────────────────────┤
│ Recent Payments             │
│ Sofía · transfer · $500     │
│ Applied to 3 charges        │
│ Mariana · cash · $180       │
│ Applied to 1 charge         │
└────────────────────────────┘
```

Rules:

- Balance is by household.
- Payer is shown as a contact.
- Do not call payer "owner" unless their role is owner.
- Charges show pets/services for context.

## Charge Card

Charge cards represent money owed.

```text
┌────────────────────────────┐
│ $180 MXN           Unpaid   │
│ Walk · Luna + Max           │
│ Casa García                 │
│ Completed Jun 3 · Rafa      │
│ Due today                   │
│ [Record Payment] [Details]  │
└────────────────────────────┘
```

Required fields:

- Amount.
- Status.
- Service/description.
- Pet(s) when applicable.
- Household.
- Date.
- Payment action.

Partial payment variant:

```text
┌────────────────────────────┐
│ $220 / $360 MXN     Partial │
│ Visit · Bruno               │
│ Casa Rivera                 │
│ Remaining: $140 MXN         │
│ [Record Payment] [Details]  │
└────────────────────────────┘
```

## Record Payment Flow

The record payment flow is a mobile sheet.

```text
┌────────────────────────────┐
│ Record Payment         X    │
│ Step 1 of 3                 │
├────────────────────────────┤
│ Who paid?                   │
│ [Sofía · payer        v]    │
│ [Unknown payer]             │
├────────────────────────────┤
│ Amount                      │
│ [$500.00 MXN             ]  │
├────────────────────────────┤
│ Method                      │
│ [Cash] [Transfer] [Other]   │
├────────────────────────────┤
│ Received                    │
│ Today · 14:30               │
├────────────────────────────┤
│ [Next]                      │
└────────────────────────────┘
```

Steps:

1. Payment details.
2. Allocate to charges.
3. Confirm.

## Step 1: Payment Details

Fields:

- Payer contact.
- Unknown payer toggle.
- Amount.
- Currency.
- Method.
- Received date/time.
- Reference.
- Notes.

Methods:

- Cash.
- Bank transfer.
- Card external.
- Other.

Rules:

- Payer is optional.
- Amount is required.
- Method is required.
- Currency defaults to business currency.
- Transfer reference is optional but encouraged.

## Step 2: Allocate To Charges

```text
┌────────────────────────────┐
│ Allocate Payment       X    │
│ Payment: $500 MXN          │
│ Remaining: $140 MXN        │
├────────────────────────────┤
│ Suggested charges           │
│ ┌────────────────────────┐ │
│ │ [x] Walk · Luna + Max   │ │
│ │ Casa García · $180      │ │
│ │ Apply: [$180]           │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ [x] Visit · Luna        │ │
│ │ Casa García · $180      │ │
│ │ Apply: [$180]           │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ [ ] Visit · Bruno       │ │
│ │ Casa Rivera · $360      │ │
│ │ Apply: [$0]             │ │
│ └────────────────────────┘ │
├────────────────────────────┤
│ [Back]              [Next]  │
└────────────────────────────┘
```

Suggested charge order:

1. Charges from the selected household.
2. Charges linked to selected booking.
3. Oldest unpaid charges.
4. Charges for pets mentioned in context.
5. Other open charges.

Rules:

- Payment can be allocated across multiple charges.
- Allocation amount can be partial.
- Remaining amount can stay unallocated.
- Currency must match selected charges.
- Charge status updates after save.

### Cross-Household Allocation

If the payer covers another household:

```text
Other open charges
[Show charges from other households]
```

When enabled:

```text
Casa Rivera
Bruno · Visit · $360 unpaid
```

The UI should make cross-household allocation explicit so it does not happen accidentally.

## Step 3: Confirm Payment

```text
┌────────────────────────────┐
│ Confirm Payment        X    │
├────────────────────────────┤
│ Payment                     │
│ Sofía · Transfer            │
│ $500 MXN · Today 14:30      │
│ Ref: 8821                   │
├────────────────────────────┤
│ Allocations                 │
│ Walk · Luna + Max · $180    │
│ Visit · Luna · $180         │
│ Unallocated · $140          │
├────────────────────────────┤
│ Result                      │
│ 2 charges paid              │
│ 1 payment credit remains    │
├────────────────────────────┤
│ [Back]       [Save Payment] │
└────────────────────────────┘
```

Rules:

- Confirm screen must show payer, amount, method, date, and allocations.
- If unallocated money remains, label it clearly as unallocated.
- MVP can allow unallocated amount as payment note/credit.
- Later versions can formalize credit balance.

## Payment Detail

```text
┌────────────────────────────┐
│ Payment Detail              │
│ Sofía · transfer            │
│ $500 MXN                    │
│ Jun 3 · Ref 8821            │
├────────────────────────────┤
│ Allocated To                │
│ $180 · Walk · Luna + Max    │
│ $180 · Visit · Luna         │
├────────────────────────────┤
│ Unallocated                 │
│ $140 MXN                    │
├────────────────────────────┤
│ Notes                       │
│ Paid for Luna and Max       │
└────────────────────────────┘
```

Actions:

- Edit notes.
- Add allocation.
- Remove allocation.
- Void payment, later/admin-only.

MVP can defer edit/void behavior if needed.

## Charge Detail

```text
┌────────────────────────────┐
│ Charge Detail               │
│ $360 MXN · Partially paid   │
├────────────────────────────┤
│ Visit · Bruno               │
│ Casa Rivera                 │
│ Completed Jun 2             │
├────────────────────────────┤
│ Payments                    │
│ Mariana · transfer · $220   │
│ Remaining: $140             │
├────────────────────────────┤
│ [Record Payment]            │
└────────────────────────────┘
```

## Desktop Layout

Desktop expands the money workspace into columns.

```text
┌──────────────────────────────────────────────────────────────┐
│ Money                                      [+ Record Payment] │
├─────────────────────┬────────────────────────┬───────────────┤
│ Filters             │ Open Charges           │ Payment Detail│
│ Household           │ [Charge Card]          │ Sofía $500    │
│ Status              │ [Charge Card]          │ Allocations   │
│ Payer               │ [Charge Card]          │ Notes         │
│ Date range          │                        │               │
├─────────────────────┴────────────────────────┴───────────────┤
│ Recent Payments                                               │
└──────────────────────────────────────────────────────────────┘
```

Desktop rules:

- Same mobile flow for record payment.
- More filters visible.
- Selected charge/payment detail can appear in right panel.
- Do not require desktop for allocation.

## Edge Cases

### Someone Pays For Another Person's Pet

Expected behavior:

- Select payer contact.
- Select charges regardless of owner role.
- UI shows household/pet for each charge.
- Confirmation makes allocation explicit.

Example:

```text
Payer: Sofía
Allocated:
- Luna · Casa García · $180
- Bruno · Casa Rivera · $360
```

### One Payment Covers Multiple Services

Expected behavior:

- One payment record.
- Multiple payment allocations.
- Each charge updates status.

### Partial Payment

Expected behavior:

- Allocation below charge total.
- Charge status becomes partially paid.
- Remaining amount visible.

### Cash Given To Walker

Expected behavior:

- Method: cash.
- Payer contact optional.
- Notes can include "received by Rafa".
- Later version may add received_by_staff_id.

### Transfer Reference Missing

Expected behavior:

- Allow saving.
- Show optional warning.

## First Implementation Target

The first static/mock-data version should support:

- Money dashboard with unpaid and recent payments.
- Household money tab pattern.
- Record payment sheet.
- Allocation step with multiple charge cards.
- Confirm payment summary.
- Edge example where payer covers charges across households.

This validates the accounting model before building payment service logic.
