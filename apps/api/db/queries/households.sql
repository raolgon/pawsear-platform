-- name: CreateHousehold :one
INSERT INTO households (
	id,
	display_name,
	address_line1,
	address_line2,
	neighborhood,
	city,
	timezone,
	notes,
	created_at,
	updated_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetHousehold :one
SELECT *
FROM households
WHERE id = ?;

-- name: ListHouseholds :many
SELECT *
FROM households
WHERE active = 1
ORDER BY display_name;

-- name: UpdateHousehold :one
UPDATE households
SET
	display_name = ?,
	address_line1 = ?,
	address_line2 = ?,
	neighborhood = ?,
	city = ?,
	timezone = ?,
	notes = ?,
	active = ?,
	updated_at = ?
WHERE id = ?
RETURNING *;
