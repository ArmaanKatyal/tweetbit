import express, { Response, Request, NextFunction } from 'express';
import * as dotenv from 'dotenv';
import cors from 'cors';
import cookieParser from 'cookie-parser';
import { Server, ServerCredentials } from '@grpc/grpc-js';
import nodeConfig from 'config';
import { tweetRouter } from './routes/tweet.route';
import { followRouter } from './routes/follower.route';
import expressPino from 'express-pino-logger';
import pinoHttp from 'pino-http';
import logger from './utils/log.util';
import { initTracer } from './utils/opentelemetry.util';

const app = express();
dotenv.config();

app.use(cors());
if (process.env.NODE_ENV !== 'test') {
    app.use(
        pinoHttp({
            transport: {
                target: 'pino-pretty',
                options: {
                    levelFirst: true,
                    colorize: true,
                    translateTime: true,
                },
            },
        })
    );
}
app.use(
    expressPino({
        logger,
        autoLogging: true,
    })
);
app.use(express.json());
app.use(cookieParser());

try {
    initTracer();
} catch (err) {
    console.log(err);
}

const server = new Server();
if (process.env.NODE_ENV !== 'test') {
    // Start the gRPC server
    server.bindAsync(nodeConfig.get('grpc.port'), ServerCredentials.createInsecure(), () => {
        server.start();
        console.log(`gRPC server is listening on port ${nodeConfig.get('grpc.port')}`);
    });
}

app.get('/health', (_: Request, res: Response) => {
    res.status(200).send('OK');
});

app.use('/api/tweet', tweetRouter);
app.use('/api/user', followRouter);

export { app };
