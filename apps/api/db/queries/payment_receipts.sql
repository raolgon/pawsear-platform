-- name: CreatePaymentReceipt :one
INSERT INTO payment_receipts (
	id,
	payment_id,
	receipt_number,
	snapshot_json,
	issued_at,
	created_at
)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetPaymentReceipt :one
SELECT *
FROM payment_receipts
WHERE id = ?;

-- name: GetPaymentReceiptByPayment :one
SELECT *
FROM payment_receipts
WHERE payment_id = ?;

-- name: GetPaymentReceiptSource :one
SELECT
	p.id AS payment_id,
	p.payer_contact_id,
	c.display_name AS payer_name,
	p.received_at,
	p.amount_minor,
	p.currency,
	p.method,
	p.reference,
	p.notes
FROM payments p
LEFT JOIN contacts c ON c.id = p.payer_contact_id
WHERE p.id = ?;

-- name: ListPaymentReceiptAllocationSources :many
SELECT
	pa.charge_id,
	pa.amount_minor,
	c.description,
	h.display_name AS household_name
FROM payment_allocations pa
JOIN charges c ON c.id = pa.charge_id
JOIN households h ON h.id = c.household_id
WHERE pa.payment_id = ?
ORDER BY pa.created_at, pa.id;
