package main

import (
	"log"
	"notification-service/internal/websocket"
	"notification-service/internal/grpc" 
)

func main() {
	go websocket.StartWebSocketServer() 

	if err := grpc.StartGRPCServer(); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}
}