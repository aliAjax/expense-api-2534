package service

import (
	"strings"
	"time"
	"unicode/utf8"

	"expense-api/internal/model"
)

func validateExpenseFields(amount float64, categoryID int64, date, paymentMethod, note string) error {
	if err := model.ValidateAmount(amount); err != nil {
		return err
	}
	if categoryID <= 0 {
		return model.NewValidationError("category_id must be greater than 0")
	}
	if err := validateDate(date); err != nil {
		return err
	}

	paymentMethod = strings.TrimSpace(paymentMethod)
	if paymentMethod == "" {
		return model.NewValidationError("payment_method is required")
	}
	if utf8.RuneCountInString(paymentMethod) > 50 {
		return model.NewValidationError("payment_method must not exceed 50 characters")
	}

	if utf8.RuneCountInString(note) > 500 {
		return model.NewValidationError("note must not exceed 500 characters")
	}

	return nil
}

func validateDate(date string) error {
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return model.NewValidationError("date must use YYYY-MM-DD format")
	}
	if parsed.Format("2006-01-02") != date {
		return model.NewValidationError("date must use YYYY-MM-DD format")
	}
	return nil
}

func validateMonth(month string) error {
	parsed, err := time.Parse("2006-01", month)
	if err != nil {
		return model.NewValidationError("month must use YYYY-MM format")
	}
	if parsed.Format("2006-01") != month {
		return model.NewValidationError("month must use YYYY-MM format")
	}
	return nil
}
