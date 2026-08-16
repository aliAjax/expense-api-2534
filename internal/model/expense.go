package model

type Expense struct {
	ID            int64   `json:"id"`
	AmountCents   int64   `json:"amount_cents"`
	Amount        float64 `json:"amount"`
	CategoryID    int64   `json:"category_id"`
	CategoryName  string  `json:"category"`
	Date          string  `json:"date"`
	PaymentMethod string  `json:"payment_method"`
	Note          string  `json:"note"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type CreateExpenseRequest struct {
	Amount        float64 `json:"amount"`
	CategoryID    int64   `json:"category_id"`
	Date          string  `json:"date"`
	PaymentMethod string  `json:"payment_method"`
	Note          string  `json:"note"`
}

type UpdateExpenseRequest struct {
	Amount        float64 `json:"amount"`
	CategoryID    int64   `json:"category_id"`
	Date          string  `json:"date"`
	PaymentMethod string  `json:"payment_method"`
	Note          string  `json:"note"`
}

type ExpenseFilter struct {
	CategoryID int64
	Category   string
}
