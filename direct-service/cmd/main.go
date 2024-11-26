package main

import (
	"context"
	"database/sql"
	"direct-service/config"
	"direct-service/internal/handlers"
	"direct-service/internal/service"
	pb "direct-service/pkg/proto"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"direct-service/internal/redis"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	postgresUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", conf.DB_USER, conf.DB_PASSWORD, conf.DB_HOST, conf.DB_PORT, conf.DB_NAME)
	db, err := sql.Open("postgres", postgresUrl)
	if err != nil {
		log.Fatalf("Ma'lumotlar bazasiga ulanishda xato: %v", err)
	}
	defer db.Close()

	redisClient := redis.NewRedisClient()

	directService := service.NewDirectService(db, redisClient)
	directHandler := handlers.NewDirectHandler(*directService)

	grpcServer := grpc.NewServer()
	pb.RegisterDirectServiceServer(grpcServer, directHandler)
	reflection.Register(grpcServer)
	lis, err := net.Listen("tcp", ":"+conf.SERVER_PORT)
	if err != nil {
		log.Fatalf("Portni tinglashda xato: %v", err)
	}

	go func() {
		log.Println("gRPC server :" + conf.SERVER_PORT + " portida ishga tushdi")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC serverni ishga tushirishda xato: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Serverni to'xtatish boshlanmoqda...")

	grpcServer.GracefulStop()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := db.Close(); err != nil {
		log.Printf("Ma'lumotlar bazasi ulanishini yopishda xato: %v", err)
	}

	if err := redisClient.Close(); err != nil {
		log.Printf("Redis ulanishini yopishda xato: %v", err)
	}
	<-ctx.Done()

	log.Println("Server to'xtatildi")
}
