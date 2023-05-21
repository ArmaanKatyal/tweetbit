process.env.NODE_ENV = 'test';

import chai from 'chai';
import chaiHttp from 'chai-http';
import { app } from '../app';
import { MongoMemoryServer } from 'mongodb-memory-server';
import mongoose from 'mongoose';
import nodeConfig from 'config';
import { User } from '../models/user.model';
import jwt from 'jsonwebtoken';

let accessToken: string;
let refreshToken: string;

chai.use(chaiHttp);
describe('Auth Test', () => {
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
    it('should register a new user', async () => {
        const res = await chai.request(app).post('/api/auth/register').send({
            email: 'test@test.com',
            password: 'test',
            username: 'test',
            first_name: 'test',
            last_name: 'test',
        });

        chai.expect(res.status).to.equal(200);
        chai.expect(res.body).to.have.property('message');
    });
    it('should not register a new user with the same email', async () => {
        const res = await chai.request(app).post('/api/auth/register').send({
            email: 'test@test.com',
            password: 'test',
            username: 'test1',
            first_name: 'test',
            last_name: 'test',
        });

        chai.expect(res.status).to.equal(400);
        chai.expect(res.body).to.have.property('error');
        chai.expect(res.body.error).to.equal(nodeConfig.get('error_codes.USER_ALREADY_EXISTS'));
    });
    it('should not register a new user with invalid input', async () => {
        const res = await chai.request(app).post('/api/auth/register').send({
            email: 'test1@gmail.com',
            password: 12345,
            username: 'test2',
            first_name: 'test',
            last_name: 'test',
        });

        chai.expect(res.status).to.equal(400);
        chai.expect(res.body).to.have.property('error');
        chai.expect(res.body.error).to.equal(nodeConfig.get('error_codes.INVALID_INPUT'));
    });
    it('should cause an internal server error', async () => {
        const res = await chai.request(app).post('/api/auth/register').send({
            email: 'test2@test.com',
            password: 'test',
            username: 'test',
            first_name: 'test',
            last_name: 'test',
        });

        chai.expect(res.status).to.equal(500);
        chai.expect(res.body).to.have.property('error');
        chai.expect(res.body.error).to.equal(nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR'));
    });
    it('should login a user successfully', async () => {
        const res = await chai.request(app).post('/api/auth/login').send({
            email: 'test@test.com',
            password: 'test',
        });

        // NOTE: In Future the response can contain more data which is why we are not checking for the exact response
        // but rather checking for the status code and the presence of the access_token and user object
        // but we might want add more checks in the future for the additional data

        chai.expect(res.status).to.equal(200);
        chai.expect(res.body).to.have.property('access_token');
        chai.expect(res.body).to.have.property('user');
        chai.expect(res.body.user.email).to.equal('test@test.com');
        chai.expect(res.body.user.username).to.equal('test');
        chai.expect(res.body.user.first_name).to.equal('test');
        chai.expect(res.body.user.last_name).to.equal('test');

        // Save the access token for future tests
        accessToken = res.body.access_token;
        refreshToken = res.body.refresh_token;
    });
    it('should login and contain a refresh_token as a cookie', async () => {
        const res = await chai.request(app).post('/api/auth/login').send({
            email: 'test@test.com',
            password: 'test',
        });

        chai.expect(res.status).to.equal(200);
        chai.expect(res.body).to.have.property('access_token');
        chai.expect(res.header).to.have.property('set-cookie');
        chai.expect(res.header['set-cookie'][0]).to.contain('refresh_token');
    });
    it('should not login a user with invalid input', async () => {
        const res = await chai.request(app).post('/api/auth/login').send({
            email: 'test@test.com',
            password: 12345,
        });

        chai.expect(res.status).to.equal(400);
        chai.expect(res.body).to.have.property('error');
        chai.expect(res.body.error).to.equal(nodeConfig.get('error_codes.INVALID_INPUT'));
    });
    it('should not login as a user that does not exist', async () => {
        const res = await chai.request(app).post('/api/auth/login').send({
            email: 'abc@abc.com',
            password: 'test',
        });

        chai.expect(res.status).to.equal(400);
        chai.expect(res.body).to.have.property('error');
        chai.expect(res.body.error).to.equal(nodeConfig.get('error_codes.USER_NOT_FOUND'));
    });
    it('should not login as a user with incorrect password', async () => {
        const res = await chai.request(app).post('/api/auth/login').send({
            email: 'test@test.com',
            password: 'test1',
        });

        chai.expect(res.status).to.equal(400);
        chai.expect(res.body).to.have.property('error');
        chai.expect(res.body.error).to.equal(nodeConfig.get('error_codes.INVALID_PASSWORD'));
    });
    it('should not login as auth exists but user does not', async () => {
        const findUser = await User.findOne({ email: 'test@test.com' });
        let { _id, ...newPayload } = (findUser as any)._doc;
        const deletedUser = await User.deleteOne({ email: 'test@test.com' });
        const res = await chai.request(app).post('/api/auth/login').send({
            email: 'test@test.com',
            password: 'test',
        });

        chai.expect(res.status).to.equal(400);
        chai.expect(res.body).to.have.property('error');
        chai.expect(res.body.error).to.equal(nodeConfig.get('error_codes.USER_NOT_FOUND'));
        chai.expect(deletedUser.deletedCount).to.equal(1);

        await User.create(newPayload);
    });
    it('should logout a user successfully', async () => {
        const res = await chai.request(app).post('/api/auth/logout').send({
            email: 'test@test.com',
            password: 'test',
        });

        chai.expect(res.status).to.equal(200);
        chai.expect(res.body).to.have.property('message');
    });
    it('should refresh a token successfully', async () => {
        const res = await chai
            .request(app)
            .get('/api/auth/refresh')
            .set('Cookie', `refresh_token=${refreshToken}`);

        chai.expect(res.status).to.equal(200);
        chai.expect(res.body).to.have.property('access_token');
    });
    it('should not refresh a token with an invalid refresh token', async () => {
        const res = await chai
            .request(app)
            .get('/api/auth/refresh')
            .set('Cookie', 'refresh_token=invalid_token');

        chai.expect(res.status).to.equal(401);
        chai.expect(res.body).to.have.property('error');
        chai.expect(res.body.error).to.equal(nodeConfig.get('error_codes.INVALID_TOKEN'));
    });
    it('should not refresh a token with an expired refresh token', async () => {
        let expiredRefreshToken = await jwt.sign(
            { email: 'test@test.com', type: 'refresh' },
            nodeConfig.get('jwt.refresh_token_secret'),
            { expiresIn: '1ms' }
        );
        const res = await chai
            .request(app)
            .get('/api/auth/refresh')
            .set('Cookie', `refresh_token=${expiredRefreshToken}`);

        chai.expect(res.status).to.equal(401);
        chai.expect(res.body).to.have.property('error');
        chai.expect(res.body.error).to.equal(nodeConfig.get('error_codes.TOKEN_EXPIRED'));
    });
    it('should not refresh a token as access_token was given inplace of refresh_token', async () => {
        const res = await chai
            .request(app)
            .get('/api/auth/refresh')
            .set('Cookie', `refresh_token=${accessToken}`);

        chai.expect(res.status).to.equal(401);
        chai.expect(res.body).to.have.property('error');
        chai.expect(res.body.error).to.equal(nodeConfig.get('error_codes.INVALID_TOKEN'));
    });
});
