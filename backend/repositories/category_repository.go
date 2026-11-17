package repositories

import (
	"database/sql"
	"projet-forum/backend/models"
)

type CategoryRepository struct {
	DB *sql.DB
}

// GetAllCategories récupère la liste complète des catégories (les franchises NBA) en base de données.
func (r *CategoryRepository) GetAllCategories() ([]models.Category, error) {
	// Exécution de la requête SQL directe
	rows, err := r.DB.Query("SELECT id, name, description FROM Categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Sécurité : on s'assure de fermer la connexion à la fin

	var categories []models.Category

	// On parcourt chaque ligne retournée par la base de données
	for rows.Next() {
		var cat models.Category
		// On remplit notre structure 'cat' avec les données de la ligne actuelle
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description); err != nil {
			return nil, err
		}
		// On l'ajoute à notre liste finale
		categories = append(categories, cat)
	}

	return categories, nil
}
