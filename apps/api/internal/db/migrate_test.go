package db

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
)

func TestMigrateCreatesInitialSchema(t *testing.T) {
	database := openMigratedTestDB(t)

	tables := []string{
		"households",
		"contacts",
		"household_contacts",
		"pets",
		"bookings",
		"care_tasks",
		"charges",
		"payments",
		"payment_allocations",
		"messages",
		"detected_requests",
		"booking_sources",
	}

	for _, table := range tables {
		t.Run(table, func(t *testing.T) {
			if !tableExists(t, database, table) {
				t.Fatalf("expected table %q to exist", table)
			}
		})
	}
}

func TestMigrateIsIdempotent(t *testing.T) {
	database := openMigratedTestDB(t)

	if err := Migrate(context.Background(), database); err != nil {
		t.Fatalf("run migrations again: %v", err)
	}

	var applied int
	err := database.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", "0001_initial_schema.sql").Scan(&applied)
	if err != nil {
		t.Fatalf("count applied migration: %v", err)
	}
	if applied != 1 {
		t.Fatalf("expected migration recorded once, got %d", applied)
	}
}

func TestOpenEnablesForeignKeys(t *testing.T) {
	database := openMigratedTestDB(t)

	_, err := database.Exec(`
		INSERT INTO pets (
			id,
			household_id,
			name,
			species,
			created_at,
			updated_at
		)
		VALUES ('pet_missing_household', 'missing_household', 'Luna', 'dog', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z')
	`)
	if err == nil {
		t.Fatal("expected foreign key violation")
	}
}

func openMigratedTestDB(t *testing.T) *sql.DB {
	t.Helper()

	database, err := Open(filepath.Join(t.TempDir(), "nested", "pawsear-test.db"))
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() {
		_ = database.Close()
	})

	if err := Migrate(context.Background(), database); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	return database
}

func tableExists(t *testing.T, database *sql.DB, table string) bool {
	t.Helper()

	var name string
	err := database.QueryRow(
		"SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?",
		table,
	).Scan(&name)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		t.Fatalf("check table %q: %v", table, err)
	}

	return name == table
}
