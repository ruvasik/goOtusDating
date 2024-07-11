package main

import (
    "log"
    "net/http"
    "database/sql"  // Add this import
    "github.com/gorilla/mux"
    "github.com/ruvasik/goOtusDating/internal/database"
    "github.com/ruvasik/goOtusDating/internal/handlers"
)

func main() {
    // Initialize the database connections
    database.InitDB()
    defer database.CloseDB()

    // Create a new router
    r := mux.NewRouter()

    // Set up routes with the router, master, and slave databases
    handlers.SetupRoutes(r, database.DBMaster, []*sql.DB{database.DBSlave})

    // Start the server
    log.Fatal(http.ListenAndServe(":8080", r))
}
