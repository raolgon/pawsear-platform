package httpapi

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	dbqueries "github.com/pawsear/pawsear-platform/apps/api/internal/db/queries"
)

type detectedRequestBookingInput struct {
	HouseholdID          string   `json:"householdId"`
	ServiceType          string   `json:"serviceType"`
	Status               string   `json:"status"`
	StartAt              string   `json:"startAt"`
	EndAt                *string  `json:"endAt"`
	LocationType         string   `json:"locationType"`
	AddressSnapshot      *string  `json:"addressSnapshot"`
	RequestedByContactID *string  `json:"requestedByContactId"`
	AssignedStaffID      *string  `json:"assignedStaffId"`
	Notes                *string  `json:"notes"`
	PetIDs               []string `json:"petIds"`
	ReviewNotes          *string  `json:"reviewNotes"`
	SourceNote           *string  `json:"sourceNote"`
}

func (h *mvpHandler) convertDetectedRequest(w http.ResponseWriter, r *http.Request) {
	var input detectedRequestBookingInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	if err := validateDetectedRequestBookingInput(input); err != nil {
		writeInvalid(w, err)
		return
	}
	detected, err := h.queries.GetDetectedRequestDetail(r.Context(), r.PathValue("id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if detected.Status == "converted_to_booking" {
		writeConflict(w, fmt.Errorf("detected request was already converted"))
		return
	}
	if detected.Status == "ignored" {
		writeConflict(w, fmt.Errorf("ignored detected request must be reopened before conversion"))
		return
	}

	bookingInput := bookingInputFromDetectedRequest(input, detected)
	if err := validateConversionSchedule(bookingInput.StartAt, bookingInput.EndAt); err != nil {
		writeInvalid(w, err)
		return
	}
	params, err := h.bookingCreateParams(bookingInput)
	if err != nil {
		writeInvalid(w, err)
		return
	}
	if err := h.validatePetIDs(r, input.PetIDs); err != nil {
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
	created, err := queries.CreateBooking(r.Context(), params)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if err := addBookingPets(r, queries, created.ID, input.PetIDs); err != nil {
		writeStoreError(w, err)
		return
	}
	sourceID, err := newRecordID()
	if err != nil {
		writeStoreError(w, err)
		return
	}
	now := timestamp(h.now)
	if err := queries.CreateBookingSource(r.Context(), dbqueries.CreateBookingSourceParams{
		ID: sourceID, BookingID: created.ID, MessageID: nullableID(detected.MessageID),
		DetectedRequestID: nullableID(detected.ID), SourceNote: conversionSourceNote(input.SourceNote, detected.Channel), CreatedAt: now,
	}); err != nil {
		writeStoreError(w, err)
		return
	}
	converted, err := queries.ConvertDetectedRequest(r.Context(), dbqueries.ConvertDetectedRequestParams{
		ConvertedBookingID: nullableID(created.ID), ReviewNotes: optionalText(input.ReviewNotes), UpdatedAt: now, ID: detected.ID,
	})
	if err == sql.ErrNoRows {
		writeConflict(w, fmt.Errorf("detected request cannot be converted from its current state"))
		return
	}
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if err := tx.Commit(); err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"booking": bookingResponse(created), "detectedRequest": detectedRequestResponse(converted),
	})
}

func validateConversionSchedule(startAt string, endAt *string) error {
	startAt = strings.TrimSpace(startAt)
	if startAt == "" {
		return nil
	}
	if err := validateOptionalTimestamp(&startAt, "startAt"); err != nil {
		return err
	}
	if err := validateOptionalTimestamp(endAt, "endAt"); err != nil {
		return err
	}
	if endAt == nil || strings.TrimSpace(*endAt) == "" {
		return nil
	}
	start, _ := time.Parse(time.RFC3339, startAt)
	end, _ := time.Parse(time.RFC3339, strings.TrimSpace(*endAt))
	if !end.After(start) {
		return fmt.Errorf("endAt must be after startAt")
	}
	return nil
}

func validateDetectedRequestBookingInput(input detectedRequestBookingInput) error {
	if input.Status != "" && !allowed(input.Status, "requested", "confirmed") {
		return fmt.Errorf("status must be requested or confirmed")
	}
	return nil
}

func bookingInputFromDetectedRequest(input detectedRequestBookingInput, detected dbqueries.GetDetectedRequestDetailRow) bookingInput {
	householdID := input.HouseholdID
	if householdID == "" && detected.HouseholdID.Valid {
		householdID = detected.HouseholdID.String
	}
	serviceType := input.ServiceType
	if serviceType == "" && detected.DetectedServiceType.Valid {
		serviceType = detected.DetectedServiceType.String
	}
	startAt := input.StartAt
	if startAt == "" && detected.DetectedStartAt.Valid {
		startAt = detected.DetectedStartAt.String
	}
	endAt := input.EndAt
	if endAt == nil && detected.DetectedEndAt.Valid {
		endAt = &detected.DetectedEndAt.String
	}
	requestedBy := input.RequestedByContactID
	if requestedBy == nil && detected.ContactID.Valid {
		requestedBy = &detected.ContactID.String
	}
	source := detected.Channel
	if !allowed(source, "whatsapp", "telegram") {
		source = "import"
	}
	return bookingInput{
		HouseholdID: householdID, ServiceType: serviceType, Status: input.Status, StartAt: startAt,
		EndAt: endAt, LocationType: input.LocationType, AddressSnapshot: input.AddressSnapshot,
		RequestedByContactID: requestedBy, AssignedStaffID: input.AssignedStaffID, Source: source,
		Notes: input.Notes, PetIDs: input.PetIDs,
	}
}

func addBookingPets(r *http.Request, queries *dbqueries.Queries, bookingID string, petIDs []string) error {
	for _, petID := range petIDs {
		if err := queries.AddBookingPet(r.Context(), dbqueries.AddBookingPetParams{
			BookingID: bookingID, PetID: petID,
		}); err != nil {
			return err
		}
	}
	return nil
}

func conversionSourceNote(value *string, channel string) sql.NullString {
	if note := optionalText(value); note.Valid {
		return note
	}
	return nullableString("Converted from reviewed " + channel + " message")
}
