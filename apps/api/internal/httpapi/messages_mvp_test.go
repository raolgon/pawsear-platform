package httpapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMessageImportIsIdempotentAndMatchesKnownContact(t *testing.T) {
	database := openTestDB(t)
	router := NewRouter(database)
	household := requestJSON(t, router, http.MethodPost, "/api/households", `{"displayName":"Casa Luna"}`, http.StatusCreated)
	householdID := stringField(t, household, "id")
	contact := requestJSON(t, router, http.MethodPost, "/api/contacts", `{"displayName":"Elena","telegramId":"telegram-42"}`, http.StatusCreated)
	contactID := stringField(t, contact, "id")
	requestJSON(t, router, http.MethodPost, "/api/households/"+householdID+"/contacts", `{"contactId":"`+contactID+`","role":"owner","isPrimary":true}`, http.StatusCreated)

	body := `{
		"channel":"telegram",
		"externalConversationId":"chat-42",
		"externalMessageId":"message-100",
		"senderExternalId":"telegram-42",
		"body":"¿Puedes pasear a Luna mañana a las 8?",
		"sentAt":"2026-06-20T14:00:00Z",
		"suggestion":{"serviceType":"walk","confidence":"medium"}
	}`
	created := requestJSON(t, router, http.MethodPost, "/api/message-imports", body, http.StatusCreated)
	if created["duplicate"] != false {
		t.Fatalf("expected a new import, got %#v", created)
	}
	detected := created["detectedRequest"].(map[string]any)
	if detected["householdId"] != householdID || detected["contactId"] != contactID {
		t.Fatalf("expected matched household and contact, got %#v", detected)
	}

	duplicate := requestJSON(t, router, http.MethodPost, "/api/message-imports", body, http.StatusOK)
	if duplicate["duplicate"] != true {
		t.Fatalf("expected duplicate import, got %#v", duplicate)
	}
	if duplicate["message"].(map[string]any)["id"] != created["message"].(map[string]any)["id"] {
		t.Fatal("expected duplicate import to return the original message")
	}
	otherConversation := `{
		"channel":"telegram",
		"externalConversationId":"chat-99",
		"externalMessageId":"message-100",
		"body":"El mismo ID puede existir en otro chat"
	}`
	other := requestJSON(t, router, http.MethodPost, "/api/message-imports", otherConversation, http.StatusCreated)
	if other["message"].(map[string]any)["id"] == created["message"].(map[string]any)["id"] {
		t.Fatal("expected the same external message ID in another conversation to create a new message")
	}

	queue := requestJSON(t, router, http.MethodGet, "/api/detected-requests?status=needs_review", "", http.StatusOK)
	requests := queue["detectedRequests"].([]any)
	if len(requests) != 2 {
		t.Fatalf("expected two review items, got %d", len(requests))
	}
	detectedID := detected["id"].(string)
	updated := requestJSON(t, router, http.MethodPatch, "/api/detected-requests/"+detectedID, `{"status":"needs_more_info","reviewNotes":"Confirm the date"}`, http.StatusOK)
	if updated["status"] != "needs_more_info" {
		t.Fatalf("expected needs_more_info, got %#v", updated)
	}
	detail := requestJSON(t, router, http.MethodGet, "/api/detected-requests/"+detectedID, "", http.StatusOK)
	if detail["body"] != "¿Puedes pasear a Luna mañana a las 8?" || detail["channel"] != "telegram" {
		t.Fatalf("expected source message detail, got %#v", detail)
	}
	pet := requestJSON(t, router, http.MethodPost, "/api/pets", `{"householdId":"`+householdID+`","name":"Luna","species":"dog"}`, http.StatusCreated)
	petID := stringField(t, pet, "id")
	requestJSON(t, router, http.MethodPost, "/api/detected-requests/"+detectedID+"/bookings", `{"status":"completed","startAt":"2026-06-22T14:00:00Z"}`, http.StatusBadRequest)
	requestJSON(t, router, http.MethodPost, "/api/detected-requests/"+detectedID+"/bookings", `{"startAt":"tomorrow"}`, http.StatusBadRequest)
	requestJSON(t, router, http.MethodPost, "/api/detected-requests/"+detectedID+"/bookings", `{"startAt":"2026-06-22T14:00:00Z","endAt":"2026-06-22T13:00:00Z"}`, http.StatusBadRequest)
	conversion := requestJSON(t, router, http.MethodPost, "/api/detected-requests/"+detectedID+"/bookings", `{
		"status":"confirmed",
		"startAt":"2026-06-22T14:00:00Z",
		"petIds":["`+petID+`"],
		"reviewNotes":"Confirmed with Elena"
	}`, http.StatusCreated)
	booking := conversion["booking"].(map[string]any)
	if booking["householdId"] != householdID || booking["source"] != "telegram" || booking["status"] != "confirmed" {
		t.Fatalf("expected detected values in converted booking, got %#v", booking)
	}
	bookingID := stringField(t, booking, "id")
	convertedDetail := requestJSON(t, router, http.MethodGet, "/api/detected-requests/"+detectedID, "", http.StatusOK)
	if convertedDetail["status"] != "converted_to_booking" || convertedDetail["convertedBookingId"] != bookingID {
		t.Fatalf("expected terminal converted request, got %#v", convertedDetail)
	}
	if convertedDetail["convertedBookingStartAt"] != "2026-06-22T14:00:00Z" || convertedDetail["convertedBookingHouseholdId"] != householdID {
		t.Fatalf("expected converted booking navigation context, got %#v", convertedDetail)
	}
	requestJSON(t, router, http.MethodPatch, "/api/detected-requests/"+detectedID, `{"status":"needs_review"}`, http.StatusConflict)
	requestJSON(t, router, http.MethodPost, "/api/detected-requests/"+detectedID+"/bookings", `{"startAt":"2026-06-23T14:00:00Z"}`, http.StatusConflict)

	var sourceCount int
	if err := database.QueryRow(`SELECT COUNT(*) FROM booking_sources WHERE booking_id = ? AND message_id = ? AND detected_request_id = ?`, bookingID, detected["messageId"], detectedID).Scan(&sourceCount); err != nil {
		t.Fatalf("query booking source: %v", err)
	}
	if sourceCount != 1 {
		t.Fatalf("expected one booking source, got %d", sourceCount)
	}
}

func TestMessageImportRejectsMissingExternalIdentifiers(t *testing.T) {
	router := NewRouter(openTestDB(t))
	requestJSON(t, router, http.MethodPost, "/api/message-imports", `{"channel":"telegram","body":"Hola"}`, http.StatusBadRequest)
}

func TestLinkDetectedRequestSenderRecognizesFutureMessages(t *testing.T) {
	database := openTestDB(t)
	router := NewRouter(database)
	household := requestJSON(t, router, http.MethodPost, "/api/households", `{"displayName":"Casa Sol"}`, http.StatusCreated)
	householdID := stringField(t, household, "id")
	contact := requestJSON(t, router, http.MethodPost, "/api/contacts", `{"displayName":"Mariana"}`, http.StatusCreated)
	contactID := stringField(t, contact, "id")
	requestJSON(t, router, http.MethodPost, "/api/households/"+householdID+"/contacts", `{"contactId":"`+contactID+`","role":"owner","isPrimary":true}`, http.StatusCreated)

	first := requestJSON(t, router, http.MethodPost, "/api/message-imports", `{
		"channel":"telegram",
		"externalConversationId":"chat-new",
		"externalMessageId":"message-new-1",
		"senderExternalId":"sender-new",
		"body":"Necesito un paseo"
	}`, http.StatusCreated)
	detectedID := stringField(t, first["detectedRequest"].(map[string]any), "id")
	linked := requestJSON(t, router, http.MethodPost, "/api/detected-requests/"+detectedID+"/contact-link", `{"contactId":"`+contactID+`"}`, http.StatusOK)
	if linked["contactId"] != contactID || linked["householdId"] != householdID {
		t.Fatalf("expected linked request context, got %#v", linked)
	}

	second := requestJSON(t, router, http.MethodPost, "/api/message-imports", `{
		"channel":"telegram",
		"externalConversationId":"chat-new",
		"externalMessageId":"message-new-2",
		"senderExternalId":"sender-new",
		"body":"También agrega una visita"
	}`, http.StatusCreated)
	detected := second["detectedRequest"].(map[string]any)
	if detected["contactId"] != contactID || detected["householdId"] != householdID {
		t.Fatalf("expected future message to match linked sender, got %#v", detected)
	}
}

func TestLinkDetectedRequestCanCreateSenderContact(t *testing.T) {
	database := openTestDB(t)
	router := NewRouter(database)
	household := requestJSON(t, router, http.MethodPost, "/api/households", `{"displayName":"Casa Tomy"}`, http.StatusCreated)
	householdID := stringField(t, household, "id")
	message := requestJSON(t, router, http.MethodPost, "/api/message-imports", `{
		"channel":"telegram",
		"externalConversationId":"chat-inline",
		"externalMessageId":"message-inline",
		"senderExternalId":"sender-inline",
		"body":"Paseo para Tomy mañana a las 10"
	}`, http.StatusCreated)
	detectedID := stringField(t, message["detectedRequest"].(map[string]any), "id")
	linked := requestJSON(t, router, http.MethodPost, "/api/detected-requests/"+detectedID+"/contact-link", `{
		"displayName":"Rafael",
		"householdId":"`+householdID+`",
		"role":"owner"
	}`, http.StatusOK)
	if linked["contactName"] != "Rafael" || linked["householdId"] != householdID {
		t.Fatalf("expected inline sender creation and household link, got %#v", linked)
	}
}

func TestMessageImportRequiresConfiguredAutomationToken(t *testing.T) {
	router := NewRouterWithAutomationToken(openTestDB(t), "local-secret")
	body := []byte(`{"channel":"manual","body":"Necesito un paseo"}`)

	unauthorized := httptest.NewRecorder()
	router.ServeHTTP(unauthorized, httptest.NewRequest(http.MethodPost, "/api/message-imports", bytes.NewReader(body)))
	if unauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized without token, got %d", unauthorized.Code)
	}

	authorizedRequest := httptest.NewRequest(http.MethodPost, "/api/message-imports", bytes.NewReader(body))
	authorizedRequest.Header.Set("Authorization", "Bearer local-secret")
	authorized := httptest.NewRecorder()
	router.ServeHTTP(authorized, authorizedRequest)
	if authorized.Code != http.StatusCreated {
		t.Fatalf("expected import with token, got %d: %s", authorized.Code, authorized.Body.String())
	}
}
