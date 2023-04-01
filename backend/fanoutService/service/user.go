package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/helpers"
	pb "github.com/ArmaanKatyal/tweetbit/backend/fanoutService/proto"
)

type IFollowUser struct {
	UserId     string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	FollowerId string `protobuf:"bytes,2,opt,name=follower_id,json=followerId,proto3" json:"follower_id,omitempty"`
}

func (server *FanoutServer) FollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb.FollowUserResponse, error) {
	log.Printf("FollowUser: %v", req.String())

	if helpers.StringToBool(helpers.GetConfigValue("featureFlag.enableKafka")) && helpers.StringToBool(helpers.GetConfigValue("featureFlag.enableFollowUser")) {
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
	}
	return &pb.FollowUserResponse{Success: true}, nil
}

func (server *FanoutServer) UnfollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb.FollowUserResponse, error) {
	log.Printf("UnFollowUser: %v", req.String())

	if helpers.StringToBool(helpers.GetConfigValue("featureFlag.enableKafka")) && helpers.StringToBool(helpers.GetConfigValue("featureFlag.enableUnfollowUser")) {
		go func() {
			topic := "unfollowUser"
			message := &IFollowUser{req.UserId, req.FollowerId}
			json_message, err := json.Marshal(message)
			if err != nil {
				log.Fatalf("failed to marshal unfollow user: %s", err)
			}
			PublishMessage(ctx, topic, json_message)
		}()
	}
	return &pb.FollowUserResponse{Success: true}, nil
}
