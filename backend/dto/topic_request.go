package dto

type CreateTopicRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	CategoryID int    `json:"category_id"`
	UserID     int    `json:"user_id"`
}
