package routes

import (
	"encoding/json"
	"net/http"
	"warehouse-backend/controllers"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UserRoutes(router chi.Router) {
	// Создать пользователя
	router.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		var user struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		// Декодируем тело запроса в структуру user
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Создаём пользователя через контроллер
		createdUser, err := controllers.CreateUser(user.Name, user.Email, user.Password)
		if err != nil {
			http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Возвращаем успешный ответ
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(createdUser)
	})

	// Получить пользователя по ID
	router.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		// Получаем параметр id из URL
		id := chi.URLParam(r, "id")
		objectID, err := primitive.ObjectIDFromHex(id) // Преобразуем строку в ObjectID
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Получаем пользователя через контроллер
		user, err := controllers.GetUserByID(objectID.Hex())
		if err != nil {
			http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
			return
		}

		// Возвращаем успешный ответ с данными пользователя
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})
}
