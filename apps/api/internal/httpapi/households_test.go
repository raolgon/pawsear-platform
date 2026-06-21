package httpapi

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/pawsear/pawsear-platform/apps/api/internal/db"
)

func TestHouseholdEndpoints(t *testing.T) {
	database := openTestDB(t)
	router := NewRouter(database)

	createBody := []byte(`{
		"displayName": "Casa de Luna y Max",
		"neighborhood": "Roma Norte",
		"city": "CDMX",
		"timezone": "America/Mexico_City"
	}`)

	createReq := httptest.NewRequest(http.MethodPost, "/api/households", bytes.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRes := httptest.NewRecorder()
	router.ServeHTTP(createRes, createReq)

	if createRes.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d: %s", http.StatusCreated, createRes.Code, createRes.Body.String())
	}

	var created struct {
		ID          string `json:"id"`
		DisplayName string `json:"displayName"`
	}
	if err := json.NewDecoder(createRes.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected created household id")
	}
	if created.DisplayName != "Casa de Luna y Max" {
		t.Fatalf("unexpected displayName %q", created.DisplayName)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/households/"+created.ID, nil)
	getRes := httptest.NewRecorder()
	router.ServeHTTP(getRes, getReq)

	if getRes.Code != http.StatusOK {
		t.Fatalf("expected get status %d, got %d: %s", http.StatusOK, getRes.Code, getRes.Body.String())
	}

	patchBody := []byte(`{"notes":"Cliente frecuente"}`)
	patchReq := httptest.NewRequest(http.MethodPatch, "/api/households/"+created.ID, bytes.NewReader(patchBody))
	patchReq.Header.Set("Content-Type", "application/json")
	patchRes := httptest.NewRecorder()
	router.ServeHTTP(patchRes, patchReq)

	if patchRes.Code != http.StatusOK {
		t.Fatalf("expected patch status %d, got %d: %s", http.StatusOK, patchRes.Code, patchRes.Body.String())
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/households", nil)
	listRes := httptest.NewRecorder()
	router.ServeHTTP(listRes, listReq)

	if listRes.Code != http.StatusOK {
		t.Fatalf("expected list status %d, got %d: %s", http.StatusOK, listRes.Code, listRes.Body.String())
	}

	var list struct {
		Households []struct {
			ID string `json:"id"`
		} `json:"households"`
	}
	if err := json.NewDecoder(listRes.Body).Decode(&list); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if len(list.Households) != 1 {
		t.Fatalf("expected 1 household, got %d", len(list.Households))
	}
}

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	path := filepath.Join(t.TempDir(), "pawsear-test.db")
	database, err := db.Open(path)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() {
		_ = database.Close()
	})

	if err := db.Migrate(context.Background(), database); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	return database
}
