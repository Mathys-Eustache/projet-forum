package main

import (
	"fmt"
	"net/http"
	"projet-forum/backend/config"
	"projet-forum/backend/controllers"
	"projet-forum/backend/repositories"
	"projet-forum/backend/routers"
	"projet-forum/backend/services"
)

// enableCORS permet à ton frontend (port 8081) de communiquer avec ton backend (port 8080)
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
	// 1. Initialisation de la base de données
	db := config.InitDB()
	defer db.Close()

	// 2. Initialisation des composants (Architecture MVC / Service)
	userRepo := repositories.InitUserRepository(db)
	authService := services.InitAuthService(userRepo)
	authController := controllers.InitAuthController(authService)

	categoryRepo := &repositories.CategoryRepository{DB: db}
	categoryService := &services.CategoryService{Repo: categoryRepo}
	categoryController := &controllers.CategoryController{Service: categoryService}

	topicRepo := &repositories.TopicRepository{DB: db}
	topicService := services.InitTopicService(topicRepo)
	topicController := controllers.InitTopicController(topicService)

	// 3. Configuration du routeur API
	mux := http.NewServeMux()

	// Routes API (uniquement de la donnée JSON)
	routers.SetupAuthRoutes(mux, authController)
	routers.SetupCategoryRoutes(mux, categoryController)
	routers.SetupTopicRoutes(mux, topicController)

	// 4. Lancement du serveur Backend
	fmt.Println("Serveur API Backend lancé sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", enableCORS(mux)); err != nil {
		fmt.Printf("Erreur lors du lancement du serveur API: %v\n", err)
	}
}
