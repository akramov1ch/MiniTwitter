# Dockerfile

# 1. Go muhitini o'rnatish
FROM golang:1.23.1 AS builder

# 2. Ishchi katalogni yaratish
WORKDIR /app

# 3. Modullarni yuklash
COPY go.mod go.sum ./
RUN go mod download

# 4. Manba kodini yuklash
COPY . .

# 5. Dastur qurilishi
RUN CGO_ENABLED=0 GOOS=linux go build -o notification-service ./cmd/main.go

# 6. Minimal konteynerni yaratish
FROM alpine:latest   

# 7. Kerakli kutubxonalarni o'rnatish
RUN apk --no-cache add ca-certificates

# 8. Qurilgan dastur faylini ko'chirish
WORKDIR /root/
COPY --from=builder /app/notification-service .

EXPOSE 5056
# 9. Dastur ishga tushirish
CMD ["./notification-service"]