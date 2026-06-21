package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMVPBackendFlow(t *testing.T) {
	router := NewRouter(openTestDB(t))

	household := requestJSON(t, router, http.MethodPost, "/api/households", `{"displayName":"Casa de Mora"}`, http.StatusCreated)
	householdID := stringField(t, household, "id")

	contact := requestJSON(t, router, http.MethodPost, "/api/contacts", `{"displayName":"Ana","phone":"+525500000000"}`, http.StatusCreated)
	contactID := stringField(t, contact, "id")

	requestJSON(t, router, http.MethodPost, "/api/households/"+householdID+"/contacts", `{"contactId":"`+contactID+`","role":"owner","isPrimary":true}`, http.StatusCreated)

	pet := requestJSON(t, router, http.MethodPost, "/api/pets", `{"householdId":"`+householdID+`","name":"Mora","species":"dog"}`, http.StatusCreated)
	petID := stringField(t, pet, "id")

	staff := requestJSON(t, router, http.MethodPost, "/api/staff", `{"displayName":"Rafa","role":"walker"}`, http.StatusCreated)
	staffID := stringField(t, staff, "id")

	booking := requestJSON(t, router, http.MethodPost, "/api/bookings", `{
		"householdId":"`+householdID+`",
		"serviceType":"walk",
		"status":"confirmed",
		"startAt":"2026-06-05T15:00:00Z",
		"assignedStaffId":"`+staffID+`",
		"petIds":["`+petID+`"]
	}`, http.StatusCreated)
	bookingID := stringField(t, booking, "id")

	task := requestJSON(t, router, http.MethodPost, "/api/care-tasks", `{
		"bookingId":"`+bookingID+`",
		"householdId":"`+householdID+`",
		"petId":"`+petID+`",
		"taskType":"walk",
		"title":"Paseo de tarde",
		"dueAt":"2026-06-05T15:00:00Z"
	}`, http.StatusCreated)
	if stringField(t, task, "status") != "pending" {
		t.Fatalf("expected pending task, got %v", task["status"])
	}
	requestJSON(t, router, http.MethodPatch, "/api/bookings/"+bookingID, `{"status":"in_progress"}`, http.StatusOK)
	completedBooking := requestJSON(t, router, http.MethodPatch, "/api/bookings/"+bookingID, `{"status":"completed"}`, http.StatusOK)
	if completedBooking["completedAt"] == nil {
		t.Fatalf("expected completedAt to be set")
	}
	completedTask := requestJSON(t, router, http.MethodPatch, "/api/care-tasks/"+stringField(t, task, "id"), `{"status":"completed"}`, http.StatusOK)
	if completedTask["completedAt"] == nil {
		t.Fatalf("expected task completedAt to be set")
	}
	skippedTask := requestJSON(t, router, http.MethodPost, "/api/care-tasks", `{
		"bookingId":"`+bookingID+`",
		"householdId":"`+householdID+`",
		"taskType":"photo_update",
		"title":"Enviar foto"
	}`, http.StatusCreated)
	skippedTaskID := stringField(t, skippedTask, "id")
	requestJSON(t, router, http.MethodPatch, "/api/care-tasks/"+skippedTaskID, `{"status":"skipped"}`, http.StatusBadRequest)
	skippedTask = requestJSON(t, router, http.MethodPatch, "/api/care-tasks/"+skippedTaskID, `{"status":"skipped","skippedReason":"Client declined"}`, http.StatusOK)
	if skippedTask["status"] != "skipped" {
		t.Fatalf("expected skipped task, got %#v", skippedTask)
	}

	charge := requestJSON(t, router, http.MethodPost, "/api/charges", `{
		"householdId":"`+householdID+`",
		"bookingId":"`+bookingID+`",
		"description":"Paseo de Mora",
		"amountMinor":20000
	}`, http.StatusCreated)
	chargeID := stringField(t, charge, "id")

	payment := requestJSON(t, router, http.MethodPost, "/api/payments", `{
		"payerContactId":"`+contactID+`",
		"receivedAt":"2026-06-05T18:00:00Z",
		"amountMinor":8000,
		"method":"cash",
		"allocations":[{"chargeId":"`+chargeID+`","amountMinor":8000}]
	}`, http.StatusCreated)
	allocations := payment["allocations"].([]any)
	if len(allocations) != 1 {
		t.Fatalf("expected one allocation, got %d", len(allocations))
	}

	charges := requestJSON(t, router, http.MethodGet, "/api/charges?householdId="+householdID, "", http.StatusOK)
	firstCharge := charges["charges"].([]any)[0].(map[string]any)
	if firstCharge["status"] != "partially_paid" || firstCharge["outstandingMinor"] != float64(12000) {
		t.Fatalf("expected partially paid charge with 12000 outstanding, got %#v", firstCharge)
	}
	requestJSON(t, router, http.MethodPost, "/api/payments", `{
		"amountMinor":13000,
		"method":"cash",
		"allocations":[{"chargeId":"`+chargeID+`","amountMinor":13000}]
	}`, http.StatusBadRequest)
	requestJSON(t, router, http.MethodPost, "/api/payments", `{
		"amountMinor":12000,
		"method":"cash",
		"allocations":[{"chargeId":"`+chargeID+`","amountMinor":12000}]
	}`, http.StatusCreated)
	charges = requestJSON(t, router, http.MethodGet, "/api/charges?householdId="+householdID, "", http.StatusOK)
	firstCharge = charges["charges"].([]any)[0].(map[string]any)
	if firstCharge["status"] != "paid" || firstCharge["outstandingMinor"] != float64(0) {
		t.Fatalf("expected paid charge with no outstanding balance, got %#v", firstCharge)
	}

	dashboard := requestJSON(t, router, http.MethodGet, "/api/dashboard/today?date=2026-06-05", "", http.StatusOK)
	if len(dashboard["bookings"].([]any)) != 1 {
		t.Fatalf("expected one dashboard booking")
	}
	if len(dashboard["careTasks"].([]any)) != 1 {
		t.Fatalf("expected one dashboard care task")
	}
}

func requestJSON(t *testing.T, router http.Handler, method string, path string, body string, wantStatus int) map[string]any {
	t.Helper()

	var reader *bytes.Reader
	if body == "" {
		reader = bytes.NewReader(nil)
	} else {
		reader = bytes.NewReader([]byte(body))
	}
	req := httptest.NewRequest(method, path, reader)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != wantStatus {
		t.Fatalf("%s %s expected status %d, got %d: %s", method, path, wantStatus, res.Code, res.Body.String())
	}

	var payload map[string]any
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response for %s %s: %v", method, path, err)
	}
	return payload
}

func stringField(t *testing.T, payload map[string]any, field string) string {
	t.Helper()
	value, ok := payload[field].(string)
	if !ok || value == "" {
		t.Fatalf("expected string field %q in %#v", field, payload)
	}
	return value
}
