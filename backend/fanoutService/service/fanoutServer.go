package service

import (
	"log"
	"net"

	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/helpers"
	pb "github.com/ArmaanKatyal/tweetbit/backend/fanoutService/proto"
	"google.golang.org/grpc"
)

var (
	PORT = helpers.GetConfigValue("server.port")
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
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterTweetServiceServer(s, server)
	pb.RegisterUserServiceServer(s, server)

	log.Printf("Server listening on port %s", PORT)
	return s.Serve(lis)
}
