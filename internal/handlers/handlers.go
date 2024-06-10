package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/ruvasik/goOtusDating/internal/models"
)

var db *sql.DB

func SetupRoutes(r *mux.Router, database *sql.DB) {
	db = database
	r.HandleFunc("/user/register", RegisterUser).Methods("POST")
	r.HandleFunc("/user/get/{id}", GetUser).Methods("GET")
	r.HandleFunc("/login", Login).Methods("POST")
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Реализуйте обработчик для регистрации пользователя
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Реализуйте обработчик для получения пользователя по ID
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Простой ответ для теста
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}
