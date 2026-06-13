package services

import "projet-forum/backend/repositories"

type PostService struct {
	Repo *repositories.PostRepository
}

func InitPostService(repo *repositories.PostRepository) *PostService {
	return &PostService{Repo: repo}
}

func (s *PostService) GetPostsByTopic(topicID int, limit int, offset int) ([]repositories.PostResponse, error) {
	return s.Repo.GetPostsByTopic(topicID, limit, offset)
}

func (s *PostService) CreatePost(content string, topicID int, pseudo string) error {
	var userID int
	err := s.Repo.DB.QueryRow("SELECT id FROM Users WHERE username = ?", pseudo).Scan(&userID)
	if err != nil {
		return err
	}
	return s.Repo.CreatePost(content, topicID, userID)
}

func (s *PostService) DeletePost(id int, pseudo string) error {
	return s.Repo.DeletePost(id, pseudo)
}
