package handlers

import (
    "database/sql"
    "encoding/json"
    "log" 
    "net/http"
    "golang.org/x/crypto/bcrypt"
    "github.com/gorilla/mux"
    "github.com/ruvasik/goOtusDating/internal/models"
    "github.com/ruvasik/goOtusDating/internal/database"
)

var masterDB *sql.DB

func SetupRoutes(r *mux.Router, master *sql.DB, slaves []*sql.DB) {
    masterDB = master
    database.SlaveDBs = slaves

    r.HandleFunc("/user/register", RegisterUser).Methods("POST")
    r.HandleFunc("/user/get/{id}", GetUser).Methods("GET")
    r.HandleFunc("/user/search", SearchUsers).Methods("GET")
    r.HandleFunc("/login", Login).Methods("POST")
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }
    user.Password = string(hashedPassword)

    query := `INSERT INTO users (first_name, last_name, birth_date, gender, interests, city, username, password)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
    err = masterDB.QueryRow(query, user.FirstName, user.LastName, user.BirthDate, user.Gender, user.Interests, user.City, user.Username, user.Password).Scan(&user.ID)
    if err != nil {
        http.Error(w, "Error creating user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
    db := database.GetSlaveDB()
    params := mux.Vars(r)
    id := params["id"]

    var user models.User
    query := `SELECT id, first_name, last_name, birth_date, gender, interests, city, username FROM users WHERE id = $1`
    err := db.QueryRow(query, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.BirthDate, &user.Gender, &user.Interests, &user.City, &user.Username)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Printf("User not found: %v", id)
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }
        log.Printf("Error fetching user: %v", err)
        http.Error(w, "Error fetching user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(user)
}

func SearchUsers(w http.ResponseWriter, r *http.Request) {
    db := database.GetSlaveDB()
    firstName := r.URL.Query().Get("first_name")
    lastName := r.URL.Query().Get("last_name")

    log.Printf("Searching for users with first name: %s and last name: %s", firstName, lastName)

    rows, err := db.Query("SELECT id, first_name, last_name, birth_date, city FROM users WHERE first_name LIKE $1 AND last_name LIKE $2 ORDER BY id",
        firstName+"%", lastName+"%")
    if err != nil {
        log.Printf("Error executing query: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var u models.User
        if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.BirthDate, &u.City); err != nil {
            log.Printf("Error scanning row: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        users = append(users, u)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Error iterating rows: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("Found %d users", len(users))

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func Login(w http.ResponseWriter, r *http.Request) {
    var creds models.User
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    var storedUser models.User
    query := `SELECT id, password FROM users WHERE username = $1`
    err = masterDB.QueryRow(query, creds.Username).Scan(&storedUser.ID, &storedUser.Password)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "User not found", http.StatusUnauthorized)
            return
        }
        http.Error(w, "Error fetching user", http.StatusInternalServerError)
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(creds.Password))
    if err != nil {
        http.Error(w, "Invalid password", http.StatusUnauthorized)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Login successful"))
}
