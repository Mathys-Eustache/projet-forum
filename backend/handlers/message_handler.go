package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

type Message struct {
	ID     int    `json:"id"`
	Heure  string `json:"heure"`
	Pseudo string `json:"pseudo"`
	Texte  string `json:"texte"`
}

func HandleMessages(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")

		if r.Method == "OPTIONS" {
			return
		}

		switch r.Method {
		case "GET":
			rows, err := db.Query("SELECT p.id, DATE_FORMAT(p.created_at, '%H:%i'), u.username, p.content FROM Posts p JOIN Users u ON p.user_id = u.id ORDER BY p.created_at ASC")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var messages []Message
			for rows.Next() {
				var msg Message
				err := rows.Scan(&msg.ID, &msg.Heure, &msg.Pseudo, &msg.Texte)
				if err != nil {
					continue
				}
				messages = append(messages, msg)
			}

			json.NewEncoder(w).Encode(messages)

		case "POST":
			var req struct {
				Texte  string `json:"texte"`
				Pseudo string `json:"pseudo"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var userID int
			err := db.QueryRow("SELECT id FROM Users WHERE username = ?", req.Pseudo).Scan(&userID)
			if err != nil {
				userID = 1 // Sécurité : permet au message de s'enregistrer même si le pseudo n'est pas encore bien synchronisé
			}

			_, err = db.Exec("INSERT INTO Posts (user_id, topic_id, content) VALUES (?, 1, ?)", userID, req.Texte)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)

		default:
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		}
	}
}

func HandleDeleteMessage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")

		if r.Method == "OPTIONS" {
			return
		}

		if r.Method == "DELETE" {
			parts := strings.Split(r.URL.Path, "/")
			id := parts[len(parts)-1]

			_, err := db.Exec("DELETE FROM Posts WHERE id = ?", id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		}

	}
}
