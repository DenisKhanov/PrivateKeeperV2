# Stage 1: Builder
FROM golang:1.22.2-alpine AS builder

# Создаем и переходим в директорию приложения
WORKDIR /app

# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./

# Копируем server tlsconfig.
COPY internal/tlsconfig/cert/server ./

# Загружаем все зависимости
RUN go mod download

# Копируем исходный код в рабочую директорию контейнера
COPY ./ ./

# Копируем файл .env
COPY server.env ./

# Сборка бинарного файла
RUN CGO_ENABLED=0 GOOS=linux go build -a -o privatekeeperv2 ./cmd/private_keeper__server/server.go

# Stage 2: Runner
FROM alpine:3.19

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем собранный бинарник и файл .env из предыдущего этапа
COPY --from=builder /app/privatekeeperv2 .
COPY --from=builder /app/server.env ./
COPY --from=builder /app/internal/tlsconfig/cert/server /app/internal/tlsconfig/cert/server/

# Проверка наличия сертификатов
RUN ls -la /app/internal/tlsconfig/cert/server/

# Запускаем бинарник
CMD ["./privatekeeperv2"]