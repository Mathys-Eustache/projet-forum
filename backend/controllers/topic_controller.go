package controllers

import (
	"encoding/json"
	"net/http"
	"projet-forum/backend/dto"
	"projet-forum/backend/repositories"
	"projet-forum/backend/services"
	"strconv"
	"strings"
)

type TopicController struct {
	Service *services.TopicService
}

func InitTopicController(service *services.TopicService) *TopicController {
	return &TopicController{Service: service}
}

func (c *TopicController) CreateTopic(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pseudo := r.Header.Get("Pseudo")
	if pseudo == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := c.Service.CreateTopic(req, pseudo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *TopicController) GetTopics(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := r.URL.Query().Get("category")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sort") // Récupération du tri

	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	offset := 0
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}

	var topics []repositories.TopicResponse
	var err error

	if categoryIDStr != "" {
		categoryID, errConv := strconv.Atoi(categoryIDStr)
		if errConv == nil {
			topics, err = c.Service.GetTopicsByCategory(categoryID, limit, offset, search, sortBy)
		} else {
			topics, err = c.Service.GetAllTopics(limit, offset, search, sortBy)
		}
	} else {
		topics, err = c.Service.GetAllTopics(limit, offset, search, sortBy)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if topics == nil {
		topics = []repositories.TopicResponse{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topics)
}

func (c *TopicController) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/topics/")
	topicID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	pseudo := r.Header.Get("Pseudo")
	if pseudo == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = c.Service.DeleteTopic(topicID, pseudo)
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

func (c *TopicController) UpdateTopic(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/topics/")
	topicID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	pseudo := r.Header.Get("Pseudo")
	if pseudo == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = c.Service.UpdateTopic(topicID, req.Content, pseudo)
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

func (c *TopicController) UpdateTopicStatus(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/topics/status/")
	topicID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	pseudo := r.Header.Get("Pseudo")
	err = c.Service.UpdateTopicStatus(topicID, req.Status, pseudo)
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

func (c *TopicController) ReactTopic(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/topics/react/")
	topicID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	var req struct {
		Action string `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	pseudo := r.Header.Get("Pseudo")
	if pseudo == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = c.Service.ReactTopic(topicID, req.Action, pseudo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
