package main

import (
    "bufio"
    "database/sql"
    "fmt"
    "github.com/ruvasik/goOtusDating/internal/database"
    "log"
    "math/rand"
    "os"
    "strings"
    "time"
    "unicode/utf8"
)

const (
    numUsers = 1000000
    csvFile  = "/app/people.v2.csv" // Используем абсолютный путь
)

func main() {
    db := database.Connect()
    defer db.Close()

    namesSurnames, err := fetchNamesSurnames(csvFile)
    if err != nil {
        log.Fatalf("Error fetching names and surnames: %v", err)
    }

    if len(namesSurnames) == 0 {
        log.Fatalf("No names and surnames were fetched from the CSV.")
    }

    err = generateUsers(db, namesSurnames, numUsers)
    if err != nil {
        log.Fatalf("Error generating users: %v", err)
    }

    fmt.Println("Data insertion completed successfully.")
}

func fetchNamesSurnames(fileName string) ([][3]string, error) {
    file, err := os.Open(fileName)
    if err != nil {
        return nil, fmt.Errorf("Error opening CSV file: %v", err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    var namesSurnames [][3]string

    for scanner.Scan() {
        line := scanner.Text()
        if !utf8.ValidString(line) {
            log.Printf("Skipping invalid UTF-8 line: %s", line)
            continue
        }
        parts := strings.Split(line, ",")
        if len(parts) == 3 {
            namesSurnames = append(namesSurnames, [3]string{parts[0], parts[1], parts[2]})
        } else {
            log.Printf("Skipping malformed line: %s", line)
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("Error reading CSV: %v", err)
    }

    log.Printf("Fetched %d names and surnames from the CSV.", len(namesSurnames))
    return namesSurnames, nil
}

func generateUsers(db *sql.DB, namesSurnames [][3]string, numUsers int) error {
    rand.Seed(time.Now().UnixNano())

    stmt, err := db.Prepare("INSERT INTO users (first_name, last_name, birth_date, city) VALUES ($1, $2, $3, $4)")
    if err != nil {
        return fmt.Errorf("Error preparing SQL statement: %v", err)
    }
    defer stmt.Close()

    for i := 0; i < numUsers; i++ {
        nameSurname := namesSurnames[rand.Intn(len(namesSurnames))]
        _, err := stmt.Exec(nameSurname[0], nameSurname[1], "2000-01-01", nameSurname[2]) // Используем дату по умолчанию
        if err != nil {
            return fmt.Errorf("Error inserting user: %v", err)
        }

        if i%10000 == 0 {
            fmt.Printf("Inserted %d users...\n", i)
        }
    }

    return nil
}
