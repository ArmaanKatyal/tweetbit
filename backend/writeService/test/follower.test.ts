import chai from 'chai';
import * as sinon from 'sinon';
import { app } from '../src/app';
import { userClient } from '../src/services/user.service';
import dotenv from 'dotenv';
import request from 'supertest';
import { PrismaClient } from '@prisma/client';
import { checkIfUserExists } from '../src/helpers/verifyUser.helper';

declare global {
    namespace NodeJS {
        interface Global {}
    }
}

interface CustomNodeJsGlobal extends NodeJS.Global {
    prisma: PrismaClient;
}

declare const global: CustomNodeJsGlobal;
const prisma = global.prisma || new PrismaClient();

dotenv.config();
const test_token = process.env.TEST_TOKEN || '';

describe('/api/user', async () => {
    let userEmail: string = 'test@test.com'
    let test_user: any;
    before(async () => {
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
            }
        })
    })
    describe('[POST] /follow/:userEmail', () => {
        beforeEach(() => {
            sinon.createSandbox();
            // mock the followUser function
            sinon.mock(userClient).expects('FollowUser').returns({Success: true});
        });
        afterEach(() => {
            sinon.restore();
        });
        after(async () => {
            await prisma.$transaction([
                prisma.user_Followers.deleteMany(),
            ]);
        });
        it('should follow the user with the given email', async () => {
            const { status, body } = await request(app).post(`/api/user/follow/${userEmail}`).set(
                'Authorization', 'Bearer ' + test_token
            )
            chai.expect(status).to.equal(200);
            chai.expect(body).to.have.property('newFollower');
            chai.expect(body).to.have.property('user');
            chai.expect(body).to.have.property('user_followed');
        });
    });

    describe('[POST] /unfollow/:userEmail', () => {
        before(async () => {
            let [_, user_id] = await checkIfUserExists(userEmail);
            await prisma.user_Followers.create({
                data: {
                    follower: {
                        connect: {
                            id: test_user.id,
                        },
                    },
                    user_id: user_id!,
                }
            })
        })
        beforeEach(() => {
            sinon.createSandbox();
            // mock the followUser function
            sinon.mock(userClient).expects('UnfollowUser').returns({Success: true});
        });
        afterEach(() => {
            sinon.restore();
        });
        after(async () => {
            await prisma.$transaction([
                prisma.user_Followers.deleteMany(),
            ]);
        });
        it('should unfollow the user with the given email', async () => {
            const {status, body} = await request(app).post(`/api/user/unfollow/${userEmail}`).set(
                'Authorization', 'Bearer ' + test_token
            )
            chai.expect(status).to.equal(200);
            chai.expect(body).to.have.property('deletedFollower');
            chai.expect(body).to.have.property('user');
            chai.expect(body).to.have.property('user_unfollowed');
        });
    });
});
