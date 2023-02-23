// import { afterAll, beforeEach, describe, expect, it } from "vitest";
import request from 'supertest'
import {app} from '../src/app'
import dotenv from 'dotenv';
import { PrismaClient } from '@prisma/client';
import chai from 'chai';
import * as sinon from 'sinon';
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
describe('/api/tweet', async () => {
    describe('[POST] /api/tweet/create', () => {
        beforeEach(async () => {
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
});
