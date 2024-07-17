package database

import (
    "database/sql"
    _ "github.com/lib/pq"
    "log"
    "os"
    "strings"
    "fmt"
)

var (
    DBMaster *sql.DB
    DBSlaves []*sql.DB
)

func InitDB() {
    var err error

    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    masterConnStr := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_MASTER"), dbUser, dbPassword, dbName)
    DBMaster, err = sql.Open("postgres", masterConnStr)
    if err != nil {
        log.Fatal(err)
    }

    // Проверка соединения с мастер базой данных
    err = DBMaster.Ping()
    if err != nil {
        log.Fatalf("Failed to ping master database: %v", err)
    }
    log.Println("Successfully connected to master database")

    slaveHosts := strings.Split(os.Getenv("DB_SLAVES"), ",")
    for _, host := range slaveHosts {
        connStr := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", host, dbUser, dbPassword, dbName)
        db, err := sql.Open("postgres", connStr)
        if err != nil {
            log.Fatalf("Failed to connect to slave database %s: %v", host, err)
        }

        // Проверка соединения с каждым slave
        err = db.Ping()
        if err != nil {
            log.Fatalf("Failed to ping slave database %s: %v", host, err)
        }
        log.Printf("Successfully connected to slave database %s", host)

        DBSlaves = append(DBSlaves, db)
    }
}

func CloseDB() {
    if DBMaster != nil {
        DBMaster.Close()
    }
    for _, db := range DBSlaves {
        if db != nil {
            db.Close()
        }
    }
}

func GetSlaveDB() *sql.DB {
    // Use a simple round-robin mechanism to distribute the load across slaves
    if len(DBSlaves) == 0 {
        return DBMaster
    }
    slave := DBSlaves[0]
    DBSlaves = append(DBSlaves[1:], slave)
    return slave
}
