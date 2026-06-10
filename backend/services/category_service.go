package services

import (
	"projet-forum/backend/models"
	"projet-forum/backend/repositories"
)

type CategoryService struct {
	Repo *repositories.CategoryRepository
}

func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	return s.Repo.GetAllCategories()
}
