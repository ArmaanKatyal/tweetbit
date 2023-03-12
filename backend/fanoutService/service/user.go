package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/helpers"
	pb "github.com/ArmaanKatyal/tweetbit/backend/fanoutService/proto"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type IFollowUser struct {
	UserId     string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	FollowerId string `protobuf:"bytes,2,opt,name=follower_id,json=followerId,proto3" json:"follower_id,omitempty"`
}

func (server *FanoutServer) FollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb.FollowUserResponse, error) {
	log.Printf("Received: %v", req.String())

	if helpers.StringToBool(helpers.GetConfigValue("featureFlag.enableKafka")) && helpers.StringToBool(helpers.GetConfigValue("featureFlag.enableFollowUser")) {
		// publish to kafka
		go func() {
			p, err := kafka.NewProducer(&kafka.ConfigMap{
				"bootstrap.servers": helpers.GetConfigValue("kafka.bootstrap.servers"),
				"client.id":         helpers.GetConfigValue("kafka.client.id"),
				"acks":              helpers.GetConfigValue("kafka.acks"),
			})
			if err != nil {
				log.Fatalf("failed to create producer: %s", err)
			}
			topic := "followUser"
			op := NewKafka(p, topic)
			message := &IFollowUser{req.UserId, req.FollowerId}
			json_message, err := json.Marshal(message)
			if err != nil {
				log.Fatalf("failed to marshal follow user: %s", err)
			}
			if err = op.PublishMessage(json_message); err != nil {
				log.Fatalf("failed to place follow user: %s", err)
			}
		}()
	}
	return &pb.FollowUserResponse{Success: true}, nil
}

func (server *FanoutServer) UnfollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb.FollowUserResponse, error) {
	log.Printf("Received: %v", req.String())

	if helpers.StringToBool(helpers.GetConfigValue("featureFlag.enableKafka")) && helpers.StringToBool(helpers.GetConfigValue("featureFlag.enableUnfollowUser")) {
		go func() {
			p, err := kafka.NewProducer(&kafka.ConfigMap{
				"bootstrap.servers": helpers.GetConfigValue("kafka.bootstrap.servers"),
				"client.id":         helpers.GetConfigValue("kafka.client.id"),
				"acks":              helpers.GetConfigValue("kafka.acks"),
			})
			if err != nil {
				log.Fatalf("failed to create producer: %s", err)
			}
			topic := "unfollowUser"
			op := NewKafka(p, topic)
			message := &IFollowUser{req.UserId, req.FollowerId}
			json_message, err := json.Marshal(message)
			if err != nil {
				log.Fatalf("failed to marshal unfollow user: %s", err)
			}
			if err = op.PublishMessage(json_message); err != nil {
				log.Fatalf("failed to place unfollow user: %s", err)
			}
		}()
	}
	return &pb.FollowUserResponse{Success: true}, nil
}
