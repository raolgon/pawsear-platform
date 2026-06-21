package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func SeedDemo(ctx context.Context, database *sql.DB) error {
	var count int
	if err := database.QueryRowContext(ctx, "SELECT COUNT(*) FROM households").Scan(&count); err != nil {
		return fmt.Errorf("count demo households: %w", err)
	}
	if count > 0 {
		return nil
	}

	now := time.Now().UTC()
	today := now.Format(time.DateOnly)
	ts := now.Format(time.RFC3339Nano)
	at := func(clock string) string {
		return today + "T" + clock + ":00Z"
	}

	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin demo seed: %w", err)
	}
	defer tx.Rollback()

	statements := []string{
		`INSERT INTO households (id, display_name, address_line1, neighborhood, city, timezone, notes, created_at, updated_at) VALUES
			('demo-household-garcia', 'Casa Garcia', 'Calle Tonala 120', 'Roma Norte', 'CDMX', 'America/Mexico_City', 'Gate code 4421. Max should use the front clip harness.', ?, ?),
			('demo-household-torres', 'Casa Torres', 'Av Amsterdam 45', 'Condesa', 'CDMX', 'America/Mexico_City', 'Use side entrance. Mia needs medicine after food.', ?, ?),
			('demo-household-rivera', 'Casa Rivera', 'Zempoala 77', 'Narvarte', 'CDMX', 'America/Mexico_City', 'Confirm requested walk time before booking.', ?, ?)`,
		`INSERT INTO contacts (id, display_name, phone, whatsapp_id, email, notes, created_at, updated_at) VALUES
			('demo-contact-mariana', 'Mariana Garcia', '+525500000101', 'whatsapp-mariana', 'mariana@example.test', 'Primary contact for Casa Garcia.', ?, ?),
			('demo-contact-sofia', 'Sofia Garcia', '+525500000102', 'whatsapp-sofia', 'sofia@example.test', 'Payer for Casa Garcia.', ?, ?),
			('demo-contact-daniela', 'Daniela Torres', '+525500000201', 'whatsapp-daniela', 'daniela@example.test', 'Primary contact for Casa Torres.', ?, ?),
			('demo-contact-pablo', 'Pablo Torres', '+525500000202', NULL, 'pablo@example.test', 'Emergency contact and payer.', ?, ?),
			('demo-contact-claudia', 'Claudia Rivera', '+525500000301', 'whatsapp-claudia', 'claudia@example.test', 'Payer for Casa Rivera.', ?, ?)`,
		`INSERT INTO household_contacts (household_id, contact_id, role, is_primary, notes, created_at) VALUES
			('demo-household-garcia', 'demo-contact-mariana', 'owner', 1, 'Primary WhatsApp contact.', ?),
			('demo-household-garcia', 'demo-contact-sofia', 'payer', 1, 'Pays weekly.', ?),
			('demo-household-torres', 'demo-contact-daniela', 'owner', 1, 'Primary WhatsApp contact.', ?),
			('demo-household-torres', 'demo-contact-pablo', 'emergency_contact', 1, 'Also payer.', ?),
			('demo-household-rivera', 'demo-contact-claudia', 'payer', 1, 'Confirm payment by transfer.', ?)`,
		`INSERT INTO pets (id, household_id, name, species, breed, size, sex, behavior_notes, medical_notes, feeding_notes, vet_notes, created_at, updated_at) VALUES
			('demo-pet-luna', 'demo-household-garcia', 'Luna', 'dog', 'Mixed', 'medium', 'female', 'Nervous around bikes.', NULL, '1 cup morning and night.', 'Vet: Central Vet.', ?, ?),
			('demo-pet-max', 'demo-household-garcia', 'Max', 'dog', 'Labrador mix', 'large', 'male', 'Pulls on leash.', 'Joint supplement nightly.', '1.5 cups morning.', 'Use front clip harness.', ?, ?),
			('demo-pet-mia', 'demo-household-torres', 'Mia', 'cat', 'Domestic shorthair', 'small', 'female', 'Hides under bed when nervous.', 'Half tablet after food.', 'Wet food morning/night.', 'Indoor only.', ?, ?),
			('demo-pet-bruno', 'demo-household-rivera', 'Bruno', 'dog', 'Terrier mix', 'small', 'male', 'Reactive near scooters.', NULL, '1 cup after walks.', 'Confirm walk timing.', ?, ?)`,
		`INSERT INTO staff_members (id, display_name, phone, role, created_at, updated_at) VALUES
			('demo-staff-rafa', 'Rafa', '+525500001001', 'walker', ?, ?),
			('demo-staff-vale', 'Vale', '+525500001002', 'sitter', ?, ?)`,
	}

	for _, statement := range statements {
		if _, err := tx.ExecContext(ctx, statement, repeated(ts, countPlaceholders(statement))...); err != nil {
			return fmt.Errorf("insert demo records: %w", err)
		}
	}

	if err := insertDemoOperations(ctx, tx, ts, at); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit demo seed: %w", err)
	}
	return nil
}

func insertDemoOperations(ctx context.Context, tx *sql.Tx, ts string, at func(string) string) error {
	statements := []struct {
		sql  string
		args []any
	}{
		{`INSERT INTO bookings (id, household_id, service_type, status, start_at, end_at, location_type, requested_by_contact_id, assigned_staff_id, source, notes, created_at, updated_at) VALUES
			('demo-booking-garcia-walk', 'demo-household-garcia', 'walk', 'confirmed', ?, ?, 'household_home', 'demo-contact-mariana', 'demo-staff-rafa', 'manual', 'Weekday walk for Luna and Max.', ?, ?),
			('demo-booking-rivera-visit', 'demo-household-rivera', 'visit', 'requested', ?, ?, 'household_home', 'demo-contact-claudia', NULL, 'whatsapp', 'Needs review: client asked if 08:00 works.', ?, ?)`,
			[]any{at("09:00"), at("10:00"), ts, ts, at("16:30"), at("17:00"), ts, ts}},
		{`INSERT INTO booking_pets (booking_id, pet_id, notes) VALUES
			('demo-booking-garcia-walk', 'demo-pet-luna', NULL),
			('demo-booking-garcia-walk', 'demo-pet-max', NULL),
			('demo-booking-rivera-visit', 'demo-pet-bruno', 'Confirm whether this should become a walk.')`, nil},
		{`INSERT INTO care_tasks (id, booking_id, household_id, pet_id, task_type, title, instructions, due_at, status, assigned_staff_id, created_at, updated_at) VALUES
			('demo-task-mia-medicine', NULL, 'demo-household-torres', 'demo-pet-mia', 'medicine', 'Mia medicine', 'Give half tablet after food.', ?, 'pending', 'demo-staff-vale', ?, ?),
			('demo-task-garcia-water', 'demo-booking-garcia-walk', 'demo-household-garcia', NULL, 'water', 'Refresh water after walk', 'Leave bowls full before locking up.', ?, 'pending', 'demo-staff-rafa', ?, ?)`,
			[]any{at("08:00"), ts, ts, at("10:00"), ts, ts}},
		{`INSERT INTO charges (id, household_id, booking_id, description, amount_minor, currency, status, due_date, created_at, updated_at) VALUES
			('demo-charge-garcia-walk', 'demo-household-garcia', 'demo-booking-garcia-walk', 'Walk for Luna and Max', 20000, 'MXN', 'unpaid', ?, ?, ?),
			('demo-charge-rivera-visit', 'demo-household-rivera', 'demo-booking-rivera-visit', 'Visit for Bruno', 18000, 'MXN', 'partially_paid', ?, ?, ?)`,
			[]any{at("23:59"), ts, ts, at("23:59"), ts, ts}},
		{`INSERT INTO payments (id, payer_contact_id, received_at, amount_minor, currency, method, reference, notes, created_at, updated_at) VALUES
			('demo-payment-rivera-partial', 'demo-contact-claudia', ?, 9000, 'MXN', 'bank_transfer', 'DEMO-8821', 'Partial payment for Bruno visit.', ?, ?)`,
			[]any{at("12:00"), ts, ts}},
		{`INSERT INTO payment_allocations (id, payment_id, charge_id, amount_minor, created_at) VALUES
			('demo-allocation-rivera-partial', 'demo-payment-rivera-partial', 'demo-charge-rivera-visit', 9000, ?)`,
			[]any{ts}},
	}
	for _, statement := range statements {
		if _, err := tx.ExecContext(ctx, statement.sql, statement.args...); err != nil {
			return fmt.Errorf("insert demo operations: %w", err)
		}
	}
	return nil
}

func countPlaceholders(statement string) int {
	count := 0
	for _, char := range statement {
		if char == '?' {
			count++
		}
	}
	return count
}

func repeated(value string, count int) []any {
	values := make([]any, count)
	for index := range values {
		values[index] = value
	}
	return values
}
