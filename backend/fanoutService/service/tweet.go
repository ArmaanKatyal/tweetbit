package service

import (
	"context"
	"encoding/json"
	"log"
	"net"

	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/helpers"
	pb "github.com/ArmaanKatyal/tweetbit/backend/fanoutService/proto"
	"github.com/confluentinc/confluent-kafka-go/kafka"
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

// Run the server
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

type TweetPlacer struct {
	producer      *kafka.Producer
	topic         string
	delivery_chan chan kafka.Event
}

func NewTweetPlacer(p *kafka.Producer, t string) *TweetPlacer {
	return &TweetPlacer{
		producer:      p,
		topic:         t,
		delivery_chan: make(chan kafka.Event),
	}
}

type TweetI struct {
	Id            string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Content       string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	UserId        string `protobuf:"bytes,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Uuid          string `protobuf:"bytes,4,opt,name=uuid,proto3" json:"uuid,omitempty"`
	CreatedAt     string `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	LikesCount    int32  `protobuf:"varint,6,opt,name=likes_count,json=likesCount,proto3" json:"likes_count,omitempty"`
	RetweetsCount int32  `protobuf:"varint,7,opt,name=retweets_count,json=retweetsCount,proto3" json:"retweets_count,omitempty"`
}

// PlaceTweet places a tweet in the kafka topic
func (op *TweetPlacer) PlaceTweet(tweet *pb.CreateTweetRequest) error {
	value := &TweetI{tweet.Id, tweet.Content, tweet.UserId, tweet.Uuid, tweet.CreatedAt, tweet.LikesCount, tweet.RetweetsCount}
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = op.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &op.topic, Partition: kafka.PartitionAny},
		Value:          []byte(jsonValue),
	}, op.delivery_chan)
	if err != nil {
		return err
	}

	e := <-op.delivery_chan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}
	return nil
}

func (sever *FanoutServer) CreateTweet(_ context.Context, in *pb.CreateTweetRequest) (*pb.CreateTweetResponse, error) {
	log.Printf("Received: %v", in.String())
	// followers := make(chan []*pb.User)
	if helpers.StringToBool(helpers.GetConfigValue("featureFlag.enableKafka")) && helpers.StringToBool(helpers.GetConfigValue("featureFlag.enableCreateTweet")) {
		// make a concurrent call to user graph service
		// go GetFollowers(User{in.GetId(), in.GetUuid()}, followers)
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
			op := NewTweetPlacer(p, topic)
			if err = op.PlaceTweet(in); err != nil {
				log.Fatalf("failed to place tweet: %s", err)
			}
		}()
	}
	return &pb.CreateTweetResponse{}, nil
}
