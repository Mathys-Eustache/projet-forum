package models

type Topic struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CategoryID int    `json:"category_id"`
	UserID     int    `json:"user_id"`
	CreatedAt  string `json:"created_at"`
}
