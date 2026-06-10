package controllers

import (
	"encoding/json"
	"net/http"
	"projet-forum/backend/services"
)

type CategoryController struct {
	Service *services.CategoryService
}

func (c *CategoryController) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := c.Service.GetAllCategories()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des catégories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(categories)
}
