package service

import (
	"context"

	"expense-api/internal/model"
	"expense-api/internal/repository"
)

type SummaryService struct {
	repository *repository.SummaryRepository
}

func NewSummaryService(summaryRepository *repository.SummaryRepository) *SummaryService {
	return &SummaryService{repository: summaryRepository}
}

func (s *SummaryService) Daily(ctx context.Context, date string) (*model.DailySummary, error) {
	if err := validateDate(date); err != nil {
		return nil, err
	}

	summary, err := s.repository.SummaryByDate(ctx, date)
	if err != nil {
		return nil, model.NewInternalError("failed to summarize daily expenses")
	}
	return summary, nil
}

func (s *SummaryService) Monthly(ctx context.Context, month string) (*model.MonthlySummary, error) {
	if err := validateMonth(month); err != nil {
		return nil, err
	}

	summary, err := s.repository.SummaryByMonth(ctx, month)
	if err != nil {
		return nil, model.NewInternalError("failed to summarize monthly expenses")
	}
	return summary, nil
}
