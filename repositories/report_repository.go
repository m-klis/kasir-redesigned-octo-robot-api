package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) TodayReport() (*models.TodayReport, error) {
	todayReport := &models.TodayReport{}

	query := `
			SELECT 
			COUNT(*) AS total_transactions,
			COALESCE(SUM(total_amount), 0) AS total_revenue
			FROM transaction
			WHERE created_at >= CURRENT_DATE;
		`

	err := repo.db.QueryRow(query).Scan(&todayReport.TotalTransaction, &todayReport.TotalRevenue)
	if err == sql.ErrNoRows {
		return nil, errors.New("report tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	query = `
			SELECT
				p.name,
				SUM(td.quantity) AS total_sold
			FROM transaction_detail td
			JOIN product p ON td.product_id = p.id
			JOIN transaction t ON td.transaction_id = t.id
			WHERE t.created_at >= CURRENT_DATE
			GROUP BY p.id, p.name
			ORDER BY total_sold DESC
			LIMIT 1
		`

	err = repo.db.QueryRow(query).Scan(&todayReport.BestSelling.Name, &todayReport.BestSelling.Quantity)
	if err == sql.ErrNoRows {
		return nil, errors.New("produk terlaris tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return todayReport, nil
}
