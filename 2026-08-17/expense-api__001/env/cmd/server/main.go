package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"expense-api/internal/config"
	"expense-api/internal/handler"
	"expense-api/internal/repository"
	"expense-api/internal/server"
	"expense-api/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.Load()

	if err := ensureDatabaseDirectory(cfg.DBPath); err != nil {
		return err
	}

	db, err := repository.Open(cfg.DBPath)
	if err != nil {
		return err
	}
	defer db.Close()

	categoryRepository := repository.NewCategoryRepository(db)
	expenseRepository := repository.NewExpenseRepository(db)
	summaryRepository := repository.NewSummaryRepository(db)

	categoryService := service.NewCategoryService(categoryRepository)
	expenseService := service.NewExpenseService(expenseRepository, categoryRepository)
	summaryService := service.NewSummaryService(summaryRepository)

	categoryHandler := handler.NewCategoryHandler(categoryService)
	expenseHandler := handler.NewExpenseHandler(expenseService)
	summaryHandler := handler.NewSummaryHandler(summaryService)

	httpServer := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           server.NewRouter(expenseHandler, categoryHandler, summaryHandler),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("expense-api listening on :%s", cfg.Port)
		serverErrors <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-ctx.Done():
		log.Println("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return err
	}

	if err := <-serverErrors; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func ensureDatabaseDirectory(path string) error {
	if path == "" || path == ":memory:" {
		return nil
	}

	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return nil
}
