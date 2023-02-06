process.env.NODE_ENV = 'test';

import chai from 'chai';
import chaiHttp from 'chai-http';
import { app } from '../app';
import {MongoMemoryServer} from 'mongodb-memory-server';
import mongoose from 'mongoose';

chai.use(chaiHttp);
describe('test setup', () => {
    let mongoServer: MongoMemoryServer;
    before(async () => {
        mongoServer = await MongoMemoryServer.create();
        const mongoUri = await mongoServer.getUri();
        await mongoose.connect(mongoUri);
    });
    after(async () => {
        await mongoose.disconnect();
        await mongoServer.stop();
    });
    it('should pass', () => {
        chai.request(app).get('/').then((res) => {
            chai.expect(res.status).to.equal(200);
        })
    });
    
});