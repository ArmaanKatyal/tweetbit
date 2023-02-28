import { PrismaClient } from '@prisma/client';
import chai from 'chai';
import nodeConfig from 'config';
import dotenv from 'dotenv';
import * as sinon from 'sinon';
import request from 'supertest';
import { app } from '../src/app';
import { checkIfUserExists } from '../src/helpers/verifyUser.helper';
import { tweetClient } from '../src/services/tweet.service';

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
const test_token_2 = process.env.TEST_TOKEN_2 || '';
describe('/api/tweet', async () => {
    let test_user: any;
    before(async () => {
        // create a test user
        test_user = await prisma.user.create({
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
                id: test_user.id,
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
                // the order is important as we have foreign key constraints
                prisma.tweet_Likes.deleteMany(),
                prisma.tweet_Comments.deleteMany(),
                prisma.tweet.deleteMany(),
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
                            id: test_user.id,
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
                // the order is important as we have foreign key constraints
                prisma.tweet_Likes.deleteMany(),
                prisma.tweet_Comments.deleteMany(),
                prisma.tweet.deleteMany(),
            ]);
        });

        it('should delete a tweet', async () => {
            const response = await request(app)
                .post(`/api/tweet/delete/${tweet.id}`)
                .set('Authorization', 'Bearer ' + test_token);
            chai.expect(response.status).to.equal(200);
            chai.expect(response.body.tweet.id).to.equal(tweet.id);
            chai.expect(response.body.tweet.content).to.equal('DELETED_TEST');
        });

        it('should not delete a tweet as prisma fails', async () => {
            sinon.restore();
            const response = await request(app)
                .post(`/api/tweet/delete/${tweet.id}213`)
                .set('Authorization', 'Bearer ' + test_token);
            chai.expect(response.status).to.equal(500);
            chai.expect(response.body).to.have.property('error');
            chai.expect(response.body.error).to.equal(
                nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR')
            );
        });
    });

    describe('[POST] /api/tweet/like/:tweetId', async () => {
        let existingTweet: any;
        before(async () => {
            existingTweet = await prisma.tweet.create({
                data: {
                    uuid: process.env.TEST_UUID!,
                    content: 'TEST',
                    user: {
                        connect: {
                            id: test_user.id,
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
                // the order is important as we have foreign key constraints
                prisma.tweet_Likes.deleteMany(),
                prisma.tweet_Comments.deleteMany(),
                prisma.tweet.deleteMany(),
            ]);
        });

        it('should like a tweet', async () => {
            const res = await request(app)
                .post(`/api/tweet/like/${existingTweet.id}`)
                .set('Authorization', 'Bearer ' + test_token);
            let [_, user_id] = await checkIfUserExists('test@abc.com');
            let likedTweet = await prisma.tweet_Likes.findFirst({
                where: {
                    user_id: user_id!,
                    tweet_id: existingTweet.id,
                },
            });
            chai.expect(res.status).to.equal(200);
            chai.expect(res.body).to.have.property('id');
            chai.expect(res.body).to.have.property('user_id');
            chai.expect(res.body.likes_count).to.equal(1);
            chai.expect(likedTweet).to.not.be.null;
            chai.expect(likedTweet!.tweet_id).to.equal(existingTweet.id);
            chai.expect(likedTweet!.user_id).to.equal(user_id);
        });

        it('should not like a tweet if user already liked it', async () => {
            const res = await request(app)
                .post(`/api/tweet/like/${existingTweet.id}`)
                .set('Authorization', 'Bearer ' + test_token);
            chai.expect(res.status).to.equal(400);
            chai.expect(res.body).to.have.property('error');
            chai.expect(res.body.error).to.equal(nodeConfig.get('error_codes.TWEET_ALREADY_LIKED'));
        });
    });

    describe('[POST] /api/tweet/unlike/:tweetId', async () => {
        let existingTweet: any;
        before(async () => {
            existingTweet = await prisma.tweet.create({
                data: {
                    uuid: process.env.TEST_UUID!,
                    content: 'TEST',
                    likes_count: 2,
                    user: {
                        connect: {
                            email: 'test@abc.com',
                        },
                    },
                },
            });
            await prisma.tweet_Likes.create({
                data: {
                    user: {
                        connect: {
                            id: test_user.id,
                        },
                    },
                    tweet: {
                        connect: {
                            id: existingTweet.id,
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
                // the order is important as we have foreign key constraints
                prisma.tweet_Likes.deleteMany(),
                prisma.tweet_Comments.deleteMany(),
                prisma.tweet.deleteMany(),
            ]);
        });

        it('should unlike a tweet', async () => {
            let [_, user_id] = await checkIfUserExists('test@abc.com');
            const res = await request(app)
                .post(`/api/tweet/unlike/${existingTweet.id}`)
                .set('Authorization', 'Bearer ' + test_token);
            let dislikedTweet = await prisma.tweet_Likes.findFirst({
                where: {
                    user_id: user_id!,
                    tweet_id: existingTweet.id,
                },
            });
            chai.expect(res.status).to.equal(200);
            chai.expect(res.body.tweet).to.have.property('id');
            chai.expect(res.body.tweet).to.have.property('user_id');
            chai.expect(res.body.tweet.likes_count).to.equal(1);
            chai.expect(res.body.deleted_like.count).to.equal(1);
            chai.expect(dislikedTweet).to.be.null;
        });

        it('should not unlike a tweet if user has not liked it', async () => {
            const res = await request(app)
                .post(`/api/tweet/unlike/${existingTweet.id}`)
                .set('Authorization', 'Bearer ' + test_token);
            chai.expect(res.status).to.equal(400);
            chai.expect(res.body).to.have.property('error');
            chai.expect(res.body.error).to.equal(nodeConfig.get('error_codes.TWEET_NOT_LIKED'));
        });
    });

    describe('[POST] /api/tweet/retweet/:tweetId', async () => {
        let test_user_retweet: any;
        let existingTweet: any;
        before(async () => {
            // create a new user to retweet
            test_user_retweet = await prisma.user.create({
                data: {
                    email: 'test@abc1.com',
                    uuid: '12345',
                    first_name: 'test',
                    last_name: 'test',
                },
            });
            // create a tweet to retweet
            existingTweet = await prisma.tweet.create({
                data: {
                    uuid: process.env.TEST_UUID!,
                    content: 'TEST',
                    user: {
                        connect: {
                            id: test_user.id,
                        },
                    },
                },
            });
        });
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
                // the order is important as we have foreign key constraints
                prisma.tweet_Likes.deleteMany(),
                prisma.tweet_Comments.deleteMany(),
                prisma.tweet.deleteMany(),
                prisma.user.delete({
                    where: {
                        id: test_user_retweet.id,
                    },
                }),
            ]);
        });

        it('should retweet a tweet', async () => {
            const res = await request(app)
                .post(`/api/tweet/retweet/${existingTweet.id}`)
                .set('Authorization', 'Bearer ' + test_token_2);
            let retweetedTweet = await prisma.tweet.findFirst({
                where: {
                    user_id: test_user_retweet.id,
                },
            });

            chai.expect(res.status).to.equal(200);
            chai.expect(res.body.retweets_count).to.equal(1);
            chai.expect(retweetedTweet).to.not.be.null;
        });

        it('should not retweet a tweet if user is the owner of the tweet', async () => {
            const res = await request(app)
                .post(`/api/tweet/retweet/${existingTweet.id}`)
                .set('Authorization', 'Bearer ' + test_token);
            chai.expect(res.status).to.equal(400);
            chai.expect(res.body).to.have.property('error');
            chai.expect(res.body.error).to.equal(nodeConfig.get('error_codes.TWEET_OWNER'));
        });
        it('should not retweet a tweet if user already retweeted it', async () => {
            const res = await request(app)
                .post(`/api/tweet/retweet/${existingTweet.id}`)
                .set('Authorization', 'Bearer ' + test_token_2);
            chai.expect(res.status).to.equal(400);
            chai.expect(res.body).to.have.property('error');
            chai.expect(res.body.error).to.equal(
                nodeConfig.get('error_codes.TWEET_ALREADY_RETWEETED')
            );
        });
    });

    describe('[POST] /api/tweet/comment/:tweetId', async () => {
        let existingTweet: any;
        before(async () => {
            existingTweet = await prisma.tweet.create({
                data: {
                    uuid: process.env.TEST_UUID!,
                    content: 'TEST',
                    user: {
                        connect: {
                            id: test_user.id,
                        },
                    },
                }
            });
        });
        beforeEach(async () => {
            sinon.createSandbox();
            // mock the createTweet function
            // sinon.mock(tweetClient).expects('createTweet').returns({});
        });
        afterEach(async () => {
            sinon.restore();
        });
        after(async () => {
            await prisma.$transaction([
                // the order is important as we have foreign key constraints
                prisma.tweet_Likes.deleteMany(),
                prisma.tweet_Comments.deleteMany(),
                prisma.tweet.deleteMany(),
            ]);
        });

        it('should comment on a tweet', async () => {
            const res = await request(app).post(`/api/tweet/comment/${existingTweet.id}`).send({
                content: 'test comment',
            }).set('Authorization', 'Bearer ' + test_token);
            chai.expect(res.status).to.equal(200);
            chai.expect(res.body.tweet_id).to.equal(existingTweet.id);
            chai.expect(res.body.user_id).to.equal(test_user.id);
        });
    });

    describe('helper/verify_user', () => {
        it('should return true and the user id', async () => {
            const [exists, _] = await checkIfUserExists('test@abc.com');
            chai.expect(exists).to.be.true;
        });
        it('should return false and null', async () => {
            const [exists, id] = await checkIfUserExists('someRandomText');
            chai.expect(exists).to.be.false;
            chai.expect(id).to.be.null;
        });
    });
});
