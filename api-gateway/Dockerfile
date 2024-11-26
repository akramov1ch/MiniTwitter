FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY . .
COPY tls/ ./tls/
COPY .env ./

RUN go mod download

RUN go build -o main cmd/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/tls ./tls/
COPY --from=builder /app/.env ./

EXPOSE 5050

CMD ["./main"]