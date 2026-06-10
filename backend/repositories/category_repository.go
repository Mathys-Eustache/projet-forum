package repositories

import (
	"database/sql"
	"projet-forum/backend/models"
)

type CategoryRepository struct {
	DB *sql.DB
}

func (r *CategoryRepository) GetAllCategories() ([]models.Category, error) {
	rows, err := r.DB.Query("SELECT id, name, description FROM Categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	return categories, nil
}
