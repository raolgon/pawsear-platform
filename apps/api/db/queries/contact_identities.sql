-- name: CreateContactChannelIdentity :one
INSERT INTO contact_channel_identities (
	id, contact_id, channel, external_user_id, created_at, updated_at
)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT (contact_id, channel, external_user_id) DO UPDATE SET
	updated_at = excluded.updated_at
RETURNING *;

-- name: GetContactByChannelIdentity :one
SELECT c.*
FROM contacts AS c
JOIN contact_channel_identities AS identity ON identity.contact_id = c.id
WHERE identity.channel = ? AND identity.external_user_id = ? AND c.active = 1
LIMIT 1;

-- name: ListHouseholdsForContactResolution :many
SELECT DISTINCT h.*
FROM households AS h
JOIN household_contacts AS hc ON hc.household_id = h.id
WHERE hc.contact_id = ? AND h.active = 1
ORDER BY hc.is_primary DESC, hc.created_at, h.display_name;

-- name: ContactBelongsToHousehold :one
SELECT EXISTS (
	SELECT 1
	FROM household_contacts
	WHERE contact_id = ? AND household_id = ?
) AS belongs;

-- name: ListContactHouseholdPets :many
SELECT p.household_id, p.id AS pet_id, p.name AS pet_name
FROM pets AS p
JOIN household_contacts AS hc ON hc.household_id = p.household_id
JOIN households AS h ON h.id = p.household_id
WHERE hc.contact_id = ? AND p.active = 1 AND h.active = 1
ORDER BY p.household_id, p.name;

-- name: ListContactHouseholdResolutionOptions :many
SELECT DISTINCT hc.contact_id, h.id AS household_id, h.display_name AS household_name
FROM household_contacts AS hc
JOIN households AS h ON h.id = hc.household_id
WHERE h.active = 1
ORDER BY h.display_name;

-- name: SetConversationHouseholdByMessage :exec
UPDATE conversations
SET household_id = ?, updated_at = ?
WHERE id = (SELECT conversation_id FROM messages WHERE messages.id = ?);

-- name: SetDetectedRequestHousehold :exec
UPDATE detected_requests
SET household_id = ?, updated_at = ?
WHERE id = ?;
