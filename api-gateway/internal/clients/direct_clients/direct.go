package direct_clients

import (
	"api-gateway/config"
	proto "api-gateway/protos/direct-proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DialDirectGrpc() proto.DirectServiceClient {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}
	conn, err := grpc.NewClient(conf.DIRECT_SERVER_NAME + ":" + conf.DIRECT_SERVER_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to dial grpc client Direct:", err)
	}
	return proto.NewDirectServiceClient(conn)
}
