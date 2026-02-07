package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
	"time"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) TodayReport() (*models.TodayReport, error) {
	return s.repo.TodayReport()
}

func (s *ReportService) ReportByDate(startDate, endDate time.Time) (*models.TodayReport, error) {
	return s.repo.ReportByDate(startDate, endDate)
}
