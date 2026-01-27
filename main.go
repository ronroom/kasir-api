package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	// Category routes
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

	// Product routes
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProdukByID(w, r)
		} else if r.Method == "PUT" {
			updateProduk(w, r)
		} else if r.Method == "DELETE" {
			deleteProduk(w, r)
		}
	})

	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			json.NewEncoder(w).Encode(produk)
		} else if r.Method == "POST" {
			var newProduk Produk
			err := json.NewDecoder(r.Body).Decode(&newProduk)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			
			newProduk.ID = len(produk) + 1
			produk = append(produk, newProduk)

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newProduk)
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

	port := os.Getenv("PORT")
	if port == "" {
		// Try different port environment variables for various platforms
		port = os.Getenv("ZEABUR_PORT")
	}
	if port == "" {
		port = os.Getenv("SERVER_PORT")
	}
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s\n", port)
	fmt.Printf("Environment PORT: %s\n", os.Getenv("PORT"))
	fmt.Printf("Environment ZEABUR_PORT: %s\n", os.Getenv("ZEABUR_PORT"))
	fmt.Println("Endpoints:")
	fmt.Println("GET /api/produk - Get all products")
	fmt.Println("GET /categories - Get all categories")
	fmt.Println("GET /health - Health check")
	
	err := http.ListenAndServe("0.0.0.0:"+port, nil)
	if err != nil {
		fmt.Println("gagal running server!")
	}
}