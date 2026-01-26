package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{ID: 1, Name: "Makanan", Description: "Kategori untuk produk makanan"},
	{ID: 2, Name: "Minuman", Description: "Kategori untuk produk minuman"},
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