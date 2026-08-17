package service_test

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"expense-api/internal/model"
	"expense-api/internal/repository"
	"expense-api/internal/service"
)

func openServiceTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := repository.Open(filepath.Join(t.TempDir(), "expense.db"))
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func TestExpenseCreateListAndAmount(t *testing.T) {
	db := openServiceTestDB(t)
	categories := repository.NewCategoryRepository(db)
	expenses := repository.NewExpenseRepository(db)
	summary := repository.NewSummaryRepository(db)
	categoryService := service.NewCategoryService(categories)
	expenseService := service.NewExpenseService(expenses, categories)
	summaryService := service.NewSummaryService(summary)

	ctx := context.Background()
	category, err := categoryService.Create(ctx, model.CreateCategoryRequest{Name: "Food"})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}

	created, err := expenseService.Create(ctx, model.CreateExpenseRequest{
		Amount:        12.34,
		CategoryID:    category.ID,
		Date:          "2026-08-16",
		PaymentMethod: "cash",
		Note:          "lunch",
	})
	if err != nil {
		t.Fatalf("create expense: %v", err)
	}
	if created.AmountCents != 1234 {
		t.Fatalf("created amount cents = %d, want 1234", created.AmountCents)
	}
	if created.Amount != 12.34 {
		t.Fatalf("created amount = %v, want 12.34", created.Amount)
	}

	got, err := expenseService.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("get expense: %v", err)
	}
	if got.AmountCents != 1234 || got.Amount != 12.34 || got.CategoryName != "Food" {
		t.Fatalf("unexpected expense after get: %+v", got)
	}

	list, err := expenseService.List(ctx, model.ExpenseFilter{})
	if err != nil {
		t.Fatalf("list expenses: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("list length = %d, want 1", len(list))
	}
	if list[0].ID != created.ID || list[0].AmountCents != 1234 || list[0].Amount != 12.34 {
		t.Fatalf("unexpected listed expense: %+v", list[0])
	}

	daily, err := summaryService.Daily(ctx, "2026-08-16")
	if err != nil {
		t.Fatalf("daily summary: %v", err)
	}
	if daily.Count != 1 || daily.TotalCents != 1234 || daily.Total != 12.34 || len(daily.ByCategory) != 1 {
		t.Fatalf("unexpected daily summary: %+v", daily)
	}
}

func TestExpenseFilterByNameCaseInsensitive(t *testing.T) {
	db := openServiceTestDB(t)
	categories := repository.NewCategoryRepository(db)
	expenses := repository.NewExpenseRepository(db)
	categoryService := service.NewCategoryService(categories)
	expenseService := service.NewExpenseService(expenses, categories)

	ctx := context.Background()
	food, err := categoryService.Create(ctx, model.CreateCategoryRequest{Name: "Food"})
	if err != nil {
		t.Fatalf("create food category: %v", err)
	}
	transport, err := categoryService.Create(ctx, model.CreateCategoryRequest{Name: "Transport"})
	if err != nil {
		t.Fatalf("create transport category: %v", err)
	}

	if _, err := expenseService.Create(ctx, model.CreateExpenseRequest{
		Amount: 10, CategoryID: food.ID, Date: "2026-08-16", PaymentMethod: "cash", Note: "breakfast",
	}); err != nil {
		t.Fatalf("create food expense: %v", err)
	}
	if _, err := expenseService.Create(ctx, model.CreateExpenseRequest{
		Amount: 20, CategoryID: transport.ID, Date: "2026-08-16", PaymentMethod: "card", Note: "taxi",
	}); err != nil {
		t.Fatalf("create transport expense: %v", err)
	}

	list, err := expenseService.List(ctx, model.ExpenseFilter{Category: "food"})
	if err != nil {
		t.Fatalf("list by category: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("filtered list length = %d, want 1", len(list))
	}
	if list[0].CategoryName != "Food" {
		t.Fatalf("filtered category = %q, want Food", list[0].CategoryName)
	}
}

func TestExpenseUpdateAndGet(t *testing.T) {
	db := openServiceTestDB(t)
	categories := repository.NewCategoryRepository(db)
	expenses := repository.NewExpenseRepository(db)
	categoryService := service.NewCategoryService(categories)
	expenseService := service.NewExpenseService(expenses, categories)

	ctx := context.Background()
	food, err := categoryService.Create(ctx, model.CreateCategoryRequest{Name: "Food"})
	if err != nil {
		t.Fatalf("create food category: %v", err)
	}
	book, err := categoryService.Create(ctx, model.CreateCategoryRequest{Name: "Books"})
	if err != nil {
		t.Fatalf("create books category: %v", err)
	}
	created, err := expenseService.Create(ctx, model.CreateExpenseRequest{
		Amount: 30, CategoryID: food.ID, Date: "2026-08-16", PaymentMethod: "cash", Note: "snack",
	})
	if err != nil {
		t.Fatalf("create expense: %v", err)
	}

	updated, err := expenseService.Update(ctx, created.ID, model.UpdateExpenseRequest{
		Amount: 45.67, CategoryID: book.ID, Date: "2026-08-17", PaymentMethod: "card", Note: "book",
	})
	if err != nil {
		t.Fatalf("update expense: %v", err)
	}
	if updated.AmountCents != 4567 || updated.Amount != 45.67 || updated.CategoryName != "Books" || updated.Date != "2026-08-17" {
		t.Fatalf("unexpected updated expense: %+v", updated)
	}
	got, err := expenseService.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("get updated expense: %v", err)
	}
	if got.AmountCents != 4567 || got.Amount != 45.67 || got.CategoryID != book.ID {
		t.Fatalf("unexpected persisted expense: %+v", got)
	}
}

func TestSummaryDailyMonthlyByCategory(t *testing.T) {
	db := openServiceTestDB(t)
	categories := repository.NewCategoryRepository(db)
	expenses := repository.NewExpenseRepository(db)
	summary := repository.NewSummaryRepository(db)
	categoryService := service.NewCategoryService(categories)
	expenseService := service.NewExpenseService(expenses, categories)
	summaryService := service.NewSummaryService(summary)

	ctx := context.Background()
	food, err := categoryService.Create(ctx, model.CreateCategoryRequest{Name: "Food"})
	if err != nil {
		t.Fatalf("create food category: %v", err)
	}
	transport, err := categoryService.Create(ctx, model.CreateCategoryRequest{Name: "Transport"})
	if err != nil {
		t.Fatalf("create transport category: %v", err)
	}
	if _, err := expenseService.Create(ctx, model.CreateExpenseRequest{
		Amount: 12.5, CategoryID: food.ID, Date: "2026-08-16", PaymentMethod: "cash", Note: "lunch",
	}); err != nil {
		t.Fatalf("create food expense: %v", err)
	}
	if _, err := expenseService.Create(ctx, model.CreateExpenseRequest{
		Amount: 7.25, CategoryID: transport.ID, Date: "2026-08-16", PaymentMethod: "card", Note: "bus",
	}); err != nil {
		t.Fatalf("create transport expense: %v", err)
	}
	if _, err := expenseService.Create(ctx, model.CreateExpenseRequest{
		Amount: 3.25, CategoryID: food.ID, Date: "2026-08-17", PaymentMethod: "cash", Note: "coffee",
	}); err != nil {
		t.Fatalf("create next-day food expense: %v", err)
	}

	daily, err := summaryService.Daily(ctx, "2026-08-16")
	if err != nil {
		t.Fatalf("daily summary: %v", err)
	}
	if daily.Count != 2 || daily.TotalCents != 1975 || daily.Total != 19.75 {
		t.Fatalf("unexpected daily total: %+v", daily)
	}
	if len(daily.ByCategory) != 2 {
		t.Fatalf("daily by-category length = %d, want 2", len(daily.ByCategory))
	}

	monthly, err := summaryService.Monthly(ctx, "2026-08")
	if err != nil {
		t.Fatalf("monthly summary: %v", err)
	}
	if monthly.Count != 3 || monthly.TotalCents != 2300 || monthly.Total != 23 {
		t.Fatalf("unexpected monthly total: %+v", monthly)
	}
	if len(monthly.ByCategory) != 2 {
		t.Fatalf("monthly by-category length = %d, want 2", len(monthly.ByCategory))
	}
}

func TestCategoryUpdateDeleteNotFoundAndForeignKey(t *testing.T) {
	db := openServiceTestDB(t)
	categories := repository.NewCategoryRepository(db)
	expenses := repository.NewExpenseRepository(db)
	categoryService := service.NewCategoryService(categories)
	expenseService := service.NewExpenseService(expenses, categories)

	ctx := context.Background()
	created, err := categoryService.Create(ctx, model.CreateCategoryRequest{Name: "Old Name"})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}
	updated, err := categoryService.Update(ctx, created.ID, model.UpdateCategoryRequest{Name: "New Name"})
	if err != nil {
		t.Fatalf("update category: %v", err)
	}
	if updated.Name != "New Name" {
		t.Fatalf("updated name = %q, want New Name", updated.Name)
	}

	if _, err := categoryService.Update(ctx, 999999, model.UpdateCategoryRequest{Name: "Missing"}); err == nil {
		t.Fatal("update missing category: expected error, got nil")
	} else {
		var apiErr *model.APIError
		if !errors.As(err, &apiErr) || apiErr.Code != model.CodeNotFound {
			t.Fatalf("update missing category error = %T %v, want not found", err, err)
		}
	}

	if _, err := expenseService.Create(ctx, model.CreateExpenseRequest{
		Amount: 1, CategoryID: updated.ID, Date: "2026-08-16", PaymentMethod: "cash", Note: "x",
	}); err != nil {
		t.Fatalf("create expense for category: %v", err)
	}
	if err := categoryService.Delete(ctx, updated.ID); err == nil {
		t.Fatal("delete used category: expected error, got nil")
	}

	if err := categoryService.Delete(ctx, 999999); err == nil {
		t.Fatal("delete missing category: expected error, got nil")
	}

	unused, err := categoryService.Create(ctx, model.CreateCategoryRequest{Name: "Unused"})
	if err != nil {
		t.Fatalf("create unused category: %v", err)
	}
	if err := categoryService.Delete(ctx, unused.ID); err != nil {
		t.Fatalf("delete unused category: %v", err)
	}
}

func TestServiceContextCancellation(t *testing.T) {
	db := openServiceTestDB(t)
	categories := repository.NewCategoryRepository(db)
	expenses := repository.NewExpenseRepository(db)
	summary := repository.NewSummaryRepository(db)
	categoryService := service.NewCategoryService(categories)
	expenseService := service.NewExpenseService(expenses, categories)
	summaryService := service.NewSummaryService(summary)

	ctx := context.Background()
	category, err := categoryService.Create(ctx, model.CreateCategoryRequest{Name: "Food"})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}
	if _, err := expenseService.Create(ctx, model.CreateExpenseRequest{
		Amount: 5, CategoryID: category.ID, Date: "2026-08-16", PaymentMethod: "cash", Note: "x",
	}); err != nil {
		t.Fatalf("create expense: %v", err)
	}

	canceled, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := categoryService.List(canceled); err == nil {
		t.Fatal("list categories with canceled context: expected error, got nil")
	}
	if _, err := expenseService.List(canceled, model.ExpenseFilter{}); err == nil {
		t.Fatal("list expenses with canceled context: expected error, got nil")
	}
	if _, err := summaryService.Daily(canceled, "2026-08-16"); err == nil {
		t.Fatal("daily summary with canceled context: expected error, got nil")
	}

	select {
	case <-ctx.Done():
		t.Fatal("background context unexpectedly canceled")
	case <-time.After(10 * time.Millisecond):
	}
}
