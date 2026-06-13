package routers

import (
	"net/http"
	"projet-forum/backend/controllers"
	"projet-forum/backend/handlers"
)

func SetupPostRoutes(mux *http.ServeMux, controller *controllers.PostController) {
	mux.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			controller.HandlePosts(w, r)
		} else if r.Method == "POST" {
			handlers.AuthMiddleware(http.HandlerFunc(controller.HandlePosts))(w, r)
		}
	})

	mux.HandleFunc("/api/posts/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			handlers.AuthMiddleware(http.HandlerFunc(controller.DeletePost))(w, r)
		}
	})
}
