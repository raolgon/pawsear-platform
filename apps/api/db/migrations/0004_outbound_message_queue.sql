CREATE TABLE outbound_messages (
	id TEXT PRIMARY KEY,
	detected_request_id TEXT NOT NULL REFERENCES detected_requests(id) ON DELETE RESTRICT,
	channel TEXT NOT NULL CHECK (channel IN ('telegram')),
	recipient_external_id TEXT NOT NULL,
	template_key TEXT NOT NULL CHECK (template_key IN ('request_details', 'booking_confirmed', 'request_declined')),
	body TEXT NOT NULL,
	status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'failed', 'cancelled')),
	attempts INTEGER NOT NULL DEFAULT 0 CHECK (attempts >= 0),
	last_error TEXT,
	created_at TEXT NOT NULL,
	sent_at TEXT,
	updated_at TEXT NOT NULL
);

CREATE INDEX idx_outbound_messages_status_created ON outbound_messages(status, created_at);
CREATE INDEX idx_outbound_messages_detected_request ON outbound_messages(detected_request_id, created_at);
