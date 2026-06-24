-- name: CreateConversation :one
INSERT INTO conversations (
	id,
	channel,
	external_conversation_id,
	primary_contact_id,
	household_id,
	last_message_at,
	created_at,
	updated_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetConversationByExternalID :one
SELECT *
FROM conversations
WHERE channel = ? AND external_conversation_id = ?
ORDER BY created_at
LIMIT 1;

-- name: UpdateConversationFromImport :one
UPDATE conversations
SET
	primary_contact_id = COALESCE(primary_contact_id, ?),
	household_id = COALESCE(household_id, ?),
	last_message_at = ?,
	updated_at = ?
WHERE id = ?
RETURNING *;

-- name: CreateMessage :one
INSERT INTO messages (
	id,
	conversation_id,
	sender_contact_id,
	direction,
	body,
	sent_at,
	imported_at,
	external_message_id,
	sender_external_id
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetImportedMessageByExternalID :one
SELECT m.*
FROM messages AS m
JOIN conversations AS c ON c.id = m.conversation_id
WHERE c.channel = ? AND c.external_conversation_id = ? AND m.external_message_id = ?
ORDER BY m.imported_at
LIMIT 1;

-- name: CreateDetectedRequest :one
INSERT INTO detected_requests (
	id,
	message_id,
	household_id,
	contact_id,
	detected_service_type,
	detected_start_at,
	detected_end_at,
	confidence,
	status,
	raw_payload_json,
	created_at,
	updated_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'needs_review', ?, ?, ?)
RETURNING *;

-- name: GetDetectedRequestByMessage :one
SELECT *
FROM detected_requests
WHERE message_id = ?
ORDER BY created_at
LIMIT 1;

-- name: ListDetectedRequests :many
SELECT
	dr.*,
	m.body AS message_body,
	m.sent_at AS message_sent_at,
	m.external_message_id,
	m.sender_external_id,
	c.channel,
	c.external_conversation_id,
	b.start_at AS converted_booking_start_at,
	b.household_id AS converted_booking_household_id,
	COALESCE(h.display_name, '') AS household_name,
	COALESCE(ct.display_name, '') AS contact_name
FROM detected_requests AS dr
JOIN messages AS m ON m.id = dr.message_id
JOIN conversations AS c ON c.id = m.conversation_id
LEFT JOIN households AS h ON h.id = dr.household_id
LEFT JOIN contacts AS ct ON ct.id = dr.contact_id
LEFT JOIN bookings AS b ON b.id = dr.converted_booking_id
WHERE (? = '' OR dr.status = ?)
ORDER BY dr.created_at DESC;

-- name: GetDetectedRequestDetail :one
SELECT
	dr.*,
	m.body AS message_body,
	m.sent_at AS message_sent_at,
	m.external_message_id,
	m.sender_external_id,
	c.channel,
	c.external_conversation_id,
	b.start_at AS converted_booking_start_at,
	b.household_id AS converted_booking_household_id,
	COALESCE(h.display_name, '') AS household_name,
	COALESCE(ct.display_name, '') AS contact_name
FROM detected_requests AS dr
JOIN messages AS m ON m.id = dr.message_id
JOIN conversations AS c ON c.id = m.conversation_id
LEFT JOIN households AS h ON h.id = dr.household_id
LEFT JOIN contacts AS ct ON ct.id = dr.contact_id
LEFT JOIN bookings AS b ON b.id = dr.converted_booking_id
WHERE dr.id = ?;

-- name: UpdateDetectedRequestStatus :one
UPDATE detected_requests
SET status = ?, review_notes = ?, updated_at = ?
WHERE id = ?
RETURNING *;

-- name: ConvertDetectedRequest :one
UPDATE detected_requests
SET
	status = 'converted_to_booking',
	converted_booking_id = ?,
	review_notes = ?,
	updated_at = ?
WHERE id = ? AND converted_booking_id IS NULL AND status != 'ignored'
RETURNING *;

-- name: CreateBookingSource :exec
INSERT INTO booking_sources (
	id,
	booking_id,
	message_id,
	detected_request_id,
	source_note,
	created_at
)
VALUES (?, ?, ?, ?, ?, ?);

-- name: GetContactByChannelID :one
SELECT *
FROM contacts
WHERE
	(? = 'telegram' AND telegram_id = ?)
	OR (? = 'whatsapp' AND whatsapp_id = ?)
LIMIT 1;

-- name: LinkContactTelegramIdentity :one
UPDATE contacts
SET telegram_id = ?, updated_at = ?
WHERE id = ?
RETURNING *;

-- name: LinkContactWhatsappIdentity :one
UPDATE contacts
SET whatsapp_id = ?, updated_at = ?
WHERE id = ?
RETURNING *;

-- name: LinkMessageSenderContact :exec
UPDATE messages
SET sender_contact_id = ?
WHERE id = ?;

-- name: LinkConversationContactContext :exec
UPDATE conversations
SET primary_contact_id = ?, household_id = ?, updated_at = ?
WHERE conversations.id = (
	SELECT messages.conversation_id
	FROM messages
	WHERE messages.id = ?
);

-- name: LinkDetectedRequestContext :exec
UPDATE detected_requests
SET contact_id = ?, household_id = ?, updated_at = ?
WHERE id = ?;

-- name: GetPrimaryHouseholdForContact :one
SELECT h.*
FROM households AS h
JOIN household_contacts AS hc ON hc.household_id = h.id
WHERE hc.contact_id = ? AND h.active = 1
ORDER BY hc.is_primary DESC, hc.created_at
LIMIT 1;
