package controllers

import (
	"encoding/json"
	"net/http"

	"projet-forum/backend/dto"
	"projet-forum/backend/services"
)

type TopicController struct {
	service *services.TopicService
}

func InitTopicController(service *services.TopicService) *TopicController {
	return &TopicController{service: service}
}

func (c *TopicController) GetAllTopicsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	topics, err := c.service.GetAllTopics()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"erreur": "Erreur serveur"})
		return
	}

	json.NewEncoder(w).Encode(topics)
}

func (c *TopicController) CreateTopicHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req dto.CreateTopicRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erreur": "Format des donnees invalide"})
		return
	}

	id, err := c.service.CreateTopic(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erreur": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Sujet cree avec succes",
		"id":      id,
	})
}
