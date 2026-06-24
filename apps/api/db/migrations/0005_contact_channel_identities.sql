CREATE TABLE contact_channel_identities (
	id TEXT PRIMARY KEY,
	contact_id TEXT NOT NULL REFERENCES contacts(id) ON DELETE RESTRICT,
	channel TEXT NOT NULL CHECK (channel IN ('telegram', 'whatsapp')),
	external_user_id TEXT NOT NULL,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	UNIQUE (channel, external_user_id),
	UNIQUE (contact_id, channel, external_user_id)
);

CREATE INDEX idx_contact_channel_identities_contact
ON contact_channel_identities(contact_id, channel);

INSERT INTO contact_channel_identities (
	id, contact_id, channel, external_user_id, created_at, updated_at
)
SELECT lower(hex(randomblob(16))), id, 'telegram', telegram_id, created_at, updated_at
FROM contacts
WHERE telegram_id IS NOT NULL AND telegram_id != ''
UNION ALL
SELECT lower(hex(randomblob(16))), id, 'whatsapp', whatsapp_id, created_at, updated_at
FROM contacts
WHERE whatsapp_id IS NOT NULL AND whatsapp_id != '';
