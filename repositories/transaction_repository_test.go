package repositories

import (
	"kasir-api/models"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateTransaction_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewTransactionRepository(db)

	items := []models.CheckoutItem{
		{ProductID: 1, Quantity: 2},
	}

	mock.ExpectBegin()

	// Mock product query
	rows := sqlmock.NewRows([]string{"name", "price", "stock"}).
		AddRow("Test Product", 1000, 10)
	mock.ExpectQuery("SELECT name, price, stock FROM products WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	// Mock update stock
	mock.ExpectExec("UPDATE products SET stock = stock - \\$1 WHERE id = \\$2 AND stock >= \\$1").
		WithArgs(2, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock insert transaction
	mock.ExpectQuery("INSERT INTO transactions").
		WithArgs(2000, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Mock insert transaction details
	mock.ExpectQuery("INSERT INTO transaction_details").
		WithArgs(1, 1, 2, 2000).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectCommit()

	tx, err := repo.CreateTransaction(items, false)
	if err != nil {
		t.Errorf("error was not expected while creating transaction: %s", err)
	}

	if tx == nil {
		t.Errorf("expected transaction, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetSalesSummary_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewTransactionRepository(db)
	now := time.Now()

	// Mock Revenue Query
	mock.ExpectQuery("SELECT COUNT\\(\\*\\), COALESCE\\(SUM\\(total_amount\\), 0\\) FROM transactions").
		WithArgs(now, now).
		WillReturnRows(sqlmock.NewRows([]string{"count", "revenue"}).AddRow(5, 50000))

	// Mock Best Seller Query
	// Note: We use regexp for complex query matching
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.name, COALESCE(SUM(td.quantity), 0) as total_qty FROM transaction_details`)).
		WithArgs(now, now).
		WillReturnRows(sqlmock.NewRows([]string{"name", "total_qty"}).AddRow("Best Product", 10))

	summary, err := repo.GetSalesSummary(now, now)
	if err != nil {
		t.Errorf("error was not expected while getting summary: %s", err)
	}

	if summary.TotalRevenue != 50000 {
		t.Errorf("expected revenue 50000, got %d", summary.TotalRevenue)
	}
	if summary.ProdukTerlaris.Nama != "Best Product" {
		t.Errorf("expected best seller 'Best Product', got '%s'", summary.ProdukTerlaris.Nama)
	}
}
