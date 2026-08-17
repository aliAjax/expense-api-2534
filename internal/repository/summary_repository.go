package repository

import (
	"context"
	"database/sql"
	"time"

	"expense-api/internal/model"
)

type SummaryRepository struct {
	db *sql.DB
}

func NewSummaryRepository(db *sql.DB) *SummaryRepository {
	return &SummaryRepository{db: db}
}

func (r *SummaryRepository) SummaryByDate(ctx context.Context, date string) (*model.DailySummary, error) {
	summary := &model.DailySummary{
		Date:       date,
		ByCategory: make([]model.CategorySummary, 0),
	}

	if err := r.db.QueryRow(
		`SELECT COALESCE(SUM(amount_cents), 0), COUNT(*)
		 FROM expenses WHERE expense_date = ?`,
		date,
	).Scan(&summary.TotalCents, &summary.Count); err != nil {
		return nil, err
	}
	summary.Total = model.CentsToAmount(summary.TotalCents)

	breakdown, err := r.categoryBreakdown(ctx, `e.expense_date = ?`, date)
	if err != nil {
		return nil, err
	}
	summary.ByCategory = breakdown

	return summary, nil
}

func (r *SummaryRepository) SummaryByMonth(ctx context.Context, month string) (*model.MonthlySummary, error) {
	start, err := time.Parse("2006-01", month)
	if err != nil {
		return nil, err
	}
	startText := start.Format("2006-01-02")
	endText := start.AddDate(0, 1, 0).Format("2006-01-02")

	summary := &model.MonthlySummary{
		Month:      month,
		ByCategory: make([]model.CategorySummary, 0),
	}

	if err := r.db.QueryRow(
		`SELECT COALESCE(SUM(amount_cents), 0), COUNT(*)
		 FROM expenses WHERE expense_date >= ? AND expense_date < ?`,
		startText,
		endText,
	).Scan(&summary.TotalCents, &summary.Count); err != nil {
		return nil, err
	}
	summary.Total = model.CentsToAmount(summary.TotalCents)

	breakdown, err := r.categoryBreakdown(
		ctx,
		`e.expense_date >= ? AND e.expense_date < ?`,
		startText,
		endText,
	)
	if err != nil {
		return nil, err
	}
	summary.ByCategory = breakdown

	return summary, nil
}

func (r *SummaryRepository) categoryBreakdown(ctx context.Context, where string, args ...any) ([]model.CategorySummary, error) {
	query := `SELECT
		c.id, c.name, COALESCE(SUM(e.amount_cents), 0), COUNT(*)
		FROM categories c
		JOIN expenses e ON e.category_id = c.id
		WHERE ` + where + `
		GROUP BY c.id, c.name
		ORDER BY c.name COLLATE NOCASE`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.CategorySummary, 0)
	for rows.Next() {
		var item model.CategorySummary
		if err := rows.Scan(&item.CategoryID, &item.Category, &item.TotalCents, &item.Count); err != nil {
			return nil, err
		}
		item.Total = model.CentsToAmount(item.TotalCents)
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
