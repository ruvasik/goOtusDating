# Используем базовый образ для сборки и выполнения
FROM golang:1.21-alpine

RUN apk update && apk add --no-cache postgresql-client

WORKDIR /app

# Копируем файлы go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Компилируем бинарные файлы
RUN go build -o main ./cmd/api
RUN go build -o generate ./cmd/generate

RUN chmod +x /app/main /app/generate

COPY people.v2.csv /app/people.v2.csv

# Указываем команду для запуска приложения
CMD ["./main"]
