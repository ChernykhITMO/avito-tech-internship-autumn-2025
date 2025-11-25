FROM golang:1.25.1 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/app /app/app

EXPOSE 8080

CMD ["./app"]
