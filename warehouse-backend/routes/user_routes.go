package routes

import (
	"encoding/json"
	"net/http"
	"warehouse-backend/controllers"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UserRoutes(router chi.Router) {

	router.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		var user struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
			return
		}

		createdUser, err := controllers.CreateUser(user.Name, user.Email, user.Password)
		if err != nil {
			http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(createdUser)
	})

	router.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")
		objectID, err := primitive.ObjectIDFromHex(id) // Преобразуем строку в ObjectID
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		user, err := controllers.GetUserByID(objectID.Hex())
		if err != nil {
			http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})
}
