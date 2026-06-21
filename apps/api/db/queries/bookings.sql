-- name: CreateBooking :one
INSERT INTO bookings (
	id,
	household_id,
	service_type,
	status,
	start_at,
	end_at,
	location_type,
	address_snapshot,
	requested_by_contact_id,
	assigned_staff_id,
	source,
	notes,
	created_at,
	updated_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetBooking :one
SELECT *
FROM bookings
WHERE id = ?;

-- name: ListBookings :many
SELECT *
FROM bookings
ORDER BY start_at, created_at;

-- name: ListBookingsByRange :many
SELECT *
FROM bookings
WHERE start_at >= ? AND start_at < ?
ORDER BY start_at, created_at;

-- name: ListBookingsByHousehold :many
SELECT *
FROM bookings
WHERE household_id = ?
ORDER BY start_at DESC, created_at DESC;

-- name: UpdateBooking :one
UPDATE bookings
SET
	household_id = ?,
	service_type = ?,
	status = ?,
	start_at = ?,
	end_at = ?,
	location_type = ?,
	address_snapshot = ?,
	requested_by_contact_id = ?,
	assigned_staff_id = ?,
	source = ?,
	notes = ?,
	completed_at = ?,
	cancelled_at = ?,
	updated_at = ?
WHERE id = ?
RETURNING *;

-- name: AddBookingPet :exec
INSERT INTO booking_pets (
	booking_id,
	pet_id,
	notes
)
VALUES (?, ?, ?)
ON CONFLICT (booking_id, pet_id) DO UPDATE SET
	notes = excluded.notes;

-- name: DeleteBookingPets :exec
DELETE FROM booking_pets
WHERE booking_id = ?;

-- name: ListBookingPets :many
SELECT
	bp.booking_id,
	bp.pet_id,
	bp.notes,
	p.name,
	p.species
FROM booking_pets AS bp
JOIN pets AS p ON p.id = bp.pet_id
WHERE bp.booking_id = ?
ORDER BY p.name;
