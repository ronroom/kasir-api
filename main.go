package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

var produk = []Produk{
	{ID: 1, Nama: "Indomie", Harga: 3500, Stok: 10},
	{ID: 2, Nama: "Vit 1000ml", Harga: 3000, Stok: 40},
	{ID: 2, Nama: "Kecap", Harga: 12000, Stok: 20},
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

// PUT localhost:8080/api/produk/{id}
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
			w.WriteHeader(http.StatusNoContent) //204
			{
				json.NewEncoder(w).Encode(map[string]string{"message": "Produk deleted successfully"})
			}
			return
		}
	}
	http.Error(w, "Produk not found", http.StatusNotFound)
}

func main() {
	// GET localhost:8080/api/produk/{id}
	// PUT localhost:8080/api/produk/{id}
	// DELETE localhost:8080/api/produk/{id}
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProdukByID(w, r)
		} else if r.Method == "PUT" {
			updateProduk(w, r)
		} else if r.Method == "DELETE" {
			// delete logic here
			deleteProduk(w, r)
		}

	})

	// GET localhost:8080/api/produk
	// POST localhost:8080/api/produk
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			json.NewEncoder(w).Encode(produk)
		} else if r.Method == "POST" {
			// baca data dari request body
			var newProduk Produk
			err := json.NewDecoder(r.Body).Decode(&newProduk)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// masukin data ke variable Produk
			newProduk.ID = len(produk) + 1
			produk = append(produk, newProduk)

			w.WriteHeader(http.StatusCreated) //201
			json.NewEncoder(w).Encode(newProduk)
		}
	})

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "Ok",
			"message": "API Running"})
	})
	fmt.Println("Server running di localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server!")
	}
}
