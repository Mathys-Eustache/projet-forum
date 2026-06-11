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

func main() {
	// 1. Connexion à la base de données
	db := config.InitDB()
	defer db.Close()

	// 2. Initialisation du module AUTH (Utilisateurs)
	userRepo := repositories.InitUserRepository(db)
	authService := services.InitAuthService(userRepo)
	authController := controllers.InitAuthController(authService)

	// 3. Initialisation du module CATEGORIES
	categoryRepo := &repositories.CategoryRepository{DB: db}
	categoryService := &services.CategoryService{Repo: categoryRepo}
	categoryController := &controllers.CategoryController{Service: categoryService}

	// 4. Configuration des Routes sur le Mux
	mux := http.NewServeMux()
	routers.SetupAuthRoutes(mux, authController)
	routers.SetupCategoryRoutes(mux, categoryController)

	// 5. Lancement du serveur
	fmt.Println("Serveur lance sur http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
