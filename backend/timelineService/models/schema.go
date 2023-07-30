package models

import "time"

type User struct {
	Id              uint       `json:"id"`
	Uuid            string     `json:"uuid"`
	Email           string     `json:"email"`
	First_name      string     `json:"first_name"`
	Last_name       string     `json:"last_name"`
	Created_at      *time.Time `json:"created_at"`
	Favorites_count uint       `json:"favorites_count"`
	Followers_count uint       `json:"followers_count"`
	Following_count uint       `json:"following_count"`
}

type Tweet struct {
	Id             uint       `json:"id"`
	Uuid           string     `json:"uuid"`
	User_id        uint       `json:"user_id"`
	Content        string     `json:"content"`
	Created_at     *time.Time `json:"created_at"`
	Likes_count    uint       `json:"likes_count"`
	Retweets_count uint       `json:"retweets_count"`
}

type User_Followers struct {
	Id          uint `json:"id"`
	User_id     uint `json:"user_id"`
	Follower_id uint `json:"follower_id"`
}

type Tweet_likes struct {
	Id         uint       `json:"id"`
	Tweet_id   uint       `json:"tweet_id"`
	User_id    uint       `json:"user_id"`
	Created_at *time.Time `json:"created_at"`
}

type Tweet_Comments struct {
	Id         uint       `json:"id"`
	Tweet_id   uint       `json:"tweet_id"`
	User_id    uint       `json:"user_id"`
	Content    string     `json:"content"`
	Created_at *time.Time `json:"created_at"`
}

type Tweet_Redis struct {
	Id             string     `json:"id"`
	Uuid           string     `json:"uuid"`
	User_id        string     `json:"user_id"`
	Content        string     `json:"content"`
	Created_at     *time.Time `json:"created_at"`
	Likes_count    string     `json:"likes_count"`
	Retweets_count string     `json:"retweets_count"`
}
