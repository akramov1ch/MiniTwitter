package main

import (
	"api-gateway/internal/router"
	_ "api-gateway/cmd/docs" // This line is important for swagger
)

// @title API Gateway
// @version 1.0
// @description API for managing users, tweets, comments, likes, and direct messages.
// @host localhost:5050
// @BasePath /
func main() {
	router.Router()
}
