package httpapi

import (
	"database/sql"
	"fmt"
	"net/http"

	dbqueries "github.com/pawsear/pawsear-platform/apps/api/internal/db/queries"
)

type contactInput struct {
	DisplayName string  `json:"displayName"`
	Phone       *string `json:"phone"`
	WhatsappID  *string `json:"whatsappId"`
	TelegramID  *string `json:"telegramId"`
	Email       *string `json:"email"`
	Notes       *string `json:"notes"`
	Active      *bool   `json:"active"`
}

type householdContactInput struct {
	ContactID string  `json:"contactId"`
	Role      string  `json:"role"`
	IsPrimary bool    `json:"isPrimary"`
	Notes     *string `json:"notes"`
}

func (h *mvpHandler) listContacts(w http.ResponseWriter, r *http.Request) {
	contacts, err := h.queries.ListContacts(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"contacts": contactListResponse(contacts)})
}

func (h *mvpHandler) createContact(w http.ResponseWriter, r *http.Request) {
	var input contactInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}

	params, err := h.contactCreateParams(input)
	if err != nil {
		writeInvalid(w, err)
		return
	}
	created, err := h.queries.CreateContact(r.Context(), params)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, contactResponse(created))
}

func (h *mvpHandler) getContact(w http.ResponseWriter, r *http.Request) {
	contact, err := h.queries.GetContact(r.Context(), r.PathValue("id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, contactResponse(contact))
}

func (h *mvpHandler) updateContact(w http.ResponseWriter, r *http.Request) {
	current, err := h.queries.GetContact(r.Context(), r.PathValue("id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}

	var input contactInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	params, err := h.contactUpdateParams(current, input)
	if err != nil {
		writeInvalid(w, err)
		return
	}
	updated, err := h.queries.UpdateContact(r.Context(), params)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, contactResponse(updated))
}

func (h *mvpHandler) listHouseholdContacts(w http.ResponseWriter, r *http.Request) {
	rows, err := h.queries.ListHouseholdContacts(r.Context(), r.PathValue("id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"contacts": householdContactListResponse(rows)})
}

func (h *mvpHandler) linkHouseholdContact(w http.ResponseWriter, r *http.Request) {
	var input householdContactInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	role, err := validateHouseholdContact(input)
	if err != nil {
		writeInvalid(w, err)
		return
	}

	err = h.queries.LinkHouseholdContact(r.Context(), dbqueries.LinkHouseholdContactParams{
		HouseholdID: r.PathValue("id"),
		ContactID:   input.ContactID,
		Role:        role,
		IsPrimary:   intFlag(input.IsPrimary),
		Notes:       optionalText(input.Notes),
		CreatedAt:   timestamp(h.now),
	})
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"status": "linked"})
}

func (h *mvpHandler) contactCreateParams(input contactInput) (dbqueries.CreateContactParams, error) {
	displayName, err := requiredText(input.DisplayName, "displayName")
	if err != nil {
		return dbqueries.CreateContactParams{}, err
	}
	recordID, err := newRecordID()
	if err != nil {
		return dbqueries.CreateContactParams{}, err
	}
	now := timestamp(h.now)
	return dbqueries.CreateContactParams{
		ID:          recordID,
		DisplayName: displayName,
		Phone:       optionalText(input.Phone),
		WhatsappID:  optionalText(input.WhatsappID),
		TelegramID:  optionalText(input.TelegramID),
		Email:       optionalText(input.Email),
		Notes:       optionalText(input.Notes),
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (h *mvpHandler) contactUpdateParams(current dbqueries.Contact, input contactInput) (dbqueries.UpdateContactParams, error) {
	displayName := current.DisplayName
	if input.DisplayName != "" {
		value, err := requiredText(input.DisplayName, "displayName")
		if err != nil {
			return dbqueries.UpdateContactParams{}, err
		}
		displayName = value
	}
	active := boolFlag(current.Active)
	if input.Active != nil {
		active = *input.Active
	}
	return dbqueries.UpdateContactParams{
		DisplayName: displayName,
		Phone:       patchText(input.Phone, current.Phone),
		WhatsappID:  patchText(input.WhatsappID, current.WhatsappID),
		TelegramID:  patchText(input.TelegramID, current.TelegramID),
		Email:       patchText(input.Email, current.Email),
		Notes:       patchText(input.Notes, current.Notes),
		Active:      intFlag(active),
		UpdatedAt:   timestamp(h.now),
		ID:          current.ID,
	}, nil
}

func validateHouseholdContact(input householdContactInput) (string, error) {
	contactID, err := requiredText(input.ContactID, "contactId")
	if err != nil {
		return "", err
	}
	role := defaultText(input.Role, "owner")
	if !allowed(role, "owner", "partner", "family", "domestic_worker", "payer", "emergency_contact", "vet", "other") {
		return "", fmt.Errorf("role is not supported")
	}
	input.ContactID = contactID
	return role, nil
}

func patchText(value *string, current sql.NullString) sql.NullString {
	if value == nil {
		return current
	}
	return optionalText(value)
}

func contactListResponse(rows []dbqueries.Contact) []map[string]any {
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, contactResponse(row))
	}
	return items
}

func contactResponse(row dbqueries.Contact) map[string]any {
	return map[string]any{
		"id":          row.ID,
		"displayName": row.DisplayName,
		"phone":       textValue(row.Phone),
		"whatsappId":  textValue(row.WhatsappID),
		"telegramId":  textValue(row.TelegramID),
		"email":       textValue(row.Email),
		"notes":       textValue(row.Notes),
		"active":      boolFlag(row.Active),
		"createdAt":   row.CreatedAt,
		"updatedAt":   row.UpdatedAt,
	}
}

func householdContactListResponse(rows []dbqueries.ListHouseholdContactsRow) []map[string]any {
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, map[string]any{
			"householdId": row.HouseholdID,
			"contactId":   row.ContactID,
			"role":        row.Role,
			"isPrimary":   boolFlag(row.IsPrimary),
			"notes":       textValue(row.Notes),
			"displayName": row.DisplayName,
			"phone":       textValue(row.Phone),
			"whatsappId":  textValue(row.WhatsappID),
			"telegramId":  textValue(row.TelegramID),
			"email":       textValue(row.Email),
			"createdAt":   row.CreatedAt,
		})
	}
	return items
}
