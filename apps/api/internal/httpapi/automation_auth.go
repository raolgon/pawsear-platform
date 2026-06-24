package httpapi

import (
	"crypto/subtle"
	"errors"
	"net/http"
	"strings"
)

var errAutomationToken = errors.New("a valid automation token is required")

func requireAutomationToken(expected string, next http.Handler) http.Handler {
	expected = strings.TrimSpace(expected)
	if expected == "" {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		provided, hasBearer := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		provided = strings.TrimSpace(provided)
		if !hasBearer || subtle.ConstantTimeCompare([]byte(provided), []byte(expected)) != 1 {
			writeAPIError(w, http.StatusUnauthorized, "unauthorized", errAutomationToken)
			return
		}
		next.ServeHTTP(w, r)
	})
}
