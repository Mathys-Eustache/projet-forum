package repositories

import (
	"database/sql"
	"projet-forum/backend/models"
)

type TopicRepository struct {
	DB *sql.DB
}

func (r *TopicRepository) GetAllTopics() ([]models.Topic, error) {
	rows, err := r.DB.Query("SELECT id, title, content, category_id, user_id, created_at FROM Topics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []models.Topic

	for rows.Next() {
		var t models.Topic
		if err := rows.Scan(&t.ID, &t.Title, &t.Content, &t.CategoryID, &t.UserID, &t.CreatedAt); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}

	return topics, nil
}

func (r *TopicRepository) CreateTopic(t models.Topic) (int, error) {
	query := "INSERT INTO Topics (title, content, category_id, user_id) VALUES (?, ?, ?, ?)"
	result, err := r.DB.Exec(query, t.Title, t.Content, t.CategoryID, t.UserID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
