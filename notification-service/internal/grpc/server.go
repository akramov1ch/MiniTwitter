package grpc

import (
	"log"
	"net"
	"google.golang.org/grpc"
	"notification-service/proto"
)

func StartGRPCServer() error {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		return err
	}

	s := grpc.NewServer()
	proto.RegisterNotificationServiceServer(s, &NotificationServiceServer{})

	log.Println("gRPC server is running on port :50051")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
		return err
	}
	return nil
}
