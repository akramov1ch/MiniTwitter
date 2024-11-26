# Makefile

# O'zgaruvchilar
APP_NAME = notification-service
DOCKER_IMAGE = $(APP_NAME):latest

# Qurish
build:
	go build -o $(APP_NAME) ./cmd/main.go

# Docker konteynerini yaratish
docker:
	docker build -t $(DOCKER_IMAGE) .

# Docker konteynerini ishga tushirish
run:
	docker run --rm -p 5056:5056 $(DOCKER_IMAGE)

# Tozalash
clean:
	rm -f $(APP_NAME)

# Barcha jarayonlarni bajarish
all: clean build docker run