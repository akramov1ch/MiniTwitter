package tweet_clients

import (
	"api-gateway/config"
	proto "api-gateway/protos/tweet-proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DialTweetGrpc() proto.TweetServiceClient {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}
	conn, err := grpc.NewClient(conf.TWEET_SERVER_NAME + ":" + conf.TWEET_SERVER_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to dial grpc client Tweet:", err)
	}
	return proto.NewTweetServiceClient(conn)
}
