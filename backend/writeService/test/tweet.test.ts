import request from 'supertest';
import { app } from '../src/app';
import dotenv from 'dotenv';
import nodeConfig from 'config';
import { PrismaClient } from '@prisma/client';
import chai from 'chai';
import * as sinon from 'sinon';
import { tweetClient } from '../src/services/tweet.service';
import { checkIfUserExists } from '../src/helpers/verifyUser.helper';

declare global {
    namespace NodeJS {
        interface Global {}
    }
}

// add prisma to the NodeJS global type
interface CustomNodeJsGlobal extends NodeJS.Global {
    prisma: PrismaClient;
}

declare const global: CustomNodeJsGlobal;
const prisma = global.prisma || new PrismaClient();

dotenv.config();
const test_token = process.env.TEST_TOKEN || '';
describe('/api/tweet', async () => {
    before(async () => {
        // create a test user
        await prisma.user.create({
            data: {
                email: 'test@abc.com',
                uuid: process.env.TEST_UUID!,
                first_name: 'test',
                last_name: 'test',
            },
        });
    });
    after(async () => {
        await prisma.user.delete({
            where: {
                email: 'test@abc.com',
            },
        });
    });
    describe('[POST] /api/tweet/create', () => {
        beforeEach(async () => {
            sinon.createSandbox();
            // mock the createTweet function
            sinon.mock(tweetClient).expects('createTweet').returns({});
        });
        afterEach(async () => {
            sinon.restore();
        });
        after(async () => {
            await prisma.$transaction([
                prisma.tweet.deleteMany(),
                prisma.tweet_Likes.deleteMany(),
                prisma.tweet_Comments.deleteMany(),
            ]);
        });
        it('should return 201 and the created tweet', async () => {
            const { status, body } = await request(app)
                .post('/api/tweet/create')
                .send({
                    content: 'TEST',
                })
                .set('Authorization', 'Bearer ' + test_token);
            chai.expect(status).to.equal(201);
            chai.expect(body).to.have.property('id');
            chai.expect(body).to.have.property('content');
            chai.expect(body).to.have.property('user_id');
            chai.expect(body).to.have.property('uuid');
            chai.expect(body).to.have.property('created_at');
            chai.expect(body).to.have.property('likes_count');
            chai.expect(body).to.have.property('retweets_count');
            chai.expect(body.content).to.equal('TEST');
        });

        it('should not create a tweet as user is not found', async () => {
            let wrong_user_token = process.env.USER_NOT_FOUND_TOKEN || '';
            const { status, body } = await request(app)
                .post('/api/tweet/create')
                .send({
                    content: 'TEST',
                })
                .set('Authorization', `Bearer ${wrong_user_token}`);
            chai.expect(status).to.equal(400);
            chai.expect(body).to.have.property('error');
            chai.expect(body.error).to.equal(nodeConfig.get('error_codes.USER_NOT_FOUND'));
        });

        it('should create a tweet but grpc call fails', async () => {
            sinon.restore();
            // mock the createTweet function that throws an error
            sinon.mock(tweetClient).expects('createTweet').throws(new Error('some error'));
            const { status, body } = await request(app)
                .post('/api/tweet/create')
                .send({
                    content: 'TEST',
                })
                .set('Authorization', 'Bearer ' + test_token);
            chai.expect(status).to.equal(500);
            chai.expect(body).to.have.property('error');
            chai.expect(body.error).to.equal(nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR'));
        });
    });

    describe('[POST] /delete/:tweetId', async () => {
        let tweet: any;
        before(async () => {
            // create a fake tweet
            tweet = await prisma.tweet.create({
                data: {
                    uuid: process.env.TEST_UUID!,
                    content: 'DELETED_TEST',
                    user: {
                        connect: {
                            email: 'test@abc.com',
                        },
                    },
                },
            });
        });
        beforeEach(async () => {
            sinon.createSandbox();
            // mock the deleteTweet function
            // sinon.mock(tweetClient).expects('deleteTweet').returns({});
        });
        afterEach(async () => {
            sinon.restore();
        });
        after(async () => {
            await prisma.$transaction([
                prisma.tweet.deleteMany(),
                prisma.tweet_Likes.deleteMany(),
                prisma.tweet_Comments.deleteMany(),
            ]);
        });

        it('should delete a tweet', async () => {
            const response = await request(app)
                .post(`/api/tweet/delete/${tweet.id}`)
                .set('Authorization', 'Bearer ' + test_token);
            chai.expect(response.status).to.equal(200);
            chai.expect(response.body.id).to.equal(tweet.id);
            chai.expect(response.body.content).to.equal('DELETED_TEST');
        });

        it('should not delete a tweet as prisma fails', async () => {
            sinon.restore();
            const response = await request(app)
                .post(`/api/tweet/delete/${tweet.id}213`)
                .set('Authorization', 'Bearer ' + test_token);
            chai.expect(response.status).to.equal(500);
            chai.expect(response.body).to.have.property('error');
            chai.expect(response.body.error).to.equal(nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR'));
        });
    });

    describe('[POST] /api/tweet/like/:tweetId', async () => {});

    describe('[POST] /api/tweet/unlike/:tweetId', async () => {});

    describe('[POST] /api/tweet/retweet/:tweetId', async () => {});

    describe('[POST] /api/tweet/comment/:tweetId', async () => {});

    describe('helper/verify_user', () => {
        it('should return true and the user id', async () => {
            const [exists, id] = await checkIfUserExists('test@abc.com');
            chai.expect(exists).to.be.true;
        });
        it('should return false and null', async () => {
            const [exists, id] = await checkIfUserExists('someRandomText');
            chai.expect(exists).to.be.false;
            chai.expect(id).to.be.null;
        });
    });
});
