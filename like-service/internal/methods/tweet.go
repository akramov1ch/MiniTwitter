package methods

import (
	grp "like-service/pkg/proto"
	"like-service/config"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectTweet() grp.TweetServiceClient {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	conn, err := grpc.NewClient(conf.TWEET_SERVER_NAME+":"+conf.TWEET_SERVER_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error connect tweet micro...", err)
	}

	client := grp.NewTweetServiceClient(conn)
	return client
}
