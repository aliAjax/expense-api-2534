package service

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"expense-api/internal/model"
	"expense-api/internal/repository"
)

type CategoryService struct {
	repository *repository.CategoryRepository
}

func NewCategoryService(categoryRepository *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repository: categoryRepository}
}

func (s *CategoryService) Create(ctx context.Context, request model.CreateCategoryRequest) (*model.Category, error) {
	name := strings.TrimSpace(request.Name)
	if name == "" {
		return nil, model.NewValidationError("name is required")
	}
	if utf8.RuneCountInString(name) > 50 {
		return nil, model.NewValidationError("name must not exceed 50 characters")
	}

	createdAt := time.Now().UTC().Format(time.RFC3339)
	category, err := s.repository.Create(ctx, name, createdAt)
	if err != nil {
		if repository.IsUniqueViolation(err) {
			return nil, model.NewConflictError("category already exists")
		}
		return nil, model.NewInternalError("failed to create category")
	}

	return category, nil
}

func (s *CategoryService) List(ctx context.Context) ([]model.Category, error) {
	categories, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, model.NewInternalError("failed to list categories")
	}
	return categories, nil
}

func (s *CategoryService) Update(ctx context.Context, id int64, request model.UpdateCategoryRequest) (*model.Category, error) {
	name := strings.TrimSpace(request.Name)
	if name == "" {
		return nil, model.NewValidationError("name is required")
	}
	if utf8.RuneCountInString(name) > 50 {
		return nil, model.NewValidationError("name must not exceed 50 characters")
	}

	category, err := s.repository.Update(ctx, id, name)
	if err != nil {
		return nil, model.NewInternalError("failed to update category")
	}

	return category, nil
}

func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		if repository.IsForeignKeyViolation(err) {
			return model.NewConflictError("category is in use by expenses")
		}
		if repository.IsNotFound(err) {
			return model.NewNotFoundError("category not found")
		}
		return model.NewInternalError("failed to delete category")
	}
	return nil
}
