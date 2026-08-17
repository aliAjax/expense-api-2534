package service

import (
	"context"
	"strings"
	"time"

	"expense-api/internal/model"
	"expense-api/internal/repository"
)

type ExpenseService struct {
	expenseRepository  *repository.ExpenseRepository
	categoryRepository *repository.CategoryRepository
}

func NewExpenseService(expenseRepository *repository.ExpenseRepository, categoryRepository *repository.CategoryRepository) *ExpenseService {
	return &ExpenseService{
		expenseRepository:  expenseRepository,
		categoryRepository: categoryRepository,
	}
}

func (s *ExpenseService) Create(ctx context.Context, request model.CreateExpenseRequest) (*model.Expense, error) {
	if err := validateExpenseFields(request.Amount, request.CategoryID, request.Date, request.PaymentMethod, request.Note); err != nil {
		return nil, err
	}

	category, err := s.categoryRepository.FindByID(ctx, request.CategoryID)
	if err != nil {
		if repository.IsNotFound(err) {
			return nil, model.NewNotFoundError("category not found")
		}
		return nil, model.NewInternalError("failed to find category")
	}

	now := time.Now().UTC().Format(time.RFC3339)
	expense := &model.Expense{
		AmountCents:   model.AmountToCents(request.Amount),
		CategoryID:    request.CategoryID,
		CategoryName:  category.Name,
		Date:          request.Date,
		PaymentMethod: strings.TrimSpace(request.PaymentMethod),
		Note:          strings.TrimSpace(request.Note),
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	id, err := s.expenseRepository.Create(ctx, expense)
	if err != nil {
		return nil, model.NewInternalError("failed to create expense")
	}

	expense.ID = id
	expense.Amount = model.CentsToAmount(expense.AmountCents)
	return expense, nil
}

func (s *ExpenseService) List(ctx context.Context, filter model.ExpenseFilter) ([]model.Expense, error) {
	expenses, err := s.expenseRepository.FindAll(ctx, filter)
	if err != nil {
		return nil, model.NewInternalError("failed to list expenses")
	}
	return expenses, nil
}

func (s *ExpenseService) Get(ctx context.Context, id int64) (*model.Expense, error) {
	expense, err := s.expenseRepository.FindByID(ctx, id)
	if err != nil {
		if repository.IsNotFound(err) {
			return nil, model.NewNotFoundError("expense not found")
		}
		return nil, model.NewInternalError("failed to get expense")
	}
	return expense, nil
}

func (s *ExpenseService) Update(ctx context.Context, id int64, request model.UpdateExpenseRequest) (*model.Expense, error) {
	if err := validateExpenseFields(request.Amount, request.CategoryID, request.Date, request.PaymentMethod, request.Note); err != nil {
		return nil, err
	}

	category, err := s.categoryRepository.FindByID(ctx, request.CategoryID)
	if err != nil {
		if repository.IsNotFound(err) {
			return nil, model.NewNotFoundError("category not found")
		}
		return nil, model.NewInternalError("failed to find category")
	}

	expense := &model.Expense{
		AmountCents:   model.AmountToCents(request.Amount),
		CategoryID:    request.CategoryID,
		CategoryName:  category.Name,
		Date:          request.Date,
		PaymentMethod: strings.TrimSpace(request.PaymentMethod),
		Note:          strings.TrimSpace(request.Note),
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
	}

	if err := s.expenseRepository.Update(ctx, id, expense); err != nil {
		if repository.IsNotFound(err) {
			return nil, model.NewNotFoundError("expense not found")
		}
		return nil, model.NewInternalError("failed to update expense")
	}

	updated, err := s.expenseRepository.FindByID(ctx, id)
	if err != nil {
		if repository.IsNotFound(err) {
			return nil, model.NewNotFoundError("expense not found")
		}
		return nil, model.NewInternalError("failed to get expense")
	}
	return updated, nil
}

func (s *ExpenseService) Delete(ctx context.Context, id int64) error {
	if err := s.expenseRepository.Delete(ctx, id); err != nil {
		if repository.IsNotFound(err) {
			return model.NewNotFoundError("expense not found")
		}
		return model.NewInternalError("failed to delete expense")
	}
	return nil
}
