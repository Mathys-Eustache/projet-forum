package routers

import (
	"net/http"
	"projet-forum/backend/controllers"
)

func SetupCategoryRoutes(mux *http.ServeMux, categoryController *controllers.CategoryController) {
	mux.HandleFunc("/api/categories", categoryController.GetCategories)
}
