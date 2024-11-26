package like_clients

import (
	"api-gateway/config"
	proto "api-gateway/protos/like-proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DialLikeGrpc() proto.LikeServiceClient {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}
	conn, err := grpc.NewClient(conf.LIKE_SERVER_NAME + ":" + conf.LIKE_SERVER_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to dial grpc client Like:", err)
	}
	return proto.NewLikeServiceClient(conn)
}
