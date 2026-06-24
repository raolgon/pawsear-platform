package httpapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	dbqueries "github.com/pawsear/pawsear-platform/apps/api/internal/db/queries"
)

type paymentReceiptSnapshot struct {
	BusinessName     string                     `json:"businessName"`
	ReceiptNumber    string                     `json:"receiptNumber"`
	PaymentID        string                     `json:"paymentId"`
	IssuedAt         string                     `json:"issuedAt"`
	ReceivedAt       string                     `json:"receivedAt"`
	PayerName        string                     `json:"payerName"`
	AmountMinor      int64                      `json:"amountMinor"`
	Currency         string                     `json:"currency"`
	Method           string                     `json:"method"`
	Reference        *string                    `json:"reference"`
	Allocations      []paymentReceiptAllocation `json:"allocations"`
	AllocatedMinor   int64                      `json:"allocatedMinor"`
	UnallocatedMinor int64                      `json:"unallocatedMinor"`
}

type paymentReceiptAllocation struct {
	ChargeID      string `json:"chargeId"`
	Description   string `json:"description"`
	HouseholdName string `json:"householdName"`
	AmountMinor   int64  `json:"amountMinor"`
}

func (h *mvpHandler) issuePaymentReceipt(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("id")
	if existing, err := h.queries.GetPaymentReceiptByPayment(r.Context(), paymentID); err == nil {
		h.writePaymentReceipt(w, http.StatusOK, existing)
		return
	} else if err != sql.ErrNoRows {
		writeStoreError(w, err)
		return
	}

	source, err := h.queries.GetPaymentReceiptSource(r.Context(), paymentID)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	allocationRows, err := h.queries.ListPaymentReceiptAllocationSources(r.Context(), paymentID)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	receiptID, err := newRecordID()
	if err != nil {
		writeStoreError(w, err)
		return
	}
	issuedAt := timestamp(h.now)
	receiptNumber := receiptNumber(issuedAt, receiptID)
	snapshot := receiptSnapshot(source, allocationRows, receiptNumber, issuedAt)
	snapshotJSON, err := json.Marshal(snapshot)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	created, err := h.queries.CreatePaymentReceipt(r.Context(), dbqueries.CreatePaymentReceiptParams{
		ID: receiptID, PaymentID: paymentID, ReceiptNumber: receiptNumber,
		SnapshotJson: string(snapshotJSON), IssuedAt: issuedAt, CreatedAt: issuedAt,
	})
	if err != nil {
		writeStoreError(w, err)
		return
	}
	h.writePaymentReceipt(w, http.StatusCreated, created)
}

func (h *mvpHandler) getPaymentReceiptByPayment(w http.ResponseWriter, r *http.Request) {
	receipt, err := h.queries.GetPaymentReceiptByPayment(r.Context(), r.PathValue("id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	h.writePaymentReceipt(w, http.StatusOK, receipt)
}

func (h *mvpHandler) downloadPaymentReceipt(w http.ResponseWriter, r *http.Request) {
	receipt, err := h.queries.GetPaymentReceiptByPayment(r.Context(), r.PathValue("id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	snapshot, err := decodeReceiptSnapshot(receipt)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	format := r.PathValue("format")
	var contents []byte
	var contentType string
	switch format {
	case "png":
		contents, err = renderPaymentReceiptPNG(snapshot)
		contentType = "image/png"
	case "pdf":
		contents, err = renderPaymentReceiptPDF(snapshot)
		contentType = "application/pdf"
	default:
		writeNotFound(w, fmt.Errorf("receipt format was not found"))
		return
	}
	if err != nil {
		writeStoreError(w, err)
		return
	}

	disposition := "inline"
	if r.URL.Query().Get("download") == "1" {
		disposition = "attachment"
	}
	filename := "recibo-" + receipt.ReceiptNumber + "." + format
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`%s; filename="%s"`, disposition, filename))
	w.Header().Set("Cache-Control", "private, max-age=300")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(contents)
}

func (h *mvpHandler) writePaymentReceipt(w http.ResponseWriter, status int, receipt dbqueries.PaymentReceipt) {
	snapshot, err := decodeReceiptSnapshot(receipt)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, status, map[string]any{
		"id": receipt.ID, "paymentId": receipt.PaymentID, "receiptNumber": receipt.ReceiptNumber,
		"issuedAt": receipt.IssuedAt, "snapshot": snapshot,
	})
}

func decodeReceiptSnapshot(receipt dbqueries.PaymentReceipt) (paymentReceiptSnapshot, error) {
	var snapshot paymentReceiptSnapshot
	if err := json.Unmarshal([]byte(receipt.SnapshotJson), &snapshot); err != nil {
		return paymentReceiptSnapshot{}, fmt.Errorf("decode receipt snapshot: %w", err)
	}
	return snapshot, nil
}

func receiptSnapshot(source dbqueries.GetPaymentReceiptSourceRow, rows []dbqueries.ListPaymentReceiptAllocationSourcesRow, number string, issuedAt string) paymentReceiptSnapshot {
	payerName := "Pagador no registrado"
	if source.PayerName.Valid {
		payerName = source.PayerName.String
	}
	allocations := make([]paymentReceiptAllocation, 0, len(rows))
	var allocatedMinor int64
	for _, row := range rows {
		allocations = append(allocations, paymentReceiptAllocation{
			ChargeID: row.ChargeID, Description: row.Description,
			HouseholdName: row.HouseholdName, AmountMinor: row.AmountMinor,
		})
		allocatedMinor += row.AmountMinor
	}
	unallocatedMinor := source.AmountMinor - allocatedMinor
	if unallocatedMinor < 0 {
		unallocatedMinor = 0
	}
	return paymentReceiptSnapshot{
		BusinessName: "Pawsear", ReceiptNumber: number, PaymentID: source.PaymentID,
		IssuedAt: issuedAt, ReceivedAt: source.ReceivedAt, PayerName: payerName,
		AmountMinor: source.AmountMinor, Currency: source.Currency, Method: source.Method,
		Reference: textValue(source.Reference), Allocations: allocations,
		AllocatedMinor: allocatedMinor, UnallocatedMinor: unallocatedMinor,
	}
}

func receiptNumber(issuedAt string, id string) string {
	issued, err := time.Parse(time.RFC3339Nano, issuedAt)
	if err != nil {
		issued = time.Now().UTC()
	}
	suffix := id
	if len(suffix) > 8 {
		suffix = suffix[len(suffix)-8:]
	}
	return fmt.Sprintf("REC-%s-%s", issued.Format("20060102"), strings.ToUpper(suffix))
}
