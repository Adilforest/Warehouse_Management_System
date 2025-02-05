package routes

import (
	"net/http"
	"warehouse-backend/controllers"
	"warehouse-backend/middleware"
)

func ContactRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/contact", middleware.LoggingMiddleware(controllers.ContactController))
}
