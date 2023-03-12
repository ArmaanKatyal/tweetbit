package test

import (
	"context"
	"testing"

	pb "github.com/ArmaanKatyal/tweetbit/backend/fanoutService/proto"
	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/service"
)

func TestFollowUser(t *testing.T) {
	response, error := service.NewFanoutServer().FollowUser(
		context.Background(),
		&pb.FollowUserRequest{
			UserId:     "1",
			FollowerId: "2",
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

func TestUnfollowUser(t *testing.T) {
	response, error := service.NewFanoutServer().UnfollowUser(
		context.Background(),
		&pb.FollowUserRequest{
			UserId:     "1",
			FollowerId: "2",
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
