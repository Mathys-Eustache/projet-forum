package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"projet-forum/backend/dto"
	"projet-forum/backend/services"
)

type PostController struct {
	service *services.PostService
}

func InitPostController(service *services.PostService) *PostController {
	return &PostController{service: service}
}

func (c *PostController) GetPostsByTopicHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	topicIDStr := r.URL.Query().Get("topic_id")
	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil || topicID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erreur": "topic_id manquant ou invalide"})
		return
	}

	posts, err := c.service.GetPostsByTopic(topicID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"erreur": "Erreur serveur"})
		return
	}

	json.NewEncoder(w).Encode(posts)
}

func (c *PostController) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req dto.CreatePostRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erreur": "Format des donnees invalide"})
		return
	}

	id, err := c.service.CreatePost(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erreur": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Message ajoute avec succes",
		"id":      id,
	})
}
