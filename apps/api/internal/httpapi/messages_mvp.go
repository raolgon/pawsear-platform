package httpapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"

	dbqueries "github.com/pawsear/pawsear-platform/apps/api/internal/db/queries"
)

const maxImportedMessageLength = 20_000

type messageImportInput struct {
	Channel                string                   `json:"channel"`
	ExternalConversationID string                   `json:"externalConversationId"`
	ExternalMessageID      string                   `json:"externalMessageId"`
	SenderExternalID       string                   `json:"senderExternalId"`
	Direction              string                   `json:"direction"`
	Body                   string                   `json:"body"`
	SentAt                 *string                  `json:"sentAt"`
	Suggestion             *detectedSuggestionInput `json:"suggestion"`
}

type detectedSuggestionInput struct {
	HouseholdID string  `json:"householdId"`
	ContactID   string  `json:"contactId"`
	ServiceType string  `json:"serviceType"`
	StartAt     *string `json:"startAt"`
	EndAt       *string `json:"endAt"`
	Confidence  string  `json:"confidence"`
}

type detectedRequestUpdateInput struct {
	Status      string  `json:"status"`
	ReviewNotes *string `json:"reviewNotes"`
}

func (h *mvpHandler) importMessage(w http.ResponseWriter, r *http.Request) {
	var input messageImportInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	if err := validateMessageImport(&input); err != nil {
		writeInvalid(w, err)
		return
	}

	if existing, found, err := h.existingMessageImport(r.Context(), input); err != nil {
		writeStoreError(w, err)
		return
	} else if found {
		writeJSON(w, http.StatusOK, existing)
		return
	}

	response, err := h.createMessageImport(r.Context(), input)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, response)
}

func (h *mvpHandler) createMessageImport(ctx context.Context, input messageImportInput) (map[string]any, error) {
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	queries := dbqueries.New(tx)
	input.Suggestion = mergeMessageSuggestion(input.Suggestion, inferredMessageSuggestion(input.Body, h.now()))

	contactID, householdID, err := matchMessageContext(ctx, queries, input)
	if err != nil {
		return nil, err
	}
	contactID, householdID = applySuggestionContext(input.Suggestion, contactID, householdID)

	conversation, err := findOrCreateConversation(ctx, queries, input, contactID, householdID, h.now)
	if err != nil {
		return nil, err
	}
	now := timestamp(h.now)
	messageID, err := newRecordID()
	if err != nil {
		return nil, err
	}
	message, err := queries.CreateMessage(ctx, dbqueries.CreateMessageParams{
		ID: messageID, ConversationID: conversation.ID, SenderContactID: nullableID(contactID),
		Direction: input.Direction, Body: input.Body, SentAt: optionalText(input.SentAt), ImportedAt: now,
		ExternalMessageID: nullableString(input.ExternalMessageID), SenderExternalID: nullableString(input.SenderExternalID),
	})
	if err != nil {
		return nil, err
	}
	detectedID, err := newRecordID()
	if err != nil {
		return nil, err
	}
	detected, err := queries.CreateDetectedRequest(ctx, detectedRequestParams(detectedID, message.ID, input, contactID, householdID, now))
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return messageImportResponse(conversation, message, detected, false), nil
}

func (h *mvpHandler) existingMessageImport(ctx context.Context, input messageImportInput) (map[string]any, bool, error) {
	if input.ExternalMessageID == "" {
		return nil, false, nil
	}
	message, err := h.queries.GetImportedMessageByExternalID(ctx, dbqueries.GetImportedMessageByExternalIDParams{
		Channel: input.Channel, ExternalConversationID: nullableString(input.ExternalConversationID),
		ExternalMessageID: sql.NullString{String: input.ExternalMessageID, Valid: true},
	})
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	conversation, err := h.queries.GetConversationByExternalID(ctx, dbqueries.GetConversationByExternalIDParams{
		Channel: input.Channel, ExternalConversationID: nullableString(input.ExternalConversationID),
	})
	if err != nil {
		return nil, false, err
	}
	detected, err := h.queries.GetDetectedRequestByMessage(ctx, message.ID)
	if err != nil {
		return nil, false, err
	}
	return messageImportResponse(conversation, message, detected, true), true, nil
}

func (h *mvpHandler) listDetectedRequests(w http.ResponseWriter, r *http.Request) {
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	if status != "" && !allowed(status, "needs_review", "confirmed", "ignored", "needs_more_info", "converted_to_booking") {
		writeInvalid(w, fmt.Errorf("status is not supported"))
		return
	}
	rows, err := h.queries.ListDetectedRequests(r.Context(), dbqueries.ListDetectedRequestsParams{Column1: status, Status: status})
	if err != nil {
		writeStoreError(w, err)
		return
	}
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, detectedRequestListResponse(row))
	}
	writeJSON(w, http.StatusOK, map[string]any{"detectedRequests": items})
}

func (h *mvpHandler) getDetectedRequest(w http.ResponseWriter, r *http.Request) {
	row, err := h.queries.GetDetectedRequestDetail(r.Context(), r.PathValue("id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, detectedRequestDetailResponse(row))
}

func (h *mvpHandler) updateDetectedRequest(w http.ResponseWriter, r *http.Request) {
	var input detectedRequestUpdateInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	if !allowed(input.Status, "needs_review", "confirmed", "ignored", "needs_more_info") {
		writeInvalid(w, fmt.Errorf("status is not supported"))
		return
	}
	current, err := h.queries.GetDetectedRequestDetail(r.Context(), r.PathValue("id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if !detectedRequestTransitionAllowed(current.Status, input.Status) {
		writeConflict(w, fmt.Errorf("detected request status transition is not allowed"))
		return
	}
	updated, err := h.queries.UpdateDetectedRequestStatus(r.Context(), dbqueries.UpdateDetectedRequestStatusParams{
		Status: input.Status, ReviewNotes: optionalText(input.ReviewNotes), UpdatedAt: timestamp(h.now), ID: r.PathValue("id"),
	})
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, detectedRequestResponse(updated))
}

func detectedRequestTransitionAllowed(current string, next string) bool {
	if current == next {
		return true
	}
	switch current {
	case "needs_review", "needs_more_info", "confirmed":
		return allowed(next, "needs_review", "needs_more_info", "confirmed", "ignored")
	case "ignored":
		return next == "needs_review"
	default:
		return false
	}
}

func validateMessageImport(input *messageImportInput) error {
	input.Channel = strings.ToLower(strings.TrimSpace(input.Channel))
	if !allowed(input.Channel, "telegram", "whatsapp", "manual") {
		return fmt.Errorf("channel is not supported")
	}
	input.ExternalConversationID = strings.TrimSpace(input.ExternalConversationID)
	input.ExternalMessageID = strings.TrimSpace(input.ExternalMessageID)
	input.SenderExternalID = strings.TrimSpace(input.SenderExternalID)
	input.Direction = defaultText(input.Direction, "inbound")
	if !allowed(input.Direction, "inbound", "outbound", "system") {
		return fmt.Errorf("direction is not supported")
	}
	input.Body = strings.TrimSpace(input.Body)
	if input.Body == "" {
		return fmt.Errorf("body is required")
	}
	if len(input.Body) > maxImportedMessageLength {
		return fmt.Errorf("body must not exceed %d bytes", maxImportedMessageLength)
	}
	if input.Channel != "manual" && (input.ExternalConversationID == "" || input.ExternalMessageID == "") {
		return fmt.Errorf("externalConversationId and externalMessageId are required for channel imports")
	}
	if err := validateOptionalTimestamp(input.SentAt, "sentAt"); err != nil {
		return err
	}
	if input.Suggestion != nil {
		if err := validateSuggestion(input.Suggestion); err != nil {
			return err
		}
	}
	return nil
}

func validateSuggestion(input *detectedSuggestionInput) error {
	input.Confidence = defaultText(input.Confidence, "unknown")
	if !allowed(input.Confidence, "unknown", "low", "medium", "high") {
		return fmt.Errorf("suggestion.confidence is not supported")
	}
	if input.ServiceType != "" && !allowed(input.ServiceType, "walk", "client_home_sitting", "boarding", "visit", "transport", "other") {
		return fmt.Errorf("suggestion.serviceType is not supported")
	}
	if err := validateOptionalTimestamp(input.StartAt, "suggestion.startAt"); err != nil {
		return err
	}
	return validateOptionalTimestamp(input.EndAt, "suggestion.endAt")
}

func validateOptionalTimestamp(value *string, field string) error {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil
	}
	if _, err := time.Parse(time.RFC3339, strings.TrimSpace(*value)); err != nil {
		return fmt.Errorf("%s must use RFC3339", field)
	}
	return nil
}

func matchMessageContext(ctx context.Context, queries *dbqueries.Queries, input messageImportInput) (string, string, error) {
	if input.SenderExternalID == "" || input.Channel == "manual" {
		return "", "", nil
	}
	contact, err := contactForMessageIdentity(ctx, queries, input)
	if err == sql.ErrNoRows {
		return "", "", nil
	}
	if err != nil {
		return "", "", err
	}
	if householdID, found, err := conversationHouseholdForContact(ctx, queries, input, contact.ID); err != nil {
		return "", "", err
	} else if found {
		return contact.ID, householdID, nil
	}
	households, err := queries.ListHouseholdsForContactResolution(ctx, contact.ID)
	if err != nil {
		return "", "", err
	}
	if len(households) == 0 {
		return contact.ID, "", nil
	}
	if len(households) == 1 {
		return contact.ID, households[0].ID, nil
	}
	householdID, err := householdFromUniquePetMention(ctx, queries, contact.ID, input.Body)
	if err != nil {
		return "", "", err
	}
	return contact.ID, householdID, nil
}

func contactForMessageIdentity(ctx context.Context, queries *dbqueries.Queries, input messageImportInput) (dbqueries.Contact, error) {
	contact, err := queries.GetContactByChannelIdentity(ctx, dbqueries.GetContactByChannelIdentityParams{
		Channel: input.Channel, ExternalUserID: input.SenderExternalID,
	})
	if err != sql.ErrNoRows {
		return contact, err
	}
	return queries.GetContactByChannelID(ctx, dbqueries.GetContactByChannelIDParams{
		Column1: input.Channel, TelegramID: nullableString(input.SenderExternalID),
		Column3: input.Channel, WhatsappID: nullableString(input.SenderExternalID),
	})
}

func conversationHouseholdForContact(ctx context.Context, queries *dbqueries.Queries, input messageImportInput, contactID string) (string, bool, error) {
	conversation, err := queries.GetConversationByExternalID(ctx, dbqueries.GetConversationByExternalIDParams{
		Channel: input.Channel, ExternalConversationID: nullableString(input.ExternalConversationID),
	})
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	if !conversation.HouseholdID.Valid {
		return "", false, nil
	}
	belongs, err := queries.ContactBelongsToHousehold(ctx, dbqueries.ContactBelongsToHouseholdParams{
		ContactID: contactID, HouseholdID: conversation.HouseholdID.String,
	})
	return conversation.HouseholdID.String, belongs != 0, err
}

func householdFromUniquePetMention(ctx context.Context, queries *dbqueries.Queries, contactID string, body string) (string, error) {
	pets, err := queries.ListContactHouseholdPets(ctx, contactID)
	if err != nil {
		return "", err
	}
	messageWords := normalizedMessageWords(body)
	matches := make(map[string]bool)
	for _, pet := range pets {
		if wordsContainName(messageWords, normalizedMessageWords(pet.PetName)) {
			matches[pet.HouseholdID] = true
		}
	}
	if len(matches) != 1 {
		return "", nil
	}
	for householdID := range matches {
		return householdID, nil
	}
	return "", nil
}

func normalizedMessageWords(value string) map[string]bool {
	words := strings.FieldsFunc(normalizeMessageText(value), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	result := make(map[string]bool, len(words))
	for _, word := range words {
		result[word] = true
	}
	return result
}

func wordsContainName(messageWords map[string]bool, nameWords map[string]bool) bool {
	if len(nameWords) == 0 {
		return false
	}
	for word := range nameWords {
		if !messageWords[word] {
			return false
		}
	}
	return true
}

func findOrCreateConversation(ctx context.Context, queries *dbqueries.Queries, input messageImportInput, contactID string, householdID string, now func() time.Time) (dbqueries.Conversation, error) {
	externalID := input.ExternalConversationID
	if input.Channel == "manual" && externalID == "" {
		generated, err := newRecordID()
		if err != nil {
			return dbqueries.Conversation{}, err
		}
		externalID = "manual-" + generated
	}
	params := dbqueries.GetConversationByExternalIDParams{Channel: input.Channel, ExternalConversationID: nullableString(externalID)}
	conversation, err := queries.GetConversationByExternalID(ctx, params)
	if err == nil {
		lastMessageAt := timestamp(now)
		if input.SentAt != nil {
			lastMessageAt = strings.TrimSpace(*input.SentAt)
		}
		return queries.UpdateConversationFromImport(ctx, dbqueries.UpdateConversationFromImportParams{
			PrimaryContactID: nullableID(contactID), HouseholdID: nullableID(householdID), LastMessageAt: nullableString(lastMessageAt), UpdatedAt: timestamp(now), ID: conversation.ID,
		})
	}
	if err != sql.ErrNoRows {
		return dbqueries.Conversation{}, err
	}
	id, err := newRecordID()
	if err != nil {
		return dbqueries.Conversation{}, err
	}
	createdAt := timestamp(now)
	return queries.CreateConversation(ctx, dbqueries.CreateConversationParams{
		ID: id, Channel: input.Channel, ExternalConversationID: nullableString(externalID), PrimaryContactID: nullableID(contactID),
		HouseholdID: nullableID(householdID), LastMessageAt: optionalText(input.SentAt), CreatedAt: createdAt, UpdatedAt: createdAt,
	})
}

func detectedRequestParams(id string, messageID string, input messageImportInput, contactID string, householdID string, now string) dbqueries.CreateDetectedRequestParams {
	suggestion := input.Suggestion
	params := dbqueries.CreateDetectedRequestParams{
		ID: id, MessageID: messageID, HouseholdID: nullableID(householdID), ContactID: nullableID(contactID), Confidence: "unknown",
		RawPayloadJson: nullableJSON(input), CreatedAt: now, UpdatedAt: now,
	}
	if suggestion != nil {
		params.DetectedServiceType = nullableString(suggestion.ServiceType)
		params.DetectedStartAt = optionalText(suggestion.StartAt)
		params.DetectedEndAt = optionalText(suggestion.EndAt)
		params.Confidence = suggestion.Confidence
	}
	return params
}

func applySuggestionContext(suggestion *detectedSuggestionInput, contactID string, householdID string) (string, string) {
	if suggestion == nil {
		return contactID, householdID
	}
	if suggestion.ContactID != "" {
		contactID = strings.TrimSpace(suggestion.ContactID)
	}
	if suggestion.HouseholdID != "" {
		householdID = strings.TrimSpace(suggestion.HouseholdID)
	}
	return contactID, householdID
}

func nullableID(value string) sql.NullString { return nullableString(strings.TrimSpace(value)) }

func nullableString(value string) sql.NullString {
	return sql.NullString{String: value, Valid: value != ""}
}

func nullableJSON(value any) sql.NullString {
	payload, err := json.Marshal(value)
	if err != nil {
		return sql.NullString{}
	}
	return nullableString(string(payload))
}

func messageImportResponse(conversation dbqueries.Conversation, message dbqueries.Message, detected dbqueries.DetectedRequest, duplicate bool) map[string]any {
	return map[string]any{
		"duplicate":       duplicate,
		"conversation":    map[string]any{"id": conversation.ID, "channel": conversation.Channel, "externalConversationId": textValue(conversation.ExternalConversationID)},
		"message":         map[string]any{"id": message.ID, "body": message.Body, "externalMessageId": textValue(message.ExternalMessageID), "sentAt": textValue(message.SentAt)},
		"detectedRequest": detectedRequestResponse(detected),
	}
}

func detectedRequestResponse(row dbqueries.DetectedRequest) map[string]any {
	return map[string]any{
		"id": row.ID, "messageId": row.MessageID, "householdId": textValue(row.HouseholdID), "contactId": textValue(row.ContactID),
		"serviceType": textValue(row.DetectedServiceType), "startAt": textValue(row.DetectedStartAt), "endAt": textValue(row.DetectedEndAt),
		"confidence": row.Confidence, "status": row.Status, "convertedBookingId": textValue(row.ConvertedBookingID),
		"reviewNotes": textValue(row.ReviewNotes), "createdAt": row.CreatedAt, "updatedAt": row.UpdatedAt,
	}
}

func detectedRequestListResponse(row dbqueries.ListDetectedRequestsRow) map[string]any {
	return map[string]any{
		"id": row.ID, "messageId": row.MessageID, "householdId": textValue(row.HouseholdID), "householdName": row.HouseholdName,
		"contactId": textValue(row.ContactID), "contactName": row.ContactName, "channel": row.Channel, "body": row.MessageBody,
		"sentAt": textValue(row.MessageSentAt), "externalMessageId": textValue(row.ExternalMessageID),
		"senderExternalId":       textValue(row.SenderExternalID),
		"externalConversationId": textValue(row.ExternalConversationID), "serviceType": textValue(row.DetectedServiceType),
		"convertedBookingStartAt": textValue(row.ConvertedBookingStartAt), "convertedBookingHouseholdId": textValue(row.ConvertedBookingHouseholdID),
		"startAt": textValue(row.DetectedStartAt), "endAt": textValue(row.DetectedEndAt), "confidence": row.Confidence,
		"status": row.Status, "convertedBookingId": textValue(row.ConvertedBookingID), "reviewNotes": textValue(row.ReviewNotes),
		"createdAt": row.CreatedAt, "updatedAt": row.UpdatedAt,
	}
}

func detectedRequestDetailResponse(row dbqueries.GetDetectedRequestDetailRow) map[string]any {
	return map[string]any{
		"id": row.ID, "messageId": row.MessageID, "householdId": textValue(row.HouseholdID), "householdName": row.HouseholdName,
		"contactId": textValue(row.ContactID), "contactName": row.ContactName, "channel": row.Channel, "body": row.MessageBody,
		"sentAt": textValue(row.MessageSentAt), "externalMessageId": textValue(row.ExternalMessageID),
		"senderExternalId":       textValue(row.SenderExternalID),
		"externalConversationId": textValue(row.ExternalConversationID), "serviceType": textValue(row.DetectedServiceType),
		"convertedBookingStartAt": textValue(row.ConvertedBookingStartAt), "convertedBookingHouseholdId": textValue(row.ConvertedBookingHouseholdID),
		"startAt": textValue(row.DetectedStartAt), "endAt": textValue(row.DetectedEndAt), "confidence": row.Confidence,
		"status": row.Status, "convertedBookingId": textValue(row.ConvertedBookingID), "reviewNotes": textValue(row.ReviewNotes),
		"createdAt": row.CreatedAt, "updatedAt": row.UpdatedAt,
	}
}
