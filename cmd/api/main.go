package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/ruvasik/goOtusDating/internal/database"
    _ "github.com/lib/pq"
)

type User struct {
    ID        int    `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    BirthDate string `json:"birth_date"`
    City      string `json:"city"`
}

var db *sql.DB

func main() {
    var err error
    db = database.Connect()
    if err != nil {
        log.Fatal(err)
    }

    r := mux.NewRouter()
    r.HandleFunc("/user/search", searchUsers).Methods("GET")

    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}

func searchUsers(w http.ResponseWriter, r *http.Request) {
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

    var users []User
    for rows.Next() {
        var u User
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
