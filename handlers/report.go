package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"santori/linkchecker/models"

	"github.com/jung-kurt/gofpdf"
)

func GenerateReport(w http.ResponseWriter, r *http.Request) {
	var req models.PDFReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request")
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 12)
	pdf.AddPage()
	pdf.Cell(0, 10, fmt.Sprintf("PDF with info"))
	pdf.Ln(12)

	models.Mu.Lock()
	defer models.Mu.Unlock()

	for _, num := range req.LinksNum {
		if group, ok := models.LinksStorage[num-1]; ok {
			pdf.Cell(0, 8, fmt.Sprintf("Request %d:", num))
			pdf.Ln(8)
			for url, status := range group {
				pdf.Cell(0, 6, fmt.Sprintf(" - %s : %s", url, status))
				pdf.Ln(6)
			}
			pdf.Ln(4)
		} else {
			pdf.Cell(0, 6, fmt.Sprintf("Request %d: doesn't exist", num))
			pdf.Ln(6)
		}
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=report.pdf")

	if err := pdf.Output(w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error generating pdf")
		return
	}
}
