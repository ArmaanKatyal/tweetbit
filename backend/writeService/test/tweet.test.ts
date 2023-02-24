// import { afterAll, beforeEach, describe, expect, it } from "vitest";
import request from 'supertest'
import {app} from '../src/app'
import dotenv from 'dotenv';
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
            }
        });
    });
    after(async () => {
        await prisma.user.delete({
                where: {
                    email: 'test@abc.com',
                }
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

            const {status, body} = await request(app).post('/api/tweet/create').send({
                content: 'TEST'
            }).set('Authorization', 'Bearer ' + test_token);
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
    });
    
    describe('[POST] /delete/:tweetId', async () => {});

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
