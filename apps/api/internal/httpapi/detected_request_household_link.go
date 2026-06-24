package httpapi

import (
	"fmt"
	"net/http"

	dbqueries "github.com/pawsear/pawsear-platform/apps/api/internal/db/queries"
)

type detectedRequestHouseholdLinkInput struct {
	HouseholdID string `json:"householdId"`
}

func (h *mvpHandler) linkDetectedRequestHousehold(w http.ResponseWriter, r *http.Request) {
	var input detectedRequestHouseholdLinkInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	householdID, err := requiredText(input.HouseholdID, "householdId")
	if err != nil {
		writeInvalid(w, err)
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
	if !detected.ContactID.Valid {
		writeConflict(w, fmt.Errorf("identify the sender before choosing a household"))
		return
	}
	household, err := queries.GetHousehold(r.Context(), householdID)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if household.Active == 0 {
		writeConflict(w, fmt.Errorf("household is inactive"))
		return
	}
	belongs, err := queries.ContactBelongsToHousehold(r.Context(), dbqueries.ContactBelongsToHouseholdParams{
		ContactID: detected.ContactID.String, HouseholdID: householdID,
	})
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if belongs == 0 {
		writeConflict(w, fmt.Errorf("contact is not linked to this household"))
		return
	}
	now := timestamp(h.now)
	if err := queries.SetConversationHouseholdByMessage(r.Context(), dbqueries.SetConversationHouseholdByMessageParams{
		HouseholdID: nullableID(householdID), UpdatedAt: now, ID: detected.MessageID,
	}); err != nil {
		writeStoreError(w, err)
		return
	}
	if err := queries.SetDetectedRequestHousehold(r.Context(), dbqueries.SetDetectedRequestHouseholdParams{
		HouseholdID: nullableID(householdID), UpdatedAt: now, ID: detected.ID,
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
