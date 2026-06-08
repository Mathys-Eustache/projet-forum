package services

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"strings"
	"unicode"

	"projet-forum/backend/dto"
	"projet-forum/backend/models"
	"projet-forum/backend/repositories"
)

type AuthService struct {
	repo *repositories.UserRepository
}

func InitAuthService(repo *repositories.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(req dto.RegisterRequest) (int, error) {
	if len(req.Password) < 12 {
		return -1, errors.New("le mot de passe doit contenir au moins 12 caractères")
	}

	hasUpper := false
	hasSpecial := false
	for _, char := range req.Password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
		}
	}

	if !hasUpper || !hasSpecial {
		return -1, errors.New("le mot de passe doit contenir au moins une majuscule et un caractère spécial")
	}

	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Email) == "" {
		return -1, errors.New("le pseudo et l'email sont obligatoires")
	}

	hasher := sha512.New()
	hasher.Write([]byte(req.Password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	newUser := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	return s.repo.CreateUser(newUser)
}
