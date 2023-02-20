// Original file: proto/tweet.proto

import type * as grpc from '@grpc/grpc-js'
import type { MethodDefinition } from '@grpc/proto-loader'
import type { CreateTweetRequest as _fanout_CreateTweetRequest, CreateTweetRequest__Output as _fanout_CreateTweetRequest__Output } from '../fanout/CreateTweetRequest';
import type { CreateTweetResponse as _fanout_CreateTweetResponse, CreateTweetResponse__Output as _fanout_CreateTweetResponse__Output } from '../fanout/CreateTweetResponse';

export interface TweetServiceClient extends grpc.Client {
  CreateTweet(argument: _fanout_CreateTweetRequest, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_fanout_CreateTweetResponse__Output>): grpc.ClientUnaryCall;
  CreateTweet(argument: _fanout_CreateTweetRequest, metadata: grpc.Metadata, callback: grpc.requestCallback<_fanout_CreateTweetResponse__Output>): grpc.ClientUnaryCall;
  CreateTweet(argument: _fanout_CreateTweetRequest, options: grpc.CallOptions, callback: grpc.requestCallback<_fanout_CreateTweetResponse__Output>): grpc.ClientUnaryCall;
  CreateTweet(argument: _fanout_CreateTweetRequest, callback: grpc.requestCallback<_fanout_CreateTweetResponse__Output>): grpc.ClientUnaryCall;
  createTweet(argument: _fanout_CreateTweetRequest, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_fanout_CreateTweetResponse__Output>): grpc.ClientUnaryCall;
  createTweet(argument: _fanout_CreateTweetRequest, metadata: grpc.Metadata, callback: grpc.requestCallback<_fanout_CreateTweetResponse__Output>): grpc.ClientUnaryCall;
  createTweet(argument: _fanout_CreateTweetRequest, options: grpc.CallOptions, callback: grpc.requestCallback<_fanout_CreateTweetResponse__Output>): grpc.ClientUnaryCall;
  createTweet(argument: _fanout_CreateTweetRequest, callback: grpc.requestCallback<_fanout_CreateTweetResponse__Output>): grpc.ClientUnaryCall;
  
}

export interface TweetServiceHandlers extends grpc.UntypedServiceImplementation {
  CreateTweet: grpc.handleUnaryCall<_fanout_CreateTweetRequest__Output, _fanout_CreateTweetResponse>;
  
}

export interface TweetServiceDefinition extends grpc.ServiceDefinition {
  CreateTweet: MethodDefinition<_fanout_CreateTweetRequest, _fanout_CreateTweetResponse, _fanout_CreateTweetRequest__Output, _fanout_CreateTweetResponse__Output>
}
