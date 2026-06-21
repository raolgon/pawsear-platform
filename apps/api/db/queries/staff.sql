-- name: CreateStaffMember :one
INSERT INTO staff_members (
	id,
	display_name,
	phone,
	role,
	created_at,
	updated_at
)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetStaffMember :one
SELECT *
FROM staff_members
WHERE id = ?;

-- name: ListStaffMembers :many
SELECT *
FROM staff_members
WHERE active = 1
ORDER BY display_name;

-- name: UpdateStaffMember :one
UPDATE staff_members
SET
	display_name = ?,
	phone = ?,
	role = ?,
	active = ?,
	updated_at = ?
WHERE id = ?
RETURNING *;
