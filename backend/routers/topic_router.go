package routers

import (
	"net/http"

	"projet-forum/backend/controllers"
)

func SetupTopicRoutes(mux *http.ServeMux, topicController *controllers.TopicController) {
	mux.HandleFunc("/api/topics", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			topicController.GetAllTopicsHandler(w, r)
		case http.MethodPost:
			topicController.CreateTopicHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"erreur": "Methode non autorisee"}`))
		}
	})
}
