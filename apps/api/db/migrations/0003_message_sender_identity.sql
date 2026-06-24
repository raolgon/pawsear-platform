ALTER TABLE messages ADD COLUMN sender_external_id TEXT;

UPDATE messages
SET sender_external_id = (
	SELECT json_extract(detected_requests.raw_payload_json, '$.senderExternalId')
	FROM detected_requests
	WHERE detected_requests.message_id = messages.id
	LIMIT 1
)
WHERE sender_external_id IS NULL;

CREATE INDEX idx_messages_sender_external_id ON messages(sender_external_id);
CREATE UNIQUE INDEX idx_contacts_telegram_id_unique
ON contacts(telegram_id)
WHERE telegram_id IS NOT NULL AND telegram_id != '';
CREATE UNIQUE INDEX idx_contacts_whatsapp_id_unique
ON contacts(whatsapp_id)
WHERE whatsapp_id IS NOT NULL AND whatsapp_id != '';
