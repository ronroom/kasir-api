package services

import (
	"fmt"
	"kasir-api/models"
	"kasir-api/repositories"
	"time"
)

type TransactionService struct {
	repo repositories.TransactionRepository
}

func NewTransactionService(repo repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items, useLock)
}

func (s *TransactionService) GetDailyReport() (*models.SalesSummary, error) {
	// Start of Day
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// End of Day (until next day 00:00:00 basically, or just use now if we want real-time up to now)
	endOfDay := startOfDay.Add(24 * time.Hour)

	return s.repo.GetSalesSummary(startOfDay, endOfDay)
}

func (s *TransactionService) GetReport(startDateStr, endDateStr string) (*models.SalesSummary, error) {
	layout := "2006-01-02"
	startDate, err := time.ParseInLocation(layout, startDateStr, time.Local)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format (YYYY-MM-DD)")
	}

	endDate, err := time.ParseInLocation(layout, endDateStr, time.Local)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date format (YYYY-MM-DD)")
	}

	// Adjust endDate to include the full day
	endDate = endDate.Add(24 * time.Hour)

	return s.repo.GetSalesSummary(startDate, endDate)
}
