package controllers

import (
	"encoding/json"
	"net/http"

	"projet-forum/backend/dto"
	"projet-forum/backend/services"
)

type AuthController struct {
	service *services.AuthService
}

func InitAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service: service}
}

// RegisterHandler gère l'inscription d'un nouvel utilisateur
func (c *AuthController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Configuration des en-têtes CORS pour autoriser le frontend
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Gestion de la requête de pré-vérification (Preflight)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"erreur": "Méthode non autorisée"}`, http.StatusMethodNotAllowed)
		return
	}

	// On décode les données envoyées par le frontend (JSON)
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"erreur": "Format des données invalide"}`, http.StatusBadRequest)
		return
	}

	// Appel du service pour créer l'utilisateur en base
	id, err := c.service.Register(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erreur": err.Error()})
		return
	}

	// Réponse de succès
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Inscription réussie",
		"id":      id,
	})
}

// LoginHandler gère la connexion d'un utilisateur et génère un token
func (c *AuthController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Configuration CORS
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"erreur": "Méthode non autorisée"}`, http.StatusMethodNotAllowed)
		return
	}

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"erreur": "Format des données invalide"}`, http.StatusBadRequest)
		return
	}

	// Vérification des identifiants et récupération du token
	token, username, role, err := c.service.Login(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"erreur": err.Error()})
		return
	}

	// Renvoi du token au frontend
	json.NewEncoder(w).Encode(map[string]string{
		"token":    token,
		"username": username,
		"role":     role,
	})
}
