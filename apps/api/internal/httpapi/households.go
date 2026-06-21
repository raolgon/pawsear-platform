package httpapi

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/pawsear/pawsear-platform/apps/api/internal/household"
)

type householdHandler struct {
	service *household.Service
}

func newHouseholdHandler(database *sql.DB) *householdHandler {
	return &householdHandler{
		service: household.NewService(database),
	}
}

func (h *householdHandler) list(w http.ResponseWriter, r *http.Request) {
	households, err := h.service.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"households": households,
	})
}

func (h *householdHandler) create(w http.ResponseWriter, r *http.Request) {
	var input household.CreateInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}

	created, err := h.service.Create(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *householdHandler) get(w http.ResponseWriter, r *http.Request) {
	found, err := h.service.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, found)
}

func (h *householdHandler) update(w http.ResponseWriter, r *http.Request) {
	var input household.UpdateInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}

	updated, err := h.service.Update(r.Context(), r.PathValue("id"), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func householdStatus(err error) (int, string) {
	switch {
	case errors.Is(err, household.ErrNotFound):
		return http.StatusNotFound, "not_found"
	case errors.Is(err, household.ErrInvalid):
		return http.StatusBadRequest, "invalid_household"
	case errors.Is(err, household.ErrNoChanges):
		return http.StatusBadRequest, "no_changes"
	default:
		return http.StatusInternalServerError, "internal_error"
	}
}
