package main

import (
	"fmt"
	"net/http"
	"projet-forum/backend/config"
	"projet-forum/backend/controllers"
	"projet-forum/backend/handlers"
	"projet-forum/backend/repositories"
	"projet-forum/backend/routers"
	"projet-forum/backend/services"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Pseudo")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	db := config.InitDB()
	defer db.Close()

	userRepo := repositories.InitUserRepository(db)
	authService := services.InitAuthService(userRepo)
	authController := controllers.InitAuthController(authService)

	categoryRepo := &repositories.CategoryRepository{DB: db}
	categoryService := &services.CategoryService{Repo: categoryRepo}
	categoryController := &controllers.CategoryController{Service: categoryService}

	topicRepo := &repositories.TopicRepository{DB: db}
	topicService := services.InitTopicService(topicRepo)
	topicController := controllers.InitTopicController(topicService)

	postRepo := &repositories.PostRepository{DB: db}
	postService := services.InitPostService(postRepo)
	postController := controllers.InitPostController(postService)

	mux := http.NewServeMux()

	// 1. Autoriser Go à lire les fichiers statiques (CSS, images, JS)
	fs := http.FileServer(http.Dir("frontend/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// 2. La seule route dont tu as besoin pour afficher tes 30 équipes !
	mux.HandleFunc("/team", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/templates/team.html")
	})

	// Routes API
	routers.SetupAuthRoutes(mux, authController)
	routers.SetupCategoryRoutes(mux, categoryController)
	routers.SetupTopicRoutes(mux, topicController)
	routers.SetupPostRoutes(mux, postController)

	mux.HandleFunc("/api/messages", handlers.HandleMessages(db))
	mux.HandleFunc("/api/messages/", handlers.HandleDeleteMessage(db))

	fmt.Println("Serveur lance sur http://localhost:8080")
	http.ListenAndServe(":8080", enableCORS(mux))
}
