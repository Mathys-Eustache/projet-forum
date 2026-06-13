package repositories

import (
	"database/sql"
	"errors"
)

type PostRepository struct {
	DB *sql.DB
}

type PostResponse struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	CreatedAt string `json:"created_at"`
}

func (r *PostRepository) GetPostsByTopic(topicID int, limit int, offset int) ([]PostResponse, error) {
	query := `
		SELECT p.id, p.content, u.username, DATE_FORMAT(p.created_at, '%d/%m/%Y %H:%i')
		FROM Posts p 
		JOIN Users u ON p.user_id = u.id 
		WHERE p.topic_id = ? 
		ORDER BY p.created_at ASC
		LIMIT ? OFFSET ?
	`
	rows, err := r.DB.Query(query, topicID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []PostResponse
	for rows.Next() {
		var p PostResponse
		if err := rows.Scan(&p.ID, &p.Content, &p.Author, &p.CreatedAt); err != nil {
			continue
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) CreatePost(content string, topicID int, userID int) error {
	_, err := r.DB.Exec("INSERT INTO Posts (content, topic_id, user_id) VALUES (?, ?, ?)", content, topicID, userID)
	return err
}

func (r *PostRepository) DeletePost(id int, pseudo string) error {
	var author string
	err := r.DB.QueryRow("SELECT u.username FROM Posts p JOIN Users u ON p.user_id = u.id WHERE p.id = ?", id).Scan(&author)
	if err != nil {
		return err
	}
	if author != pseudo {
		return errors.New("forbidden")
	}
	_, err = r.DB.Exec("DELETE FROM Posts WHERE id = ?", id)
	return err
}
