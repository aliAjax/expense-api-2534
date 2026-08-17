package handler

import (
	"net/http"
	"strings"

	"expense-api/internal/model"
	"expense-api/internal/service"
)

type SummaryHandler struct {
	service *service.SummaryService
}

func NewSummaryHandler(summaryService *service.SummaryService) *SummaryHandler {
	return &SummaryHandler{service: summaryService}
}

func (h *SummaryHandler) Daily(w http.ResponseWriter, r *http.Request) {
	date := strings.TrimSpace(r.URL.Query().Get("date"))
	if date == "" {
		writeError(w, model.NewValidationError("date query parameter is required"))
		return
	}

	summary, err := h.service.Daily(r.Context(), date)
	if err != nil {
		writeError(w, err)
		return
	}
	writeSuccess(w, http.StatusOK, summary)
}

func (h *SummaryHandler) Monthly(w http.ResponseWriter, r *http.Request) {
	month := strings.TrimSpace(r.URL.Query().Get("month"))
	if month == "" {
		writeError(w, model.NewValidationError("month query parameter is required"))
		return
	}

	summary, err := h.service.Monthly(r.Context(), month)
	if err != nil {
		writeError(w, err)
		return
	}
	writeSuccess(w, http.StatusOK, summary)
}
