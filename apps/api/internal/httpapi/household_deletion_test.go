package httpapi

import (
	"net/http"
	"testing"
)

func TestDeleteHouseholdRemovesDependentDataAndPreservesPaymentLedger(t *testing.T) {
	database := openTestDB(t)
	router := NewRouter(database)
	target := requestJSON(t, router, http.MethodPost, "/api/households", `{"displayName":"Casa Borrar"}`, http.StatusCreated)
	targetID := stringField(t, target, "id")
	requestJSON(t, router, http.MethodPost, "/api/households", `{"displayName":"Casa Conservar"}`, http.StatusCreated)
	contact := requestJSON(t, router, http.MethodPost, "/api/contacts", `{"displayName":"Pagador","telegramId":"delete-sender"}`, http.StatusCreated)
	contactID := stringField(t, contact, "id")
	requestJSON(t, router, http.MethodPost, "/api/households/"+targetID+"/contacts", `{"contactId":"`+contactID+`","role":"owner","isPrimary":true}`, http.StatusCreated)
	orphan := requestJSON(t, router, http.MethodPost, "/api/contacts", `{"displayName":"Contacto huérfano"}`, http.StatusCreated)
	requestJSON(t, router, http.MethodPost, "/api/households/"+targetID+"/contacts", `{"contactId":"`+stringField(t, orphan, "id")+`","role":"family"}`, http.StatusCreated)
	pet := requestJSON(t, router, http.MethodPost, "/api/pets", `{"householdId":"`+targetID+`","name":"Tomy","species":"dog"}`, http.StatusCreated)
	petID := stringField(t, pet, "id")
	message := requestJSON(t, router, http.MethodPost, "/api/message-imports", `{
		"channel":"telegram","externalConversationId":"delete-chat",
		"externalMessageId":"delete-message","senderExternalId":"delete-sender",
		"body":"Paseo para Tomy mañana a las 10"
	}`, http.StatusCreated)
	detectedID := stringField(t, message["detectedRequest"].(map[string]any), "id")
	conversion := requestJSON(t, router, http.MethodPost, "/api/detected-requests/"+detectedID+"/bookings", `{
		"status":"confirmed","startAt":"2026-06-23T10:00:00-06:00","petIds":["`+petID+`"]
	}`, http.StatusCreated)
	bookingID := stringField(t, conversion["booking"].(map[string]any), "id")
	requestJSON(t, router, http.MethodPost, "/api/care-tasks", `{
		"bookingId":"`+bookingID+`","householdId":"`+targetID+`","petId":"`+petID+`",
		"taskType":"walk","title":"Paseo"
	}`, http.StatusCreated)
	charge := requestJSON(t, router, http.MethodPost, "/api/charges", `{
		"householdId":"`+targetID+`","bookingId":"`+bookingID+`",
		"description":"Paseo","amountMinor":10000
	}`, http.StatusCreated)
	requestJSON(t, router, http.MethodPost, "/api/payments", `{
		"payerContactId":"`+contactID+`","amountMinor":10000,"method":"cash",
		"allocations":[{"chargeId":"`+stringField(t, charge, "id")+`","amountMinor":10000}]
	}`, http.StatusCreated)

	requestJSON(t, router, http.MethodDelete, "/api/households/"+targetID, `{"confirmationName":"incorrecto"}`, http.StatusBadRequest)
	requestJSON(t, router, http.MethodDelete, "/api/households/"+targetID, `{"confirmationName":"Casa Borrar"}`, http.StatusOK)

	for _, table := range []string{"pets", "bookings", "care_tasks", "charges", "payment_allocations", "conversations", "messages", "detected_requests", "booking_sources", "outbound_messages"} {
		var count int
		if err := database.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
		if count != 0 {
			t.Fatalf("expected %s to be empty, got %d", table, count)
		}
	}
	var households, contacts, payments int
	if err := database.QueryRow("SELECT COUNT(*) FROM households").Scan(&households); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow("SELECT COUNT(*) FROM contacts").Scan(&contacts); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow("SELECT COUNT(*) FROM payments").Scan(&payments); err != nil {
		t.Fatal(err)
	}
	if households != 1 || contacts != 1 || payments != 1 {
		t.Fatalf("expected other household and payment ledger preserved; households=%d contacts=%d payments=%d", households, contacts, payments)
	}
}
