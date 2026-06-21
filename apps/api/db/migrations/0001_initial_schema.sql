CREATE TABLE households (
	id TEXT PRIMARY KEY,
	display_name TEXT NOT NULL,
	address_line1 TEXT,
	address_line2 TEXT,
	neighborhood TEXT,
	city TEXT,
	timezone TEXT,
	notes TEXT,
	active INTEGER NOT NULL DEFAULT 1,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE INDEX idx_households_display_name ON households(display_name);
CREATE INDEX idx_households_active ON households(active);

CREATE TABLE contacts (
	id TEXT PRIMARY KEY,
	display_name TEXT NOT NULL,
	phone TEXT,
	whatsapp_id TEXT,
	telegram_id TEXT,
	email TEXT,
	notes TEXT,
	active INTEGER NOT NULL DEFAULT 1,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE INDEX idx_contacts_display_name ON contacts(display_name);
CREATE INDEX idx_contacts_phone ON contacts(phone);
CREATE INDEX idx_contacts_whatsapp_id ON contacts(whatsapp_id);
CREATE INDEX idx_contacts_telegram_id ON contacts(telegram_id);

CREATE TABLE household_contacts (
	household_id TEXT NOT NULL REFERENCES households(id) ON DELETE RESTRICT,
	contact_id TEXT NOT NULL REFERENCES contacts(id) ON DELETE RESTRICT,
	role TEXT NOT NULL CHECK (role IN ('owner', 'partner', 'family', 'domestic_worker', 'payer', 'emergency_contact', 'vet', 'other')),
	is_primary INTEGER NOT NULL DEFAULT 0,
	notes TEXT,
	created_at TEXT NOT NULL,
	PRIMARY KEY (household_id, contact_id, role)
);

CREATE INDEX idx_household_contacts_contact_id ON household_contacts(contact_id);
CREATE INDEX idx_household_contacts_role ON household_contacts(role);

CREATE TABLE pets (
	id TEXT PRIMARY KEY,
	household_id TEXT NOT NULL REFERENCES households(id) ON DELETE RESTRICT,
	name TEXT NOT NULL,
	species TEXT NOT NULL CHECK (species IN ('dog', 'cat', 'other')),
	breed TEXT,
	size TEXT CHECK (size IS NULL OR size IN ('small', 'medium', 'large', 'giant', 'unknown')),
	sex TEXT CHECK (sex IS NULL OR sex IN ('female', 'male', 'unknown')),
	birthdate TEXT,
	color_markings TEXT,
	behavior_notes TEXT,
	medical_notes TEXT,
	feeding_notes TEXT,
	vet_notes TEXT,
	active INTEGER NOT NULL DEFAULT 1,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE INDEX idx_pets_household_id ON pets(household_id);
CREATE INDEX idx_pets_name ON pets(name);
CREATE INDEX idx_pets_active ON pets(active);

CREATE TABLE pet_medications (
	id TEXT PRIMARY KEY,
	pet_id TEXT NOT NULL REFERENCES pets(id) ON DELETE RESTRICT,
	name TEXT NOT NULL,
	dosage TEXT,
	instructions TEXT,
	schedule_note TEXT,
	active INTEGER NOT NULL DEFAULT 1,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE INDEX idx_pet_medications_pet_id ON pet_medications(pet_id);

CREATE TABLE pet_diets (
	id TEXT PRIMARY KEY,
	pet_id TEXT NOT NULL REFERENCES pets(id) ON DELETE RESTRICT,
	food_name TEXT,
	amount TEXT,
	schedule_note TEXT,
	instructions TEXT,
	active INTEGER NOT NULL DEFAULT 1,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE INDEX idx_pet_diets_pet_id ON pet_diets(pet_id);

CREATE TABLE staff_members (
	id TEXT PRIMARY KEY,
	display_name TEXT NOT NULL,
	phone TEXT,
	role TEXT NOT NULL CHECK (role IN ('owner_operator', 'walker', 'sitter', 'admin')),
	active INTEGER NOT NULL DEFAULT 1,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE INDEX idx_staff_members_active ON staff_members(active);

CREATE TABLE bookings (
	id TEXT PRIMARY KEY,
	household_id TEXT NOT NULL REFERENCES households(id) ON DELETE RESTRICT,
	service_type TEXT NOT NULL CHECK (service_type IN ('walk', 'client_home_sitting', 'boarding', 'visit', 'transport', 'other')),
	status TEXT NOT NULL CHECK (status IN ('requested', 'confirmed', 'in_progress', 'completed', 'cancelled')),
	start_at TEXT NOT NULL,
	end_at TEXT,
	location_type TEXT NOT NULL CHECK (location_type IN ('household_home', 'caregiver_home', 'other')),
	address_snapshot TEXT,
	requested_by_contact_id TEXT REFERENCES contacts(id) ON DELETE RESTRICT,
	assigned_staff_id TEXT REFERENCES staff_members(id) ON DELETE RESTRICT,
	source TEXT NOT NULL DEFAULT 'manual' CHECK (source IN ('manual', 'whatsapp', 'telegram', 'import')),
	notes TEXT,
	completed_at TEXT,
	cancelled_at TEXT,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE INDEX idx_bookings_household_id ON bookings(household_id);
CREATE INDEX idx_bookings_start_at ON bookings(start_at);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_requested_by_contact_id ON bookings(requested_by_contact_id);
CREATE INDEX idx_bookings_assigned_staff_id ON bookings(assigned_staff_id);

CREATE TABLE booking_pets (
	booking_id TEXT NOT NULL REFERENCES bookings(id) ON DELETE RESTRICT,
	pet_id TEXT NOT NULL REFERENCES pets(id) ON DELETE RESTRICT,
	notes TEXT,
	PRIMARY KEY (booking_id, pet_id)
);

CREATE INDEX idx_booking_pets_pet_id ON booking_pets(pet_id);

CREATE TABLE care_routines (
	id TEXT PRIMARY KEY,
	household_id TEXT NOT NULL REFERENCES households(id) ON DELETE RESTRICT,
	pet_id TEXT REFERENCES pets(id) ON DELETE RESTRICT,
	task_type TEXT NOT NULL CHECK (task_type IN ('food', 'medicine', 'walk', 'water', 'cleaning', 'pickup', 'dropoff', 'photo_update', 'other')),
	title TEXT NOT NULL,
	instructions TEXT,
	schedule_note TEXT,
	active INTEGER NOT NULL DEFAULT 1,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE INDEX idx_care_routines_household_id ON care_routines(household_id);
CREATE INDEX idx_care_routines_pet_id ON care_routines(pet_id);
CREATE INDEX idx_care_routines_active ON care_routines(active);

CREATE TABLE care_tasks (
	id TEXT PRIMARY KEY,
	booking_id TEXT REFERENCES bookings(id) ON DELETE RESTRICT,
	household_id TEXT NOT NULL REFERENCES households(id) ON DELETE RESTRICT,
	pet_id TEXT REFERENCES pets(id) ON DELETE RESTRICT,
	task_type TEXT NOT NULL CHECK (task_type IN ('food', 'medicine', 'walk', 'water', 'cleaning', 'pickup', 'dropoff', 'photo_update', 'other')),
	title TEXT NOT NULL,
	instructions TEXT,
	due_at TEXT,
	status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'skipped', 'cancelled')),
	assigned_staff_id TEXT REFERENCES staff_members(id) ON DELETE RESTRICT,
	completed_at TEXT,
	completed_by_staff_id TEXT REFERENCES staff_members(id) ON DELETE RESTRICT,
	skipped_reason TEXT,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE INDEX idx_care_tasks_due_at ON care_tasks(due_at);
CREATE INDEX idx_care_tasks_status ON care_tasks(status);
CREATE INDEX idx_care_tasks_booking_id ON care_tasks(booking_id);
CREATE INDEX idx_care_tasks_pet_id ON care_tasks(pet_id);

CREATE TABLE charges (
	id TEXT PRIMARY KEY,
	household_id TEXT NOT NULL REFERENCES households(id) ON DELETE RESTRICT,
	booking_id TEXT REFERENCES bookings(id) ON DELETE RESTRICT,
	description TEXT NOT NULL,
	amount_minor INTEGER NOT NULL CHECK (amount_minor > 0),
	currency TEXT NOT NULL DEFAULT 'MXN',
	status TEXT NOT NULL DEFAULT 'unpaid' CHECK (status IN ('unpaid', 'partially_paid', 'paid', 'waived', 'void')),
	due_date TEXT,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE INDEX idx_charges_household_id ON charges(household_id);
CREATE INDEX idx_charges_booking_id ON charges(booking_id);
CREATE INDEX idx_charges_status ON charges(status);
CREATE INDEX idx_charges_due_date ON charges(due_date);

CREATE TABLE payments (
	id TEXT PRIMARY KEY,
	payer_contact_id TEXT REFERENCES contacts(id) ON DELETE RESTRICT,
	received_at TEXT NOT NULL,
	amount_minor INTEGER NOT NULL CHECK (amount_minor > 0),
	currency TEXT NOT NULL DEFAULT 'MXN',
	method TEXT NOT NULL CHECK (method IN ('cash', 'bank_transfer', 'card_external', 'other')),
	reference TEXT,
	notes TEXT,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE INDEX idx_payments_payer_contact_id ON payments(payer_contact_id);
CREATE INDEX idx_payments_received_at ON payments(received_at);

CREATE TABLE payment_allocations (
	id TEXT PRIMARY KEY,
	payment_id TEXT NOT NULL REFERENCES payments(id) ON DELETE RESTRICT,
	charge_id TEXT NOT NULL REFERENCES charges(id) ON DELETE RESTRICT,
	amount_minor INTEGER NOT NULL CHECK (amount_minor > 0),
	created_at TEXT NOT NULL,
	UNIQUE (payment_id, charge_id)
);

CREATE INDEX idx_payment_allocations_payment_id ON payment_allocations(payment_id);
CREATE INDEX idx_payment_allocations_charge_id ON payment_allocations(charge_id);

CREATE TABLE conversations (
	id TEXT PRIMARY KEY,
	channel TEXT NOT NULL CHECK (channel IN ('whatsapp', 'telegram', 'manual')),
	external_conversation_id TEXT,
	primary_contact_id TEXT REFERENCES contacts(id) ON DELETE RESTRICT,
	household_id TEXT REFERENCES households(id) ON DELETE RESTRICT,
	last_message_at TEXT,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE INDEX idx_conversations_channel_external ON conversations(channel, external_conversation_id);
CREATE INDEX idx_conversations_primary_contact_id ON conversations(primary_contact_id);
CREATE INDEX idx_conversations_household_id ON conversations(household_id);

CREATE TABLE messages (
	id TEXT PRIMARY KEY,
	conversation_id TEXT NOT NULL REFERENCES conversations(id) ON DELETE RESTRICT,
	sender_contact_id TEXT REFERENCES contacts(id) ON DELETE RESTRICT,
	direction TEXT NOT NULL CHECK (direction IN ('inbound', 'outbound', 'system')),
	body TEXT NOT NULL,
	sent_at TEXT,
	imported_at TEXT NOT NULL,
	external_message_id TEXT
);

CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX idx_messages_sender_contact_id ON messages(sender_contact_id);
CREATE INDEX idx_messages_sent_at ON messages(sent_at);

CREATE TABLE detected_requests (
	id TEXT PRIMARY KEY,
	message_id TEXT NOT NULL REFERENCES messages(id) ON DELETE RESTRICT,
	household_id TEXT REFERENCES households(id) ON DELETE RESTRICT,
	contact_id TEXT REFERENCES contacts(id) ON DELETE RESTRICT,
	detected_service_type TEXT CHECK (detected_service_type IS NULL OR detected_service_type IN ('walk', 'client_home_sitting', 'boarding', 'visit', 'transport', 'other')),
	detected_start_at TEXT,
	detected_end_at TEXT,
	confidence TEXT NOT NULL DEFAULT 'unknown' CHECK (confidence IN ('low', 'medium', 'high', 'unknown')),
	status TEXT NOT NULL DEFAULT 'needs_review' CHECK (status IN ('needs_review', 'confirmed', 'ignored', 'needs_more_info', 'converted_to_booking')),
	converted_booking_id TEXT REFERENCES bookings(id) ON DELETE RESTRICT,
	raw_payload_json TEXT,
	review_notes TEXT,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE INDEX idx_detected_requests_status ON detected_requests(status);
CREATE INDEX idx_detected_requests_message_id ON detected_requests(message_id);
CREATE INDEX idx_detected_requests_converted_booking_id ON detected_requests(converted_booking_id);

CREATE TABLE booking_sources (
	id TEXT PRIMARY KEY,
	booking_id TEXT NOT NULL REFERENCES bookings(id) ON DELETE RESTRICT,
	message_id TEXT REFERENCES messages(id) ON DELETE RESTRICT,
	detected_request_id TEXT REFERENCES detected_requests(id) ON DELETE RESTRICT,
	source_note TEXT,
	created_at TEXT NOT NULL,
	CHECK (message_id IS NOT NULL OR detected_request_id IS NOT NULL)
);

CREATE INDEX idx_booking_sources_booking_id ON booking_sources(booking_id);
CREATE INDEX idx_booking_sources_message_id ON booking_sources(message_id);
CREATE INDEX idx_booking_sources_detected_request_id ON booking_sources(detected_request_id);
