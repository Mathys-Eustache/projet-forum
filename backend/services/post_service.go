package services

import (
	"errors"
	"strings"

	"projet-forum/backend/dto"
	"projet-forum/backend/models"
	"projet-forum/backend/repositories"
)

type PostService struct {
	repo *repositories.PostRepository
}

func InitPostService(repo *repositories.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) GetPostsByTopic(topicID int) ([]models.Post, error) {
	return s.repo.GetPostsByTopic(topicID)
}

func (s *PostService) CreatePost(req dto.CreatePostRequest) (int, error) {
	if strings.TrimSpace(req.Content) == "" {
		return 0, errors.New("le contenu du message est obligatoire")
	}

	if req.TopicID <= 0 || req.UserID <= 0 {
		return 0, errors.New("identifiants de sujet ou utilisateur invalides")
	}

	newPost := models.Post{
		Content: req.Content,
		TopicID: req.TopicID,
		UserID:  req.UserID,
	}

	return s.repo.CreatePost(newPost)
}
