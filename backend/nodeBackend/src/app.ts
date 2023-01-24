import express from 'express';
import * as dotenv from 'dotenv';
import cors from 'cors';
import { connect } from 'mongoose';

// import routes
import { authRouter } from './routes/auth';

const app = express();
dotenv.config();

const run = async () => {
    if (process.env.VERSION === 'dev') {
        await connect(process.env.MONGO_URI_DEV!);
    } else if (process.env.VERSION === 'prod') {
        await connect(process.env.MONGO_URI_PROD!);
    }
    else {
        console.error('No version specified');
    }
}
run();

app.use(cors()) // allow cross-origin requests
app.use(express.json());    // parse requests of content-type - application/json

app.get('/', (req, res) => {
    res.send('Hello World');
});

app.use('/api/auth', authRouter);

export  { app };