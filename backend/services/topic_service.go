package services

import (
	"errors"
	"strings"

	"projet-forum/backend/dto"
	"projet-forum/backend/models"
	"projet-forum/backend/repositories"
)

type TopicService struct {
	repo *repositories.TopicRepository
}

func InitTopicService(repo *repositories.TopicRepository) *TopicService {
	return &TopicService{repo: repo}
}

func (s *TopicService) GetAllTopics() ([]models.Topic, error) {
	return s.repo.GetAllTopics()
}

func (s *TopicService) CreateTopic(req dto.CreateTopicRequest) (int, error) {
	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Content) == "" {
		return 0, errors.New("le titre et le contenu sont obligatoires")
	}

	if req.CategoryID <= 0 || req.UserID <= 0 {
		return 0, errors.New("identifiants de categorie ou utilisateur invalides")
	}

	newTopic := models.Topic{
		Title:      req.Title,
		Content:    req.Content,
		CategoryID: req.CategoryID,
		UserID:     req.UserID,
	}

	return s.repo.CreateTopic(newTopic)
}
