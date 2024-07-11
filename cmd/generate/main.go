package main

import (
    "fmt"
    "log"
    "encoding/csv"
    "os"
    "strings"

    "github.com/ruvasik/goOtusDating/internal/database"
)

type User struct {
	LastName   string
	FirstName  string
	BirthDate  string
	City       string
}

func main() {
    db, _ := database.Connect()
    defer db.Close()


	// Устанавливаем количество пользователей, которых нужно создать
	targetUserCount := 1000000 // пример значения

	// Открываем CSV файл
	csvFile, err := os.Open("/app/people.v2.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	// Создаем CSV reader
	reader := csv.NewReader(csvFile)

	// Читаем все строки
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Инициализируем счетчик созданных пользователей
	createdUserCount := 0

	// Обрабатываем и вставляем данные
	for _, record := range records {
		if createdUserCount >= targetUserCount {
			break
		}

		// Разделяем фамилию и имя
		nameParts := strings.Split(record[0], " ")
		if len(nameParts) != 2 {
			log.Fatalf("Неправильный формат имени: %s", record[0])
		}

		user := User{
			LastName:  nameParts[0],
			FirstName: nameParts[1],
			BirthDate: record[1],
			City:      record[2],
		}

		_, err := db.Exec("INSERT INTO users (last_name, first_name, birth_date, city) VALUES ($1, $2, $3, $4)",
			user.LastName, user.FirstName, user.BirthDate, user.City)
		if err != nil {
			log.Fatal(err)
		}

		// Увеличиваем счетчик созданных пользователей
		createdUserCount++

		// Выводим прогресс каждые 10 000 записей
		if createdUserCount%10000 == 0 {
			fmt.Printf("Создано %d пользователей\n", createdUserCount)
		}
	}

	// Выводим итоговое количество созданных пользователей
	fmt.Printf("Всего создано пользователей: %d\n", createdUserCount)
}