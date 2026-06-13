package services

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"strings"
	"time"
	"unicode"

	"projet-forum/backend/dto"
	"projet-forum/backend/models"
	"projet-forum/backend/repositories"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte("cle_secrete_ultra_robuste_nba_forum_2026")

type AuthService struct {
	repo *repositories.UserRepository
}

func InitAuthService(repo *repositories.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(req dto.RegisterRequest) (int, error) {
	if len(req.Password) < 12 {
		return -1, errors.New("le mot de passe doit contenir au moins 12 caracteres")
	}

	hasUpper, hasSpecial := false, false
	for _, char := range req.Password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
		}
	}

	if !hasUpper || !hasSpecial {
		return -1, errors.New("le mot de passe doit contenir au moins une majuscule et un caractere special")
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

func (s *AuthService) Login(req dto.LoginRequest) (string, string, error) {
	u, err := s.repo.GetUserByLogin(req.Email)
	if err != nil {
		return "", "", errors.New("identifiants invalides")
	}

	hasher := sha512.New()
	hasher.Write([]byte(req.Password))
	if hex.EncodeToString(hasher.Sum(nil)) != u.Password {
		return "", "", errors.New("identifiants invalides")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  u.ID,
		"username": u.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", errors.New("erreur interne")
	}

	return tokenString, u.Username, nil
}
