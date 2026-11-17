package routers

import (
	"net/http"
	"projet-forum/backend/controllers"
)

func SetupAuthRoutes(mux *http.ServeMux, authController *controllers.AuthController) {
	mux.HandleFunc("/api/register", authController.RegisterHandler)
	mux.HandleFunc("/api/login", authController.LoginHandler)
}
