package controllers

import (
	"encoding/json"
	"net/http"
	"projet-forum/backend/repositories"
	"projet-forum/backend/services"
	"strconv"
	"strings"
)

type PostController struct {
	Service *services.PostService
}

func InitPostController(service *services.PostService) *PostController {
	return &PostController{Service: service}
}

func (c *PostController) HandlePosts(w http.ResponseWriter, r *http.Request) {
	topicIDStr := r.URL.Query().Get("topic_id")
	topicID, _ := strconv.Atoi(topicIDStr)

	if r.Method == "GET" {
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		limit := 10
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}

		offset := 0
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}

		posts, err := c.Service.GetPostsByTopic(topicID, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if posts == nil {
			posts = []repositories.PostResponse{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
		return
	}

	if r.Method == "POST" {
		var req struct {
			Content string `json:"content"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		pseudo := r.Header.Get("Pseudo")
		if pseudo == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err := c.Service.CreatePost(req.Content, topicID, pseudo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (c *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	pseudo := r.Header.Get("Pseudo")
	if pseudo == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = c.Service.DeletePost(postID, pseudo)
	if err != nil {
		if err.Error() == "forbidden" {
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
