package httpapi

import (
	"fmt"
	"net/http"

	dbqueries "github.com/pawsear/pawsear-platform/apps/api/internal/db/queries"
)

type bookingInput struct {
	HouseholdID          string   `json:"householdId"`
	ServiceType          string   `json:"serviceType"`
	Status               string   `json:"status"`
	StartAt              string   `json:"startAt"`
	EndAt                *string  `json:"endAt"`
	LocationType         string   `json:"locationType"`
	AddressSnapshot      *string  `json:"addressSnapshot"`
	RequestedByContactID *string  `json:"requestedByContactId"`
	AssignedStaffID      *string  `json:"assignedStaffId"`
	Source               string   `json:"source"`
	Notes                *string  `json:"notes"`
	CompletedAt          *string  `json:"completedAt"`
	CancelledAt          *string  `json:"cancelledAt"`
	PetIDs               []string `json:"petIds"`
}

func (h *mvpHandler) listBookings(w http.ResponseWriter, r *http.Request) {
	var rows []dbqueries.Booking
	var err error
	switch {
	case r.URL.Query().Get("householdId") != "":
		rows, err = h.queries.ListBookingsByHousehold(r.Context(), r.URL.Query().Get("householdId"))
	case r.URL.Query().Get("date") != "":
		start, end, rangeErr := parseDayRange(r.URL.Query().Get("date"), h.now)
		if rangeErr != nil {
			writeInvalid(w, rangeErr)
			return
		}
		rows, err = h.queries.ListBookingsByRange(r.Context(), dbqueries.ListBookingsByRangeParams{StartAt: start, StartAt_2: end})
	default:
		rows, err = h.queries.ListBookings(r.Context())
	}
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"bookings": bookingListResponse(rows)})
}

func (h *mvpHandler) createBooking(w http.ResponseWriter, r *http.Request) {
	var input bookingInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	params, err := h.bookingCreateParams(input)
	if err != nil {
		writeInvalid(w, err)
		return
	}
	if err := h.validatePetIDs(r, input.PetIDs); err != nil {
		writeInvalid(w, err)
		return
	}
	created, err := h.queries.CreateBooking(r.Context(), params)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if err := h.replaceBookingPets(r, created.ID, input.PetIDs); err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, bookingResponse(created))
}

func (h *mvpHandler) getBooking(w http.ResponseWriter, r *http.Request) {
	booking, err := h.queries.GetBooking(r.Context(), r.PathValue("id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	pets, err := h.queries.ListBookingPets(r.Context(), booking.ID)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	response := bookingResponse(booking)
	response["pets"] = bookingPetListResponse(pets)
	writeJSON(w, http.StatusOK, response)
}

func (h *mvpHandler) updateBooking(w http.ResponseWriter, r *http.Request) {
	current, err := h.queries.GetBooking(r.Context(), r.PathValue("id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	var input bookingInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	params, err := h.bookingUpdateParams(current, input)
	if err != nil {
		writeInvalid(w, err)
		return
	}
	if input.PetIDs != nil {
		if err := h.validatePetIDs(r, input.PetIDs); err != nil {
			writeInvalid(w, err)
			return
		}
	}
	updated, err := h.queries.UpdateBooking(r.Context(), params)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if input.PetIDs != nil {
		if err := h.replaceBookingPets(r, updated.ID, input.PetIDs); err != nil {
			writeStoreError(w, err)
			return
		}
	}
	writeJSON(w, http.StatusOK, bookingResponse(updated))
}

func (h *mvpHandler) bookingCreateParams(input bookingInput) (dbqueries.CreateBookingParams, error) {
	householdID, serviceType, startAt, err := validateBookingRequired(input)
	if err != nil {
		return dbqueries.CreateBookingParams{}, err
	}
	recordID, err := newRecordID()
	if err != nil {
		return dbqueries.CreateBookingParams{}, err
	}
	now := timestamp(h.now)
	return dbqueries.CreateBookingParams{
		ID: recordID, HouseholdID: householdID, ServiceType: serviceType,
		Status: defaultBookingStatus(input.Status), StartAt: startAt, EndAt: optionalText(input.EndAt),
		LocationType: defaultLocationType(input.LocationType), AddressSnapshot: optionalText(input.AddressSnapshot),
		RequestedByContactID: optionalText(input.RequestedByContactID), AssignedStaffID: optionalText(input.AssignedStaffID),
		Source: defaultBookingSource(input.Source), Notes: optionalText(input.Notes), CreatedAt: now, UpdatedAt: now,
	}, nil
}

func (h *mvpHandler) bookingUpdateParams(current dbqueries.Booking, input bookingInput) (dbqueries.UpdateBookingParams, error) {
	status := defaultText(input.Status, current.Status)
	if !allowed(status, "requested", "confirmed", "in_progress", "completed", "cancelled") {
		return dbqueries.UpdateBookingParams{}, fmt.Errorf("status is not supported")
	}
	if !bookingStatusTransitionAllowed(current.Status, status) {
		return dbqueries.UpdateBookingParams{}, fmt.Errorf("booking status transition is not allowed")
	}
	serviceType := defaultText(input.ServiceType, current.ServiceType)
	if !allowed(serviceType, "walk", "client_home_sitting", "boarding", "visit", "transport", "other") {
		return dbqueries.UpdateBookingParams{}, fmt.Errorf("serviceType is not supported")
	}
	completedAt := patchText(input.CompletedAt, current.CompletedAt)
	cancelledAt := patchText(input.CancelledAt, current.CancelledAt)
	if status == "completed" && !completedAt.Valid {
		completedAt = optionalText(ptr(timestamp(h.now)))
	}
	if status == "cancelled" && !cancelledAt.Valid {
		cancelledAt = optionalText(ptr(timestamp(h.now)))
	}
	return dbqueries.UpdateBookingParams{
		HouseholdID: defaultText(input.HouseholdID, current.HouseholdID), ServiceType: serviceType, Status: status,
		StartAt: defaultText(input.StartAt, current.StartAt), EndAt: patchText(input.EndAt, current.EndAt),
		LocationType:         defaultText(input.LocationType, current.LocationType),
		AddressSnapshot:      patchText(input.AddressSnapshot, current.AddressSnapshot),
		RequestedByContactID: patchText(input.RequestedByContactID, current.RequestedByContactID),
		AssignedStaffID:      patchText(input.AssignedStaffID, current.AssignedStaffID),
		Source:               defaultText(input.Source, current.Source), Notes: patchText(input.Notes, current.Notes),
		CompletedAt: completedAt,
		CancelledAt: cancelledAt,
		UpdatedAt:   timestamp(h.now), ID: current.ID,
	}, nil
}

func bookingStatusTransitionAllowed(current string, next string) bool {
	if current == next {
		return true
	}
	switch current {
	case "requested":
		return next == "confirmed" || next == "cancelled"
	case "confirmed":
		return next == "in_progress" || next == "cancelled"
	case "in_progress":
		return next == "completed" || next == "cancelled"
	default:
		return false
	}
}

func validateBookingRequired(input bookingInput) (string, string, string, error) {
	householdID, err := requiredText(input.HouseholdID, "householdId")
	if err != nil {
		return "", "", "", err
	}
	startAt, err := requiredText(input.StartAt, "startAt")
	if err != nil {
		return "", "", "", err
	}
	serviceType := defaultText(input.ServiceType, "walk")
	if !allowed(serviceType, "walk", "client_home_sitting", "boarding", "visit", "transport", "other") {
		return "", "", "", fmt.Errorf("serviceType is not supported")
	}
	return householdID, serviceType, startAt, nil
}

func defaultBookingStatus(value string) string {
	status := defaultText(value, "requested")
	if allowed(status, "requested", "confirmed", "in_progress", "completed", "cancelled") {
		return status
	}
	return "requested"
}

func defaultLocationType(value string) string {
	locationType := defaultText(value, "household_home")
	if allowed(locationType, "household_home", "caregiver_home", "other") {
		return locationType
	}
	return "household_home"
}

func defaultBookingSource(value string) string {
	source := defaultText(value, "manual")
	if allowed(source, "manual", "whatsapp", "telegram", "import") {
		return source
	}
	return "manual"
}

func (h *mvpHandler) replaceBookingPets(r *http.Request, bookingID string, petIDs []string) error {
	if err := h.queries.DeleteBookingPets(r.Context(), bookingID); err != nil {
		return err
	}
	for _, petID := range petIDs {
		petID, err := requiredText(petID, "petIds")
		if err != nil {
			return err
		}
		err = h.queries.AddBookingPet(r.Context(), dbqueries.AddBookingPetParams{BookingID: bookingID, PetID: petID})
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *mvpHandler) validatePetIDs(r *http.Request, petIDs []string) error {
	for _, petID := range petIDs {
		cleanID, err := requiredText(petID, "petIds")
		if err != nil {
			return err
		}
		if _, err := h.queries.GetPet(r.Context(), cleanID); err != nil {
			return fmt.Errorf("petId %s was not found", cleanID)
		}
	}
	return nil
}

func bookingListResponse(rows []dbqueries.Booking) []map[string]any {
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, bookingResponse(row))
	}
	return items
}

func bookingResponse(row dbqueries.Booking) map[string]any {
	return map[string]any{
		"id": row.ID, "householdId": row.HouseholdID, "serviceType": row.ServiceType, "status": row.Status,
		"startAt": row.StartAt, "endAt": textValue(row.EndAt), "locationType": row.LocationType,
		"addressSnapshot": textValue(row.AddressSnapshot), "requestedByContactId": textValue(row.RequestedByContactID),
		"assignedStaffId": textValue(row.AssignedStaffID), "source": row.Source, "notes": textValue(row.Notes),
		"completedAt": textValue(row.CompletedAt), "cancelledAt": textValue(row.CancelledAt),
		"createdAt": row.CreatedAt, "updatedAt": row.UpdatedAt,
	}
}

func bookingPetListResponse(rows []dbqueries.ListBookingPetsRow) []map[string]any {
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, map[string]any{
			"bookingId": row.BookingID, "petId": row.PetID, "name": row.Name, "species": row.Species, "notes": textValue(row.Notes),
		})
	}
	return items
}
