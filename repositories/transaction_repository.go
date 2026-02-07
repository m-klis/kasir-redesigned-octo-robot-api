package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"strings"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM product WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE product SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transaction (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	err = bulkInsertDetails(tx, transactionID, details)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		CreatedAt:   time.Now(),
		Details:     details,
	}, nil
}

// function for bulk insert details
func bulkInsertDetails(tx *sql.Tx, transactionID int, details []models.TransactionDetail) error {
	const batchSize = 500

	for i := 0; i < len(details); i += batchSize {
		end := i + batchSize
		if end > len(details) {
			end = len(details)
		}

		// Process the current chunk
		chunk := details[i:end]

		valueStrings := make([]string, 0, len(chunk))
		valueArgs := make([]interface{}, 0, len(chunk)*4)

		for j, d := range chunk {
			n := j * 4
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", n+1, n+2, n+3, n+4))
			valueArgs = append(valueArgs, transactionID, d.ProductID, d.Quantity, d.Subtotal)
		}

		stmt := fmt.Sprintf("INSERT INTO transaction_detail (transaction_id, product_id, quantity, subtotal) VALUES %s",
			strings.Join(valueStrings, ","))

		if _, err := tx.Exec(stmt, valueArgs...); err != nil {
			return err
		}
	}
	return nil
}
