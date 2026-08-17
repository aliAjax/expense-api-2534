package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"expense-api/internal/model"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, name, createdAt string) (*model.Category, error) {
	result, err := r.db.Exec(
		`INSERT INTO categories (name, created_at) VALUES (?, ?)`,
		name,
		createdAt,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("get category id: %w", err)
	}

	return &model.Category{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
	}, nil
}

func (r *CategoryRepository) FindAll(ctx context.Context) ([]model.Category, error) {
	rows, err := r.db.Query(`SELECT id, name, created_at FROM categories ORDER BY name COLLATE NOCASE`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]model.Category, 0)
	for rows.Next() {
		var category model.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) FindByID(ctx context.Context, id int64) (*model.Category, error) {
	var category model.Category
	err := r.db.QueryRow(
		`SELECT id, name, created_at FROM categories WHERE id = ?`,
		id,
	).Scan(&category.ID, &category.Name, &category.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) FindByName(ctx context.Context, name string) (*model.Category, error) {
	var category model.Category
	err := r.db.QueryRow(
		`SELECT id, name, created_at FROM categories WHERE name = ? COLLATE NOCASE`,
		name,
	).Scan(&category.ID, &category.Name, &category.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) Update(ctx context.Context, id int64, name string) (*model.Category, error) {
	result, err := r.db.Exec(`UPDATE categories SET name = ? WHERE id = ?`, name, id)
	if err != nil {
		return nil, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, sql.ErrNoRows
	}

	return r.FindByID(ctx, id)
}

func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.Exec(`DELETE FROM categories WHERE id = ?`, id)
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

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
