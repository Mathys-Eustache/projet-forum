package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte("cle_secrete_ultra_robuste_nba_forum_2026")

// AuthMiddleware est un "videur" qui protège les routes sensibles de l'API.
// Il vérifie si la requête contient un jeton (token) valide avant d'autoriser l'accès.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. On cherche l'en-tête "Authorization" envoyé par le navigateur
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Non autorisé", http.StatusUnauthorized)
			return
		}

		// 2. On vérifie que le format est bien "Bearer <le_token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Token mal formaté", http.StatusUnauthorized)
			return
		}

		// 3. On décrypte et on valide le token avec notre clé secrète
		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("méthode de signature invalide")
			}
			return jwtSecretKey, nil
		})

		// Si le token est faux, expiré ou modifié, on bloque l'accès
		if err != nil || !token.Valid {
			http.Error(w, "Token invalide ou expiré", http.StatusUnauthorized)
			return
		}

		// 4. Si tout est bon, on laisse la requête continuer vers le Controller prévu
		next.ServeHTTP(w, r)
	}
}
