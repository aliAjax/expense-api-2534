package server

import (
	"net/http"

	"expense-api/internal/handler"
)

func NewRouter(expenseHandler *handler.ExpenseHandler, categoryHandler *handler.CategoryHandler, summaryHandler *handler.SummaryHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", health)

	mux.HandleFunc("POST /api/v1/categories", categoryHandler.Create)
	mux.HandleFunc("GET /api/v1/categories", categoryHandler.List)
	mux.HandleFunc("PUT /api/v1/categories/{id}", categoryHandler.Update)
	mux.HandleFunc("DELETE /api/v1/categories/{id}", categoryHandler.Delete)

	mux.HandleFunc("POST /api/v1/expenses", expenseHandler.Create)
	mux.HandleFunc("GET /api/v1/expenses", expenseHandler.List)
	mux.HandleFunc("GET /api/v1/expenses/{id}", expenseHandler.Get)
	mux.HandleFunc("PUT /api/v1/expenses/{id}", expenseHandler.Update)
	mux.HandleFunc("DELETE /api/v1/expenses/{id}", expenseHandler.Delete)

	mux.HandleFunc("GET /api/v1/summary/daily", summaryHandler.Daily)
	mux.HandleFunc("GET /api/v1/summary/monthly", summaryHandler.Monthly)

	return withLogging(mux)
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"code":0,"message":"success","data":{"status":"ok"}}`))
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
