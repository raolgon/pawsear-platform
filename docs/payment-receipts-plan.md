# Payment Receipt Prototype

## Goal

Allow an operator to issue and inspect a local payment receipt, then download the same
receipt as PDF or PNG. Telegram, WhatsApp, n8n, fiscal invoicing, and automatic delivery
are intentionally outside this prototype.

## Tasks

- [x] Add an immutable receipt record linked one-to-one with a payment.
- [x] Snapshot payer, payment, allocation, charge, and household labels at issue time.
- [x] Generate a stable human-readable receipt number.
- [x] Render one shared receipt design as PDF and PNG.
- [x] Add API endpoints to issue, inspect, and download receipt artifacts.
- [x] Add a payment detail view with receipt preview and download actions.
- [x] Add receipt actions to recent payment history.
- [x] Cover missing payer, partial allocations, and repeated issue requests in tests.
- [x] Verify API tests, frontend checks, and rendered PDF layout.

## Product Rules

- A receipt confirms money received; it is not a charge, invoice, fiscal receipt, or CFDI.
- Issuing a receipt must not change payment or allocation totals.
- Receipt content is immutable after issue, even if related records later change.
- Saving a payment succeeds independently from receipt generation.
- Repeating the issue action returns the existing receipt instead of creating duplicates.
- Sensitive care notes, addresses, and medical information never appear on the receipt.

## Deferred

- Business profile and custom logo.
- Telegram, WhatsApp, email, and n8n delivery.
- Receipt cancellation or replacement.
- Fiscal numbering, tax fields, electronic signatures, and CFDI integration.
