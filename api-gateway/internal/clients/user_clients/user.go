package user_clients

import (
	"api-gateway/config"
	proto "api-gateway/protos/user-proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DialUserGrpc() proto.UserServiceClient {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}
	conn, err := grpc.NewClient(conf.USER_SERVER_NAME + ":" + conf.USER_SERVER_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to dial grpc client User:", err)
	}
	return proto.NewUserServiceClient(conn)
}
