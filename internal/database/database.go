package database

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "log"
    "os"
)

var SlaveDBs []*sql.DB

func Connect() (*sql.DB, []*sql.DB) {
    masterDBUser := os.Getenv("DB_USER")
    masterDBPassword := os.Getenv("DB_PASSWORD")
    masterDBName := os.Getenv("DB_NAME")
    masterDBHost := os.Getenv("DB_HOST")
    masterDBPort := os.Getenv("DB_PORT")

    masterConnStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", masterDBUser, masterDBPassword, masterDBName, masterDBHost, masterDBPort)
    masterDB, err := sql.Open("postgres", masterConnStr)
    if err != nil {
        log.Fatal(err)
    }

    slaveDBUser := os.Getenv("DB_USER")
    slaveDBPassword := os.Getenv("DB_PASSWORD")
    slaveDBName := os.Getenv("DB_NAME")
    slaveHost := os.Getenv("DB_SLAVE_HOST")
    slavePort := os.Getenv("DB_SLAVE_PORT")

    slaveConnStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", slaveDBUser, slaveDBPassword, slaveDBName, slaveHost, slavePort)
    slaveDB, err := sql.Open("postgres", slaveConnStr)
    if err != nil {
        log.Fatal(err)
    }
    SlaveDBs = append(SlaveDBs, slaveDB)

    return masterDB, SlaveDBs
}

func GetSlaveDB() *sql.DB {
    // Simple round-robin load balancing
    return SlaveDBs[0] // You can implement more advanced load balancing
}
