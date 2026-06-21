-- name: CreateContact :one
INSERT INTO contacts (
	id,
	display_name,
	phone,
	whatsapp_id,
	telegram_id,
	email,
	notes,
	created_at,
	updated_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetContact :one
SELECT *
FROM contacts
WHERE id = ?;

-- name: ListContacts :many
SELECT *
FROM contacts
WHERE active = 1
ORDER BY display_name;

-- name: UpdateContact :one
UPDATE contacts
SET
	display_name = ?,
	phone = ?,
	whatsapp_id = ?,
	telegram_id = ?,
	email = ?,
	notes = ?,
	active = ?,
	updated_at = ?
WHERE id = ?
RETURNING *;

-- name: LinkHouseholdContact :exec
INSERT INTO household_contacts (
	household_id,
	contact_id,
	role,
	is_primary,
	notes,
	created_at
)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT (household_id, contact_id, role) DO UPDATE SET
	is_primary = excluded.is_primary,
	notes = excluded.notes;

-- name: ListHouseholdContacts :many
SELECT
	hc.household_id,
	hc.contact_id,
	hc.role,
	hc.is_primary,
	hc.notes,
	hc.created_at,
	c.display_name,
	c.phone,
	c.whatsapp_id,
	c.telegram_id,
	c.email
FROM household_contacts AS hc
JOIN contacts AS c ON c.id = hc.contact_id
WHERE hc.household_id = ?
ORDER BY hc.is_primary DESC, c.display_name;
