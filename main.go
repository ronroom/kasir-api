package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/models"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	PORT   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

var db *sql.DB

func loadConfig() Config {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		PORT:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}
	return config
}

func createTablesAndData() {
	if db == nil {
		return
	}

	// Create tables
	categoryTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT
	);`

	productTable := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		price INTEGER NOT NULL,
		stock INTEGER NOT NULL,
		category_id INTEGER REFERENCES categories(id)
	);`

	_, err := db.Exec(categoryTable)
	if err != nil {
		fmt.Printf("Failed to create categories table: %v\n", err)
		return
	}

	_, err = db.Exec(productTable)
	if err != nil {
		fmt.Printf("Failed to create products table: %v\n", err)
		return
	}

	// Add category_id column if not exists
	_, err = db.Exec("ALTER TABLE products ADD COLUMN IF NOT EXISTS category_id INTEGER REFERENCES categories(id)")
	if err != nil {
		fmt.Printf("Failed to add category_id column: %v\n", err)
		return
	}

	fmt.Println("Database tables created successfully")

	// Insert sample data
	var count int
	db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)

	if count == 0 {
		products := []models.Product{
			{Name: "Indomie", Price: 3500, Stock: 10, CategoryID: 1},
			{Name: "Vit 1000ml", Price: 3000, Stock: 40, CategoryID: 2},
			{Name: "Kecap", Price: 12000, Stock: 20, CategoryID: 1},
		}

		for _, prod := range products {
			_, err := db.Exec("INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4)", prod.Name, prod.Price, prod.Stock, prod.CategoryID)
			if err != nil {
				fmt.Printf("Failed to insert product %s: %v\n", prod.Name, err)
			}
		}
		fmt.Println("Sample products inserted")
	}

	// Insert sample categories
	db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count)
	if count == 0 {
		categories := []models.Category{
			{Name: "Makanan", Description: "Kategori untuk produk makanan"},
			{Name: "Minuman", Description: "Kategori untuk produk minuman"},
			{Name: "Elektronik", Description: "Kategori untuk produk elektronik"},
		}

		for _, cat := range categories {
			_, err := db.Exec("INSERT INTO categories (name, description) VALUES ($1, $2)", cat.Name, cat.Description)
			if err != nil {
				fmt.Printf("Failed to insert category %s: %v\n", cat.Name, err)
			}
		}
		fmt.Println("Sample categories inserted")
	}
}

func main() {
	config := loadConfig()

	// Debug: Print loaded config
	fmt.Printf("Loaded config - PORT: %s, DB_CONN: %s\n", config.PORT, config.DBConn)

	// Setup database
	if config.DBConn != "" {
		var err error
		db, err = database.InitDB(config.DBConn)
		if err != nil {
			log.Fatal("Failed to initialize database:", err)
		}
		defer db.Close()

		// Create tables and insert sample data
		createTablesAndData()
	} else {
		fmt.Println("No database connection configured, running without database")
	}

	// Dependency Injection
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Category routes with dependency injection
	http.HandleFunc("/categories/", categoryHandler.HandleCategoryByID)
	http.HandleFunc("/categories", categoryHandler.HandleCategories)

	// Product routes with dependency injection
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)
	http.HandleFunc("/api/produk", productHandler.HandleProducts)

	// Root endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Kasir API is running",
			"config": map[string]string{
				"port":         config.PORT,
				"db_connected": fmt.Sprintf("%t", db != nil),
			},
			"endpoints": []string{
				"GET /health",
				"GET /api/produk",
				"POST /api/produk",
				"GET /api/produk/{id}",
				"PUT /api/produk/{id}",
				"DELETE /api/produk/{id}",
				"GET /categories",
				"POST /categories",
				"GET /categories/{id}",
				"PUT /categories/{id}",
				"DELETE /categories/{id}",
			},
		})
	})

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "Ok",
			"message": "API Running"})
	})

	port := config.PORT
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
	}

	addr := "0.0.0.0:" + port
	fmt.Println("Server running di", addr)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("gagal running server", err)
	}
}
