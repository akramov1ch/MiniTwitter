APP_NAME=main
DOCKER_IMAGE=api-gateway
DOCKERFILE=Dockerfile

run:
	go run cmd/main.go

swag:
	swag init -g cmd/main.go -o ./cmd/docs --parseDependency --parseInternal --parseDepth 2

build:
	docker build -t $(DOCKER_IMAGE) -f $(DOCKERFILE) .

run-docker:
	docker run -p 5051:5050 --env-file .env $(DOCKER_IMAGE)

clean:
	docker rmi $(DOCKER_IMAGE)

proto:
	find . -name "*.proto" -type f -exec protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative {} \;
