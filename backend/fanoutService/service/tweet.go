package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/helpers"
	pb "github.com/ArmaanKatyal/tweetbit/backend/fanoutService/proto"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ITweet struct {
	Id            string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Content       string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	UserId        string `protobuf:"bytes,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Uuid          string `protobuf:"bytes,4,opt,name=uuid,proto3" json:"uuid,omitempty"`
	CreatedAt     string `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	LikesCount    string `protobuf:"varint,6,opt,name=likes_count,json=likesCount,proto3" json:"likes_count,omitempty"`
	RetweetsCount string `protobuf:"varint,7,opt,name=retweets_count,json=retweetsCount,proto3" json:"retweets_count,omitempty"`
}

func (sever *FanoutServer) CreateTweet(_ context.Context, req *pb.CreateTweetRequest) (*pb.CreateTweetResponse, error) {
	log.Printf("CreateTweet: %v", req.String())
	// followers := make(chan []*pb.User)
	if helpers.StringToBool(helpers.GetConfigValue("featureFlag.enableKafka")) && helpers.StringToBool(helpers.GetConfigValue("featureFlag.enableCreateTweet")) {
		go func() {
			// publish the tweet to kafka to be consumed by other services
			p, err := kafka.NewProducer(&kafka.ConfigMap{
				"bootstrap.servers": helpers.GetConfigValue("kafka.bootstrap.servers"),
				"client.id":         helpers.GetConfigValue("kafka.client.id"),
				"acks":              helpers.GetConfigValue("kafka.acks"),
			})
			if err != nil {
				log.Fatalf("failed to create producer: %s", err)
			}
			topic := "createTweet"
			op := NewKafka(p, topic)
			value := &ITweet{req.Id, req.Content, req.UserId, req.Uuid, req.CreatedAt, req.LikesCount, req.RetweetsCount}
			jsonValue, err := json.Marshal(value)
			if err != nil {
				log.Fatalf("failed to marshal tweet: %s", err)
			}
			if err = op.PublishMessage(jsonValue); err != nil {
				log.Fatalf("failed to place tweet: %s", err)
			}
		}()
	}
	return &pb.CreateTweetResponse{Success: true}, nil
}
