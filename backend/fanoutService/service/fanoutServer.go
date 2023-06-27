package service

import (
	"log"
	"net"

	pb "github.com/ArmaanKatyal/tweetbit/backend/fanoutService/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type FanoutServer struct {
	pb.UnimplementedTweetServiceServer
	pb.UnimplementedUserServiceServer
}

func NewFanoutServer() *FanoutServer {
	return &FanoutServer{}
}

// Run the server
func (server *FanoutServer) Run() error {
	port := viper.GetString("server.port")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterTweetServiceServer(s, server)
	pb.RegisterUserServiceServer(s, server)

	log.Printf("Server listening on port %s", port)
	return s.Serve(lis)
}
