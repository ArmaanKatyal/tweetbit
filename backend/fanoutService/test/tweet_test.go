package test

import (
	"context"
	"testing"

	pb "github.com/ArmaanKatyal/tweetbit/backend/fanoutService/proto"
	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/service"
)

func TestCreateTweet(t *testing.T) {
	response, error := service.NewFanoutServer().CreateTweet(
		context.Background(),
		&pb.CreateTweetRequest{
			Id:            "1",
			Content:       "TEST",
			UserId:        "1",
			Uuid:          "SomeUUID",
			CreatedAt:     "SomeTime",
			LikesCount:    "0",
			RetweetsCount: "0",
		},
	)

	if error != nil {
		t.Errorf("Error: %v", error)
	}

	// check if test passed
	if response.Success != true {
		t.Errorf("Test failed")
	}
}
