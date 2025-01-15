package routes

import (
	"net/http"
	"warehouse-backend/controllers"
	"warehouse-backend/middleware"
)

// ContactRoutes настраивает маршруты для обработки запросов, связанных с формой "Contact Us".
func ContactRoutes(mux *http.ServeMux) {
	// Регистрируем обработчик для пути /api/contact с middleware для логирования
	mux.HandleFunc("/api/contact", middleware.LoggingMiddleware(controllers.ContactController))
}
