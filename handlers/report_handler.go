package handlers

import (
	"encoding/json"
	"fmt"
	"kasir-api/services"
	"net/http"
	"time"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) TodayReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	transaction, err := h.service.TodayReport()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func (h *ReportHandler) ReportByDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startDate, err := h.formDateParam(r.URL.Query().Get("start_date"))
	if err != nil {
		http.Error(w, "format start_date tidak sesuai", http.StatusBadRequest)
		return
	}

	endDate, err := h.formDateParam(r.URL.Query().Get("end_date"))
	if err != nil {
		http.Error(w, "format end_date tidak sesuai", http.StatusBadRequest)
		return
	}

	if endDate.Before(startDate) {
		http.Error(w, "range date tidak sesuai", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.ReportByDate(startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func (h *ReportHandler) formDateParam(dateParam string) (time.Time, error) {
	// Define the layout (Standard YYYY-MM-DD)
	layout := "2006-01-02"

	// Parse the string
	t, err := time.Parse(layout, dateParam)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return t, err
	}
	return t, nil
}
