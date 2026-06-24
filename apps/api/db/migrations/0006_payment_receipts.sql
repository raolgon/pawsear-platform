CREATE TABLE payment_receipts (
	id TEXT PRIMARY KEY,
	payment_id TEXT NOT NULL REFERENCES payments(id) ON DELETE RESTRICT,
	receipt_number TEXT NOT NULL UNIQUE,
	snapshot_json TEXT NOT NULL,
	issued_at TEXT NOT NULL,
	created_at TEXT NOT NULL,
	UNIQUE (payment_id)
);

CREATE INDEX idx_payment_receipts_issued_at ON payment_receipts(issued_at);
