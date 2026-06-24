package httpapi

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	dbqueries "github.com/pawsear/pawsear-platform/apps/api/internal/db/queries"
)

type outboundReplyInput struct {
	TemplateKey string `json:"templateKey"`
}

type outboundDeliveryInput struct {
	Status string  `json:"status"`
	Error  *string `json:"error"`
}

func (h *mvpHandler) queueOutboundReply(w http.ResponseWriter, r *http.Request) {
	var input outboundReplyInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	if !allowed(input.TemplateKey, "request_details", "booking_confirmed", "request_declined") {
		writeInvalid(w, fmt.Errorf("templateKey is not supported"))
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
	if detected.Channel != "telegram" || !detected.ExternalConversationID.Valid {
		writeInvalid(w, fmt.Errorf("only Telegram requests can receive replies in this version"))
		return
	}
	if _, err := queries.GetPendingOutboundMessageForRequest(r.Context(), detected.ID); err == nil {
		writeConflict(w, fmt.Errorf("request already has a pending reply"))
		return
	} else if err != sql.ErrNoRows {
		writeStoreError(w, err)
		return
	}
	body, nextStatus, err := outboundTemplate(detected, input.TemplateKey)
	if err != nil {
		writeConflict(w, err)
		return
	}
	id, err := newRecordID()
	if err != nil {
		writeStoreError(w, err)
		return
	}
	now := timestamp(h.now)
	created, err := queries.CreateOutboundMessage(r.Context(), dbqueries.CreateOutboundMessageParams{
		ID: id, DetectedRequestID: detected.ID, Channel: "telegram",
		RecipientExternalID: detected.ExternalConversationID.String, TemplateKey: input.TemplateKey,
		Body: body, CreatedAt: now, UpdatedAt: now,
	})
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if nextStatus != "" {
		if _, err := queries.UpdateDetectedRequestStatus(r.Context(), dbqueries.UpdateDetectedRequestStatusParams{
			Status: nextStatus, ReviewNotes: detected.ReviewNotes, UpdatedAt: now, ID: detected.ID,
		}); err != nil {
			writeStoreError(w, err)
			return
		}
	}
	if err := tx.Commit(); err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, outboundMessageResponse(created))
}

func (h *mvpHandler) listOutboundMessages(w http.ResponseWriter, r *http.Request) {
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	if status != "" && !allowed(status, "pending", "sent", "failed", "cancelled") {
		writeInvalid(w, fmt.Errorf("status is not supported"))
		return
	}
	limit := 200
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 || parsed > 200 {
			writeInvalid(w, fmt.Errorf("limit must be between 1 and 200"))
			return
		}
		limit = parsed
	}
	rows, err := h.queries.ListOutboundMessages(r.Context(), dbqueries.ListOutboundMessagesParams{
		StatusFilter: status, RequestFilter: strings.TrimSpace(r.URL.Query().Get("detectedRequestId")), ResultLimit: int64(limit),
	})
	if err != nil {
		writeStoreError(w, err)
		return
	}
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, outboundMessageResponse(row))
	}
	writeJSON(w, http.StatusOK, map[string]any{"outboundMessages": items})
}

func (h *mvpHandler) updateOutboundDelivery(w http.ResponseWriter, r *http.Request) {
	var input outboundDeliveryInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	now := timestamp(h.now)
	var updated dbqueries.OutboundMessage
	var err error
	switch input.Status {
	case "sent":
		updated, err = h.queries.MarkOutboundMessageSent(r.Context(), dbqueries.MarkOutboundMessageSentParams{
			SentAt: nullableString(now), UpdatedAt: now, ID: r.PathValue("id"),
		})
	case "failed":
		errorText := optionalText(input.Error)
		if errorText.Valid && len(errorText.String) > 500 {
			errorText.String = errorText.String[:500]
		}
		updated, err = h.queries.MarkOutboundMessageFailed(r.Context(), dbqueries.MarkOutboundMessageFailedParams{
			LastError: errorText, UpdatedAt: now, ID: r.PathValue("id"),
		})
	default:
		writeInvalid(w, fmt.Errorf("status must be sent or failed"))
		return
	}
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, outboundMessageResponse(updated))
}

func outboundTemplate(detected dbqueries.GetDetectedRequestDetailRow, templateKey string) (string, string, error) {
	switch templateKey {
	case "request_details":
		if detected.Status == "converted_to_booking" {
			return "", "", fmt.Errorf("converted requests cannot ask for booking details")
		}
		return "Gracias por escribirnos. ¿Nos confirmas la fecha, hora, servicio y mascota para revisar tu solicitud?", "needs_more_info", nil
	case "request_declined":
		if detected.Status == "converted_to_booking" {
			return "", "", fmt.Errorf("converted requests cannot be declined")
		}
		return "Gracias por escribirnos. Por ahora no podemos aceptar esta solicitud. Si quieres, envíanos otra fecha u horario.", "ignored", nil
	case "booking_confirmed":
		if detected.Status != "converted_to_booking" || !detected.ConvertedBookingStartAt.Valid {
			return "", "", fmt.Errorf("only converted requests can send a booking confirmation")
		}
		return "Tu servicio quedó confirmado para " + formatReplyTime(detected.ConvertedBookingStartAt.String) + ". Si necesitas cambiar algo, responde a este mensaje.", "", nil
	default:
		return "", "", fmt.Errorf("template is not supported")
	}
}

func formatReplyTime(value string) string {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return value
	}
	location, err := time.LoadLocation("America/Mexico_City")
	if err == nil {
		parsed = parsed.In(location)
	}
	return parsed.Format("02/01/2006 a las 15:04")
}

func outboundMessageResponse(row dbqueries.OutboundMessage) map[string]any {
	return map[string]any{
		"id": row.ID, "detectedRequestId": row.DetectedRequestID, "channel": row.Channel,
		"recipientExternalId": row.RecipientExternalID, "templateKey": row.TemplateKey, "body": row.Body,
		"status": row.Status, "attempts": row.Attempts, "lastError": textValue(row.LastError),
		"createdAt": row.CreatedAt, "sentAt": textValue(row.SentAt), "updatedAt": row.UpdatedAt,
	}
}
