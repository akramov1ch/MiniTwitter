# O'zgaruvchilar
BINARY_NAME=like-service
DOCKER_IMAGE_NAME=like-service
DOCKER_TAG=latest

# Go buyruqlari
GO=go
GOTEST=$(GO) test
GOBUILD=$(GO) build

# Docker buyruqlari
DOCKER=docker

# Migratsiya yaratish uchun o'zgaruvchilar
MIGRATION_NAME ?= new_migration
MIGRATION_DIR := /pkg/migrations

.PHONY: all build test clean run docker-build docker-run docker-push create-migration migrate-up migrate-down

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/main.go

test:
	$(GOTEST) -v ./...

clean:
	$(GO) clean
	rm -f $(BINARY_NAME)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/main.go
	./$(BINARY_NAME)

docker-build:
	$(DOCKER) build -t $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) .

docker-run:
	$(DOCKER) run -p 8080:8080 \
	-e DB_HOST=POSTGRES \
	-e DB_PORT=5432 \
	-e DB_USER=postgres \
	-e DB_PASSWORD=vakhaboff \
	-e DB_NAME=likes \
	-e REDIS_HOST=REDIS \
	-e REDIS_PORT=6379 \
	$(DOCKER_IMAGE_NAME):$(DOCKER_TAG)
	
docker-push:
	$(DOCKER) push $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

# Qo'shimcha maqsadlar
lint:
	golangci-lint run

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

# Barcha testlarni va lintlarni o'tkazish
check: fmt vet lint test

# Migratsiya yaratish buyrug'i
create-migration:
	@read -p "Migratsiya nomi: " name; \
	timestamp=$$(date +%Y%m%d%H%M%S); \
	up_file="$(MIGRATION_DIR)/$${timestamp}_$${name}.up.sql"; \
	down_file="$(MIGRATION_DIR)/$${timestamp}_$${name}.down.sql"; \
	touch $$up_file $$down_file; \
	echo "Migratsiya fayllari yaratildi:"; \
	echo "- $$up_file"; \
	echo "- $$down_file"

# Migratsiyalarni bajarish buyrug'i
migrate-up:
	migrate -database "postgresql://postgres:vakhaboff@localhost:5432/likes?sslmode=disable" -path ./$(MIGRATION_DIR) up

# Migratsiyalarni bekor qilish buyrug'i
migrate-down:
	migrate -database "postgresql://postgres:vakhaboff@localhost:5432/likes?sslmode=disable" -path ./$(MIGRATION_DIR) down