package repositories

import (
	"database/sql"
	"errors"
	"projet-forum/backend/models"
)

type TopicRepository struct {
	DB *sql.DB
}

type TopicResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
}

func (r *TopicRepository) CreateTopic(topic models.Topic) error {
	_, err := r.DB.Exec("INSERT INTO Topics (title, content, user_id, category_id, status) VALUES (?, ?, ?, ?, 'ouvert')",
		topic.Title, topic.Content, topic.UserID, topic.CategoryID)
	return err
}

func (r *TopicRepository) NewMethod() {}

func (r *TopicRepository) GetAllTopics(limit int, offset int, search string) ([]TopicResponse, error) {
	var rows *sql.Rows
	var err error

	if search != "" {
		searchParam := "%" + search + "%"
		rows, err = r.DB.Query("SELECT t.id, t.title, t.content, u.username, DATE_FORMAT(t.created_at, '%d/%m/%Y %H:%i'), t.status FROM Topics t JOIN Users u ON t.user_id = u.id WHERE t.title LIKE ? OR t.content LIKE ? ORDER BY t.created_at DESC LIMIT ? OFFSET ?", searchParam, searchParam, limit, offset)
	} else {
		rows, err = r.DB.Query("SELECT t.id, t.title, t.content, u.username, DATE_FORMAT(t.created_at, '%d/%m/%Y %H:%i'), t.status FROM Topics t JOIN Users u ON t.user_id = u.id ORDER BY t.created_at DESC LIMIT ? OFFSET ?", limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []TopicResponse
	for rows.Next() {
		var t TopicResponse
		if err := rows.Scan(&t.ID, &t.Title, &t.Content, &t.Author, &t.CreatedAt, &t.Status); err != nil {
			continue
		}
		topics = append(topics, t)
	}
	return topics, nil
}

func (r *TopicRepository) GetTopicsByCategory(categoryID int, limit int, offset int, search string) ([]TopicResponse, error) {
	var rows *sql.Rows
	var err error

	if search != "" {
		searchParam := "%" + search + "%"
		rows, err = r.DB.Query("SELECT t.id, t.title, t.content, u.username, DATE_FORMAT(t.created_at, '%d/%m/%Y %H:%i'), t.status FROM Topics t JOIN Users u ON t.user_id = u.id WHERE t.category_id = ? AND (t.title LIKE ? OR t.content LIKE ?) ORDER BY t.created_at DESC LIMIT ? OFFSET ?", categoryID, searchParam, searchParam, limit, offset)
	} else {
		rows, err = r.DB.Query("SELECT t.id, t.title, t.content, u.username, DATE_FORMAT(t.created_at, '%d/%m/%Y %H:%i'), t.status FROM Topics t JOIN Users u ON t.user_id = u.id WHERE t.category_id = ? ORDER BY t.created_at DESC LIMIT ? OFFSET ?", categoryID, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []TopicResponse
	for rows.Next() {
		var t TopicResponse
		if err := rows.Scan(&t.ID, &t.Title, &t.Content, &t.Author, &t.CreatedAt, &t.Status); err != nil {
			continue
		}
		topics = append(topics, t)
	}
	return topics, nil
}

func (r *TopicRepository) DeleteTopic(id int, pseudo string) error {
	var author string
	err := r.DB.QueryRow("SELECT u.username FROM Topics t JOIN Users u ON t.user_id = u.id WHERE t.id = ?", id).Scan(&author)
	if err != nil {
		return err
	}

	if author != pseudo {
		return errors.New("forbidden")
	}

	_, err = r.DB.Exec("DELETE FROM Topics WHERE id = ?", id)
	return err
}

func (r *TopicRepository) UpdateTopic(id int, content string, pseudo string) error {
	var author string
	err := r.DB.QueryRow("SELECT u.username FROM Topics t JOIN Users u ON t.user_id = u.id WHERE t.id = ?", id).Scan(&author)
	if err != nil {
		return err
	}

	if author != pseudo {
		return errors.New("forbidden")
	}

	_, err = r.DB.Exec("UPDATE Topics SET content = ? WHERE id = ?", content, id)
	return err
}

func (r *TopicRepository) UpdateTopicStatus(id int, status string, pseudo string) error {
	var author string
	err := r.DB.QueryRow("SELECT u.username FROM Topics t JOIN Users u ON t.user_id = u.id WHERE t.id = ?", id).Scan(&author)
	if err != nil {
		return err
	}

	if author != pseudo {
		return errors.New("forbidden")
	}

	_, err = r.DB.Exec("UPDATE Topics SET status = ? WHERE id = ?", status, id)
	return err
}
