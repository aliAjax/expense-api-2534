package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"expense-api/internal/model"
)

const maxRequestBodyBytes = 1 << 20

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("write json response: %v", err)
	}
}

func writeSuccess(w http.ResponseWriter, status int, data any) {
	writeJSON(w, status, model.Success(data))
}

func writeError(w http.ResponseWriter, err error) {
	var apiErr *model.APIError
	if errors.As(err, &apiErr) {
		writeJSON(w, apiErr.Status, model.Failure(apiErr))
		return
	}

	log.Printf("internal error: %v", err)
	writeJSON(w, http.StatusInternalServerError, model.Failure(model.NewInternalError("internal server error")))
}

func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		writeError(w, model.NewValidationError("invalid JSON body"))
		return false
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeError(w, model.NewValidationError("request body must contain a single JSON object"))
		return false
	}

	return true
}

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, model.NewValidationError("id must be a positive integer"))
		return 0, false
	}
	return id, true
}
