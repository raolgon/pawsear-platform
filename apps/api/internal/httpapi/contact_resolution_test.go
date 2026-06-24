package httpapi

import (
	"net/http"
	"testing"
)

func TestMessageResolutionSupportsMultipleSendersAndHouseholds(t *testing.T) {
	router := NewRouter(openTestDB(t))
	houseA := requestJSON(t, router, http.MethodPost, "/api/households", `{"displayName":"Casa A"}`, http.StatusCreated)
	houseB := requestJSON(t, router, http.MethodPost, "/api/households", `{"displayName":"Casa B"}`, http.StatusCreated)
	houseAID := stringField(t, houseA, "id")
	houseBID := stringField(t, houseB, "id")
	alice := createLinkedTelegramContact(t, router, "Alice", "alice-tg", houseAID)
	bob := createLinkedTelegramContact(t, router, "Bob", "bob-tg", houseAID)
	carol := createLinkedTelegramContact(t, router, "Carol", "carol-tg", houseBID)

	assertImportedContext(t, router, "family-a", "a-1", "alice-tg", "Paseo mañana", alice, houseAID)
	assertImportedContext(t, router, "family-a", "a-2", "bob-tg", "También confirmo", bob, houseAID)
	assertImportedContext(t, router, "carol-private", "b-1", "carol-tg", "Una visita", carol, houseBID)
}

func TestAmbiguousContactUsesPetThenRemembersConversationHousehold(t *testing.T) {
	router := NewRouter(openTestDB(t))
	houseA := requestJSON(t, router, http.MethodPost, "/api/households", `{"displayName":"Casa Luna"}`, http.StatusCreated)
	houseB := requestJSON(t, router, http.MethodPost, "/api/households", `{"displayName":"Casa Max"}`, http.StatusCreated)
	houseAID := stringField(t, houseA, "id")
	houseBID := stringField(t, houseB, "id")
	contact := requestJSON(t, router, http.MethodPost, "/api/contacts", `{"displayName":"Tutor compartido","telegramId":"multi-tg"}`, http.StatusCreated)
	contactID := stringField(t, contact, "id")
	requestJSON(t, router, http.MethodPost, "/api/households/"+houseAID+"/contacts", `{"contactId":"`+contactID+`","role":"owner"}`, http.StatusCreated)
	requestJSON(t, router, http.MethodPost, "/api/households/"+houseBID+"/contacts", `{"contactId":"`+contactID+`","role":"owner"}`, http.StatusCreated)
	requestJSON(t, router, http.MethodPost, "/api/pets", `{"householdId":"`+houseAID+`","name":"Luna","species":"dog"}`, http.StatusCreated)
	requestJSON(t, router, http.MethodPost, "/api/pets", `{"householdId":"`+houseBID+`","name":"Max","species":"dog"}`, http.StatusCreated)

	assertImportedContext(t, router, "multi-pet", "multi-1", "multi-tg", "Quiero un paseo para Max", contactID, houseBID)
	ambiguous := importTelegramRequest(t, router, "multi-ambiguous", "multi-2", "multi-tg", "Necesito un paseo")
	if ambiguous["contactId"] != contactID || ambiguous["householdId"] != nil {
		t.Fatalf("expected known contact with ambiguous household, got %#v", ambiguous)
	}
	detectedID := stringField(t, ambiguous, "id")
	linked := requestJSON(t, router, http.MethodPost, "/api/detected-requests/"+detectedID+"/household-link", `{"householdId":"`+houseAID+`"}`, http.StatusOK)
	if linked["householdId"] != houseAID {
		t.Fatalf("expected selected household, got %#v", linked)
	}
	assertImportedContext(t, router, "multi-ambiguous", "multi-3", "multi-tg", "La misma casa otra vez", contactID, houseAID)
}

func createLinkedTelegramContact(t *testing.T, router http.Handler, name string, telegramID string, householdID string) string {
	t.Helper()
	contact := requestJSON(t, router, http.MethodPost, "/api/contacts", `{"displayName":"`+name+`","telegramId":"`+telegramID+`"}`, http.StatusCreated)
	contactID := stringField(t, contact, "id")
	requestJSON(t, router, http.MethodPost, "/api/households/"+householdID+"/contacts", `{"contactId":"`+contactID+`","role":"owner"}`, http.StatusCreated)
	return contactID
}

func importTelegramRequest(t *testing.T, router http.Handler, conversationID string, messageID string, senderID string, body string) map[string]any {
	t.Helper()
	created := requestJSON(t, router, http.MethodPost, "/api/message-imports", `{
		"channel":"telegram","externalConversationId":"`+conversationID+`",
		"externalMessageId":"`+messageID+`","senderExternalId":"`+senderID+`",
		"body":"`+body+`"
	}`, http.StatusCreated)
	return created["detectedRequest"].(map[string]any)
}

func assertImportedContext(t *testing.T, router http.Handler, conversationID string, messageID string, senderID string, body string, contactID string, householdID string) {
	t.Helper()
	detected := importTelegramRequest(t, router, conversationID, messageID, senderID, body)
	if detected["contactId"] != contactID || detected["householdId"] != householdID {
		t.Fatalf("expected contact %s and household %s, got %#v", contactID, householdID, detected)
	}
}
