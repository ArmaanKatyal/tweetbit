import {app} from './app';
import * as dotenv from 'dotenv';
dotenv.config();

app.listen(process.env.PORT || 3000) // start the server