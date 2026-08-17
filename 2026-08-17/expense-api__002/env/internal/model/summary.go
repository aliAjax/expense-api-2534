package model

type CategorySummary struct {
	CategoryID int64   `json:"category_id"`
	Category   string  `json:"category"`
	TotalCents int64   `json:"total_cents"`
	Total      float64 `json:"total"`
	Count      int64   `json:"count"`
}

type DailySummary struct {
	Date       string            `json:"date"`
	TotalCents int64             `json:"total_cents"`
	Total      float64           `json:"total"`
	Count      int64             `json:"count"`
	ByCategory []CategorySummary `json:"by_category"`
}

type MonthlySummary struct {
	Month      string            `json:"month"`
	TotalCents int64             `json:"total_cents"`
	Total      float64           `json:"total"`
	Count      int64             `json:"count"`
	ByCategory []CategorySummary `json:"by_category"`
}
