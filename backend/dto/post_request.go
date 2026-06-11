package dto

type CreatePostRequest struct {
	Content string `json:"content"`
	TopicID int    `json:"topic_id"`
	UserID  int    `json:"user_id"`
}
