package httpapi

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	dbqueries "github.com/pawsear/pawsear-platform/apps/api/internal/db/queries"
)

type detectedRequestContactLinkInput struct {
	ContactID   string `json:"contactId"`
	DisplayName string `json:"displayName"`
	HouseholdID string `json:"householdId"`
	Role        string `json:"role"`
}

func (h *mvpHandler) linkDetectedRequestContact(w http.ResponseWriter, r *http.Request) {
	var input detectedRequestContactLinkInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	defer tx.Rollback()
	queries := dbqueries.New(tx)
	detected, err := queries.GetDetectedRequestDetail(r.Context(), r.PathValue("id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if !allowed(detected.Channel, "telegram", "whatsapp") || !detected.SenderExternalID.Valid {
		writeInvalid(w, fmt.Errorf("request does not contain a linkable sender identity"))
		return
	}
	contact, err := h.resolveContactForLink(r.Context(), queries, input)
	if err != nil {
		writeInvalid(w, err)
		return
	}
	if conflict, err := contactIdentityConflict(r, queries, detected.Channel, detected.SenderExternalID.String, contact); err != nil {
		writeStoreError(w, err)
		return
	} else if conflict != "" {
		writeConflict(w, fmt.Errorf("%s", conflict))
		return
	}

	now := timestamp(h.now)
	if err := linkContactIdentity(r, queries, detected.Channel, detected.SenderExternalID.String, contact.ID, now); err != nil {
		writeStoreError(w, err)
		return
	}
	householdID := ""
	households, err := queries.ListHouseholdsForContactResolution(r.Context(), contact.ID)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if len(households) == 1 {
		householdID = households[0].ID
	} else if len(households) > 1 {
		householdID, err = householdFromUniquePetMention(r.Context(), queries, contact.ID, detected.MessageBody)
		if err != nil {
			writeStoreError(w, err)
			return
		}
	}
	if err := queries.LinkMessageSenderContact(r.Context(), dbqueries.LinkMessageSenderContactParams{
		SenderContactID: nullableID(contact.ID), ID: detected.MessageID,
	}); err != nil {
		writeStoreError(w, err)
		return
	}
	if err := queries.LinkConversationContactContext(r.Context(), dbqueries.LinkConversationContactContextParams{
		PrimaryContactID: nullableID(contact.ID), HouseholdID: nullableID(householdID), UpdatedAt: now, ID: detected.MessageID,
	}); err != nil {
		writeStoreError(w, err)
		return
	}
	if err := queries.LinkDetectedRequestContext(r.Context(), dbqueries.LinkDetectedRequestContextParams{
		ContactID: nullableID(contact.ID), HouseholdID: nullableID(householdID), UpdatedAt: now, ID: detected.ID,
	}); err != nil {
		writeStoreError(w, err)
		return
	}
	if err := tx.Commit(); err != nil {
		writeStoreError(w, err)
		return
	}
	updated, err := h.queries.GetDetectedRequestDetail(r.Context(), detected.ID)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, detectedRequestDetailResponse(updated))
}

func (h *mvpHandler) resolveContactForLink(ctx context.Context, queries *dbqueries.Queries, input detectedRequestContactLinkInput) (dbqueries.Contact, error) {
	if contactID := strings.TrimSpace(input.ContactID); contactID != "" {
		return queries.GetContact(ctx, contactID)
	}
	displayName, err := requiredText(input.DisplayName, "displayName")
	if err != nil {
		return dbqueries.Contact{}, err
	}
	householdID, err := requiredText(input.HouseholdID, "householdId")
	if err != nil {
		return dbqueries.Contact{}, err
	}
	if _, err := queries.GetHousehold(ctx, householdID); err != nil {
		return dbqueries.Contact{}, fmt.Errorf("household is not available: %w", err)
	}
	params, err := h.contactCreateParams(contactInput{DisplayName: displayName})
	if err != nil {
		return dbqueries.Contact{}, err
	}
	contact, err := queries.CreateContact(ctx, params)
	if err != nil {
		return dbqueries.Contact{}, err
	}
	role, err := validateHouseholdContact(householdContactInput{ContactID: contact.ID, Role: input.Role})
	if err != nil {
		return dbqueries.Contact{}, err
	}
	if err := queries.LinkHouseholdContact(ctx, dbqueries.LinkHouseholdContactParams{
		HouseholdID: householdID, ContactID: contact.ID, Role: role, IsPrimary: 1, CreatedAt: timestamp(h.now),
	}); err != nil {
		return dbqueries.Contact{}, err
	}
	return contact, nil
}

func contactIdentityConflict(r *http.Request, queries *dbqueries.Queries, channel string, externalID string, contact dbqueries.Contact) (string, error) {
	if channel == "telegram" && contact.TelegramID.Valid && contact.TelegramID.String != externalID {
		return "contact is already linked to another Telegram account", nil
	}
	if channel == "whatsapp" && contact.WhatsappID.Valid && contact.WhatsappID.String != externalID {
		return "contact is already linked to another WhatsApp account", nil
	}
	existing, err := queries.GetContactByChannelIdentity(r.Context(), dbqueries.GetContactByChannelIdentityParams{
		Channel: channel, ExternalUserID: externalID,
	})
	if err == sql.ErrNoRows {
		existing, err = queries.GetContactByChannelID(r.Context(), dbqueries.GetContactByChannelIDParams{
			Column1: channel, TelegramID: nullableString(externalID), Column3: channel, WhatsappID: nullableString(externalID),
		})
		if err == sql.ErrNoRows {
			return "", nil
		}
	}
	if err != nil {
		return "", err
	}
	if existing.ID != contact.ID {
		return "sender identity is already linked to another contact", nil
	}
	return "", nil
}

func linkContactIdentity(r *http.Request, queries *dbqueries.Queries, channel string, externalID string, contactID string, now string) error {
	identityID, err := newRecordID()
	if err != nil {
		return err
	}
	if _, err := queries.CreateContactChannelIdentity(r.Context(), dbqueries.CreateContactChannelIdentityParams{
		ID: identityID, ContactID: contactID, Channel: channel, ExternalUserID: externalID,
		CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		return err
	}
	if channel == "telegram" {
		_, err = queries.LinkContactTelegramIdentity(r.Context(), dbqueries.LinkContactTelegramIdentityParams{
			TelegramID: nullableString(externalID), UpdatedAt: now, ID: contactID,
		})
		return err
	}
	_, err = queries.LinkContactWhatsappIdentity(r.Context(), dbqueries.LinkContactWhatsappIdentityParams{
		WhatsappID: nullableString(externalID), UpdatedAt: now, ID: contactID,
	})
	return err
}
