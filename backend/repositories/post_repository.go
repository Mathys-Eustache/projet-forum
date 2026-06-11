package repositories

import (
	"database/sql"
	"projet-forum/backend/models"
)

type PostRepository struct {
	DB *sql.DB
}

func (r *PostRepository) GetPostsByTopic(topicID int) ([]models.Post, error) {
	rows, err := r.DB.Query("SELECT id, content, topic_id, user_id, created_at FROM Posts WHERE topic_id = ?", topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.Content, &p.TopicID, &p.UserID, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (r *PostRepository) CreatePost(p models.Post) (int, error) {
	query := "INSERT INTO Posts (content, topic_id, user_id) VALUES (?, ?, ?)"
	result, err := r.DB.Exec(query, p.Content, p.TopicID, p.UserID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
