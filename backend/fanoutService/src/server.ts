import path from 'path';
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import { ProtoGrpcType } from './proto/tweet';
import { TweetServiceHandlers } from "./proto/fanout/TweetService"

const PORT = 3002;
const PROTO_FILE = './proto/tweet.proto';

const packageDef = protoLoader.loadSync(path.resolve(__dirname, `../${PROTO_FILE}`));
const grpcObj = (grpc.loadPackageDefinition(packageDef) as unknown) as ProtoGrpcType;
const tweetPackage = grpcObj.fanout;

function main() {
    const server = getServer();
    server.bindAsync(`0.0.0.0:${PORT}`, grpc.ServerCredentials.createInsecure(), (err) => {
        if (err) {
            console.error(err);
            return;
        }
        server.start();
        console.log(`Server running at port ${PORT}`);
    });
}

function getServer() {
    const server = new grpc.Server();
    server.addService(tweetPackage.TweetService.service, {
        "CreateTweet": (req, res) => {
            console.log(req.request);
            res(null, {})
        }
    } as TweetServiceHandlers)

    return server;
}

main();