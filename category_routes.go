package main

import (
	"encoding/json"
	"net/http"
)

func setupCategoryRoutes() {
	// GET /categories/{id}, PUT /categories/{id}, DELETE /categories/{id}
	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoryByID(w, r)
		} else if r.Method == "PUT" {
			updateCategory(w, r)
		} else if r.Method == "DELETE" {
			deleteCategory(w, r)
		}
	})

	// GET /categories, POST /categories
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
}