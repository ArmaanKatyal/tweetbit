import express, { Request, Response, NextFunction } from 'express';
import * as dotenv from 'dotenv';
import cors from 'cors';
import { connect } from 'mongoose';
import cookieParser from 'cookie-parser';
import logger from './utils/log.util';
import expressPino from 'express-pino-logger';
import pinoHttp from 'pino-http';

// import routes
import { authRouter } from './routes/auth.route';

const app = express();
dotenv.config();

const run = async () => {
    if (process.env.VERSION === 'dev') {
        await connect(process.env.MONGO_URI_DEV!);
    } else if (process.env.VERSION === 'prod') {
        await connect(process.env.MONGO_URI_PROD!);
    } else {
        console.error('No version specified');
    }
};
run();

app.use(cors()); // allow cross-origin requests
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
app.use(
    expressPino({
        logger,
        autoLogging: true,
    })
);

app.use(function (req: Request, res: Response, next: NextFunction) {
    res.header('Content-Type', 'application/json');
    res.header('Access-Control-Allow-Credentials', 'true');
    res.header(
        'Access-Control-Allow-Headers',
        'Origin, X-Requested-With, Content-Type, Accept, Authorization'
    );
    next();
});

app.use(express.json()); // parse requests of content-type - application/json
app.use(cookieParser());

app.get('/', (req, res) => {
    res.send('Hello World');
});

app.use('/api/auth', authRouter);

export { app };
