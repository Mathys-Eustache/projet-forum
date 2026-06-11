package models

type Post struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	TopicID   int    `json:"topic_id"`
	UserID    int    `json:"user_id"`
	CreatedAt string `json:"created_at"`
}
