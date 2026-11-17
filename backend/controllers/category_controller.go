package controllers

import (
	"encoding/json"
	"net/http"
	"projet-forum/backend/services"
)

type CategoryController struct {
	Service *services.CategoryService
}

// GetCategories renvoie la liste de toutes les catégories (franchises)
func (c *CategoryController) GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Récupération des catégories depuis la base de données
	categories, err := c.Service.GetAllCategories()
	if err != nil {
		http.Error(w, `{"erreur": "Erreur lors de la récupération des catégories"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}
