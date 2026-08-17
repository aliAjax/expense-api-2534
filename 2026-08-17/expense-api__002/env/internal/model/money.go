package model

import "math"

const maxSupportedAmount = 9.0e15

func ValidateAmount(amount float64) error {
	if math.IsNaN(amount) || math.IsInf(amount, 0) {
		return NewValidationError("amount must be a finite number")
	}
	if amount <= 0 {
		return NewValidationError("amount must be greater than 0")
	}
	if amount > maxSupportedAmount {
		return NewValidationError("amount is too large")
	}

	scaled := amount * 100
	rounded := math.Round(scaled)
	if math.Abs(scaled-rounded) > 1e-7 {
		return NewValidationError("amount must have at most 2 decimal places")
	}

	return nil
}

func AmountToCents(amount float64) int64 {
	return int64(math.Round(amount * 100))
}

func CentsToAmount(cents int64) float64 {
	return float64(cents) / 100
}
