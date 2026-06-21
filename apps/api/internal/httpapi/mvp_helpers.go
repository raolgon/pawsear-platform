package httpapi

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	dbqueries "github.com/pawsear/pawsear-platform/apps/api/internal/db/queries"
	"github.com/pawsear/pawsear-platform/apps/api/internal/platform/id"
)

type mvpHandler struct {
	db      *sql.DB
	queries *dbqueries.Queries
	now     func() time.Time
}

func newMVPHandler(database *sql.DB) *mvpHandler {
	return &mvpHandler{
		db:      database,
		queries: dbqueries.New(database),
		now:     func() time.Time { return time.Now().UTC() },
	}
}

func newRecordID() (string, error) {
	value, err := id.New()
	if err != nil {
		return "", fmt.Errorf("create id: %w", err)
	}
	return value, nil
}

func timestamp(now func() time.Time) string {
	return now().UTC().Format(time.RFC3339Nano)
}

func optionalText(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: trimmed, Valid: true}
}

func textValue(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func requiredText(value string, field string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("%s is required", field)
	}
	return trimmed, nil
}

func defaultText(value string, fallback string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fallback
	}
	return trimmed
}

func intFlag(value bool) int64 {
	if value {
		return 1
	}
	return 0
}

func boolFlag(value int64) bool {
	return value != 0
}

func writeAPIError(w http.ResponseWriter, status int, code string, err error) {
	writeJSON(w, status, errorResponse{
		Error:   code,
		Message: err.Error(),
	})
}

func writeInvalid(w http.ResponseWriter, err error) {
	writeAPIError(w, http.StatusBadRequest, "invalid_request", err)
}

func writeNotFound(w http.ResponseWriter, err error) {
	writeAPIError(w, http.StatusNotFound, "not_found", err)
}

func writeStoreError(w http.ResponseWriter, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		writeNotFound(w, err)
		return
	}
	writeAPIError(w, http.StatusInternalServerError, "internal_error", err)
}

func allowed(value string, values ...string) bool {
	for _, candidate := range values {
		if value == candidate {
			return true
		}
	}
	return false
}

func parseDayRange(raw string, now func() time.Time) (string, string, error) {
	if strings.TrimSpace(raw) == "" {
		raw = now().UTC().Format(time.DateOnly)
	}
	day, err := time.Parse(time.DateOnly, raw)
	if err != nil {
		return "", "", fmt.Errorf("date must use YYYY-MM-DD")
	}

	start := day.Format(time.DateOnly) + "T00:00:00Z"
	end := day.AddDate(0, 0, 1).Format(time.DateOnly) + "T00:00:00Z"
	return start, end, nil
}
