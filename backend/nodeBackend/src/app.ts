import express from 'express';
import * as dotenv from 'dotenv';
import cors from 'cors';

const app = express();
dotenv.config();

const run = async () => {

}
run();

app.use(cors()) // allow cross-origin requests
app.use(express.json());    // parse requests of content-type - application/json

app.get('/', (req, res) => {
    res.send('Hello World');
});

export  { app };