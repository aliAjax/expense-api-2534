package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"expense-api/internal/model"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) Create(ctx context.Context, expense *model.Expense) (int64, error) {
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO expenses
			(amount_cents, category_id, expense_date, payment_method, note, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		expense.AmountCents,
		expense.CategoryID,
		expense.Date,
		expense.PaymentMethod,
		expense.Note,
		expense.CreatedAt,
		expense.UpdatedAt,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get expense id: %w", err)
	}
	return id, nil
}

func (r *ExpenseRepository) FindAll(ctx context.Context, filter model.ExpenseFilter) ([]model.Expense, error) {
	query := `SELECT
		e.id, e.amount_cents, e.category_id, c.name, e.expense_date,
		e.payment_method, e.note, e.created_at, e.updated_at
		FROM expenses e
		JOIN categories c ON c.id = e.category_id`

	where := make([]string, 0, 2)
	args := make([]any, 0, 2)

	if filter.CategoryID > 0 {
		where = append(where, "e.category_id = ?")
		args = append(args, filter.CategoryID)
	} else if strings.TrimSpace(filter.Category) != "" {
		where = append(where, "c.name = ? COLLATE NOCASE")
		args = append(args, strings.TrimSpace(filter.Category))
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY e.expense_date DESC, e.id DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := make([]model.Expense, 0)
	for rows.Next() {
		expense, err := scanExpense(rows)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, *expense)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *ExpenseRepository) FindByID(ctx context.Context, id int64) (*model.Expense, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT
			e.id, e.amount_cents, e.category_id, c.name, e.expense_date,
			e.payment_method, e.note, e.created_at, e.updated_at
		 FROM expenses e
		 JOIN categories c ON c.id = e.category_id
		 WHERE e.id = ?`,
		id,
	)

	return scanExpense(row)
}

func (r *ExpenseRepository) Update(ctx context.Context, id int64, expense *model.Expense) error {
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE expenses
		 SET amount_cents = ?, category_id = ?, expense_date = ?, payment_method = ?, note = ?, updated_at = ?
		 WHERE id = ?`,
		expense.AmountCents,
		expense.CategoryID,
		expense.Date,
		expense.PaymentMethod,
		expense.Note,
		expense.UpdatedAt,
		id,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *ExpenseRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM expenses WHERE id = ?`, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

type expenseScanner interface {
	Scan(dest ...any) error
}

func scanExpense(scanner expenseScanner) (*model.Expense, error) {
	var expense model.Expense
	if err := scanner.Scan(
		&expense.ID,
		&expense.AmountCents,
		&expense.CategoryID,
		&expense.CategoryName,
		&expense.Date,
		&expense.PaymentMethod,
		&expense.Note,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	); err != nil {
		return nil, err
	}

	expense.Amount = model.CentsToAmount(expense.AmountCents)
	return &expense, nil
}
