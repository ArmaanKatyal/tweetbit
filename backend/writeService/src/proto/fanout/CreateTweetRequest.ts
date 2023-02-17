// Original file: proto/tweet.proto

export interface CreateTweetRequest {
    id?: string;
    content?: string;
    userId?: string;
    uuid?: string;
    createdAt?: string;
    likesCount?: number;
    retweetsCount?: number;
}

export interface CreateTweetRequest__Output {
    id?: string;
    content?: string;
    userId?: string;
    uuid?: string;
    createdAt?: string;
    likesCount?: number;
    retweetsCount?: number;
}
