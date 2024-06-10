# Указываем базовый образ
FROM golang:1.18-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Компилируем бинарный файл
RUN go build -o main ./cmd/api

# Указываем команду для запуска приложения
CMD ["./main"]
