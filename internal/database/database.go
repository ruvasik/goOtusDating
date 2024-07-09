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

    slave1Host := os.Getenv("DB_SLAVE1_HOST")
    slave1Port := os.Getenv("DB_SLAVE1_PORT")
    slave1ConnStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", slaveDBUser, slaveDBPassword, slaveDBName, slave1Host, slave1Port)
    slave1DB, err := sql.Open("postgres", slave1ConnStr)
    if err != nil {
        log.Fatal(err)
    }
    SlaveDBs = append(SlaveDBs, slave1DB)

    slave2Host := os.Getenv("DB_SLAVE2_HOST")
    slave2Port := os.Getenv("DB_SLAVE2_PORT")
    slave2ConnStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", slaveDBUser, slaveDBPassword, slaveDBName, slave2Host, slave2Port)
    slave2DB, err := sql.Open("postgres", slave2ConnStr)
    if err != nil {
        log.Fatal(err)
    }
    SlaveDBs = append(SlaveDBs, slave2DB)

    return masterDB, SlaveDBs
}

func GetSlaveDB() *sql.DB {
    // Simple round-robin load balancing
    return SlaveDBs[0] // You can implement more advanced load balancing
}
