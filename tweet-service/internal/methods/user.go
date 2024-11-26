package methods

import (
	grp "tweet-service/pkg/proto"
	"tweet-service/config"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectUser() grp.UserServiceClient {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	conn, err := grpc.NewClient(conf.USER_SERVER_NAME+":"+conf.USER_SERVER_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error connect user micro...", err)
	}

	client := grp.NewUserServiceClient(conn)
	return client
}
