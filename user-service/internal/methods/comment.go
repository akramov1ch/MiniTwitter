package methods

import (
	grp "user-service/pkg/proto"
	"user-service/config"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectComment() grp.CommentServiceClient {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	conn, err := grpc.NewClient(conf.COMMENT_SERVER_NAME+":"+conf.COMMENT_SERVER_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error connect comment micro...", err)
	}

	client := grp.NewCommentServiceClient(conn)
	return client
}
