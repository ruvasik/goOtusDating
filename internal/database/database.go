package database

import (
    "database/sql"
    _ "github.com/lib/pq"
    "log"
    "os"
    "fmt"
)

var (
    DBMaster *sql.DB
    DBSlave  *sql.DB
)

func InitDB() {
    var err error

    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    masterConnStr := fmt.Sprintf("host=db-master port=5432 user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)
    DBMaster, err = sql.Open("postgres", masterConnStr)
    if err != nil {
        log.Fatal(err)
    }

    slaveConnStr := fmt.Sprintf("host=db-slave port=5432 user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)
    DBSlave, err = sql.Open("postgres", slaveConnStr)
    if err != nil {
        log.Fatal(err)
    }
}

func CloseDB() {
    DBMaster.Close()
    DBSlave.Close()
}

func GetSlaveDB() *sql.DB {
    return DBSlave
}
