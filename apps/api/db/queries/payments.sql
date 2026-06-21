-- name: CreatePayment :one
INSERT INTO payments (
	id,
	payer_contact_id,
	received_at,
	amount_minor,
	currency,
	method,
	reference,
	notes,
	created_at,
	updated_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetPayment :one
SELECT *
FROM payments
WHERE id = ?;

-- name: ListPayments :many
SELECT *
FROM payments
ORDER BY received_at DESC, created_at DESC;

-- name: ListPaymentsByContact :many
SELECT *
FROM payments
WHERE payer_contact_id = ?
ORDER BY received_at DESC, created_at DESC;

-- name: CreatePaymentAllocation :one
INSERT INTO payment_allocations (
	id,
	payment_id,
	charge_id,
	amount_minor,
	created_at
)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: ListPaymentAllocations :many
SELECT *
FROM payment_allocations
WHERE payment_id = ?
ORDER BY created_at;

-- name: GetAllocatedTotalForCharge :one
SELECT CAST(COALESCE(SUM(amount_minor), 0) AS INTEGER)
FROM payment_allocations
WHERE charge_id = ?;
