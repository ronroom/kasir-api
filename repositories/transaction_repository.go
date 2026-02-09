package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"time"
)

type TransactionRepository interface {
	CreateTransaction(items []models.CheckoutItem, useLock bool) (*models.Transaction, error)
	GetSalesSummary(startDate, endDate time.Time) (*models.SalesSummary, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *transactionRepository {
	return &transactionRepository{db: db}
}

func (repo *transactionRepository) CreateTransaction(items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
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

		query := "SELECT name, price, stock FROM products WHERE id = $1"
		if useLock {
			query += " FOR UPDATE"
		}

		err := tx.QueryRow(query, item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		if stock < item.Quantity {
			return nil, fmt.Errorf("product %s (id: %d) has insufficient stock", productName, item.ProductID)
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		// atomic update with check
		res, err := tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}

		if rowsAffected == 0 {
			// This might happen if race condition occurred and stock wasn't locked, or if stock changed between read and update
			return nil, fmt.Errorf("failed to update stock for product %s (id: %d), possibly insufficient stock", productName, item.ProductID)
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount, created_at) VALUES ($1, $2) RETURNING id", totalAmount, models.GetCurrentTime()).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = transactionID
		err = tx.QueryRow("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING id",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal).Scan(&details[i].ID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		CreatedAt:   models.GetCurrentTime(),
		Details:     details,
	}, nil
}

func (repo *transactionRepository) GetSalesSummary(startDate, endDate time.Time) (*models.SalesSummary, error) {
	var summary models.SalesSummary

	// 1. Total Revenue & Count
	queryRevenue := "SELECT COUNT(*), COALESCE(SUM(total_amount), 0) FROM transactions WHERE created_at BETWEEN $1 AND $2"
	err := repo.db.QueryRow(queryRevenue, startDate, endDate).Scan(&summary.TotalTransaksi, &summary.TotalRevenue)
	if err != nil {
		return nil, err
	}

	// 2. Best Selling Product
	queryBestSeller := `
		SELECT p.name, COALESCE(SUM(td.quantity), 0) as total_qty
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE t.created_at BETWEEN $1 AND $2
		GROUP BY p.name
		ORDER BY total_qty DESC
		LIMIT 1`

	err = repo.db.QueryRow(queryBestSeller, startDate, endDate).Scan(&summary.ProdukTerlaris.Nama, &summary.ProdukTerlaris.QtyTerjual)
	if err == sql.ErrNoRows {
		// No transactions yet, which is fine
		summary.ProdukTerlaris = models.BestSellingProd{Nama: "-", QtyTerjual: 0}
	} else if err != nil {
		return nil, err
	}

	return &summary, nil
}
