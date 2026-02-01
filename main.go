package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	_ "github.com/lib/pq"
)

type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

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

var produk = []Produk{
	{ID: 1, Nama: "Indomie", Harga: 3500, Stok: 10},
	{ID: 2, Nama: "Vit 1000ml", Harga: 3000, Stok: 40},
	{ID: 3, Nama: "Kecap", Harga: 12000, Stok: 20},
}

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
	
	var config Config
	viper.Unmarshal(&config)
	return config
}

func initDB(connectionString string) {
	var err error
	
	// Use provided connection string or fallback to DATABASE_URL
	databaseURL := connectionString
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}
	if databaseURL == "" {
		// Skip database for local development without PostgreSQL
		fmt.Println("No database connection configured, running without database")
		return
	}

	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Printf("Failed to ping database: %v", err)
		fmt.Println("Running without database connection")
		return
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
		name VARCHAR(100) NOT NULL,
		price INTEGER NOT NULL,
		stock INTEGER NOT NULL
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
		products := []models.Product{
			{Name: "Indomie", Price: 3500, Stock: 10},
			{Name: "Vit 1000ml", Price: 3000, Stock: 40},
			{Name: "Kecap", Price: 12000, Stock: 20},
		}

		for _, prod := range products {
			_, err := db.Exec("INSERT INTO products (name, price, stock) VALUES ($1, $2, $3)", prod.Name, prod.Price, prod.Stock)
			if err != nil {
				log.Printf("Failed to insert product %s: %v", prod.Name, err)
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

func updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, c := range categories {
		if c.ID == id {
			var updatedCategory Category
			err := json.NewDecoder(r.Body).Decode(&updatedCategory)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			categories[i].Name = updatedCategory.Name
			categories[i].Description = updatedCategory.Description

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories[i])
			return
		}
	}
	http.Error(w, "Category not found", http.StatusNotFound)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, c := range categories {
		if c.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Category not found", http.StatusNotFound)
}

func getProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Produk not found", http.StatusNotFound)
}

func updateProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, p := range produk {
		if p.ID == id {
			var updatedProduk Produk
			err := json.NewDecoder(r.Body).Decode(&updatedProduk)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			produk[i].Nama = updatedProduk.Nama
			produk[i].Harga = updatedProduk.Harga
			produk[i].Stok = updatedProduk.Stok

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk[i])
			return
		}
	}
	http.Error(w, "Produk not found", http.StatusNotFound)
}

func deleteProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, p := range produk {
		if p.ID == id {
			produk = append(produk[:i], produk[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Produk not found", http.StatusNotFound)
}

func main() {
	config := loadConfig()
	
	// Initialize database
	initDB(config.DBConn)
	defer db.Close()
	
	// Initialize repositories
	productRepo := repositories.NewProductRepository(db)
	
	// Initialize services
	productService := services.NewProductService(productRepo)
	
	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService)
	
	// Setup routes with new handlers (only if database is connected)
	if db != nil {
		http.HandleFunc("/api/produk", productHandler.HandleProducts)
		http.HandleFunc("/api/produk/", productHandler.HandleProductByID)
	} else {
		fmt.Println("Database handlers disabled - no database connection")
	}
	
	// Category routes (keep old implementation for now)
	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoryByID(w, r)
		} else if r.Method == "PUT" {
			updateCategory(w, r)
		} else if r.Method == "DELETE" {
			deleteCategory(w, r)
		}
	})

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			json.NewEncoder(w).Encode(categories)
		} else if r.Method == "POST" {
			var newCategory Category
			err := json.NewDecoder(r.Body).Decode(&newCategory)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			
			newCategory.ID = len(categories) + 1
			categories = append(categories, newCategory)

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newCategory)
		}
	})

	// Root endpoint for debugging
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Kasir API is running",
			"platform": "Multi-Platform (Railway + Zeabur)",
			"endpoints": []string{
				"GET /health",
				"GET /api/produk",
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
		port = os.Getenv("ZEABUR_PORT")
	}
	if port == "" {
		port = os.Getenv("SERVER_PORT")
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