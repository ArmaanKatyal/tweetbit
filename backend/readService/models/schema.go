package models

import "time"

type User struct {
	Id              uint
	Uuid            string
	Email           string
	First_name      string
	Last_name       string
	Created_at      *time.Time
	Favorites_count uint
	Followers_count uint
	Following_count uint
}

type Tweet struct {
	Id             uint
	Uuid           string
	User_id        uint
	Content        string
	Created_at     *time.Time
	Likes_count    uint
	Retweets_count uint
}

type User_Followers struct {
	Id          uint
	User_id     uint
	Follower_id uint
}

type Tweet_likes struct {
	Id         uint
	Tweet_id   uint
	User_id    uint
	Created_at *time.Time
}

type Tweet_Comments struct {
	Id         uint
	Tweet_id   uint
	User_id    uint
	Content    string
	Created_at *time.Time
}
