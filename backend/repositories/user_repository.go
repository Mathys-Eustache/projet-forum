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

func (r *UserRepository) GetUserByLogin(login string) (models.User, error) {
	var u models.User
	query := "SELECT id, username, email, password, role FROM Users WHERE username = ? OR email = ?;"
	err := r.db.QueryRow(query, login, login).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role)
	return u, err
}
