package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	pb "github.com/ArmaanKatyal/tweetbit/backend/fanoutService/proto"
	"github.com/spf13/viper"
)

type IFollowUser struct {
	UserId     string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	FollowerId string `protobuf:"bytes,2,opt,name=follower_id,json=followerId,proto3" json:"follower_id,omitempty"`
}

func (server *FanoutServer) FollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb.FollowUserResponse, error) {
	log.Printf("FollowUser: %v", req.String())
	if followUserEnabled() {
		start := time.Now()
		// publish to kafka
		go func() {
			topic := "followUser"
			message := &IFollowUser{req.UserId, req.FollowerId}
			json_message, err := json.Marshal(message)
			if err != nil {
				log.Fatalf("failed to marshal follow user: %s", err)
			}
			PublishMessage(ctx, topic, json_message)
		}()
		end := time.Since(start).Seconds()
		server.metrics.KafkaResponseTimeHistogram.WithLabelValues("followUser").Observe(end)
		server.metrics.KafkaTransactionTotal.WithLabelValues("followUser").Inc()
	}
	return &pb.FollowUserResponse{Success: true}, nil
}

func (server *FanoutServer) UnfollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb.FollowUserResponse, error) {
	log.Printf("UnFollowUser: %v", req.String())
	if unfollowUserEnabled() {
		start := time.Now()
		go func() {
			topic := "unfollowUser"
			message := &IFollowUser{req.UserId, req.FollowerId}
			json_message, err := json.Marshal(message)
			if err != nil {
				log.Fatalf("failed to marshal unfollow user: %s", err)
			}
			PublishMessage(ctx, topic, json_message)
		}()
		end := time.Since(start).Seconds()
		server.metrics.KafkaResponseTimeHistogram.WithLabelValues("unfollowUser").Observe(end)
		server.metrics.KafkaTransactionTotal.WithLabelValues("unfollowUser").Inc()
	}
	return &pb.FollowUserResponse{Success: true}, nil
}

func followUserEnabled() bool {
	return viper.GetBool("featureFlag.enableKafka") && viper.GetBool("featureFlag.enableFollowUser")
}

func unfollowUserEnabled() bool {
	return viper.GetBool("featureFlag.enableKafka") && viper.GetBool("featureFlag.enableUnfollowUser")
}
