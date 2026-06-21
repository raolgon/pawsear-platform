package config

import "testing"

func TestLoadUsesCleanLocalDatabaseByDefault(t *testing.T) {
	t.Setenv("PAWSEAR_HTTP_ADDR", "")
	t.Setenv("PAWSEAR_DB_PATH", "")
	t.Setenv("PAWSEAR_SEED_DEMO", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.DatabasePath != "../../data/pawsear-local.db" {
		t.Fatalf("DatabasePath = %q", cfg.DatabasePath)
	}
	if cfg.SeedDemoData {
		t.Fatal("SeedDemoData should be disabled by default")
	}
}

func TestLoadEnablesDemoSeedExplicitly(t *testing.T) {
	t.Setenv("PAWSEAR_SEED_DEMO", "true")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !cfg.SeedDemoData {
		t.Fatal("SeedDemoData should be enabled when explicitly requested")
	}
}
