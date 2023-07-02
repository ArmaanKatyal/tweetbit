package service

import (
	"log"
	"net"

	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/internal"
	pb "github.com/ArmaanKatyal/tweetbit/backend/fanoutService/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type FanoutServer struct {
	pb.UnimplementedTweetServiceServer
	pb.UnimplementedUserServiceServer
	metrics *internal.PromMetrics
}

func NewFanoutServer(pm *internal.PromMetrics) *FanoutServer {
	return &FanoutServer{
		metrics: pm,
	}
}

// Run the server
func (server *FanoutServer) Run() {
	port := viper.GetString("server.port")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterTweetServiceServer(s, server)
	pb.RegisterUserServiceServer(s, server)

	log.Printf("Server listening on port %s", port)
	s.Serve(lis)
}
