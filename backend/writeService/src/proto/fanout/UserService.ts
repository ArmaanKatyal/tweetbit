// Original file: proto/user.proto

import type * as grpc from '@grpc/grpc-js';
import type { MethodDefinition } from '@grpc/proto-loader';
import type {
    FollowUserRequest as _fanout_FollowUserRequest,
    FollowUserRequest__Output as _fanout_FollowUserRequest__Output,
} from '../fanout/FollowUserRequest';
import type {
    FollowUserResponse as _fanout_FollowUserResponse,
    FollowUserResponse__Output as _fanout_FollowUserResponse__Output,
} from '../fanout/FollowUserResponse';

export interface UserServiceClient extends grpc.Client {
    FollowUser(
        argument: _fanout_FollowUserRequest,
        metadata: grpc.Metadata,
        options: grpc.CallOptions,
        callback: grpc.requestCallback<_fanout_FollowUserResponse__Output>
    ): grpc.ClientUnaryCall;
    FollowUser(
        argument: _fanout_FollowUserRequest,
        metadata: grpc.Metadata,
        callback: grpc.requestCallback<_fanout_FollowUserResponse__Output>
    ): grpc.ClientUnaryCall;
    FollowUser(
        argument: _fanout_FollowUserRequest,
        options: grpc.CallOptions,
        callback: grpc.requestCallback<_fanout_FollowUserResponse__Output>
    ): grpc.ClientUnaryCall;
    FollowUser(
        argument: _fanout_FollowUserRequest,
        callback: grpc.requestCallback<_fanout_FollowUserResponse__Output>
    ): grpc.ClientUnaryCall;
    followUser(
        argument: _fanout_FollowUserRequest,
        metadata: grpc.Metadata,
        options: grpc.CallOptions,
        callback: grpc.requestCallback<_fanout_FollowUserResponse__Output>
    ): grpc.ClientUnaryCall;
    followUser(
        argument: _fanout_FollowUserRequest,
        metadata: grpc.Metadata,
        callback: grpc.requestCallback<_fanout_FollowUserResponse__Output>
    ): grpc.ClientUnaryCall;
    followUser(
        argument: _fanout_FollowUserRequest,
        options: grpc.CallOptions,
        callback: grpc.requestCallback<_fanout_FollowUserResponse__Output>
    ): grpc.ClientUnaryCall;
    followUser(
        argument: _fanout_FollowUserRequest,
        callback: grpc.requestCallback<_fanout_FollowUserResponse__Output>
    ): grpc.ClientUnaryCall;
}

export interface UserServiceHandlers extends grpc.UntypedServiceImplementation {
    FollowUser: grpc.handleUnaryCall<_fanout_FollowUserRequest__Output, _fanout_FollowUserResponse>;
}

export interface UserServiceDefinition extends grpc.ServiceDefinition {
    FollowUser: MethodDefinition<
        _fanout_FollowUserRequest,
        _fanout_FollowUserResponse,
        _fanout_FollowUserRequest__Output,
        _fanout_FollowUserResponse__Output
    >;
}
