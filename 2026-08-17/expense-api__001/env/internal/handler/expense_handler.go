package handler

import (
	"net/http"
	"strconv"
	"strings"

	"expense-api/internal/model"
	"expense-api/internal/service"
)

type ExpenseHandler struct {
	service *service.ExpenseService
}

func NewExpenseHandler(expenseService *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: expenseService}
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request model.CreateExpenseRequest
	if !decodeJSON(w, r, &request) {
		return
	}

	expense, err := h.service.Create(r.Context(), request)
	if err != nil {
		writeError(w, err)
		return
	}
	writeSuccess(w, http.StatusCreated, expense)
}

func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
	filter, ok := parseExpenseFilter(w, r)
	if !ok {
		return
	}

	expenses, err := h.service.List(r.Context(), filter)
	if err != nil {
		writeError(w, err)
		return
	}
	writeSuccess(w, http.StatusOK, expenses)
}

func (h *ExpenseHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	expense, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeSuccess(w, http.StatusOK, expense)
}

func (h *ExpenseHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var request model.UpdateExpenseRequest
	if !decodeJSON(w, r, &request) {
		return
	}

	expense, err := h.service.Update(r.Context(), id, request)
	if err != nil {
		writeError(w, err)
		return
	}
	writeSuccess(w, http.StatusOK, expense)
}

func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}

	writeSuccess(w, http.StatusOK, map[string]any{
		"id":      id,
		"deleted": true,
	})
}

func parseExpenseFilter(w http.ResponseWriter, r *http.Request) (model.ExpenseFilter, bool) {
	var filter model.ExpenseFilter

	if raw := strings.TrimSpace(r.URL.Query().Get("category_id")); raw != "" {
		categoryID, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || categoryID <= 0 {
			writeError(w, model.NewValidationError("category_id must be a positive integer"))
			return filter, false
		}
		filter.CategoryID = categoryID
	}

	filter.Category = strings.TrimSpace(r.URL.Query().Get("category"))
	return filter, true
}
