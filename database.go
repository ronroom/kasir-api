package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	var err error
	
	// Get DATABASE_URL from environment (Railway provides this)
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// Fallback for local development
		databaseURL = "postgres://localhost/kasir_db?sslmode=disable"
	}

	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Connected to PostgreSQL database")
	createTables()
}

func createTables() {
	// Create categories table
	categoryTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT
	);`

	// Create products table
	productTable := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		nama VARCHAR(100) NOT NULL,
		harga INTEGER NOT NULL,
		stok INTEGER NOT NULL
	);`

	_, err := db.Exec(categoryTable)
	if err != nil {
		log.Fatal("Failed to create categories table:", err)
	}

	_, err = db.Exec(productTable)
	if err != nil {
		log.Fatal("Failed to create products table:", err)
	}

	fmt.Println("Database tables created successfully")
	insertSampleData()
}

func insertSampleData() {
	// Insert sample categories if table is empty
	var count int
	db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count)
	
	if count == 0 {
		categories := []Category{
			{Name: "Makanan", Description: "Kategori untuk produk makanan"},
			{Name: "Minuman", Description: "Kategori untuk produk minuman"},
			{Name: "Elektronik", Description: "Kategori untuk produk elektronik"},
			{Name: "Pakaian", Description: "Kategori untuk produk pakaian dan fashion"},
			{Name: "Kesehatan", Description: "Kategori untuk produk kesehatan dan obat-obatan"},
			{Name: "Olahraga", Description: "Kategori untuk peralatan dan perlengkapan olahraga"},
		}

		for _, cat := range categories {
			_, err := db.Exec("INSERT INTO categories (name, description) VALUES ($1, $2)", cat.Name, cat.Description)
			if err != nil {
				log.Printf("Failed to insert category %s: %v", cat.Name, err)
			}
		}
		fmt.Println("Sample categories inserted")
	}

	// Insert sample products if table is empty
	db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	
	if count == 0 {
		products := []Produk{
			{Nama: "Indomie", Harga: 3500, Stok: 10},
			{Nama: "Vit 1000ml", Harga: 3000, Stok: 40},
			{Nama: "Kecap", Harga: 12000, Stok: 20},
		}

		for _, prod := range products {
			_, err := db.Exec("INSERT INTO products (nama, harga, stok) VALUES ($1, $2, $3)", prod.Nama, prod.Harga, prod.Stok)
			if err != nil {
				log.Printf("Failed to insert product %s: %v", prod.Nama, err)
			}
		}
		fmt.Println("Sample products inserted")
	}
}