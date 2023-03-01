package service

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/ArmaanKatyal/tweetbit/backend/fanoutService/helpers"
	pb "github.com/ArmaanKatyal/tweetbit/backend/fanoutService/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type User struct {
	id   string
	uuid string
}

func GetFollowers(user User) []User {
	conn, err := grpc.Dial(helpers.GetConfigValue("userGraphService.port"), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("could not close connection: %v", err)
		}
	}(conn)

	c := pb.NewFollowerServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := c.GetUserFollowers(ctx, &pb.GetFollowerRequest{UserId: user.id, Uuid: user.uuid})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	var followers []User
	for _, follower := range r.Users {
		followers = append(followers, User{strconv.Itoa(int(follower.Id)), follower.Uuid})
	}

	return followers
}
