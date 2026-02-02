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
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Config struct {
	PORT   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

var db *sql.DB



var categories = []Category{
	{ID: 1, Name: "Makanan", Description: "Kategori untuk produk makanan"},
	{ID: 2, Name: "Minuman", Description: "Kategori untuk produk minuman"},
	{ID: 3, Name: "Elektronik", Description: "Kategori untuk produk elektronik"},
	{ID: 4, Name: "Pakaian", Description: "Kategori untuk produk pakaian dan fashion"},
	{ID: 5, Name: "Kesehatan", Description: "Kategori untuk produk kesehatan dan obat-obatan"},
	{ID: 6, Name: "Olahraga", Description: "Kategori untuk peralatan dan perlengkapan olahraga"},
}

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
		stock INTEGER NOT NULL
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

	fmt.Println("Database tables created successfully")
	
	// Insert sample data
	var count int
	db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	
	if count == 0 {
		products := []models.Product{
			{Name: "Indomie", Price: 3500, Stock: 10},
			{Name: "Vit 1000ml", Price: 3000, Stock: 40},
			{Name: "Kecap", Price: 12000, Stock: 20},
		}

		for _, prod := range products {
			_, err := db.Exec("INSERT INTO products (name, price, stock) VALUES ($1, $2, $3)", prod.Name, prod.Price, prod.Stock)
			if err != nil {
				fmt.Printf("Failed to insert product %s: %v\n", prod.Name, err)
			}
		}
		fmt.Println("Sample products inserted")
	}
}

func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, c := range categories {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
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
	
	// Category routes
	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoryByID(w, r)
		}
	})

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			json.NewEncoder(w).Encode(categories)
		}
	})

	// Product routes with dependency injection
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)
	http.HandleFunc("/api/produk", productHandler.HandleProducts)

	// Root endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Kasir API is running",
			"config": map[string]string{
				"port": config.PORT,
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