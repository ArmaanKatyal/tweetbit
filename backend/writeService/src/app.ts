import express, { Response, Request, NextFunction } from 'express';
import * as dotenv from 'dotenv';
import cors from 'cors';
import cookieParser from 'cookie-parser';

const app = express();
dotenv.config();

app.use(cors());
app.use(express.json());
app.use(cookieParser());

app.get('/', (req: Request, res: Response, next: NextFunction) => {
    res.status(200).json({ message: 'Hello World!' });
});

export { app };