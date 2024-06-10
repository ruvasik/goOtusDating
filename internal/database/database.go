package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func Connect() *sql.DB {
	connStr := "user=user password=password dbname=dating_db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
