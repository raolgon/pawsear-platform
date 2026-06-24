-- name: CreateOutboundMessage :one
INSERT INTO outbound_messages (
	id,
	detected_request_id,
	channel,
	recipient_external_id,
	template_key,
	body,
	status,
	attempts,
	created_at,
	updated_at
)
VALUES (?, ?, ?, ?, ?, ?, 'pending', 0, ?, ?)
RETURNING *;

-- name: ListOutboundMessages :many
SELECT *
FROM outbound_messages
WHERE (sqlc.arg(status_filter) = '' OR status = sqlc.arg(status_filter))
  AND (sqlc.arg(request_filter) = '' OR detected_request_id = sqlc.arg(request_filter))
ORDER BY created_at DESC
LIMIT sqlc.arg(result_limit);

-- name: GetOutboundMessage :one
SELECT * FROM outbound_messages WHERE id = ?;

-- name: GetPendingOutboundMessageForRequest :one
SELECT *
FROM outbound_messages
WHERE detected_request_id = ? AND status = 'pending'
ORDER BY created_at DESC
LIMIT 1;

-- name: MarkOutboundMessageSent :one
UPDATE outbound_messages
SET status = 'sent', attempts = attempts + 1, last_error = NULL,
	sent_at = ?, updated_at = ?
WHERE id = ? AND status IN ('pending', 'failed')
RETURNING *;

-- name: MarkOutboundMessageFailed :one
UPDATE outbound_messages
SET status = 'failed', attempts = attempts + 1, last_error = ?, updated_at = ?
WHERE id = ? AND status IN ('pending', 'failed')
RETURNING *;
