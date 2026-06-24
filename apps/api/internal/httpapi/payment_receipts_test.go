package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPaymentReceiptIssueAndDownload(t *testing.T) {
	router := NewRouter(openTestDB(t))
	household := requestJSON(t, router, http.MethodPost, "/api/households", `{"displayName":"Casa de Luna"}`, http.StatusCreated)
	householdID := stringField(t, household, "id")
	contact := requestJSON(t, router, http.MethodPost, "/api/contacts", `{"displayName":"Sofía García"}`, http.StatusCreated)
	contactID := stringField(t, contact, "id")
	charge := requestJSON(t, router, http.MethodPost, "/api/charges", `{
		"householdId":"`+householdID+`",
		"description":"Paseo de Luna",
		"amountMinor":85000
	}`, http.StatusCreated)
	chargeID := stringField(t, charge, "id")
	payment := requestJSON(t, router, http.MethodPost, "/api/payments", `{
		"payerContactId":"`+contactID+`",
		"receivedAt":"2026-06-24T22:35:00Z",
		"amountMinor":90000,
		"currency":"MXN",
		"method":"bank_transfer",
		"reference":"SPEI-483921",
		"allocations":[{"chargeId":"`+chargeID+`","amountMinor":85000}]
	}`, http.StatusCreated)
	paymentID := stringField(t, payment, "id")

	receipt := requestJSON(t, router, http.MethodPost, "/api/payments/"+paymentID+"/receipt", "", http.StatusCreated)
	receiptNumber := stringField(t, receipt, "receiptNumber")
	if !strings.HasPrefix(receiptNumber, "REC-20260624-") {
		t.Fatalf("unexpected receipt number %q", receiptNumber)
	}
	snapshot := receipt["snapshot"].(map[string]any)
	if snapshot["payerName"] != "Sofía García" || snapshot["allocatedMinor"] != float64(85000) || snapshot["unallocatedMinor"] != float64(5000) {
		t.Fatalf("unexpected receipt snapshot: %#v", snapshot)
	}

	existing := requestJSON(t, router, http.MethodPost, "/api/payments/"+paymentID+"/receipt", "", http.StatusOK)
	if existing["id"] != receipt["id"] {
		t.Fatalf("reissuing should return the existing receipt")
	}

	assertReceiptDownload(t, router, paymentID, "png", "image/png", []byte{0x89, 'P', 'N', 'G'})
	assertReceiptDownload(t, router, paymentID, "pdf", "application/pdf", []byte("%PDF"))
}

func TestPaymentReceiptSupportsUnknownPayerAndUnallocatedAmount(t *testing.T) {
	router := NewRouter(openTestDB(t))
	payment := requestJSON(t, router, http.MethodPost, "/api/payments", `{
		"receivedAt":"2026-06-24T22:35:00Z",
		"amountMinor":50000,
		"currency":"MXN",
		"method":"cash"
	}`, http.StatusCreated)
	receipt := requestJSON(t, router, http.MethodPost, "/api/payments/"+stringField(t, payment, "id")+"/receipt", "", http.StatusCreated)
	snapshot := receipt["snapshot"].(map[string]any)
	if snapshot["payerName"] != "Pagador no registrado" || snapshot["unallocatedMinor"] != float64(50000) {
		t.Fatalf("unexpected standalone receipt snapshot: %#v", snapshot)
	}
	if len(snapshot["allocations"].([]any)) != 0 {
		t.Fatalf("expected no allocations")
	}
}

func assertReceiptDownload(t *testing.T, router http.Handler, paymentID string, format string, contentType string, prefix []byte) {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, "/api/payments/"+paymentID+"/receipt/"+format+"?download=1", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("download %s returned %d: %s", format, response.Code, response.Body.String())
	}
	if response.Header().Get("Content-Type") != contentType {
		t.Fatalf("download %s content type was %q", format, response.Header().Get("Content-Type"))
	}
	if !strings.HasPrefix(string(response.Body.Bytes()), string(prefix)) {
		t.Fatalf("download %s has an invalid signature", format)
	}
	if !strings.Contains(response.Header().Get("Content-Disposition"), "attachment") {
		t.Fatalf("download %s should be an attachment", format)
	}
}
