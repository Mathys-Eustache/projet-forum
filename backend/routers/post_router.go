package routers

import (
	"net/http"

	"projet-forum/backend/controllers"
)

func SetupPostRoutes(mux *http.ServeMux, postController *controllers.PostController) {
	mux.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			postController.GetPostsByTopicHandler(w, r)
		case http.MethodPost:
			postController.CreatePostHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"erreur": "Methode non autorisee"}`))
		}
	})
}
