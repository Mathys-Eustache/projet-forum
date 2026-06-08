package repositories

import (
	"database/sql"
	"fmt"
	"projet-forum/backend/models"
)

type UserRepository struct {
	db *sql.DB
}

func InitUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(u models.User) (int, error) {
	query := "INSERT INTO Users (username, email, password) VALUES (?, ?, ?);"
	result, err := r.db.Exec(query, u.Username, u.Email, u.Password)
	if err != nil {
		return -1, fmt.Errorf("erreur lors de l'insertion: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}
