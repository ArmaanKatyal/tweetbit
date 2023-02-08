import express, { Request, Response, NextFunction } from 'express';
import * as dotenv from 'dotenv';
import cors from 'cors';
import mongoose, { connect } from 'mongoose';
import cookieParser from 'cookie-parser';
import logger from './utils/log.util';
import expressPino from 'express-pino-logger';
import pinoHttp from 'pino-http';

// import routes
import { authRouter } from './routes/auth.route';

const app = express();
dotenv.config();
mongoose.set('strictQuery', true);

const run = async () => {
    if (process.env.VERSION === 'dev') {
        await connect(process.env.MONGO_URI_DEV!);
    } else if (process.env.VERSION === 'prod') {
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
app.use(function (req: Request, res: Response, next: NextFunction) {
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

app.get('/', (req, res) => {
    res.status(200).json({
        status: mongoose.connection.readyState,
        database: mongoose.connection.name,
    });
});

app.use('/api/auth', authRouter);

export { app };
