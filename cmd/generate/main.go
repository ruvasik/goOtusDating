package main

import (
    "fmt"
    "log"
    "encoding/csv"
    "os"
    "strings"
    "database/sql"

    "github.com/ruvasik/goOtusDating/internal/database"
)

type User struct {
    LastName   string
    FirstName  string
    BirthDate  string
    City       string
}

var db *sql.DB

func main() {
    fmt.Println("Initializing database")
    database.InitDB()
    defer database.CloseDB()

    // Проверка соединения
    err := database.DBMaster.Ping()
    if err != nil {
        log.Fatalf("Failed to connect to master database: %v", err)
    }

    db = database.DBMaster

    // Устанавливаем количество пользователей, которых нужно создать
    targetUserCount := 1000000 // пример значения

    // Открываем CSV файл
    fmt.Println("Opening CSV file")
    csvFile, err := os.Open("/app/people.v2.csv")
    if err != nil {
        log.Fatal(err)
    }
    defer csvFile.Close()

    // Создаем CSV reader
    fmt.Println("Creating CSV reader")
    reader := csv.NewReader(csvFile)

    // Читаем все строки
    fmt.Println("Reading CSV file")
    records, err := reader.ReadAll()
    if err != nil {
        log.Fatal(err)
    }

    // Инициализируем счетчик созданных пользователей
    createdUserCount := 0

    fmt.Println("Start generating data")

    // Начинаем транзакцию
    tx, err := db.Begin()
    if err != nil {
        log.Fatalf("Failed to begin transaction: %v", err)
    }

    // Обрабатываем и вставляем данные
    for _, record := range records {
        if createdUserCount >= targetUserCount {
            break
        }

        // Разделяем фамилию и имя
        nameParts := strings.Split(record[0], " ")
        if len(nameParts) != 2 {
            log.Printf("Invalid name format: %s", record[0])
            err := tx.Rollback()
            if err != nil {
                log.Fatalf("Failed to rollback transaction: %v", err)
            }
            log.Fatalf("Rolled back transaction due to invalid name format")
        }

        user := User{
            LastName:  nameParts[0],
            FirstName: nameParts[1],
            BirthDate: record[1],
            City:      record[2],
        }

        fmt.Printf("Inserting user: %+v\n", user)

        // Обновленный INSERT-запрос
        _, err := tx.Exec("INSERT INTO users (first_name, last_name, birth_date, city) VALUES ($1, $2, $3, $4)",
            user.FirstName, user.LastName, user.BirthDate, user.City)
        if err != nil {
            log.Printf("Failed to insert user: %v", err)
            err := tx.Rollback()
            if err != nil {
                log.Fatalf("Failed to rollback transaction: %v", err)
            }
            log.Fatalf("Rolled back transaction due to insert error")
        }

        // Увеличиваем счетчик созданных пользователей
        createdUserCount++

        // Выводим прогресс каждые 1000 записей
        if createdUserCount % 1000 == 0 {
            fmt.Printf("Created %d users\n", createdUserCount)
        }
    }

    // Фиксируем транзакцию
    err = tx.Commit()
    if err != nil {
        log.Fatalf("Failed to commit transaction: %v", err)
    }
    log.Println("Transaction committed successfully")

    // Выводим итоговое количество созданных пользователей
    fmt.Printf("Total users created: %d\n", createdUserCount)
}
