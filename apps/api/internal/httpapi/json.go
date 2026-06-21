package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var errInvalidJSON = errors.New("invalid json")

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func readJSON(r *http.Request, target any) error {
	defer r.Body.Close()

	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("%w: %v", errInvalidJSON, err)
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("%w: request body must contain a single json object", errInvalidJSON)
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, err error) {
	status, code := householdStatus(err)
	if errors.Is(err, errInvalidJSON) {
		status = http.StatusBadRequest
		code = "invalid_json"
	}

	writeJSON(w, status, errorResponse{
		Error:   code,
		Message: err.Error(),
	})
}
