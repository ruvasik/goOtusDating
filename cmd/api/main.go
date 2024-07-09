package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/ruvasik/goOtusDating/internal/database"
    "github.com/ruvasik/goOtusDating/internal/handlers"
)

func main() {
    masterDB, slaveDBs := database.Connect()

    r := mux.NewRouter()
    handlers.SetupRoutes(r, masterDB, slaveDBs)

    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
