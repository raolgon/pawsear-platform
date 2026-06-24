package httpapi

import (
	"bytes"
	"fmt"
	"image/color"
	"math"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/go-pdf/fpdf"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
)

const receiptWidth = 1240

var (
	receiptNavy  = color.RGBA{R: 8, G: 45, B: 96, A: 255}
	receiptGreen = color.RGBA{R: 57, G: 117, B: 70, A: 255}
	receiptMuted = color.RGBA{R: 104, G: 117, B: 135, A: 255}
	receiptLine  = color.RGBA{R: 228, G: 224, B: 214, A: 255}
	receiptCream = color.RGBA{R: 255, G: 253, B: 248, A: 255}
)

func renderPaymentReceiptPNG(snapshot paymentReceiptSnapshot) ([]byte, error) {
	height := receiptCanvasHeight(len(snapshot.Allocations))
	dc := gg.NewContext(receiptWidth, height)
	dc.SetColor(receiptCream)
	dc.Clear()
	dc.SetColor(color.White)
	dc.DrawRoundedRectangle(54, 54, receiptWidth-108, float64(height-108), 34)
	dc.Fill()

	regular, err := receiptFont(goregular.TTF, 34)
	if err != nil {
		return nil, err
	}
	bold, err := receiptFont(gobold.TTF, 34)
	if err != nil {
		return nil, err
	}

	dc.SetColor(receiptNavy)
	dc.SetFontFace(receiptFontSize(bold, gobold.TTF, 42))
	dc.DrawString(snapshot.BusinessName, 108, 145)
	dc.SetFontFace(receiptFontSize(bold, gobold.TTF, 23))
	dc.SetColor(receiptGreen)
	dc.DrawRoundedRectangle(918, 98, 210, 58, 29)
	dc.Fill()
	dc.SetColor(color.White)
	dc.DrawStringAnchored("PAGO RECIBIDO", 1023, 128, 0.5, 0.5)

	dc.SetColor(receiptMuted)
	dc.SetFontFace(receiptFontSize(bold, gobold.TTF, 22))
	dc.DrawString("COMPROBANTE DE PAGO", 108, 238)
	dc.SetColor(receiptNavy)
	dc.SetFontFace(receiptFontSize(bold, gobold.TTF, 64))
	dc.DrawString(formatReceiptMoney(snapshot.AmountMinor, snapshot.Currency), 108, 330)
	dc.SetColor(receiptMuted)
	dc.SetFontFace(receiptFontSize(regular, goregular.TTF, 25))
	dc.DrawString(snapshot.ReceiptNumber, 108, 377)

	drawReceiptLine(dc, 108, 425, receiptWidth-216)
	drawReceiptField(dc, bold, regular, "Pagó", snapshot.PayerName, 108, 488)
	drawReceiptField(dc, bold, regular, "Fecha", formatReceiptDate(snapshot.ReceivedAt), 108, 570)
	drawReceiptField(dc, bold, regular, "Método", receiptMethodLabel(snapshot.Method), 650, 488)
	reference := "Sin referencia"
	if snapshot.Reference != nil && strings.TrimSpace(*snapshot.Reference) != "" {
		reference = *snapshot.Reference
	}
	drawReceiptField(dc, bold, regular, "Referencia", reference, 650, 570)

	y := 690.0
	dc.SetColor(receiptNavy)
	dc.SetFontFace(receiptFontSize(bold, gobold.TTF, 24))
	dc.DrawString("CONCEPTOS APLICADOS", 108, y)
	y += 54
	if len(snapshot.Allocations) == 0 {
		dc.SetColor(color.RGBA{R: 247, G: 244, B: 236, A: 255})
		dc.DrawRoundedRectangle(108, y, receiptWidth-216, 104, 18)
		dc.Fill()
		dc.SetColor(receiptMuted)
		dc.SetFontFace(receiptFontSize(regular, goregular.TTF, 25))
		dc.DrawString("Pago registrado sin cargos aplicados.", 136, y+62)
		y += 138
	} else {
		for _, allocation := range snapshot.Allocations {
			dc.SetColor(receiptNavy)
			dc.SetFontFace(receiptFontSize(bold, gobold.TTF, 27))
			description := truncateReceiptText(allocation.Description, 44)
			dc.DrawString(description, 108, y+34)
			dc.SetColor(receiptMuted)
			dc.SetFontFace(receiptFontSize(regular, goregular.TTF, 22))
			dc.DrawString(truncateReceiptText(allocation.HouseholdName, 54), 108, y+68)
			dc.SetColor(receiptNavy)
			dc.SetFontFace(receiptFontSize(bold, gobold.TTF, 28))
			dc.DrawStringAnchored(formatReceiptMoney(allocation.AmountMinor, snapshot.Currency), 1132, y+48, 1, 0.5)
			y += 96
			drawReceiptLine(dc, 108, y, receiptWidth-216)
			y += 24
		}
	}

	y += 22
	drawReceiptTotal(dc, bold, regular, "Total aplicado", snapshot.AllocatedMinor, snapshot.Currency, y)
	y += 62
	drawReceiptTotal(dc, bold, regular, "Saldo sin aplicar", snapshot.UnallocatedMinor, snapshot.Currency, y)
	y += 104

	dc.SetColor(color.RGBA{R: 234, G: 245, B: 230, A: 255})
	dc.DrawRoundedRectangle(108, y, receiptWidth-216, 96, 20)
	dc.Fill()
	dc.SetColor(receiptGreen)
	dc.SetFontFace(receiptFontSize(bold, gobold.TTF, 28))
	dc.DrawStringAnchored("Gracias por confiar en nosotros", receiptWidth/2, y+48, 0.5, 0.5)
	y += 150

	dc.SetColor(receiptMuted)
	dc.SetFontFace(receiptFontSize(regular, goregular.TTF, 20))
	dc.DrawStringAnchored("Comprobante interno de pago. No es factura ni CFDI.", receiptWidth/2, y, 0.5, 0.5)
	dc.DrawStringAnchored("Emitido por Pawsear", receiptWidth/2, y+34, 0.5, 0.5)

	var output bytes.Buffer
	if err := dc.EncodePNG(&output); err != nil {
		return nil, fmt.Errorf("encode receipt PNG: %w", err)
	}
	return output.Bytes(), nil
}

func renderPaymentReceiptPDF(snapshot paymentReceiptSnapshot) ([]byte, error) {
	png, err := renderPaymentReceiptPNG(snapshot)
	if err != nil {
		return nil, err
	}
	heightMM := 210 * float64(receiptCanvasHeight(len(snapshot.Allocations))) / receiptWidth
	pdf := fpdf.NewCustom(&fpdf.InitType{
		OrientationStr: "P", UnitStr: "mm", Size: fpdf.SizeType{Wd: 210, Ht: heightMM},
	})
	pdf.SetTitle("Comprobante "+snapshot.ReceiptNumber, false)
	pdf.SetAuthor("Pawsear", false)
	pdf.AddPage()
	options := fpdf.ImageOptions{ImageType: "PNG", ReadDpi: false}
	pdf.RegisterImageOptionsReader("receipt", options, bytes.NewReader(png))
	pdf.ImageOptions("receipt", 0, 0, 210, heightMM, false, options, 0, "")
	var output bytes.Buffer
	if err := pdf.Output(&output); err != nil {
		return nil, fmt.Errorf("encode receipt PDF: %w", err)
	}
	return output.Bytes(), nil
}

func receiptCanvasHeight(allocationCount int) int {
	return int(math.Max(1754, float64(1450+allocationCount*120)))
}

func receiptFont(contents []byte, size float64) (font.Face, error) {
	parsed, err := truetype.Parse(contents)
	if err != nil {
		return nil, fmt.Errorf("parse receipt font: %w", err)
	}
	return truetype.NewFace(parsed, &truetype.Options{Size: size}), nil
}

func receiptFontSize(_ font.Face, contents []byte, size float64) font.Face {
	face, err := receiptFont(contents, size)
	if err != nil {
		return nil
	}
	return face
}

func drawReceiptField(dc *gg.Context, bold font.Face, regular font.Face, label string, value string, x float64, y float64) {
	dc.SetColor(receiptMuted)
	dc.SetFontFace(receiptFontSize(bold, gobold.TTF, 20))
	dc.DrawString(strings.ToUpper(label), x, y)
	dc.SetColor(receiptNavy)
	dc.SetFontFace(receiptFontSize(regular, goregular.TTF, 27))
	dc.DrawString(truncateReceiptText(value, 34), x, y+38)
}

func drawReceiptTotal(dc *gg.Context, bold font.Face, regular font.Face, label string, amount int64, currency string, y float64) {
	dc.SetColor(receiptMuted)
	dc.SetFontFace(receiptFontSize(regular, goregular.TTF, 25))
	dc.DrawString(label, 108, y)
	dc.SetColor(receiptNavy)
	dc.SetFontFace(receiptFontSize(bold, gobold.TTF, 29))
	dc.DrawStringAnchored(formatReceiptMoney(amount, currency), 1132, y, 1, 0.5)
}

func drawReceiptLine(dc *gg.Context, x float64, y float64, width float64) {
	dc.SetColor(receiptLine)
	dc.SetLineWidth(2)
	dc.DrawLine(x, y, x+width, y)
	dc.Stroke()
}

func formatReceiptMoney(amountMinor int64, currency string) string {
	amount := amountMinor
	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}
	whole := fmt.Sprintf("%d", amount/100)
	for index := len(whole) - 3; index > 0; index -= 3 {
		whole = whole[:index] + "," + whole[index:]
	}
	return fmt.Sprintf("%s$%s.%02d %s", sign, whole, amount%100, currency)
}

func formatReceiptDate(value string) string {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return value
	}
	location, err := time.LoadLocation("America/Mexico_City")
	if err == nil {
		parsed = parsed.In(location)
	}
	months := []string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"}
	return fmt.Sprintf("%d de %s de %d · %02d:%02d", parsed.Day(), months[parsed.Month()-1], parsed.Year(), parsed.Hour(), parsed.Minute())
}

func receiptMethodLabel(method string) string {
	switch method {
	case "bank_transfer":
		return "Transferencia bancaria"
	case "cash":
		return "Efectivo"
	case "card_external":
		return "Tarjeta externa"
	default:
		return "Otro método"
	}
}

func truncateReceiptText(value string, maxRunes int) string {
	runes := []rune(strings.TrimSpace(value))
	if len(runes) <= maxRunes {
		return string(runes)
	}
	return string(runes[:maxRunes-1]) + "…"
}
