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

// CreateTopic permet à un utilisateur connecté de publier un nouveau sujet
func (c *TopicController) CreateTopic(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// On récupère le pseudo de l'auteur depuis les en-têtes (sécurisé par le middleware)
	pseudo := r.Header.Get("Pseudo")
	if pseudo == "" {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	if err := c.Service.CreateTopic(req, pseudo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetTopics renvoie la liste des sujets avec pagination, recherche et tri
func (c *TopicController) GetTopics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Récupération des paramètres de l'URL (?category=1&limit=10...)
	categoryIDStr := r.URL.Query().Get("category")
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sort")

	// Pagination par défaut (10 sujets par page)
	limit, offset := 10, 0
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
		limit = l
	}
	if o, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && o >= 0 {
		offset = o
	}

	var topics []repositories.TopicResponse
	var err error

	// Tri selon qu'une catégorie précise (franchise) est demandée ou non
	if categoryIDStr != "" {
		if categoryID, errConv := strconv.Atoi(categoryIDStr); errConv == nil {
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

	// Évite de renvoyer "null" si la base est vide (renvoie un tableau vide [])
	if topics == nil {
		topics = []repositories.TopicResponse{}
	}

	json.NewEncoder(w).Encode(topics)
}

// DeleteTopic permet de supprimer un sujet (vérifie si c'est l'auteur ou un admin)
func (c *TopicController) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	// Extraction de l'ID depuis l'URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/topics/")
	topicID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	pseudo := r.Header.Get("Pseudo")
	if pseudo == "" {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	if err := c.Service.DeleteTopic(topicID, pseudo); err != nil {
		if err.Error() == "forbidden" {
			http.Error(w, "Interdit : Vous ne pouvez pas supprimer ce sujet", http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateTopic modifie le contenu d'un sujet existant
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
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	if err := c.Service.UpdateTopic(topicID, req.Content, pseudo); err != nil {
		if err.Error() == "forbidden" {
			http.Error(w, "Interdit", http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

// UpdateTopicStatus permet aux admins de fermer ou d'ouvrir un sujet
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
	if err := c.Service.UpdateTopicStatus(topicID, req.Status, pseudo); err != nil {
		if err.Error() == "forbidden" {
			http.Error(w, "Interdit", http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

// ReactTopic gère l'ajout d'un Like ou Dislike sur un sujet
func (c *TopicController) ReactTopic(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/topics/react/")
	topicID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	var req struct {
		Action string `json:"action"` // "like" ou "dislike"
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	pseudo := r.Header.Get("Pseudo")
	if pseudo == "" {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	if err := c.Service.ReactTopic(topicID, req.Action, pseudo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
