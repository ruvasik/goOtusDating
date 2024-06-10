package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/your_project/internal/models"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	
}

func Login(w http.ResponseWriter, r *http.Request) {
	
}

func SetupRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/user/register", RegisterUser).Methods("POST")
	r.HandleFunc("/user/get/{id}", GetUser).Methods("GET")
	r.HandleFunc("/login", Login).Methods("POST")
}
