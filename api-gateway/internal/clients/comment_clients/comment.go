package comment_clients

import (
	"api-gateway/config"
	proto "api-gateway/protos/comment-proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DialCommentGrpc() proto.CommentServiceClient {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}
	conn, err := grpc.NewClient(conf.COMMENT_SERVER_NAME + ":" + conf.COMMENT_SERVER_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to dial grpc client Comment:", err)
	}
	return proto.NewCommentServiceClient(conn)
}
