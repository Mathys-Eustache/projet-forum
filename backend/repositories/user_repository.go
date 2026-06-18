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

// CreateUser gère l'insertion sécurisée d'un nouvel utilisateur lors de l'inscription
func (r *UserRepository) CreateUser(u models.User) (int, error) {
	// L'utilisation des "?" permet d'éviter les attaques par injection SQL
	query := "INSERT INTO Users (username, email, password) VALUES (?, ?, ?)"
	result, err := r.db.Exec(query, u.Username, u.Email, u.Password)
	if err != nil {
		return -1, fmt.Errorf("erreur lors de l'insertion en base: %w", err)
	}

	// Récupération de l'ID généré automatiquement par MySQL
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

// GetUserByLogin permet de retrouver un utilisateur soit par son pseudo, soit par son email
func (r *UserRepository) GetUserByLogin(login string) (models.User, error) {
	var u models.User
	query := "SELECT id, username, email, password, role FROM Users WHERE username = ? OR email = ?"

	// QueryRow est utilisé car on attend un seul résultat unique
	err := r.db.QueryRow(query, login, login).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role)

	return u, err
}
