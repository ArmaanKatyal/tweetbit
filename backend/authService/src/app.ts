import express, { Request, Response, NextFunction } from 'express';
import * as dotenv from 'dotenv';
import cors from 'cors';
import mongoose, { connect } from 'mongoose';
import cookieParser from 'cookie-parser';
import logger from './utils/log.util';
import expressPino from 'express-pino-logger';
import pinoHttp from 'pino-http';
import { authRouter } from './routes/auth.route';
import { register } from 'prom-client';
import {
    IncHttpTransaction,
    MetricsCode,
    MetricsMethod,
    ObserveHttpResponseTime,
} from './internal/prometheus';

const app = express();
dotenv.config();
mongoose.set('strictQuery', true);

const run = async () => {
    if (process.env.NODE_ENV === 'dev') {
        await connect(process.env.MONGO_URI_DEV!);
    } else if (process.env.NODE_ENV === 'prod') {
        await connect(process.env.MONGO_URI_PROD!);
    } else {
        console.error('No version specified');
    }
};

if (process.env.NODE_ENV !== 'test') {
    run();
}

app.use(cors()); // allow cross-origin requests
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

// Add headers before the routes are defined
app.use(function (_: Request, res: Response, next: NextFunction) {
    res.header('Content-Type', 'application/json'); // set default content type
    res.header('Access-Control-Allow-Credentials', 'true'); // allow cookies
    res.header(
        'Access-Control-Allow-Headers',
        'Origin, X-Requested-With, Content-Type, Accept, Authorization'
    );
    next();
});

app.use(express.json()); // parse requests of content-type - application/json
app.use(cookieParser());

app.get('/', (_, res) => {
    res.status(200).json({
        status: mongoose.connection.readyState,
        database: mongoose.connection.name,
    });
});

app.get('/health', (_: Request, res: Response) => {
    let start = Date.now();
    res.status(200).send('OK');
    IncHttpTransaction(MetricsCode.Ok, MetricsMethod.Get);
    ObserveHttpResponseTime(MetricsCode.Ok, MetricsMethod.Get, Date.now() - start);
});

app.get('/metrics', async (_: Request, res: Response) => {
    try {
        res.set('Content-Type', register.contentType);
        res.end(await register.metrics());
    } catch (ex) {
        res.status(500).end(ex);
    }
});

app.use('/api/auth', authRouter);

export { app };
