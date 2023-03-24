package models

type ITweet struct {
	Id            string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Content       string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	UserId        string `protobuf:"bytes,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Uuid          string `protobuf:"bytes,4,opt,name=uuid,proto3" json:"uuid,omitempty"`
	CreatedAt     string `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	LikesCount    string `protobuf:"varint,6,opt,name=likes_count,json=likesCount,proto3" json:"likes_count,omitempty"`
	RetweetsCount string `protobuf:"varint,7,opt,name=retweets_count,json=retweetsCount,proto3" json:"retweets_count,omitempty"`
}
