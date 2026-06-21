-- name: CreateCharge :one
INSERT INTO charges (
	id,
	household_id,
	booking_id,
	description,
	amount_minor,
	currency,
	status,
	due_date,
	created_at,
	updated_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetCharge :one
SELECT *
FROM charges
WHERE id = ?;

-- name: ListCharges :many
SELECT *
FROM charges
ORDER BY due_date IS NULL, due_date, created_at;

-- name: ListChargesByHousehold :many
SELECT *
FROM charges
WHERE household_id = ?
ORDER BY due_date IS NULL, due_date, created_at;

-- name: ListOpenCharges :many
SELECT *
FROM charges
WHERE status IN ('unpaid', 'partially_paid')
ORDER BY due_date IS NULL, due_date, created_at;

-- name: UpdateCharge :one
UPDATE charges
SET
	household_id = ?,
	booking_id = ?,
	description = ?,
	amount_minor = ?,
	currency = ?,
	status = ?,
	due_date = ?,
	updated_at = ?
WHERE id = ?
RETURNING *;

-- name: RefreshChargeStatus :one
UPDATE charges
SET
	status = CASE
		WHEN COALESCE((SELECT SUM(amount_minor) FROM payment_allocations WHERE charge_id = charges.id), 0) = 0 THEN 'unpaid'
		WHEN COALESCE((SELECT SUM(amount_minor) FROM payment_allocations WHERE charge_id = charges.id), 0) < amount_minor THEN 'partially_paid'
		ELSE 'paid'
	END,
	updated_at = ?
WHERE charges.id = ?
RETURNING *;
