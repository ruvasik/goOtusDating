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

    slaveHosts := strings.Split(os.Getenv("DB_SLAVES"), ",")
    for _, host := range slaveHosts {
        connStr := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", host, dbUser, dbPassword, dbName)
        db, err := sql.Open("postgres", connStr)
        if err != nil {
            log.Fatal(err)
        }
        DBSlaves = append(DBSlaves, db)
    }
}

func CloseDB() {
    DBMaster.Close()
    for _, db := range DBSlaves {
        db.Close()
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
