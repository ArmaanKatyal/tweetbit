import type * as grpc from '@grpc/grpc-js';
import type { MessageTypeDefinition } from '@grpc/proto-loader';

import type { TweetServiceClient as _fanout_TweetServiceClient, TweetServiceDefinition as _fanout_TweetServiceDefinition } from './fanout/TweetService';

type SubtypeConstructor<Constructor extends new (...args: any) => any, Subtype> = {
  new(...args: ConstructorParameters<Constructor>): Subtype;
};

export interface ProtoGrpcType {
  fanout: {
    CreateTweetRequest: MessageTypeDefinition
    CreateTweetResponse: MessageTypeDefinition
    TweetService: SubtypeConstructor<typeof grpc.Client, _fanout_TweetServiceClient> & { service: _fanout_TweetServiceDefinition }
  }
}

