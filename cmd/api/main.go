package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/ruvasik/goOtusDating/internal/database"
    "github.com/ruvasik/goOtusDating/internal/handlers"
)

func main() {
    db := database.Connect()
    defer db.Close()

    r := mux.NewRouter()
    handlers.SetupRoutes(r, db)

    log.Fatal(http.ListenAndServe(":8080", r))
}
