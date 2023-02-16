import express, { Response, Request, NextFunction } from 'express';
import * as dotenv from 'dotenv';
import cors from 'cors';
import cookieParser from 'cookie-parser';
import { Server, ServerCredentials } from '@grpc/grpc-js';
import nodeCofiig from 'config';
import { tweetRouter } from './routes/tweet.route';

const app = express();
dotenv.config();

app.use(cors());
app.use(express.json());
app.use(cookieParser());

const server = new Server();
if (process.env.NODE_ENV !== 'test') {
    // Start the gRPC server
    server.bindAsync(nodeCofiig.get('grpc.port'), ServerCredentials.createInsecure(), () => {
        server.start();
        console.log(`gRPC server is listening on port ${nodeCofiig.get('grpc.port')}`);
    });
}

app.get('/', (req: Request, res: Response, next: NextFunction) => {
    res.status(200).json({ message: 'Hello World' });
});

app.use('/api/tweet', tweetRouter);

export { app };