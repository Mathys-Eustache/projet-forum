package services

import (
	"projet-forum/backend/dto"
	"projet-forum/backend/models"
	"projet-forum/backend/repositories"
)

type TopicService struct {
	Repo *repositories.TopicRepository
}

func InitTopicService(repo *repositories.TopicRepository) *TopicService {
	return &TopicService{Repo: repo}
}

func (s *TopicService) CreateTopic(req dto.CreateTopicRequest, pseudo string) error {
	var userID int
	err := s.Repo.DB.QueryRow("SELECT id FROM Users WHERE username = ?", pseudo).Scan(&userID)
	if err != nil {
		return err
	}

	topic := models.Topic{
		Title:      req.Title,
		Content:    req.Content,
		UserID:     userID,
		CategoryID: req.CategoryID,
	}
	return s.Repo.CreateTopic(topic)
}

func (s *TopicService) GetAllTopics(limit int, offset int, search string, sortBy string) ([]repositories.TopicResponse, error) {
	return s.Repo.GetAllTopics(limit, offset, search, sortBy)
}

func (s *TopicService) GetTopicsByCategory(categoryID int, limit int, offset int, search string, sortBy string) ([]repositories.TopicResponse, error) {
	return s.Repo.GetTopicsByCategory(categoryID, limit, offset, search, sortBy)
}

func (s *TopicService) DeleteTopic(id int, pseudo string) error {
	return s.Repo.DeleteTopic(id, pseudo)
}

func (s *TopicService) UpdateTopic(id int, content string, pseudo string) error {
	return s.Repo.UpdateTopic(id, content, pseudo)
}

func (s *TopicService) UpdateTopicStatus(id int, status string, pseudo string) error {
	return s.Repo.UpdateTopicStatus(id, status, pseudo)
}

func (s *TopicService) ReactTopic(id int, action string, pseudo string) error {
	return s.Repo.ReactTopic(id, action, pseudo)
}
