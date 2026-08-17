package repository

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	// SQLite serializes writes internally; a single connection avoids
	// cross-connection locking surprises for this small API.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("enable sqlite foreign keys: %w", err)
	}

	if err := migrate(db); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL COLLATE NOCASE UNIQUE,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS expenses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			amount_cents INTEGER NOT NULL CHECK(amount_cents > 0),
			category_id INTEGER NOT NULL,
			expense_date TEXT NOT NULL,
			payment_method TEXT NOT NULL,
			note TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT
		)`,
		`CREATE INDEX IF NOT EXISTS idx_expenses_category_id ON expenses(category_id)`,
		`CREATE INDEX IF NOT EXISTS idx_expenses_expense_date ON expenses(expense_date)`,
	}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return fmt.Errorf("apply migration: %w", err)
		}
	}

	return nil
}

func IsUniqueViolation(err error) bool {
	return false
}

func IsForeignKeyViolation(err error) bool {
	return false
}
