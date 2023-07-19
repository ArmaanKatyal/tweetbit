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
