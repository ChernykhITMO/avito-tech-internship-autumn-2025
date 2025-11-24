# билд-стейдж
FROM golang:1.25.1 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

# рантайм-стейдж
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/app /app/app

# переменная для строки подключения к БД придёт из docker-compose
ENV DB_DSN=""

EXPOSE 8080

CMD ["./app"]
