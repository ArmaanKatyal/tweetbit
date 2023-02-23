// Original file: proto/tweet.proto

export interface CreateTweetRequest {
    id?: number;
    content?: string;
    userId?: number;
    uuid?: string;
    createdAt?: string;
    likesCount?: number;
    retweetsCount?: number;
}

export interface CreateTweetRequest__Output {
    id?: number;
    content?: string;
    userId?: number;
    uuid?: string;
    createdAt?: string;
    likesCount?: number;
    retweetsCount?: number;
}
