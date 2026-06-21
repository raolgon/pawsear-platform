package httpapi

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	dbqueries "github.com/pawsear/pawsear-platform/apps/api/internal/db/queries"
)

type careTaskInput struct {
	BookingID          *string `json:"bookingId"`
	HouseholdID        string  `json:"householdId"`
	PetID              *string `json:"petId"`
	TaskType           string  `json:"taskType"`
	Title              string  `json:"title"`
	Instructions       *string `json:"instructions"`
	DueAt              *string `json:"dueAt"`
	Status             string  `json:"status"`
	AssignedStaffID    *string `json:"assignedStaffId"`
	CompletedAt        *string `json:"completedAt"`
	CompletedByStaffID *string `json:"completedByStaffId"`
	SkippedReason      *string `json:"skippedReason"`
}

type chargeInput struct {
	HouseholdID string  `json:"householdId"`
	BookingID   *string `json:"bookingId"`
	Description string  `json:"description"`
	AmountMinor int64   `json:"amountMinor"`
	Currency    string  `json:"currency"`
	Status      string  `json:"status"`
	DueDate     *string `json:"dueDate"`
}

type paymentInput struct {
	PayerContactID *string           `json:"payerContactId"`
	ReceivedAt     string            `json:"receivedAt"`
	AmountMinor    int64             `json:"amountMinor"`
	Currency       string            `json:"currency"`
	Method         string            `json:"method"`
	Reference      *string           `json:"reference"`
	Notes          *string           `json:"notes"`
	Allocations    []allocationInput `json:"allocations"`
}

type allocationInput struct {
	ChargeID    string `json:"chargeId"`
	AmountMinor int64  `json:"amountMinor"`
}

func (h *mvpHandler) listCareTasks(w http.ResponseWriter, r *http.Request) {
	var rows []dbqueries.CareTask
	var err error
	switch {
	case r.URL.Query().Get("householdId") != "":
		rows, err = h.queries.ListCareTasksByHousehold(r.Context(), r.URL.Query().Get("householdId"))
	case r.URL.Query().Get("bookingId") != "":
		rows, err = h.queries.ListCareTasksByBooking(r.Context(), optionalText(ptr(r.URL.Query().Get("bookingId"))))
	case r.URL.Query().Get("date") != "":
		start, end, rangeErr := parseDayRange(r.URL.Query().Get("date"), h.now)
		if rangeErr != nil {
			writeInvalid(w, rangeErr)
			return
		}
		rows, err = h.queries.ListCareTasksByRange(r.Context(), dbqueries.ListCareTasksByRangeParams{DueAt: sql.NullString{String: start, Valid: true}, DueAt_2: sql.NullString{String: end, Valid: true}})
	default:
		rows, err = h.queries.ListCareTasks(r.Context())
	}
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"careTasks": careTaskListResponse(rows)})
}

func (h *mvpHandler) createCareTask(w http.ResponseWriter, r *http.Request) {
	var input careTaskInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	params, err := h.careTaskCreateParams(input)
	if err != nil {
		writeInvalid(w, err)
		return
	}
	created, err := h.queries.CreateCareTask(r.Context(), params)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, careTaskResponse(created))
}

func (h *mvpHandler) getCareTask(w http.ResponseWriter, r *http.Request) {
	task, err := h.queries.GetCareTask(r.Context(), r.PathValue("id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, careTaskResponse(task))
}

func (h *mvpHandler) updateCareTask(w http.ResponseWriter, r *http.Request) {
	current, err := h.queries.GetCareTask(r.Context(), r.PathValue("id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	var input careTaskInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	params, err := h.careTaskUpdateParams(current, input)
	if err != nil {
		writeInvalid(w, err)
		return
	}
	updated, err := h.queries.UpdateCareTask(r.Context(), params)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, careTaskResponse(updated))
}

func (h *mvpHandler) listCharges(w http.ResponseWriter, r *http.Request) {
	var rows []dbqueries.Charge
	var err error
	if r.URL.Query().Get("householdId") != "" {
		rows, err = h.queries.ListChargesByHousehold(r.Context(), r.URL.Query().Get("householdId"))
	} else if r.URL.Query().Get("status") == "open" {
		rows, err = h.queries.ListOpenCharges(r.Context())
	} else {
		rows, err = h.queries.ListCharges(r.Context())
	}
	if err != nil {
		writeStoreError(w, err)
		return
	}
	response, err := h.chargeListResponse(r.Context(), rows)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"charges": response})
}

func (h *mvpHandler) createCharge(w http.ResponseWriter, r *http.Request) {
	var input chargeInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	params, err := h.chargeCreateParams(input)
	if err != nil {
		writeInvalid(w, err)
		return
	}
	created, err := h.queries.CreateCharge(r.Context(), params)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, chargeResponse(created, 0))
}

func (h *mvpHandler) getCharge(w http.ResponseWriter, r *http.Request) {
	charge, err := h.queries.GetCharge(r.Context(), r.PathValue("id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	allocated, err := h.queries.GetAllocatedTotalForCharge(r.Context(), charge.ID)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, chargeResponse(charge, allocated))
}

func (h *mvpHandler) updateCharge(w http.ResponseWriter, r *http.Request) {
	current, err := h.queries.GetCharge(r.Context(), r.PathValue("id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	var input chargeInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	params, err := h.chargeUpdateParams(current, input)
	if err != nil {
		writeInvalid(w, err)
		return
	}
	updated, err := h.queries.UpdateCharge(r.Context(), params)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	allocated, err := h.queries.GetAllocatedTotalForCharge(r.Context(), updated.ID)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, chargeResponse(updated, allocated))
}

func (h *mvpHandler) listPayments(w http.ResponseWriter, r *http.Request) {
	var rows []dbqueries.Payment
	var err error
	if r.URL.Query().Get("payerContactId") != "" {
		rows, err = h.queries.ListPaymentsByContact(r.Context(), optionalText(ptr(r.URL.Query().Get("payerContactId"))))
	} else {
		rows, err = h.queries.ListPayments(r.Context())
	}
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"payments": paymentListResponse(rows)})
}

func (h *mvpHandler) createPayment(w http.ResponseWriter, r *http.Request) {
	var input paymentInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	created, allocations, err := h.createPaymentWithAllocations(r, input)
	if err != nil {
		writeInvalid(w, err)
		return
	}
	response := paymentResponse(created)
	response["allocations"] = paymentAllocationListResponse(allocations)
	writeJSON(w, http.StatusCreated, response)
}

func (h *mvpHandler) getPayment(w http.ResponseWriter, r *http.Request) {
	payment, err := h.queries.GetPayment(r.Context(), r.PathValue("id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	allocations, err := h.queries.ListPaymentAllocations(r.Context(), payment.ID)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	response := paymentResponse(payment)
	response["allocations"] = paymentAllocationListResponse(allocations)
	writeJSON(w, http.StatusOK, response)
}

func (h *mvpHandler) dashboardToday(w http.ResponseWriter, r *http.Request) {
	start, end, err := parseDayRange(r.URL.Query().Get("date"), h.now)
	if err != nil {
		writeInvalid(w, err)
		return
	}
	bookings, err := h.queries.ListBookingsByRange(r.Context(), dbqueries.ListBookingsByRangeParams{StartAt: start, StartAt_2: end})
	if err != nil {
		writeStoreError(w, err)
		return
	}
	tasks, err := h.queries.ListCareTasksByRange(r.Context(), dbqueries.ListCareTasksByRangeParams{DueAt: sql.NullString{String: start, Valid: true}, DueAt_2: sql.NullString{String: end, Valid: true}})
	if err != nil {
		writeStoreError(w, err)
		return
	}
	charges, err := h.queries.ListOpenCharges(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	chargeResponseItems, err := h.chargeListResponse(r.Context(), charges)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"dateStart": start, "dateEnd": end,
		"bookings":    bookingListResponse(bookings),
		"careTasks":   careTaskListResponse(tasks),
		"openCharges": chargeResponseItems,
	})
}

func (h *mvpHandler) careTaskCreateParams(input careTaskInput) (dbqueries.CreateCareTaskParams, error) {
	householdID, title, taskType, err := validateCareTask(input)
	if err != nil {
		return dbqueries.CreateCareTaskParams{}, err
	}
	recordID, err := newRecordID()
	if err != nil {
		return dbqueries.CreateCareTaskParams{}, err
	}
	now := timestamp(h.now)
	return dbqueries.CreateCareTaskParams{
		ID: recordID, BookingID: optionalText(input.BookingID), HouseholdID: householdID, PetID: optionalText(input.PetID),
		TaskType: taskType, Title: title, Instructions: optionalText(input.Instructions), DueAt: optionalText(input.DueAt),
		Status: defaultTaskStatus(input.Status), AssignedStaffID: optionalText(input.AssignedStaffID), CreatedAt: now, UpdatedAt: now,
	}, nil
}

func (h *mvpHandler) careTaskUpdateParams(current dbqueries.CareTask, input careTaskInput) (dbqueries.UpdateCareTaskParams, error) {
	status := defaultText(input.Status, current.Status)
	if !allowed(status, "pending", "completed", "skipped", "cancelled") {
		return dbqueries.UpdateCareTaskParams{}, fmt.Errorf("status is not supported")
	}
	if current.Status != status && current.Status != "pending" {
		return dbqueries.UpdateCareTaskParams{}, fmt.Errorf("care task status transition is not allowed")
	}
	skippedReason := patchText(input.SkippedReason, current.SkippedReason)
	if status == "skipped" && (!skippedReason.Valid || strings.TrimSpace(skippedReason.String) == "") {
		return dbqueries.UpdateCareTaskParams{}, fmt.Errorf("skippedReason is required when a task is skipped")
	}
	completedAt := patchText(input.CompletedAt, current.CompletedAt)
	if status == "completed" && !completedAt.Valid {
		completedAt = optionalText(ptr(timestamp(h.now)))
	}
	return dbqueries.UpdateCareTaskParams{
		BookingID: patchText(input.BookingID, current.BookingID), HouseholdID: defaultText(input.HouseholdID, current.HouseholdID),
		PetID: patchText(input.PetID, current.PetID), TaskType: defaultText(input.TaskType, current.TaskType),
		Title: defaultText(input.Title, current.Title), Instructions: patchText(input.Instructions, current.Instructions),
		DueAt: patchText(input.DueAt, current.DueAt), Status: status, AssignedStaffID: patchText(input.AssignedStaffID, current.AssignedStaffID),
		CompletedAt: completedAt, CompletedByStaffID: patchText(input.CompletedByStaffID, current.CompletedByStaffID),
		SkippedReason: skippedReason, UpdatedAt: timestamp(h.now), ID: current.ID,
	}, nil
}

func validateCareTask(input careTaskInput) (string, string, string, error) {
	householdID, err := requiredText(input.HouseholdID, "householdId")
	if err != nil {
		return "", "", "", err
	}
	title, err := requiredText(input.Title, "title")
	if err != nil {
		return "", "", "", err
	}
	taskType := defaultText(input.TaskType, "other")
	if !allowed(taskType, "food", "medicine", "walk", "water", "cleaning", "pickup", "dropoff", "photo_update", "other") {
		return "", "", "", fmt.Errorf("taskType is not supported")
	}
	return householdID, title, taskType, nil
}

func defaultTaskStatus(value string) string {
	status := defaultText(value, "pending")
	if allowed(status, "pending", "completed", "skipped", "cancelled") {
		return status
	}
	return "pending"
}

func (h *mvpHandler) chargeCreateParams(input chargeInput) (dbqueries.CreateChargeParams, error) {
	householdID, err := requiredText(input.HouseholdID, "householdId")
	if err != nil {
		return dbqueries.CreateChargeParams{}, err
	}
	description, err := requiredText(input.Description, "description")
	if err != nil {
		return dbqueries.CreateChargeParams{}, err
	}
	if input.AmountMinor <= 0 {
		return dbqueries.CreateChargeParams{}, fmt.Errorf("amountMinor must be greater than 0")
	}
	recordID, err := newRecordID()
	if err != nil {
		return dbqueries.CreateChargeParams{}, err
	}
	now := timestamp(h.now)
	return dbqueries.CreateChargeParams{
		ID: recordID, HouseholdID: householdID, BookingID: optionalText(input.BookingID), Description: description,
		AmountMinor: input.AmountMinor, Currency: defaultText(input.Currency, "MXN"), Status: defaultChargeStatus(input.Status),
		DueDate: optionalText(input.DueDate), CreatedAt: now, UpdatedAt: now,
	}, nil
}

func (h *mvpHandler) chargeUpdateParams(current dbqueries.Charge, input chargeInput) (dbqueries.UpdateChargeParams, error) {
	householdID := defaultText(input.HouseholdID, current.HouseholdID)
	description := defaultText(input.Description, current.Description)
	amountMinor := current.AmountMinor
	if input.AmountMinor > 0 {
		amountMinor = input.AmountMinor
	}
	status := defaultText(input.Status, current.Status)
	if !allowed(status, "unpaid", "partially_paid", "paid", "waived", "void") {
		return dbqueries.UpdateChargeParams{}, fmt.Errorf("status is not supported")
	}
	return dbqueries.UpdateChargeParams{
		HouseholdID: householdID, BookingID: patchText(input.BookingID, current.BookingID),
		Description: description, AmountMinor: amountMinor, Currency: defaultText(input.Currency, current.Currency),
		Status: status, DueDate: patchText(input.DueDate, current.DueDate), UpdatedAt: timestamp(h.now), ID: current.ID,
	}, nil
}

func defaultChargeStatus(value string) string {
	status := defaultText(value, "unpaid")
	if allowed(status, "unpaid", "partially_paid", "paid", "waived", "void") {
		return status
	}
	return "unpaid"
}

func (h *mvpHandler) createPaymentWithAllocations(r *http.Request, input paymentInput) (dbqueries.Payment, []dbqueries.PaymentAllocation, error) {
	if input.AmountMinor <= 0 {
		return dbqueries.Payment{}, nil, fmt.Errorf("amountMinor must be greater than 0")
	}
	total := int64(0)
	for _, allocation := range input.Allocations {
		if allocation.AmountMinor <= 0 {
			return dbqueries.Payment{}, nil, fmt.Errorf("allocation amountMinor must be greater than 0")
		}
		total += allocation.AmountMinor
	}
	if total > input.AmountMinor {
		return dbqueries.Payment{}, nil, fmt.Errorf("allocations cannot exceed payment amount")
	}

	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		return dbqueries.Payment{}, nil, err
	}
	defer tx.Rollback()

	created, allocations, err := h.insertPaymentRows(r, dbqueries.New(tx), input)
	if err != nil {
		return dbqueries.Payment{}, nil, err
	}
	if err := tx.Commit(); err != nil {
		return dbqueries.Payment{}, nil, err
	}
	return created, allocations, nil
}

func (h *mvpHandler) insertPaymentRows(r *http.Request, q *dbqueries.Queries, input paymentInput) (dbqueries.Payment, []dbqueries.PaymentAllocation, error) {
	recordID, err := newRecordID()
	if err != nil {
		return dbqueries.Payment{}, nil, err
	}
	now := timestamp(h.now)
	receivedAt := defaultText(input.ReceivedAt, now)
	method := defaultText(input.Method, "cash")
	if !allowed(method, "cash", "bank_transfer", "card_external", "other") {
		return dbqueries.Payment{}, nil, fmt.Errorf("method is not supported")
	}
	created, err := q.CreatePayment(r.Context(), dbqueries.CreatePaymentParams{
		ID: recordID, PayerContactID: optionalText(input.PayerContactID), ReceivedAt: receivedAt,
		AmountMinor: input.AmountMinor, Currency: defaultText(input.Currency, "MXN"), Method: method,
		Reference: optionalText(input.Reference), Notes: optionalText(input.Notes), CreatedAt: now, UpdatedAt: now,
	})
	if err != nil {
		return dbqueries.Payment{}, nil, err
	}
	allocations, err := h.insertAllocations(r, q, created.ID, created.Currency, input.Allocations, now)
	return created, allocations, err
}

func (h *mvpHandler) insertAllocations(r *http.Request, q *dbqueries.Queries, paymentID string, paymentCurrency string, inputs []allocationInput, now string) ([]dbqueries.PaymentAllocation, error) {
	allocations := make([]dbqueries.PaymentAllocation, 0, len(inputs))
	seenChargeIDs := make(map[string]struct{}, len(inputs))
	for _, input := range inputs {
		chargeID, err := requiredText(input.ChargeID, "chargeId")
		if err != nil {
			return nil, err
		}
		if _, exists := seenChargeIDs[chargeID]; exists {
			return nil, fmt.Errorf("chargeId %s is allocated more than once", chargeID)
		}
		seenChargeIDs[chargeID] = struct{}{}
		charge, err := q.GetCharge(r.Context(), chargeID)
		if err != nil {
			return nil, fmt.Errorf("chargeId %s was not found", chargeID)
		}
		allocated, err := q.GetAllocatedTotalForCharge(r.Context(), chargeID)
		if err != nil {
			return nil, err
		}
		if input.AmountMinor > charge.AmountMinor-allocated {
			return nil, fmt.Errorf("allocation exceeds outstanding charge balance")
		}
		if charge.Status != "unpaid" && charge.Status != "partially_paid" {
			return nil, fmt.Errorf("allocations require an open charge")
		}
		if charge.Currency != paymentCurrency {
			return nil, fmt.Errorf("payment and charge currencies must match")
		}
		allocationID, err := newRecordID()
		if err != nil {
			return nil, err
		}
		created, err := q.CreatePaymentAllocation(r.Context(), dbqueries.CreatePaymentAllocationParams{
			ID: allocationID, PaymentID: paymentID, ChargeID: chargeID, AmountMinor: input.AmountMinor, CreatedAt: now,
		})
		if err != nil {
			return nil, err
		}
		if _, err := q.RefreshChargeStatus(r.Context(), dbqueries.RefreshChargeStatusParams{UpdatedAt: now, ID: chargeID}); err != nil {
			return nil, err
		}
		allocations = append(allocations, created)
	}
	return allocations, nil
}

func ptr(value string) *string {
	return &value
}

func careTaskListResponse(rows []dbqueries.CareTask) []map[string]any {
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, careTaskResponse(row))
	}
	return items
}

func careTaskResponse(row dbqueries.CareTask) map[string]any {
	return map[string]any{
		"id": row.ID, "bookingId": textValue(row.BookingID), "householdId": row.HouseholdID, "petId": textValue(row.PetID),
		"taskType": row.TaskType, "title": row.Title, "instructions": textValue(row.Instructions), "dueAt": textValue(row.DueAt),
		"status": row.Status, "assignedStaffId": textValue(row.AssignedStaffID), "completedAt": textValue(row.CompletedAt),
		"completedByStaffId": textValue(row.CompletedByStaffID), "skippedReason": textValue(row.SkippedReason),
		"createdAt": row.CreatedAt, "updatedAt": row.UpdatedAt,
	}
}

func (h *mvpHandler) chargeListResponse(ctx context.Context, rows []dbqueries.Charge) ([]map[string]any, error) {
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		allocated, err := h.queries.GetAllocatedTotalForCharge(ctx, row.ID)
		if err != nil {
			return nil, err
		}
		items = append(items, chargeResponse(row, allocated))
	}
	return items, nil
}

func chargeResponse(row dbqueries.Charge, allocatedMinor int64) map[string]any {
	outstandingMinor := row.AmountMinor - allocatedMinor
	if outstandingMinor < 0 {
		outstandingMinor = 0
	}
	return map[string]any{
		"id": row.ID, "householdId": row.HouseholdID, "bookingId": textValue(row.BookingID),
		"description": row.Description, "amountMinor": row.AmountMinor, "currency": row.Currency,
		"allocatedMinor": allocatedMinor, "outstandingMinor": outstandingMinor, "status": row.Status,
		"dueDate": textValue(row.DueDate), "createdAt": row.CreatedAt, "updatedAt": row.UpdatedAt,
	}
}

func paymentListResponse(rows []dbqueries.Payment) []map[string]any {
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, paymentResponse(row))
	}
	return items
}

func paymentResponse(row dbqueries.Payment) map[string]any {
	return map[string]any{
		"id": row.ID, "payerContactId": textValue(row.PayerContactID), "receivedAt": row.ReceivedAt,
		"amountMinor": row.AmountMinor, "currency": row.Currency, "method": row.Method,
		"reference": textValue(row.Reference), "notes": textValue(row.Notes), "createdAt": row.CreatedAt, "updatedAt": row.UpdatedAt,
	}
}

func paymentAllocationListResponse(rows []dbqueries.PaymentAllocation) []map[string]any {
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, map[string]any{
			"id": row.ID, "paymentId": row.PaymentID, "chargeId": row.ChargeID, "amountMinor": row.AmountMinor, "createdAt": row.CreatedAt,
		})
	}
	return items
}
