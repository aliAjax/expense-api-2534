package handler

import (
	"net/http"

	"expense-api/internal/model"
	"expense-api/internal/service"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: categoryService}
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request model.CreateCategoryRequest
	if !decodeJSON(w, r, &request) {
		return
	}

	category, err := h.service.Create(r.Context(), request)
	if err != nil {
		writeError(w, err)
		return
	}

	writeSuccess(w, http.StatusCreated, category)
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeSuccess(w, http.StatusOK, categories)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var request model.UpdateCategoryRequest
	if !decodeJSON(w, r, &request) {
		return
	}

	category, err := h.service.Update(r.Context(), id, request)
	if err != nil {
		writeError(w, err)
		return
	}
	writeSuccess(w, http.StatusOK, category)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
