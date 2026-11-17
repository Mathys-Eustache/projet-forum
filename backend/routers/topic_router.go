package routers

import (
	"net/http"
	"projet-forum/backend/controllers"
	"projet-forum/backend/handlers"
)

func SetupTopicRoutes(mux *http.ServeMux, controller *controllers.TopicController) {

	// Route pour récupérer ou créer un topic
	mux.HandleFunc("/api/topics", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			controller.GetTopics(w, r)
		} else if r.Method == http.MethodPost {
			handlers.AuthMiddleware(controller.CreateTopic)(w, r)
		}
	})

	// Route pour ouvrir/fermer un topic (Admin)
	mux.HandleFunc("/api/topics/status/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			handlers.AuthMiddleware(controller.UpdateTopicStatus)(w, r)
		}
	})

	// Route pour liker/disliker un topic
	mux.HandleFunc("/api/topics/react/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			handlers.AuthMiddleware(controller.ReactTopic)(w, r)
		}
	})

	// Route pour modifier ou supprimer un topic
	mux.HandleFunc("/api/topics/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			handlers.AuthMiddleware(controller.DeleteTopic)(w, r)
		} else if r.Method == http.MethodPut {
			handlers.AuthMiddleware(controller.UpdateTopic)(w, r)
		}
	})
}
