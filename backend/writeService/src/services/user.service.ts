import path from 'path';
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import { ProtoGrpcType } from '../proto/user';
import nodeConfig from 'config';

const PROTO_FILE = './proto/user.proto';

const packageDef = protoLoader.loadSync(path.resolve(__dirname, `../../${PROTO_FILE}`));
const grpcObj = grpc.loadPackageDefinition(packageDef) as unknown as ProtoGrpcType;

export const userClient = new grpcObj.fanout.UserService(
    `0.0.0.0:${nodeConfig.get('grpc.fanout.port')}`,
    grpc.credentials.createInsecure()
);
