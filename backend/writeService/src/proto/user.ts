import type * as grpc from '@grpc/grpc-js';
import type { MessageTypeDefinition } from '@grpc/proto-loader';

import type {
    UserServiceClient as _fanout_UserServiceClient,
    UserServiceDefinition as _fanout_UserServiceDefinition,
} from './fanout/UserService';

type SubtypeConstructor<Constructor extends new (...args: any) => any, Subtype> = {
    new (...args: ConstructorParameters<Constructor>): Subtype;
};

export interface ProtoGrpcType {
    fanout: {
        FollowUserRequest: MessageTypeDefinition;
        FollowUserResponse: MessageTypeDefinition;
        UserService: SubtypeConstructor<typeof grpc.Client, _fanout_UserServiceClient> & {
            service: _fanout_UserServiceDefinition;
        };
    };
}
