package methods

import (
	grp "user-service/pkg/proto"
	"user-service/config"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectLike() grp.LikeServiceClient {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	conn, err := grpc.NewClient(conf.LIKE_SERVER_NAME+":"+conf.LIKE_SERVER_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error connect like micro...", err)
	}

	client := grp.NewLikeServiceClient(conn)
	return client
}
