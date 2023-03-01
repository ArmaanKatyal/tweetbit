package service

import (
	"context"
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
}

func NewFanoutServer() *FanoutServer {
	return &FanoutServer{}
}

func (server *FanoutServer) Run() error {
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterTweetServiceServer(s, server)
	log.Printf("Server listening on port %s", PORT)
	return s.Serve(lis)
}

func (sever *FanoutServer) CreateTweet(_ context.Context, in *pb.CreateTweetRequest) (*pb.CreateTweetResponse, error) {
	log.Printf("Received: %v", in.String())
	if helpers.StringToBool(helpers.GetConfigValue("featureFlag.enableUserGraph")) {
		go GetFollowers(User{in.GetId(), in.GetUuid()})
	}
	return &pb.CreateTweetResponse{}, nil
}
