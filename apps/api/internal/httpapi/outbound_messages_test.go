package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOutboundReplyQueueAndDelivery(t *testing.T) {
	database := openTestDB(t)
	router := NewRouter(database)
	household := requestJSON(t, router, http.MethodPost, "/api/households", `{"displayName":"Casa Telegram"}`, http.StatusCreated)
	householdID := stringField(t, household, "id")
	contact := requestJSON(t, router, http.MethodPost, "/api/contacts", `{"displayName":"Cliente","telegramId":"reply-sender"}`, http.StatusCreated)
	contactID := stringField(t, contact, "id")
	requestJSON(t, router, http.MethodPost, "/api/households/"+householdID+"/contacts", `{"contactId":"`+contactID+`","role":"owner","isPrimary":true}`, http.StatusCreated)
	message := requestJSON(t, router, http.MethodPost, "/api/message-imports", `{
		"channel":"telegram","externalConversationId":"reply-chat",
		"externalMessageId":"reply-message","senderExternalId":"reply-sender",
		"body":"Quiero reservar"
	}`, http.StatusCreated)
	detectedID := stringField(t, message["detectedRequest"].(map[string]any), "id")
	reply := requestJSON(t, router, http.MethodPost, "/api/detected-requests/"+detectedID+"/replies", `{"templateKey":"request_details"}`, http.StatusCreated)
	if reply["status"] != "pending" || reply["templateKey"] != "request_details" {
		t.Fatalf("expected pending details reply, got %#v", reply)
	}
	requestJSON(t, router, http.MethodPost, "/api/detected-requests/"+detectedID+"/replies", `{"templateKey":"request_declined"}`, http.StatusConflict)
	detected := requestJSON(t, router, http.MethodGet, "/api/detected-requests/"+detectedID, "", http.StatusOK)
	if detected["status"] != "needs_more_info" {
		t.Fatalf("expected request to need more info, got %#v", detected)
	}
	replyID := stringField(t, reply, "id")
	delivered := requestJSON(t, router, http.MethodPatch, "/api/automation/outbound-messages/"+replyID, `{"status":"sent"}`, http.StatusOK)
	if delivered["status"] != "sent" || delivered["sentAt"] == nil {
		t.Fatalf("expected sent reply, got %#v", delivered)
	}
	declined := requestJSON(t, router, http.MethodPost, "/api/detected-requests/"+detectedID+"/replies", `{"templateKey":"request_declined"}`, http.StatusCreated)
	if declined["status"] != "pending" {
		t.Fatalf("expected declined reply to be queued, got %#v", declined)
	}
	detected = requestJSON(t, router, http.MethodGet, "/api/detected-requests/"+detectedID, "", http.StatusOK)
	if detected["status"] != "ignored" {
		t.Fatalf("expected request to be ignored after approved decline, got %#v", detected)
	}
}

func TestOutboundAutomationEndpointsRequireToken(t *testing.T) {
	database := openTestDB(t)
	router := NewRouterWithAutomationToken(database, "outbound-secret")

	unauthorized := httptest.NewRecorder()
	router.ServeHTTP(unauthorized, httptest.NewRequest(http.MethodGet, "/api/automation/outbound-messages?status=pending", nil))
	if unauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized automation queue, got %d", unauthorized.Code)
	}
	authorizedRequest := httptest.NewRequest(http.MethodGet, "/api/automation/outbound-messages?status=pending", nil)
	authorizedRequest.Header.Set("Authorization", "Bearer outbound-secret")
	authorized := httptest.NewRecorder()
	router.ServeHTTP(authorized, authorizedRequest)
	if authorized.Code != http.StatusOK {
		t.Fatalf("expected authorized automation queue, got %d: %s", authorized.Code, authorized.Body.String())
	}
}
