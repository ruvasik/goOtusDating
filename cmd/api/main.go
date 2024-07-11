package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/ruvasik/goOtusDating/internal/database"
    "github.com/ruvasik/goOtusDating/internal/handlers"
)

func main() {
    r := mux.NewRouter()

    database.InitDB()
    defer database.CloseDB()

    handlers.SetupRoutes(r, database.DBMaster, database.DBSlaves)

    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
}
