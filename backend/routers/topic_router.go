package routers

import (
	"net/http"
	"projet-forum/backend/controllers"
	"projet-forum/backend/handlers"
)

func SetupTopicRoutes(mux *http.ServeMux, controller *controllers.TopicController) {
	mux.HandleFunc("/api/topics", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			controller.GetTopics(w, r)
		} else if r.Method == "POST" {
			handlers.AuthMiddleware(http.HandlerFunc(controller.CreateTopic))(w, r)
		}
	})

	mux.HandleFunc("/api/topics/status/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			handlers.AuthMiddleware(http.HandlerFunc(controller.UpdateTopicStatus))(w, r)
		}
	})

	mux.HandleFunc("/api/topics/react/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			handlers.AuthMiddleware(http.HandlerFunc(controller.ReactTopic))(w, r)
		}
	})

	mux.HandleFunc("/api/topics/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			handlers.AuthMiddleware(http.HandlerFunc(controller.DeleteTopic))(w, r)
		} else if r.Method == "PUT" {
			handlers.AuthMiddleware(http.HandlerFunc(controller.UpdateTopic))(w, r)
		}
	})
}
